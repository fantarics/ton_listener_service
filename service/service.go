package service

import (
	"encoding/json"
	"fmt"
	"github.com/BoryslavGlov/logrusx"
	"github.com/r3labs/sse/v2"
	"github.com/tonkeeper/tongo"
	"github.com/xssnick/tonutils-go/tlb"
	tele "gopkg.in/telebot.v3"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
	streamClient "tonListener/internal/stream-client"
	tonClient "tonListener/internal/ton-client"
	"tonListener/pkg/repository"
	client "tonListener/pkg/telegram"
)

type mainAccount struct {
	address tongo.AccountID
}

type transactionEventData struct {
	AccountID tongo.AccountID `json:"account_id"`
	Lt        uint64          `json:"lt"`
	TxHash    string          `json:"tx_hash"`
}

func NewTonListener(
	tonApi tonClient.TonClient,
	stream *streamClient.StreamClient,
	telegram client.TelegramAPI,
	repo repository.Repository,
	logx logrusx.Logging,
	mainAddr string,

) *TonListener {

	mainaddr, err := tongo.ParseAccountID(mainAddr)
	if err != nil {
		log.Fatal(err)
	}

	return &TonListener{tonApi: tonApi, stream: stream, telegram: telegram, repo: repo, logx: logx, mainAcc: &mainAccount{address: mainaddr}}
}

type TonListener struct {
	tonApi   tonClient.TonClient
	stream   *streamClient.StreamClient
	telegram client.TelegramAPI
	repo     repository.Repository
	logx     logrusx.Logging
	mainAcc  *mainAccount
}

func (tl *TonListener) Start() error {

	err := tl.stream.Stream.Subscribe("", func(msg *sse.Event) {
		switch string(msg.Event) {
		case "message":
			var data transactionEventData

			if err := json.Unmarshal(msg.Data, &data); err != nil {
				tl.logx.Error("json.Unmarshal() failed in Start",
					logrusx.LogField{Key: "context", Value: err},
					logrusx.LogField{Key: "msg", Value: fmt.Sprintf("%s", msg.Data)},
				)
				return
			}
			wallets, err := tl.repo.GetWallets("TON")
			if err != nil {
				tl.logx.Error("Error while trying to GetWallets in onTransaction",
					logrusx.LogField{Key: "context", Value: err},
				)
			}

			go tl.onTransaction(data, wallets)

		}
	})
	return err
}

func (tl *TonListener) onTransaction(data transactionEventData, wallets []repository.Wallet) {

	tl.logx.Info("New transaction",
		logrusx.LogField{Key: "accountID", Value: data.AccountID.ToRaw()},
		logrusx.LogField{Key: "lt", Value: data.Lt},
		logrusx.LogField{Key: "txHash", Value: data.TxHash},
	)

	if tl.mainAcc.address.ToRaw() == data.AccountID.ToRaw() {

		tx := tl.tonApi.GetTransactionByHash(data.TxHash)

		if len(tx.OutMsgs) == 0 {
			return
		}

		if tx.OutMsgs[0].Value == 0 {
			return
		}

		outMsg := tx.OutMsgs[0]

		datetime := time.Unix(tx.Utime, 0)

		var usrId uint64

		uId := strings.Split(outMsg.DecodedBody.Text, ";")
		if len(uId) == 3 {
			userId, err := strconv.ParseInt(uId[2], 10, 64)
			if err != nil {
				tl.logx.Error("Error while trying to convert uId",
					logrusx.LogField{Key: "context", Value: err},
					logrusx.LogField{Key: "decodedBody", Value: outMsg.DecodedBody.Text},
				)
			}
			usrId = uint64(userId)

		}

		humanDestination := tongo.MustParseAccountID(tx.InMsg.Destination.Address).ToHuman(false, false)

		trans := &repository.Transaction{
			TxHash:             data.TxHash,
			Hash:               outMsg.RawBody,
			DestinationAddress: humanDestination,
			IncomingValue:      0,
			WalletID:           outMsg.Source.Address,
			TransactionType:    "WITHDRAW",
			CreatedTime:        datetime,
			SendedAmount:       uint64(outMsg.Value),
			UserID:             usrId,
		}

		if err := tl.repo.CreateTrx(trans); err != nil {
			tl.logx.Error("Error while trying to CreateTrx in onTransaction",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "transaction", Value: fmt.Sprintf("%+v", trans)},
			)
		}

		user, err := tl.repo.GetUser(usrId)
		if err != nil {
			tl.logx.Error("Error while trying to CreateTrx in onTransaction",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "userId", Value: usrId},
			)
		}

		var lang string

		if user.BotLang.Valid {
			lang = user.BotLang.String
		} else {
			lang = "en"
		}

		txt := strings.Replace(transferredMap[lang], "{floatAmount}", tlb.MustFromNano(big.NewInt(int64(outMsg.Value*10)), 10).String(), 1)
		txt = strings.Replace(txt, "{currency}", "TON", 1)
		txt = strings.Replace(txt, "{destination}", humanDestination, 1)
		replyMarkup := &tele.ReplyMarkup{ResizeKeyboard: true}
		cancelButton := replyMarkup.Data("OK", "cancel", "cancel")
		replyMarkup.Inline(replyMarkup.Row(cancelButton))

		_, err = tl.telegram.SendMessage(user.ID, txt, replyMarkup)
		if err != nil {
			tl.logx.Error("Error while trying to SendMessage",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "userID", Value: user.ID},
				logrusx.LogField{Key: "text", Value: txt},
			)
		}

	}

	ok, wallet := contains(wallets, data.AccountID)
	if ok {

		tx := tl.tonApi.GetTransactionByHash(data.TxHash)

		if tx.InMsg.Value == 0 {
			return
		}

		datetime := time.Unix(tx.Utime, 0)
		humanDestination := tongo.MustParseAccountID(tx.InMsg.Destination.Address).ToHuman(false, false)

		trans := &repository.Transaction{
			TxHash:             data.TxHash,
			Hash:               tx.InMsg.RawBody,
			DestinationAddress: humanDestination,
			IncomingValue:      uint64(tx.InMsg.Value),
			WalletID:           wallet.Address,
			TransactionType:    "DEPOSIT",
			CreatedTime:        datetime,
			SendedAmount:       0,
			UserID:             wallet.UserID,
		}
		if err := tl.repo.CreateTrx(trans); err != nil {
			tl.logx.Error("Error while trying to CreateTrx in onTransaction",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "transaction", Value: fmt.Sprintf("%+v", trans)},
			)
		}
		amount := int64(tx.InMsg.Value)

		err := tl.repo.AddBalance(wallet.UserID, wallet.CoinID, amount)
		if err != nil {
			tl.logx.Error("Error while trying to AddBalance in onTransaction",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "userID", Value: wallet.UserID},
				logrusx.LogField{Key: "coin", Value: "TON"},
				logrusx.LogField{Key: "amount", Value: amount},
			)
		}

		jsonb := map[string]interface{}{
			"type":     "DEPOSIT",
			"delta":    amount,
			"currency": "TON",
		}

		body, _ := json.Marshal(jsonb)

		logg := &repository.Log{
			Action:     "change_balance",
			UserID:     wallet.UserID,
			JSONObject: body,
		}

		err = tl.repo.CreateLog(logg)
		if err != nil {
			tl.logx.Error("Error while trying to CreateLog in onTransaction",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "log", Value: fmt.Sprintf("%+v", logg)},
			)
		}

		replyMarkup := &tele.ReplyMarkup{ResizeKeyboard: true}
		cancelButton := replyMarkup.Data("OK", "cancel", "cancel")
		replyMarkup.Inline(replyMarkup.Row(cancelButton))

		var lang string

		if wallet.User.BotLang.Valid {
			lang = wallet.User.BotLang.String
		} else {
			lang = "en"
		}

		txt := strings.Replace(receivedMap[lang], "{amount}", tlb.MustFromNano(big.NewInt(int64(tx.InMsg.Value*10)), 10).String(), 1)
		txt = strings.Replace(txt, "{currency}", "TON", 1)

		_, err = tl.telegram.SendMessage(wallet.UserID, txt, replyMarkup)
		if err != nil {
			tl.logx.Error("Error while trying to SendMessage",
				logrusx.LogField{Key: "context", Value: err},
				logrusx.LogField{Key: "userID", Value: wallet.UserID},
				logrusx.LogField{Key: "text", Value: txt},
			)
		}

		_, err = tl.tonApi.TransferTokens(wallet.Mnemonics, uint64(amount), fmt.Sprintf("userID:%v | amount %v", wallet.UserID, amount))
		if err != nil {
			tl.logx.Error("Error while trying to TransferTokens", logrusx.LogField{Key: "context", Value: err})
		}

	}

}

func contains(wallets []repository.Wallet, accountId tongo.AccountID) (bool, repository.Wallet) {
	for _, wallet := range wallets {
		if wallet.IsActive {
			addr, err := tongo.ParseAccountID(wallet.Address)
			if err != nil {
				continue
			}
			if addr.ToRaw() == accountId.ToRaw() {
				return true, wallet
			}
		}
	}
	return false, repository.Wallet{}

}

package tonClient

import (
	"context"
	"fmt"
	"github.com/BoryslavGlov/logrusx"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
)

func (api *api) TransferTokens(mnemonics []string, tons uint64, message string) (string, error) {

	w, err := wallet.FromSeed(api.tonapi, mnemonics, wallet.V3R2)
	if err != nil {
		api.logx.Error("Error while trying to wallet.FromSeed in TransferTokens",
			logrusx.LogField{Key: "context", Value: err},
		)
		return "", err
	}

	addr := address.MustParseAddr(api.config.MainAddress)

	balance := tlb.MustFromNano(big.NewInt(int64(tons*10)), 10)

	msg, err := generateMsg(message)
	if err != nil {
		return "", err
	}
	transfer := &wallet.Message{
		Mode: 1 + 2 + 128,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      false,
			DstAddr:     addr,
			Amount:      balance,
			Body:        msg,
		},
	}

	tx, _, err := w.SendWaitTransaction(context.Background(), transfer)
	if err != nil {
		api.logx.Fatal("Error while trying to SendWaitTransaction",
			logrusx.LogField{Key: "context", Value: err},
			logrusx.LogField{Key: "transfer", Value: fmt.Sprintf("%+v", transfer)},
		)
	}

	return string(tx.Hash), nil

}

func generateMsg(comment string) (*cell.Cell, error) {
	var (
		body *cell.Cell
		err  error
	)
	if comment != "" {
		body, err = wallet.CreateCommentCell(comment)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil, nil

}

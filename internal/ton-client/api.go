package tonClient

import (
	"context"
	"github.com/BoryslavGlov/logrusx"
	"github.com/valyala/fasthttp"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"log"
	"tonListener/internal/config"
	"tonListener/internal/entity"
)

type TonClient interface {
	GetTransactionByHash(hash string) entity.TxResult
	TransferTokens(mnemonics []string, tons uint64, message string) (string, error)
}

type api struct {
	client  *fasthttp.Client
	request *fasthttp.Request
	config  *config.Config
	logx    logrusx.Logging
	tonapi  ton.APIClientWrapped
}

func NewApi(config *config.Config, logx logrusx.Logging) TonClient {
	req := initHeaders(config.XToken)
	client := newClient()

	tonapi, err := initTonClient(config.ConfigTON)
	if err != nil {
		log.Fatal("Error while trying to configure tonAPI ", err)
	}

	return &api{client: client, request: req, config: config, logx: logx, tonapi: tonapi}

}

func initHeaders(token string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.SetContentType("application/json")
	return req
}

func newClient() *fasthttp.Client {
	client := &fasthttp.Client{}
	return client
}

func initTonClient(url string) (ton.APIClientWrapped, error) {
	client := liteclient.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return ton.NewAPIClient(client).WithRetry(), nil
}

package tonClient

import (
	"encoding/json"
	"github.com/BoryslavGlov/logrusx"
	"github.com/valyala/fasthttp"
	"tonListener/internal/entity"
)

func (api *api) GetTransactionByHash(hash string) entity.TxResult {

	var (
		TxResult entity.TxResult
	)

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	api.request.CopyTo(req)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetRequestURI(api.config.TonURL + getTransactionUrl + hash)

	err := api.client.Do(req, res)
	if err != nil {
		api.logx.Error("Error occurred while trying to DO request in GetTransactionByHash",
			logrusx.LogField{Key: "context", Value: err},
			logrusx.LogField{Key: "Hash", Value: hash},
			logrusx.LogField{Key: "body", Value: string(res.Body())},
		)
	}

	err = json.Unmarshal(res.Body(), &TxResult)
	if err != nil {
		api.logx.Error("Error occurred while trying to Unmarshal in GetTransactionByHash",
			logrusx.LogField{Key: "context", Value: err},
			logrusx.LogField{Key: "hash", Value: hash},
			logrusx.LogField{Key: "body", Value: string(res.Body())},
		)
	}

	return TxResult
}

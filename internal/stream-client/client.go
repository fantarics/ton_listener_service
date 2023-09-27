package streamClient

import (
	"github.com/r3labs/sse/v2"
	"tonListener/internal/config"
)

type StreamClient struct {
	Stream *sse.Client
	config *config.Config
}

func NewClient(config *config.Config) *StreamClient {
	client := initHeaders(config.TonURL, config.XToken)

	return &StreamClient{config: config, Stream: client}
}

func initHeaders(host, token string) *sse.Client {
	client := sse.NewClient(host + listenerUrl)
	client.Headers["Authorization"] = "Bearer " + token
	return client
}

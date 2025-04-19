package binance

import (
	"context"
	"net/http"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

type MarketClient struct {
	Restful *binance_connector.Client
	Stream  *binance_connector.WebsocketStreamClient
}

func New(cfg *Config) *MarketClient {
	stream := binance_connector.NewWebsocketStreamClient(false)
	restful := binance_connector.NewClient(cfg.APIKey, cfg.SecretKey)
	restful.Debug = true
	restful.HTTPClient.Timeout = time.Second * 60
	restful.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	restful.NewPingService().Do(context.Background())
	return &MarketClient{Restful: restful, Stream: stream}
}

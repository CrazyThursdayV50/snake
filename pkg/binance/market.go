package binance

import (
	"context"

	"github.com/CrazyThursdayV50/goex/binance"
	"github.com/CrazyThursdayV50/goex/binance/websocket-api/api"
	"github.com/CrazyThursdayV50/pkgo/log"
)

type MarketClient struct {
	Restful *api.API
	Stream  *binance.WebSocketStreams
}

func New(ctx context.Context, logger log.Logger, cfg *Config) *MarketClient {
	restful := binance.NewWebSocketAPI(ctx, logger, cfg.APIKey, cfg.SecretKey)
	_, err := restful.Ping(ctx)
	if err != nil {
		panic(err)
	}

	stream := binance.NewWebSocketStreams()
	return &MarketClient{Restful: restful, Stream: stream}
}

package strategy

import (
	"context"
	"snake/internal/kline"
	"snake/internal/strategy"
	"snake/pkg/broadcast"
	"sync"

	"github.com/CrazyThursdayV50/pkgo/log"
)

type Service struct {
	ctx       context.Context
	logger    log.Logger
	id        int64
	broadcast *broadcast.Broadcast[*kline.Kline]
	klineRepo strategy.KlineRepository

	strategyLock sync.RWMutex
	strategies   map[int64]strategy.Strategy
}

func NewService(ctx context.Context, logger log.Logger, repo strategy.KlineRepository) *Service {
	return &Service{
		ctx:        ctx,
		logger:     logger,
		id:         0,
		broadcast:  broadcast.New[*kline.Kline](),
		klineRepo:  repo,
		strategies: make(map[int64]strategy.Strategy),
	}
}

type Response[T any] struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

func successResponse[T any](data *T) *Response[T] {
	return &Response[T]{
		Error:   "",
		Message: "",
		Data:    data,
	}
}

func failResponse[T any](err, msg string) *Response[T] {
	return &Response[T]{Error: err, Message: msg, Data: nil}
}

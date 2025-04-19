package strategy

import (
	"context"
	"snake/internal/kline"
	"snake/internal/kline/interval"
)

type KlineRepository interface {
	GetKlines(ctx context.Context, interval interval.Interval, from int64) <-chan *kline.Kline
}

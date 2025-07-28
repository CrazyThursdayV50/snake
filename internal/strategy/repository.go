package strategy

import (
	"context"
	"snake/internal/kline"
	"snake/internal/kline/interval"
)

type KlineRepository interface {
	// 获取 klines，kline 会一直推送更新
	GetKlines(ctx context.Context, interval interval.Interval, from int64) <-chan *kline.Kline
	// 获取历史 klines，kline 不会推送更新
	ListKlines(ctx context.Context, interval interval.Interval, from int64) <-chan *kline.Kline
}

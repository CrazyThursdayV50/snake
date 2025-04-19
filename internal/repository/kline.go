package repository

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
)

type KlineRepository interface {
	Insert(ctx context.Context, interval interval.Interval, klines []*models.Kline) error

	First(ctx context.Context, interval interval.Interval) (*models.Kline, error)
	Last(ctx context.Context, interval interval.Interval) (*models.Kline, error)
	List(ctx context.Context, interval interval.Interval, from, to int64) ([]*models.Kline, error)

	CheckMissing(ctx context.Context, interval interval.Interval, openTs []int64) ([]uint64, error)
}

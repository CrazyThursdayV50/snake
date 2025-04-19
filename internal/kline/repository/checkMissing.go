package repository

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	gmap "github.com/CrazyThursdayV50/pkgo/builtin/map"
	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
)

func (r *Repository) CheckMissing(ctx context.Context, interval interval.Interval, openTs []int64) ([]uint64, error) {
	var model models.Kline
	var results []int64
	db := r.db.Db(ctx).Model(&model).
		Scopes(
			models.KlineTable(interval),
			model.ColumnOpenTs().In(openTs),
		).
		Pluck(model.ColumnOpenTs().String(), &results)
	if db.Error != nil {
		return nil, db.Error
	}

	openTsMap := collector.Map(openTs, func(_ int, v int64) (bool, uint64, struct{}) {
		return true, uint64(v), struct{}{}
	})

	slice.From(results...).Iter(func(k int, v int64) (bool, error) {
		delete(openTsMap, uint64(v))
		return true, nil
	})

	return gmap.From(openTsMap).Keys().Unwrap(), nil
}

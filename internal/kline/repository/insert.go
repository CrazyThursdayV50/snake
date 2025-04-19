package repository

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	"sort"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	gmap "github.com/CrazyThursdayV50/pkgo/builtin/map"
)

func (r *Repository) Insert(ctx context.Context, interval interval.Interval, klines []*models.Kline) error {
	if len(klines) == 0 {
		return nil
	}

	klinesGroup := collector.Map(klines, func(_ int, v *models.Kline) (bool, int64, *models.Kline) {
		return true, v.OpenTs, v
	})

	var kline models.Kline
	var tempKlines []*models.Kline
	r.db.Db(ctx).Model(&kline).Scopes(models.KlineTable(interval)).Scopes(kline.ColumnOpenTs().In(gmap.From(klinesGroup).Keys().Unwrap())).FindInBatches(&tempKlines, 100, models.DefaultFindInBatchesCallback(func() {
		for _, t := range tempKlines {
			delete(klinesGroup, t.OpenTs)
		}
	}))

	klinesSlice := gmap.From(klinesGroup).Values()
	klinesSlice.WithLessFunc(func(a *models.Kline, b *models.Kline) bool {
		return a.OpenTs < b.OpenTs
	})
	sort.Sort(klinesSlice)
	if klinesSlice.Len() == 0 {
		return nil
	}
	return r.db.Db(ctx).Scopes(models.KlineTable(interval)).CreateInBatches(klinesSlice.Unwrap(), 200).Error
}

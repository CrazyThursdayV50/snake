package repository

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"

	"gorm.io/gorm/clause"
)

func (r *Repository) List(ctx context.Context, interval interval.Interval, from, to int64) ([]*models.Kline, error) {
	var model models.Kline
	var klines []*models.Kline
	db := r.db.Db(ctx).Model(&model).
		Scopes(
			models.KlineTable(interval),
			model.ColumnOpenTs().Between(from, to),
		).
		Order(clause.OrderByColumn{Column: clause.Column{Name: model.ColumnOpenTs().String()}}).Find(&klines)
	if db.Error != nil {
		return nil, db.Error
	}

	// 如果第一条数据被找到了，那么可以返回
	if len(klines) != 0 && klines[0].OpenTs == from {
		return klines, nil
	}

	// 否则，往前查询一条数据，拼到当前数据的第一条前面，再返回
	db = r.db.Db(ctx).
		Model(&model).
		Scopes(
			models.KlineTable(interval),
			model.ColumnOpenTs().LessThan(from),
		).
		Limit(1).
		Order(clause.OrderByColumn{Column: clause.Column{Name: model.ColumnOpenTs().String()}}).
		Find(&model)
	if db.Error != nil {
		return nil, db.Error
	}

	klines = append([]*models.Kline{&model}, klines...)
	return klines, nil
}

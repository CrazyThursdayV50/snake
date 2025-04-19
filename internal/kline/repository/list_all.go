package repository

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"

	"gorm.io/gorm/clause"
)

// ListAll 获取指定时间间隔的所有 kline 数据
func (r *Repository) ListAll(ctx context.Context, interval interval.Interval) ([]*models.Kline, error) {
	var klines []*models.Kline
	var model models.Kline
	db := r.db.Db(ctx).Scopes(models.KlineTable(interval)).
		Order(
			clause.OrderByColumn{
				Column: clause.Column{
					Name: model.ColumnOpenTs().String(),
				},
				Desc: false,
			},
		).Find(&klines)

	if db.Error != nil {
		return nil, db.Error
	}
	return klines, nil
}

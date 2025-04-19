package repository

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"

	"gorm.io/gorm/clause"
)

func (r *Repository) Last(ctx context.Context, interval interval.Interval) (*models.Kline, error) {
	var model models.Kline
	db := r.db.Db(ctx).Scopes(models.KlineTable(interval)).Order(clause.OrderByColumn{
		Column:  clause.Column{Name: model.ColumnOpenTs().String()},
		Desc:    true,
		Reorder: false,
	}).Limit(1).Find(&model)
	if db.Error != nil {
		return nil, db.Error
	}
	if db.RowsAffected == 0 {
		return nil, nil
	}
	return &model, nil
}

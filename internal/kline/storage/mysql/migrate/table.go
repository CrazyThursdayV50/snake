package migrate

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"

	"github.com/CrazyThursdayV50/pkgo/store/db/gorm"
)

func AutoMigrate(ctx context.Context, db *gorm.DB) {
	for _, interval := range interval.All() {
		db.Db(ctx).
			Scopes(models.KlineTable(interval)).
			AutoMigrate(new(models.Kline))
	}
}

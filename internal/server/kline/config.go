package kline

import (
	"snake/internal/service"
	"snake/pkg/binance"

	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/store/db/gorm"
)

type Config struct {
	Log     *defaultlogger.Config
	Mysql   *gorm.Config
	Binance *binance.Config
	Service *service.Config
}

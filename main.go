package main

import (
	"os"
	"snake/internal/server"
	"snake/internal/server/config"
	"snake/internal/service"
	"snake/pkg/binance"
	"time"

	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/store/db/gorm"
	g "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Config{

		Log: &defaultlogger.Config{
			Development:       true,
			Console:           true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Level:             "debug",
			CallerSkip:        1,
		},
		Service: &service.Config{Symbol: "BTCUSDT"},
		Mysql: &gorm.Config{
			Schema:        "test",
			DSN:           "alex@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
			MaxIdleConn:   10,
			MaxOpenConn:   2,
			MaxLifeTime:   3600,
			MaxIdleTime:   300,
			ServerVersion: "8.4.0",
			Gorm: g.Config{
				SkipDefaultTransaction:                   true,
				FullSaveAssociations:                     true,
				DisableForeignKeyConstraintWhenMigrating: true,
				IgnoreRelationshipsWhenMigrating:         true,
				QueryFields:                              true,
				CreateBatchSize:                          100,
				TranslateError:                           true,
				PropagateUnscoped:                        true,
			},
			Logger: logger.Config{
				SlowThreshold:        time.Millisecond * 100,
				Colorful:             true,
				ParameterizedQueries: false,
				LogLevel:             logger.Info,
			},
		},
		Binance: &binance.Config{
			APIKey:    "",
			SecretKey: "",
		},
	}

	server := server.New(&cfg)
	api, ok := os.LookupEnv("BN_APIKEY")
	if !ok {
		panic("invalid api")
	}

	secret, ok := os.LookupEnv("BN_SECRET")
	if !ok {
		panic("invalid secret")
	}

	cfg.Binance.APIKey = api
	cfg.Binance.SecretKey = secret
	server.Run()
}

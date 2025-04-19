package client

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/store/db/gorm"
	"github.com/CrazyThursdayV50/pkgo/trace"
)

func NewClient(logger log.Logger, tracer trace.Tracer, cfg *gorm.Config) *gorm.DB {
	return gorm.NewDB(logger, tracer, cfg)
}

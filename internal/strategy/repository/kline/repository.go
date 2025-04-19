package kline

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/request/resty"
)

type Repository struct {
	logger     log.Logger
	client     *resty.Client
	endpoint   string
	wsEndpoint string
}

func New(logger log.Logger, cfg *Config, client *resty.Client) *Repository {
	return &Repository{logger: logger, client: client, endpoint: cfg.Endpoint, wsEndpoint: cfg.WsEndpoint}
}

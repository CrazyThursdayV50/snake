package strategy

import (
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/client"
)

type Service struct {
	logger   log.Logger
	wsclient *client.Client
}

func NewService(logger log.Logger, wsclient *client.Client) *Service {
	return &Service{logger: logger, wsclient: wsclient}
}

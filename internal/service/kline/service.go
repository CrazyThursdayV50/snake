package kline

import (
	"snake/internal/kline"

	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/server"
	"github.com/gorilla/websocket"
)

type Repository = kline.Repository

// Service K线服务实现
type Service struct {
	logger    log.Logger
	repoKline Repository
	clients   map[*websocket.Conn]bool
	ws        *server.Server
}

// NewService 创建新的K线服务
func NewService(logger log.Logger, wsserver *server.Server, repoKline Repository) *Service {
	return &Service{
		logger:    logger,
		repoKline: repoKline,
		ws:        wsserver,
		clients:   make(map[*websocket.Conn]bool),
	}
}

type Response[T any] struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

func successResponse[T any](data *T) *Response[T] {
	return &Response[T]{
		Error:   "",
		Message: "",
		Data:    data,
	}
}

func failResponse[T any](err, msg string) *Response[T] {
	return &Response[T]{Error: err, Message: msg, Data: nil}
}

package kline

import (
	"context"

	"github.com/gorilla/websocket"
)

func HandleKlineMessage(ctx context.Context, messageType int, data []byte, err error) (int, []byte, error) {
	switch messageType {
	case websocket.PingMessage:
		return websocket.PongMessage, nil, nil

	case websocket.TextMessage:

	case websocket.BinaryMessage:

	default:
	}

	return websocket.TextMessage, nil, nil
}

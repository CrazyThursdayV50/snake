package handler

import (
	"snake/internal/kline/acl"
	"snake/internal/kline/storage/mysql/models"

	binance_connector "github.com/binance/binance-connector-go"
)

type WsKlineHandler struct {
	tempKline *binance_connector.WsKlineEvent
	canStore  bool
	uptodator func(uint64)
	storer    func(*models.Kline)
}

func NewWsKline(uptodator func(uint64), storer func(*models.Kline)) *WsKlineHandler {
	return &WsKlineHandler{
		tempKline: nil,
		canStore:  false,
		uptodator: uptodator,
		storer:    storer,
	}
}

func (h *WsKlineHandler) Handle(kline *binance_connector.WsKlineEvent) {
	if h.tempKline == nil {
		h.tempKline = kline
		h.uptodator(uint64(h.tempKline.Kline.StartTime))
		return
	}

	if h.tempKline.Kline.StartTime == kline.Kline.StartTime {
		h.tempKline = kline
		return
	}

	if h.canStore {
		model := acl.WsToDB(h.tempKline)
		h.storer(model)
		h.tempKline = kline
		return
	}

	h.canStore = true
	h.tempKline = kline
	h.uptodator(uint64(h.tempKline.Kline.StartTime))
}

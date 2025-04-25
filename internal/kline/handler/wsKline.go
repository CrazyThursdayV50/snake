package handler

import (
	"snake/internal/kline/acl"
	"snake/internal/kline/storage/mysql/models"

	binance_connector "github.com/binance/binance-connector-go"
)

type WsKlineHandler struct {
	tempKline *binance_connector.WsKlineEvent
	// canStore  bool
	uptodator func(uint64)
	storer    func(*models.Kline)
}

func NewWsKline(uptodator func(uint64), storer func(*models.Kline)) *WsKlineHandler {
	return &WsKlineHandler{
		tempKline: nil,
		// canStore:  false,
		uptodator: uptodator,
		storer:    storer,
	}
}

func (h *WsKlineHandler) Handle(kline *binance_connector.WsKlineEvent) {
	// 如果没有初始化 k 线，那么记录这个初始化的 k 线
	// 并且要调用 API 拿到此 k线前的所有更新数据
	if h.tempKline == nil {
		h.tempKline = kline
		h.uptodator(uint64(h.tempKline.Kline.StartTime))
		return
	}

	// 如果此k线为缓存k线的更新值
	// 那么更新缓存k线
	if h.tempKline.Kline.StartTime == kline.Kline.StartTime {
		h.tempKline = kline
		return
	}

	// // 默认还不能持久化
	// // 此时系统第一次拿到下一根k线
	// // 这个时候，先把缓存k线更新，然后调用API来拿到上一根k线的数据
	// if !h.canStore {
	// 	h.canStore = true
	// 	h.tempKline = kline
	// 	h.uptodator(uint64(h.tempKline.Kline.StartTime))
	// 	return
	// }

	model := acl.WsToDB(h.tempKline)
	h.storer(model)
	h.tempKline = kline
}

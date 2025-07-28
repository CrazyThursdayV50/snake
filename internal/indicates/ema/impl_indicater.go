package ema

import (
	"snake/internal/indicates"
)

// Next 计算下一个 Kline 对应的 EMA
func (e EMA) Next(kline *indicates.Kline) *EMA {
	if !e.lastKline.IsBefore(kline) {
		return nil
	}

	var ema EMA
	ema.count = e.count
	ema.alpha = e.alpha
	ema.alpha1 = e.alpha1
	ema.setLastKline(kline)
	ema.lastValue = e.Value
	ema.Value = nextEMA(e.Value, e.alpha, e.alpha1, kline.C)
	return &ema
}

// Update 用于更新当前 Kline
func (e *EMA) Update(kline *indicates.Kline) bool {
	if !e.lastKline.IsCurrent(kline) {
		return false
	}

	e.setLastKline(kline)
	e.Value = nextEMA(e.lastValue, e.alpha, e.alpha1, kline.C)
	return true
}

// Package ema
package ema

import (
	"snake/internal/indicates"
	"snake/internal/strategy/repository/kline"

	"github.com/shopspring/decimal"
)

type EMA struct {
	count     int
	alpha     decimal.Decimal
	alpha1    decimal.Decimal
	lastKline *kline.Kline
	lastValue decimal.Decimal

	Value        decimal.Decimal
	Timestamp    int64
	CurrentPrice decimal.Decimal
}

type emaBuilder struct {
	count  int
	klines []*indicates.Kline
}

func New(count int) *emaBuilder {
	return &emaBuilder{count: count}
}

func (b *emaBuilder) Klines(klines []*indicates.Kline) *emaBuilder {
	if len(klines) == b.count {
		b.klines = klines
	}
	return b
}

var emaCoefficient = decimal.NewFromFloat(2.0)

func (b *emaBuilder) Build() *EMA {
	// 检查是否有足够的 klines
	if len(b.klines) == 0 {
		return nil
	}
	
	var ema EMA
	ema.count = b.count
	ema.alpha = emaCoefficient.Div(decimal.NewFromInt(int64(b.count + 1)))
	ema.alpha1 = decimal.NewFromInt(1).Sub(ema.alpha)
	ema.setLastKline(b.klines[len(b.klines)-1])
	ema.calculate(b.klines)
	return &ema
}

func (e *EMA) setLastKline(kline *kline.Kline) {
	e.lastKline = kline
	e.Timestamp = kline.S
	e.CurrentPrice = kline.C
}

func (e *EMA) calculate(klines []*kline.Kline) {
	var lastEMA decimal.Decimal
	var currentEMA = klines[0].C
	for _, k := range klines[1:] {
		lastEMA = currentEMA
		currentEMA = nextEMA(lastEMA, e.alpha, e.alpha1, k.C)
	}

	e.Value = currentEMA
	e.lastValue = lastEMA
}

func nextEMA(currentEMA, alpha, alpha1, price decimal.Decimal) decimal.Decimal {
	return price.Mul(alpha).Add(currentEMA.Mul(alpha1))
}

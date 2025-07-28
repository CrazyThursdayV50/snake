// Package ma
package ma

import (
	"snake/internal/indicates"
	"snake/pkg/math"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/shopspring/decimal"
)

type MA struct {
	count  int
	prices []decimal.Decimal
	klines []*indicates.Kline

	// ma value
	Value decimal.Decimal
	// ma timestamp
	Timestamp int64
	// current price at this ma
	CurrentPrice decimal.Decimal // 最新价格
}

type maBuilder struct {
	count  int
	klines []*indicates.Kline
}

func New(count int) *maBuilder {

	return &maBuilder{count: count}
}

func (b *maBuilder) Klines(klines []*indicates.Kline) *maBuilder {
	if len(klines) == b.count {
		b.klines = klines
	}
	return b
}

func (m *MA) calculate() {
	m.Value = math.AverageDecimals(m.prices...)
	m.Timestamp = m.klines[m.count-1].E
	m.CurrentPrice = m.prices[m.count-1]
}

func (b *maBuilder) Build() *MA {
	var ma MA
	ma.count = b.count
	ma.prices = collector.Slice(b.klines, func(_ int, k *indicates.Kline) (bool, decimal.Decimal) {
		return true, k.C
	})

	ma.calculate()
	return &ma
}

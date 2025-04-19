package ma

import (
	"snake/internal/models"
	"snake/pkg/math"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/shopspring/decimal"
)

type MA struct {
	count     int
	prices    []decimal.Decimal
	Price     decimal.Decimal
	Timestamp int64
}

// func (m MA) Name() string { return fmt.Sprintf("%s%d", indicates.MA, m.count) }
// func (m MA) Data() *MA    { return &m }

// func (m *MA) Next(kline *models.Kline) indicates.Indicate[*MA] {
// 	var prices = slices.Delete(m.prices, 0, 0)
// 	prices = append(prices, kline.C)
// 	m.Price = math.AverageDecimals(m.prices...)
// 	return &MA{count: m.count, prices: prices, Price: decimal.Decimal}
// }

func New(klines ...*models.Kline) *MA {
	var prices = collector.Slice(klines, func(_ int, k *models.Kline) (bool, decimal.Decimal) {
		return true, k.C
	})
	var count = len(klines)
	var price = math.AverageDecimals(prices...)
	var ts = klines[count-1]
	return &MA{Price: price, count: count, Timestamp: ts, prices: prices}
}

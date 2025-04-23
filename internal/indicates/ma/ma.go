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
	var ts = klines[count-1].E
	return &MA{Price: price, count: count, Timestamp: ts, prices: prices}
}

// NextKline 计算下一个 Kline 对应的 MA
// 如果传入的 Kline 不是当前 MA 的下一个 Kline，返回 nil
func (m *MA) NextKline(kline *models.Kline) *MA {
	// 检查是否是下一个 Kline
	if kline.S <= m.Timestamp {
		return nil
	}

	// 检查时间间隔是否连续
	if kline.S != m.Timestamp+1 {
		return nil
	}

	// 移除第一个价格，添加新的价格
	var prices = make([]decimal.Decimal, len(m.prices))
	copy(prices, m.prices[1:])
	prices[len(prices)-1] = kline.C

	// 计算新的 MA
	var price = math.AverageDecimals(prices...)
	return &MA{
		count:     m.count,
		prices:    prices,
		Price:     price,
		Timestamp: kline.E,
	}
}

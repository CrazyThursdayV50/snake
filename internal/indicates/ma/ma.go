package ma

import (
	"snake/internal/kline"
	"snake/pkg/math"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/shopspring/decimal"
)

type MA struct {
	count     int
	prices    []decimal.Decimal
	Price     decimal.Decimal
	Timestamp int64
	LastPrice decimal.Decimal // 最新价格
}

// func (m MA) Name() string { return fmt.Sprintf("%s%d", indicates.MA, m.count) }
// func (m MA) Data() *MA    { return &m }

// func (m *MA) Next(kline *models.Kline) indicates.Indicate[*MA] {
// 	var prices = slices.Delete(m.prices, 0, 0)
// 	prices = append(prices, kline.C)
// 	m.Price = math.AverageDecimals(m.prices...)
// 	return &MA{count: m.count, prices: prices, Price: decimal.Decimal}
// }

func New(klines ...*kline.Kline) *MA {
	var prices = collector.Slice(klines, func(_ int, k *kline.Kline) (bool, decimal.Decimal) {
		return true, k.C
	})
	var count = len(klines)
	var price = math.AverageDecimals(prices...)
	var ts = klines[count-1].E
	var lastPrice = klines[count-1].C
	return &MA{Price: price, count: count, Timestamp: ts, prices: prices, LastPrice: lastPrice}
}

// NextKline 计算下一个 Kline 对应的 MA
// 如果传入的 Kline 时间戳小于当前 MA 的时间戳，返回 nil
func (m *MA) NextKline(kline *kline.Kline) *MA {
	// 检查是否是新的K线或更新数据
	if kline.S < m.Timestamp {
		return nil
	}

	var prices []decimal.Decimal
	// 判断是否是更新最后一条K线
	if kline.S == m.Timestamp {
		// 更新最后一条K线的价格
		prices = make([]decimal.Decimal, len(m.prices))
		copy(prices, m.prices)
		prices[len(prices)-1] = kline.C
	} else {
		// 新增K线
		prices = make([]decimal.Decimal, len(m.prices))
		copy(prices, m.prices[1:])
		prices[len(prices)-1] = kline.C
	}

	// 计算新的 MA
	var price = math.AverageDecimals(prices...)
	return &MA{
		count:     m.count,
		prices:    prices,
		Price:     price,
		Timestamp: kline.E,
		LastPrice: kline.C,
	}
}

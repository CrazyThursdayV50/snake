package bollingband

import (
	"snake/internal/kline"
	"snake/pkg/math"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/shopspring/decimal"
)

type BB struct {
	count     int
	prices    []decimal.Decimal
	MA        decimal.Decimal // 中轨
	Upper     decimal.Decimal // 上轨
	Lower     decimal.Decimal // 下轨
	Timestamp int64
	LastPrice decimal.Decimal // 最新价格
}

func New(klines ...*kline.Kline) *BB {
	var prices = collector.Slice(klines, func(_ int, k *kline.Kline) (bool, decimal.Decimal) {
		return true, k.C
	})
	var count = len(klines)
	var ma = math.AverageDecimals(prices...)
	var std = math.StandardDeviation(prices...)
	var upper = ma.Add(std.Mul(decimal.NewFromInt(2)))
	var lower = ma.Sub(std.Mul(decimal.NewFromInt(2)))
	var ts = klines[count-1].E
	var lastPrice = klines[count-1].C
	return &BB{
		count:     count,
		prices:    prices,
		MA:        ma,
		Upper:     upper,
		Lower:     lower,
		Timestamp: ts,
		LastPrice: lastPrice,
	}
}

// NextKline 计算下一个 Kline 对应的布林带
// 如果传入的 Kline 时间戳小于当前布林带的时间戳，返回 nil
func (b *BB) NextKline(kline *kline.Kline) *BB {
	// 检查是否是新的K线或更新数据
	if kline.S < b.Timestamp {
		return nil
	}

	var prices []decimal.Decimal
	// 判断是否是更新最后一条K线
	if kline.S == b.Timestamp {
		// 更新最后一条K线的价格
		prices = make([]decimal.Decimal, len(b.prices))
		copy(prices, b.prices)
		prices[len(prices)-1] = kline.C
	} else {
		// 新增K线
		prices = make([]decimal.Decimal, len(b.prices))
		copy(prices, b.prices[1:])
		prices[len(prices)-1] = kline.C
	}

	// 计算新的布林带
	var ma = math.AverageDecimals(prices...)
	var std = math.StandardDeviation(prices...)
	var upper = ma.Add(std.Mul(decimal.NewFromInt(2)))
	var lower = ma.Sub(std.Mul(decimal.NewFromInt(2)))

	return &BB{
		count:     b.count,
		prices:    prices,
		MA:        ma,
		Upper:     upper,
		Lower:     lower,
		Timestamp: kline.E,
		LastPrice: kline.C,
	}
}

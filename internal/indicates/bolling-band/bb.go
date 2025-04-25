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
	return &BB{
		count:     count,
		prices:    prices,
		MA:        ma,
		Upper:     upper,
		Lower:     lower,
		Timestamp: ts,
	}
}

// NextKline 计算下一个 Kline 对应的布林带
// 如果传入的 Kline 不是当前布林带的下一个 Kline，返回 nil
func (b *BB) NextKline(kline *kline.Kline) *BB {
	// 检查是否是下一个 Kline
	if kline.S <= b.Timestamp {
		return nil
	}

	// 检查时间间隔是否连续
	if kline.S != b.Timestamp+1 {
		return nil
	}

	// 移除第一个价格，添加新的价格
	var prices = make([]decimal.Decimal, len(b.prices))
	copy(prices, b.prices[1:])
	prices[len(prices)-1] = kline.C

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
	}
}

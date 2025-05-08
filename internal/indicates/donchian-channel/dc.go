package donchianchannel

import (
	"snake/internal/kline"

	"github.com/shopspring/decimal"
)

// DC 唐奇安通道结构体
type DC struct {
	count     int
	period    int               // 唐奇安通道周期，默认20
	prices    []decimal.Decimal // 收盘价
	highs     []decimal.Decimal // 最高价
	lows      []decimal.Decimal // 最低价
	Upper     decimal.Decimal   // 上轨 - 周期内最高价
	Middle    decimal.Decimal   // 中轨 - 上轨和下轨的中点
	Lower     decimal.Decimal   // 下轨 - 周期内最低价
	Timestamp int64
	LastPrice decimal.Decimal   // 最新价格
}

// New 创建新的唐奇安通道指标实例，使用默认20天周期
func New(klines ...*kline.Kline) *DC {
	return NewWithPeriod(20, klines...)
}

// NewWithPeriod 使用自定义周期创建唐奇安通道指标实例
func NewWithPeriod(period int, klines ...*kline.Kline) *DC {
	// 至少需要period个K线来计算唐奇安通道
	if len(klines) < period {
		return nil // 历史数据不足
	}

	// 提取价格数据
	var prices = make([]decimal.Decimal, len(klines))
	var highs = make([]decimal.Decimal, len(klines))
	var lows = make([]decimal.Decimal, len(klines))

	for i, k := range klines {
		prices[i] = k.C
		highs[i] = k.H
		lows[i] = k.L
	}

	// 计算唐奇安通道
	// 取最近period个K线的最高价和最低价
	var startIdx = len(klines) - period
	if startIdx < 0 {
		startIdx = 0
	}

	// 初始化为第一个价格
	var highest = highs[startIdx]
	var lowest = lows[startIdx]

	// 找出最近period个K线的最高价和最低价
	for i := startIdx; i < len(klines); i++ {
		if highs[i].GreaterThan(highest) {
			highest = highs[i]
		}
		if lows[i].LessThan(lowest) {
			lowest = lows[i]
		}
	}

	// 中轨是上轨和下轨的平均值
	var middle = highest.Add(lowest).Div(decimal.NewFromInt(2))

	var count = len(klines)
	var ts = klines[count-1].E
	var lastPrice = klines[count-1].C

	return &DC{
		count:     count,
		period:    period,
		prices:    prices,
		highs:     highs,
		lows:      lows,
		Upper:     highest,
		Middle:    middle,
		Lower:     lowest,
		Timestamp: ts,
		LastPrice: lastPrice,
	}
}

// NextKline 计算下一个K线对应的唐奇安通道
// 如果传入的K线时间戳小于当前唐奇安通道的时间戳，返回nil
func (d *DC) NextKline(kline *kline.Kline) *DC {
	// 检查是否是新的K线或更新数据
	if kline.S < d.Timestamp {
		return nil
	}

	var newPrices []decimal.Decimal
	var newHighs []decimal.Decimal
	var newLows []decimal.Decimal

	// 判断是否是更新最后一条K线
	if kline.S == d.Timestamp {
		// 更新最后一条K线的价格
		newPrices = make([]decimal.Decimal, len(d.prices))
		copy(newPrices, d.prices)
		newPrices[len(newPrices)-1] = kline.C

		newHighs = make([]decimal.Decimal, len(d.highs))
		copy(newHighs, d.highs)
		newHighs[len(newHighs)-1] = kline.H

		newLows = make([]decimal.Decimal, len(d.lows))
		copy(newLows, d.lows)
		newLows[len(newLows)-1] = kline.L
	} else {
		// 新增K线
		newPrices = make([]decimal.Decimal, len(d.prices)+1)
		copy(newPrices, d.prices)
		newPrices[len(newPrices)-1] = kline.C

		newHighs = make([]decimal.Decimal, len(d.highs)+1)
		copy(newHighs, d.highs)
		newHighs[len(newHighs)-1] = kline.H

		newLows = make([]decimal.Decimal, len(d.lows)+1)
		copy(newLows, d.lows)
		newLows[len(newLows)-1] = kline.L
	}

	// 计算唐奇安通道
	// 取最近period个K线的最高价和最低价
	var startIdx = len(newHighs) - d.period
	if startIdx < 0 {
		startIdx = 0
	}

	// 初始化为第一个价格
	var highest = newHighs[startIdx]
	var lowest = newLows[startIdx]

	// 找出最近period个K线的最高价和最低价
	for i := startIdx; i < len(newHighs); i++ {
		if newHighs[i].GreaterThan(highest) {
			highest = newHighs[i]
		}
		if newLows[i].LessThan(lowest) {
			lowest = newLows[i]
		}
	}

	// 中轨是上轨和下轨的平均值
	var middle = highest.Add(lowest).Div(decimal.NewFromInt(2))

	return &DC{
		count:     d.count + 1,
		period:    d.period,
		prices:    newPrices,
		highs:     newHighs,
		lows:      newLows,
		Upper:     highest,
		Middle:    middle,
		Lower:     lowest,
		Timestamp: kline.E,
		LastPrice: kline.C,
	}
}

// IsBuySignal 判断是否为买入信号（价格突破上轨）
func (d *DC) IsBuySignal(price decimal.Decimal) bool {
	return price.GreaterThanOrEqual(d.Upper)
}

// IsSellSignal 判断是否为卖出信号（价格跌破下轨）
func (d *DC) IsSellSignal(price decimal.Decimal) bool {
	return price.LessThanOrEqual(d.Lower)
}

// ChannelWidth 获取通道宽度（上轨减下轨）
func (d *DC) ChannelWidth() decimal.Decimal {
	return d.Upper.Sub(d.Lower)
}

// IsNarrowChannel 判断是否为窄通道（通道宽度小于平均值的一定比例）
func (d *DC) IsNarrowChannel(threshold decimal.Decimal) bool {
	// 计算最近period个K线的平均价格
	var startIdx = len(d.prices) - d.period
	if startIdx < 0 {
		startIdx = 0
	}

	var sum = decimal.Zero
	for i := startIdx; i < len(d.prices); i++ {
		sum = sum.Add(d.prices[i])
	}

	var avgPrice = sum.Div(decimal.NewFromInt(int64(len(d.prices) - startIdx)))

	// 计算通道宽度相对于平均价格的比例
	var widthRatio = d.ChannelWidth().Div(avgPrice)

	// 如果通道宽度比例小于阈值，则为窄通道
	return widthRatio.LessThan(threshold)
}

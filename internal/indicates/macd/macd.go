package macd

import (
	"snake/internal/kline"
	"snake/pkg/math"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/shopspring/decimal"
)

// MACD指标结构体
// MACD = EMA(12) - EMA(26)
// Signal = EMA(9) of MACD
// Histogram = MACD - Signal
type MACD struct {
	count        int
	prices       []decimal.Decimal
	fastEMA      decimal.Decimal // 快速EMA(通常是12日EMA)
	slowEMA      decimal.Decimal // 慢速EMA(通常是26日EMA)
	MACD         decimal.Decimal // MACD值 = 快速EMA - 慢速EMA
	Signal       decimal.Decimal // 信号线(通常是9日MACD的EMA)
	Histogram    decimal.Decimal // 柱状图 = MACD - 信号线
	Timestamp    int64
	fastPeriod   int               // 快速EMA周期
	slowPeriod   int               // 慢速EMA周期
	signalPeriod int               // 信号线周期
	macdValues   []decimal.Decimal // 保存MACD历史值用于计算信号线
	LastPrice    decimal.Decimal   // 最新价格
}

// New 创建新的MACD指标实例
// 默认使用12日EMA，26日EMA和9日信号线
func New(klines ...*kline.Kline) *MACD {
	return NewWithParams(12, 26, 9, klines...)
}

// NewWithParams 使用自定义参数创建MACD指标实例
func NewWithParams(fastPeriod, slowPeriod, signalPeriod int, klines ...*kline.Kline) *MACD {
	var prices = collector.Slice(klines, func(_ int, k *kline.Kline) (bool, decimal.Decimal) {
		return true, k.C
	})

	var count = len(klines)
	if count < slowPeriod {
		return nil // 历史数据不足
	}

	// 计算快速EMA（12日EMA）
	var fastEMA = calculateEMA(prices, fastPeriod)

	// 计算慢速EMA（26日EMA）
	var slowEMA = calculateEMA(prices, slowPeriod)

	// 计算MACD值 = 快速EMA - 慢速EMA
	var macdValue = fastEMA.Sub(slowEMA)

	// 获取/计算MACD的历史值
	var macdValues []decimal.Decimal
	if count >= slowPeriod+signalPeriod-1 {
		// 如果有足够的历史数据，计算过去的MACD值
		macdValues = calculateHistoricalMACD(prices, fastPeriod, slowPeriod, signalPeriod)
	} else {
		// 否则只使用当前MACD值
		macdValues = []decimal.Decimal{macdValue}
	}

	// 计算信号线（9日MACD的EMA）
	var signal decimal.Decimal
	if len(macdValues) >= signalPeriod {
		signal = calculateEMA(macdValues, signalPeriod)
	} else {
		signal = macdValue // 如果历史数据不足，信号线等于MACD
	}

	// 计算柱状图 = MACD - 信号线
	var histogram = macdValue.Sub(signal)

	var ts = klines[count-1].E
	var lastPrice = klines[count-1].C

	return &MACD{
		count:        count,
		prices:       prices,
		fastEMA:      fastEMA,
		slowEMA:      slowEMA,
		MACD:         macdValue,
		Signal:       signal,
		Histogram:    histogram,
		Timestamp:    ts,
		fastPeriod:   fastPeriod,
		slowPeriod:   slowPeriod,
		signalPeriod: signalPeriod,
		macdValues:   macdValues,
		LastPrice:    lastPrice,
	}
}

// NextKline 计算下一个K线对应的MACD
// 如果传入的K线时间戳小于当前MACD的时间戳，返回nil
func (m *MACD) NextKline(kline *kline.Kline) *MACD {
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

	// 计算新的快速EMA
	// 使用递推公式：EMA(t) = Price(t) * k + EMA(t-1) * (1-k) 其中 k = 2/(Period+1)
	var alpha = decimal.NewFromFloat(2.0 / float64(m.fastPeriod+1))
	var fastEMA = kline.C.Mul(alpha).Add(m.fastEMA.Mul(decimal.NewFromInt(1).Sub(alpha)))

	// 计算新的慢速EMA
	alpha = decimal.NewFromFloat(2.0 / float64(m.slowPeriod+1))
	var slowEMA = kline.C.Mul(alpha).Add(m.slowEMA.Mul(decimal.NewFromInt(1).Sub(alpha)))

	// 计算新的MACD值
	var macdValue = fastEMA.Sub(slowEMA)

	// 更新MACD历史值
	var macdValues []decimal.Decimal
	if len(m.macdValues) > 0 {
		macdValues = make([]decimal.Decimal, len(m.macdValues))
		if kline.S == m.Timestamp {
			// 更新最后一条MACD值
			copy(macdValues, m.macdValues)
			macdValues[len(macdValues)-1] = macdValue
		} else {
			// 新增MACD值
			copy(macdValues, m.macdValues[1:])
			macdValues[len(macdValues)-1] = macdValue
		}
	} else {
		macdValues = []decimal.Decimal{macdValue}
	}

	// 计算新的信号线
	alpha = decimal.NewFromFloat(2.0 / float64(m.signalPeriod+1))
	var signal = macdValue.Mul(alpha).Add(m.Signal.Mul(decimal.NewFromInt(1).Sub(alpha)))

	// 计算新的柱状图
	var histogram = macdValue.Sub(signal)

	return &MACD{
		count:        m.count,
		prices:       prices,
		fastEMA:      fastEMA,
		slowEMA:      slowEMA,
		MACD:         macdValue,
		Signal:       signal,
		Histogram:    histogram,
		Timestamp:    kline.E,
		fastPeriod:   m.fastPeriod,
		slowPeriod:   m.slowPeriod,
		signalPeriod: m.signalPeriod,
		macdValues:   macdValues,
		LastPrice:    kline.C,
	}
}

// calculateEMA 计算指数移动平均线
func calculateEMA(values []decimal.Decimal, period int) decimal.Decimal {
	if len(values) < period {
		// 如果数据量不足，返回简单平均值
		return math.AverageDecimals(values...)
	}

	// 使用前period个值的简单平均值作为EMA的初始值
	var ema = math.AverageDecimals(values[:period]...)

	// 计算系数 k = 2/(period+1)
	var k = decimal.NewFromFloat(2.0 / float64(period+1))
	var oneMinusK = decimal.NewFromInt(1).Sub(k)

	// 从period开始计算EMA
	for i := period; i < len(values); i++ {
		// EMA(today) = Price(today) * k + EMA(yesterday) * (1-k)
		ema = values[i].Mul(k).Add(ema.Mul(oneMinusK))
	}

	return ema
}

// calculateHistoricalMACD 计算历史MACD值
func calculateHistoricalMACD(prices []decimal.Decimal, fastPeriod, slowPeriod, signalPeriod int) []decimal.Decimal {
	if len(prices) < slowPeriod {
		return []decimal.Decimal{}
	}

	// 我们需要至少signalPeriod个MACD值来计算信号线
	var macdCount = signalPeriod
	var result = make([]decimal.Decimal, macdCount)

	// 如果数据不足，返回空
	if len(prices) < slowPeriod+macdCount-1 {
		return []decimal.Decimal{}
	}

	// 计算前macdCount个MACD值
	for i := 0; i < macdCount; i++ {
		// 取prices的子集来计算每个历史点的EMA
		endIndex := len(prices) - macdCount + i
		priceSubset := prices[:endIndex]

		// 计算该点的快速和慢速EMA
		fastEMA := calculateEMA(priceSubset, fastPeriod)
		slowEMA := calculateEMA(priceSubset, slowPeriod)

		// 计算MACD值
		result[i] = fastEMA.Sub(slowEMA)
	}

	return result
}

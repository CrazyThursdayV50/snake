package rsi

import (
	"snake/internal/kline"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/shopspring/decimal"
)

// RSI 相对强弱指标结构体
type RSI struct {
	count     int
	period    int // RSI周期，默认14
	prices    []decimal.Decimal
	changes   []decimal.Decimal // 价格变化
	gains     []decimal.Decimal // 价格上涨部分
	losses    []decimal.Decimal // 价格下跌部分（正值）
	avgGain   decimal.Decimal   // 平均上涨
	avgLoss   decimal.Decimal   // 平均下跌
	Value     decimal.Decimal   // RSI值，范围0-100
	Timestamp int64
	LastPrice decimal.Decimal   // 最新价格
}

// New 创建新的RSI指标实例，使用默认14天周期
func New(klines ...*kline.Kline) *RSI {
	return NewWithPeriod(14, klines...)
}

// NewWithPeriod 使用自定义周期创建RSI指标实例
func NewWithPeriod(period int, klines ...*kline.Kline) *RSI {
	// 至少需要period+1个K线来计算RSI
	if len(klines) < period+1 {
		return nil // 历史数据不足
	}

	// 提取收盘价
	var prices = collector.Slice(klines, func(_ int, k *kline.Kline) (bool, decimal.Decimal) {
		return true, k.C
	})

	// 计算价格变化
	changes := make([]decimal.Decimal, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		changes[i-1] = prices[i].Sub(prices[i-1])
	}

	// 分离上涨和下跌
	gains := make([]decimal.Decimal, len(changes))
	losses := make([]decimal.Decimal, len(changes))
	for i, change := range changes {
		if change.IsPositive() {
			gains[i] = change
			losses[i] = decimal.Zero
		} else {
			gains[i] = decimal.Zero
			losses[i] = change.Neg() // 转换为正值
		}
	}

	// 计算第一个RSI
	// 对前period个上涨和下跌取平均
	var sumGain = decimal.Zero
	var sumLoss = decimal.Zero
	for i := 0; i < period; i++ {
		sumGain = sumGain.Add(gains[i])
		sumLoss = sumLoss.Add(losses[i])
	}

	// 计算平均上涨和平均下跌
	var avgGain = sumGain.Div(decimal.NewFromInt(int64(period)))
	var avgLoss = sumLoss.Div(decimal.NewFromInt(int64(period)))

	// 计算RSI
	var rs decimal.Decimal
	var rsiValue decimal.Decimal

	if avgLoss.IsZero() {
		if avgGain.IsZero() {
			rsiValue = decimal.NewFromInt(50) // 如果没有变化，则RSI为50
		} else {
			rsiValue = decimal.NewFromInt(100) // 如果只有上涨，没有下跌，则RSI为100
		}
	} else {
		rs = avgGain.Div(avgLoss)
		rsiValue = rs.Div(rs.Add(decimal.NewFromInt(1))).Mul(decimal.NewFromInt(100))
	}

	var count = len(klines)
	var ts = klines[count-1].E
	var lastPrice = klines[count-1].C

	return &RSI{
		count:     count,
		period:    period,
		prices:    prices,
		changes:   changes,
		gains:     gains,
		losses:    losses,
		avgGain:   avgGain,
		avgLoss:   avgLoss,
		Value:     rsiValue,
		Timestamp: ts,
		LastPrice: lastPrice,
	}
}

// NextKline 计算下一个K线对应的RSI
// 如果传入的K线时间戳小于当前RSI的时间戳，返回nil
func (r *RSI) NextKline(kline *kline.Kline) *RSI {
	// 检查是否是新的K线或更新数据
	if kline.S < r.Timestamp {
		return nil
	}

	// 计算价格变化
	var lastPrice decimal.Decimal
	if kline.S == r.Timestamp {
		// 更新最后一条K线
		lastPrice = r.prices[len(r.prices)-2]
	} else {
		// 新增K线
		lastPrice = r.prices[len(r.prices)-1]
	}
	change := kline.C.Sub(lastPrice)

	// 计算当前的上涨和下跌
	var currentGain decimal.Decimal
	var currentLoss decimal.Decimal

	if change.IsPositive() {
		currentGain = change
		currentLoss = decimal.Zero
	} else {
		currentGain = decimal.Zero
		currentLoss = change.Neg() // 转换为正值
	}

	// 使用平滑RSI公式更新平均上涨和平均下跌
	// avgGain = ((period-1) * prevAvgGain + currentGain) / period
	// avgLoss = ((period-1) * prevAvgLoss + currentLoss) / period

	periodMinusOne := decimal.NewFromInt(int64(r.period - 1))
	periodDecimal := decimal.NewFromInt(int64(r.period))

	newAvgGain := (periodMinusOne.Mul(r.avgGain).Add(currentGain)).Div(periodDecimal)
	newAvgLoss := (periodMinusOne.Mul(r.avgLoss).Add(currentLoss)).Div(periodDecimal)

	// 计算RSI
	var rs decimal.Decimal
	var rsiValue decimal.Decimal

	if newAvgLoss.IsZero() {
		if newAvgGain.IsZero() {
			rsiValue = decimal.NewFromInt(50) // 如果没有变化，则RSI为50
		} else {
			rsiValue = decimal.NewFromInt(100) // 如果只有上涨，没有下跌，则RSI为100
		}
	} else {
		rs = newAvgGain.Div(newAvgLoss)
		rsiValue = rs.Div(rs.Add(decimal.NewFromInt(1))).Mul(decimal.NewFromInt(100))
	}

	// 更新价格数组
	var newPrices []decimal.Decimal
	if kline.S == r.Timestamp {
		// 更新最后一条K线
		newPrices = make([]decimal.Decimal, len(r.prices))
		copy(newPrices, r.prices)
		newPrices[len(newPrices)-1] = kline.C
	} else {
		// 新增K线
		newPrices = make([]decimal.Decimal, len(r.prices)+1)
		copy(newPrices, r.prices)
		newPrices[len(newPrices)-1] = kline.C
	}

	// 更新变化数组
	var newChanges []decimal.Decimal
	if kline.S == r.Timestamp {
		// 更新最后一条K线的变化
		newChanges = make([]decimal.Decimal, len(r.changes))
		copy(newChanges, r.changes)
		newChanges[len(newChanges)-1] = change
	} else {
		// 新增变化
		newChanges = make([]decimal.Decimal, len(r.changes)+1)
		copy(newChanges, r.changes)
		newChanges[len(newChanges)-1] = change
	}

	// 更新上涨和下跌数组
	var newGains []decimal.Decimal
	var newLosses []decimal.Decimal
	if kline.S == r.Timestamp {
		// 更新最后一条K线的上涨和下跌
		newGains = make([]decimal.Decimal, len(r.gains))
		copy(newGains, r.gains)
		newGains[len(newGains)-1] = currentGain

		newLosses = make([]decimal.Decimal, len(r.losses))
		copy(newLosses, r.losses)
		newLosses[len(newLosses)-1] = currentLoss
	} else {
		// 新增上涨和下跌
		newGains = make([]decimal.Decimal, len(r.gains)+1)
		copy(newGains, r.gains)
		newGains[len(newGains)-1] = currentGain

		newLosses = make([]decimal.Decimal, len(r.losses)+1)
		copy(newLosses, r.losses)
		newLosses[len(newLosses)-1] = currentLoss
	}

	return &RSI{
		count:     r.count + 1,
		period:    r.period,
		prices:    newPrices,
		changes:   newChanges,
		gains:     newGains,
		losses:    newLosses,
		avgGain:   newAvgGain,
		avgLoss:   newAvgLoss,
		Value:     rsiValue,
		Timestamp: kline.E,
		LastPrice: kline.C,
	}
}

// IsBuy 判断是否为买入信号（RSI < 30）
func (r *RSI) IsBuy() bool {
	return r.Value.LessThan(decimal.NewFromInt(30))
}

// IsSell 判断是否为卖出信号（RSI > 70）
func (r *RSI) IsSell() bool {
	return r.Value.GreaterThan(decimal.NewFromInt(70))
}

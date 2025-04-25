package rsi

import (
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRSIIndicator(t *testing.T) {
	// 创建测试数据
	startTime := time.Now().Unix() * 1000
	klines := generateTestKlines(startTime, 30)

	t.Run("New RSI with default period", func(t *testing.T) {
		// 测试默认周期(14)
		rsi := New(klines...)

		// 断言RSI实例不为空
		assert.NotNil(t, rsi)

		// 断言RSI属性已正确计算
		assert.False(t, rsi.Value.IsZero(), "RSI value should not be zero")

		// 验证RSI值是否在合理范围内(0-100)
		assert.True(t, rsi.Value.GreaterThanOrEqual(decimal.Zero), "RSI should be >= 0")
		assert.True(t, rsi.Value.LessThanOrEqual(decimal.NewFromInt(100)), "RSI should be <= 100")

		// 验证期间设置
		assert.Equal(t, 14, rsi.period)
	})

	t.Run("New RSI with custom period", func(t *testing.T) {
		// 测试自定义周期(7)
		rsi := NewWithPeriod(7, klines...)

		// 断言RSI实例不为空
		assert.NotNil(t, rsi)

		// 断言期间设置
		assert.Equal(t, 7, rsi.period)

		// 断言RSI属性已正确计算
		assert.False(t, rsi.Value.IsZero(), "RSI value should not be zero")

		// 验证RSI值是否在合理范围内(0-100)
		assert.True(t, rsi.Value.GreaterThanOrEqual(decimal.Zero), "RSI should be >= 0")
		assert.True(t, rsi.Value.LessThanOrEqual(decimal.NewFromInt(100)), "RSI should be <= 100")
	})

	t.Run("New RSI with insufficient data", func(t *testing.T) {
		// 使用少于所需的K线数量
		// RSI需要至少period+1个K线
		shortKlines := klines[:5] // 只提供5个K线
		rsi := NewWithPeriod(14, shortKlines...)

		// 断言RSI为nil，因为数据不足
		assert.Nil(t, rsi, "RSI should be nil with insufficient data")
	})

	t.Run("NextKline", func(t *testing.T) {
		// 创建初始RSI，使用前20个K线
		rsi := New(klines[:20]...)
		assert.NotNil(t, rsi, "Initial RSI should not be nil")

		// 记录初始值用于比较
		initialValue := rsi.Value

		// 使用下一个K线更新RSI
		nextKline := klines[20]
		// 确保下一个K线的开始时间是连续的
		nextKline.S = rsi.Timestamp + 1
		nextKline.E = nextKline.S + 60000 // 一分钟后结束

		updatedRSI := rsi.NextKline(nextKline)

		// 断言更新后的RSI不为空
		assert.NotNil(t, updatedRSI, "Updated RSI should not be nil")

		// 断言时间戳已更新
		assert.Equal(t, nextKline.E, updatedRSI.Timestamp)

		// 验证RSI值仍在合理范围内
		assert.True(t, updatedRSI.Value.GreaterThanOrEqual(decimal.Zero), "RSI should be >= 0")
		assert.True(t, updatedRSI.Value.LessThanOrEqual(decimal.NewFromInt(100)), "RSI should be <= 100")

		// 验证RSI值已改变
		assert.False(t, initialValue.Equal(updatedRSI.Value), "RSI value should change after updating with new kline")

		// 验证价格数组已正确更新
		assert.Equal(t, len(rsi.prices)+1, len(updatedRSI.prices))
		assert.True(t, nextKline.C.Equal(updatedRSI.prices[len(updatedRSI.prices)-1]))
	})

	t.Run("NextKline with non-sequential kline", func(t *testing.T) {
		// 创建初始RSI
		rsi := New(klines[:20]...)
		assert.NotNil(t, rsi, "Initial RSI should not be nil")

		// 创建一个不连续的K线（时间戳不连续）
		nonSequentialKline := &kline.Kline{
			O: decimal.NewFromFloat(100),
			C: decimal.NewFromFloat(101),
			H: decimal.NewFromFloat(102),
			L: decimal.NewFromFloat(99),
			V: decimal.NewFromFloat(1000),
			A: decimal.NewFromFloat(100000),
			S: rsi.Timestamp + 100, // 不连续的开始时间
			E: rsi.Timestamp + 160000,
		}

		// 应该返回nil
		updatedRSI := rsi.NextKline(nonSequentialKline)
		assert.Nil(t, updatedRSI, "RSI should be nil with non-sequential kline")
	})

	t.Run("IsBuy and IsSell signals", func(t *testing.T) {
		// 为了测试，创建一个RSI实例但手动设置其值
		testRSI := New(klines...)

		// 测试买入信号 - RSI < 30
		testRSI.Value = decimal.NewFromInt(25)
		assert.True(t, testRSI.IsBuy(), "Should be a buy signal when RSI = 25")
		assert.False(t, testRSI.IsSell(), "Should not be a sell signal when RSI = 25")

		// 测试卖出信号 - RSI > 70
		testRSI.Value = decimal.NewFromInt(75)
		assert.True(t, testRSI.IsSell(), "Should be a sell signal when RSI = 75")
		assert.False(t, testRSI.IsBuy(), "Should not be a buy signal when RSI = 75")

		// 测试中性 - 30 < RSI < 70
		testRSI.Value = decimal.NewFromInt(50)
		assert.False(t, testRSI.IsBuy(), "Should not be a buy signal when RSI = 50")
		assert.False(t, testRSI.IsSell(), "Should not be a sell signal when RSI = 50")
	})

	t.Run("RSI calculation with simple case", func(t *testing.T) {
		// 创建一个具有明确上涨趋势的K线序列，预期RSI值较高
		risingSeries := generateRisingKlines(startTime, 20, 100.0, 0.5)
		rsiRising := NewWithPeriod(14, risingSeries...)

		assert.NotNil(t, rsiRising)
		t.Logf("RSI for rising trend: %s", rsiRising.Value.String())
		assert.True(t, rsiRising.Value.GreaterThan(decimal.NewFromInt(50)),
			"RSI should be > 50 for rising trend")

		// 创建一个具有明确下跌趋势的K线序列，预期RSI值较低
		fallingSeries := generateFallingKlines(startTime, 20, 100.0, 0.5)
		rsiFalling := NewWithPeriod(14, fallingSeries...)

		assert.NotNil(t, rsiFalling)
		t.Logf("RSI for falling trend: %s", rsiFalling.Value.String())
		assert.True(t, rsiFalling.Value.LessThan(decimal.NewFromInt(50)),
			"RSI should be < 50 for falling trend")
	})
}

// 生成测试用的K线数据
func generateTestKlines(startTime int64, count int) []*kline.Kline {
	klines := make([]*kline.Kline, count)

	// 基准价格，用于生成波动的价格
	basePrice := 100.0

	for i := 0; i < count; i++ {
		// 生成波动的价格，模拟市场波动
		// 这里使用简单的正弦波函数来模拟价格波动
		priceOffset := 5.0 * float64(i%10) / 10.0
		if i%20 >= 10 {
			priceOffset = 5.0 * (1.0 - float64(i%10)/10.0)
		}

		open := basePrice + priceOffset - 0.5
		close := basePrice + priceOffset + 0.5
		high := basePrice + priceOffset + 1.0
		low := basePrice + priceOffset - 1.0

		klines[i] = &kline.Kline{
			O: decimal.NewFromFloat(open),
			C: decimal.NewFromFloat(close),
			H: decimal.NewFromFloat(high),
			L: decimal.NewFromFloat(low),
			V: decimal.NewFromFloat(1000.0 + float64(i%5)*500.0),
			A: decimal.NewFromFloat(100000.0 + float64(i%5)*50000.0),
			S: startTime + int64(i)*60000, // 每分钟一个K线
			E: startTime + int64(i+1)*60000,
		}
	}

	return klines
}

// 生成持续上涨的K线序列
func generateRisingKlines(startTime int64, count int, startPrice, increment float64) []*kline.Kline {
	klines := make([]*kline.Kline, count)

	price := startPrice
	for i := 0; i < count; i++ {
		open := price
		close := price + increment

		klines[i] = &kline.Kline{
			O: decimal.NewFromFloat(open),
			C: decimal.NewFromFloat(close),
			H: decimal.NewFromFloat(close + 0.2),
			L: decimal.NewFromFloat(open - 0.2),
			V: decimal.NewFromFloat(1000.0),
			A: decimal.NewFromFloat(100000.0),
			S: startTime + int64(i)*60000,
			E: startTime + int64(i+1)*60000,
		}

		price = close // 下一个K线的开盘价是前一个K线的收盘价
	}

	return klines
}

// 生成持续下跌的K线序列
func generateFallingKlines(startTime int64, count int, startPrice, decrement float64) []*kline.Kline {
	klines := make([]*kline.Kline, count)

	price := startPrice
	for i := 0; i < count; i++ {
		open := price
		close := price - decrement

		klines[i] = &kline.Kline{
			O: decimal.NewFromFloat(open),
			C: decimal.NewFromFloat(close),
			H: decimal.NewFromFloat(open + 0.2),
			L: decimal.NewFromFloat(close - 0.2),
			V: decimal.NewFromFloat(1000.0),
			A: decimal.NewFromFloat(100000.0),
			S: startTime + int64(i)*60000,
			E: startTime + int64(i+1)*60000,
		}

		price = close // 下一个K线的开盘价是前一个K线的收盘价
	}

	return klines
}

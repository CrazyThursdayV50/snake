package macd

import (
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMACDIndicator(t *testing.T) {
	// 创建测试数据
	startTime := time.Now().Unix() * 1000
	klines := generateTestKlines(startTime, 100)

	t.Run("New MACD with default parameters", func(t *testing.T) {
		// 测试默认参数(12,26,9)
		macd := New(klines...)

		// 断言MACD实例不为空
		assert.NotNil(t, macd)

		// 断言MACD属性已正确计算
		assert.False(t, macd.MACD.IsZero(), "MACD should not be zero")
		assert.False(t, macd.Signal.IsZero(), "Signal should not be zero")
		assert.False(t, macd.Histogram.IsZero(), "Histogram should not be zero")

		// 验证MACD = fastEMA - slowEMA
		calculatedMACD := macd.fastEMA.Sub(macd.slowEMA)
		assert.True(t, macd.MACD.Equal(calculatedMACD),
			"MACD should equal fastEMA - slowEMA")

		// 验证Histogram = MACD - Signal
		calculatedHistogram := macd.MACD.Sub(macd.Signal)
		assert.True(t, macd.Histogram.Equal(calculatedHistogram),
			"Histogram should equal MACD - Signal")
	})

	t.Run("New MACD with custom parameters", func(t *testing.T) {
		// 测试自定义参数(5,10,3)
		macd := NewWithParams(5, 10, 3, klines...)

		// 断言MACD实例不为空
		assert.NotNil(t, macd)

		// 断言参数已正确设置
		assert.Equal(t, 5, macd.fastPeriod)
		assert.Equal(t, 10, macd.slowPeriod)
		assert.Equal(t, 3, macd.signalPeriod)

		// 断言MACD属性已正确计算
		assert.False(t, macd.MACD.IsZero(), "MACD should not be zero")
		assert.False(t, macd.Signal.IsZero(), "Signal should not be zero")
		assert.False(t, macd.Histogram.IsZero(), "Histogram should not be zero")
	})

	t.Run("NextKline", func(t *testing.T) {
		// 创建初始MACD，使用前90个K线
		macd := New(klines[:90]...)
		assert.NotNil(t, macd, "Initial MACD should not be nil")

		// 使用下一个K线更新MACD
		nextKline := klines[90]
		// 确保下一个K线的开始时间是连续的
		nextKline.S = macd.Timestamp + 1
		nextKline.E = nextKline.S + 60000 // 一分钟后结束

		updatedMACD := macd.NextKline(nextKline)

		// 断言更新后的MACD不为空
		assert.NotNil(t, updatedMACD, "Updated MACD should not be nil")

		// 断言时间戳已更新
		assert.Equal(t, nextKline.E, updatedMACD.Timestamp)

		// 断言价格列表长度相同
		assert.Equal(t, len(macd.prices), len(updatedMACD.prices))

		// 验证更新的MACD = fastEMA - slowEMA
		calculatedMACD := updatedMACD.fastEMA.Sub(updatedMACD.slowEMA)
		assert.True(t, updatedMACD.MACD.Equal(calculatedMACD),
			"Updated MACD should equal fastEMA - slowEMA")

		// 验证更新的Histogram = MACD - Signal
		calculatedHistogram := updatedMACD.MACD.Sub(updatedMACD.Signal)
		assert.True(t, updatedMACD.Histogram.Equal(calculatedHistogram),
			"Updated Histogram should equal MACD - Signal")
	})

	t.Run("NextKline with non-sequential kline", func(t *testing.T) {
		// 创建初始MACD
		macd := New(klines[:90]...)
		assert.NotNil(t, macd, "Initial MACD should not be nil")

		// 创建一个不连续的K线（时间戳不连续）
		nonSequentialKline := &kline.Kline{
			O: decimal.NewFromFloat(100),
			C: decimal.NewFromFloat(101),
			H: decimal.NewFromFloat(102),
			L: decimal.NewFromFloat(99),
			V: decimal.NewFromFloat(1000),
			A: decimal.NewFromFloat(100000),
			S: macd.Timestamp + 100, // 不连续的开始时间
			E: macd.Timestamp + 101,
		}

		// 应该返回nil
		updatedMACD := macd.NextKline(nonSequentialKline)
		assert.Nil(t, updatedMACD)
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

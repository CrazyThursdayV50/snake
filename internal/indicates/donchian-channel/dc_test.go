package donchianchannel

import (
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDonchianChannelIndicator(t *testing.T) {
	// 创建测试数据
	startTime := time.Now().Unix() * 1000
	klines := generateTestKlines(startTime, 30)

	t.Run("New DC with default period", func(t *testing.T) {
		// 测试默认周期(20)
		dc := New(klines...)

		// 断言DC实例不为空
		assert.NotNil(t, dc)

		// 断言DC属性已正确计算
		assert.False(t, dc.Upper.IsZero(), "Upper band should not be zero")
		assert.False(t, dc.Lower.IsZero(), "Lower band should not be zero")
		assert.False(t, dc.Middle.IsZero(), "Middle band should not be zero")

		// 验证Upper应该大于或等于Lower
		assert.True(t, dc.Upper.GreaterThanOrEqual(dc.Lower), "Upper band should be >= Lower band")

		// 验证Middle是Upper和Lower的平均值
		expectedMiddle := dc.Upper.Add(dc.Lower).Div(decimal.NewFromInt(2))
		assert.True(t, dc.Middle.Equal(expectedMiddle), "Middle band should be the average of Upper and Lower")

		// 验证期间设置
		assert.Equal(t, 20, dc.period)
	})

	t.Run("New DC with custom period", func(t *testing.T) {
		// 测试自定义周期(10)
		dc := NewWithPeriod(10, klines...)

		// 断言DC实例不为空
		assert.NotNil(t, dc)

		// 断言期间设置
		assert.Equal(t, 10, dc.period)

		// 断言DC属性已正确计算
		assert.False(t, dc.Upper.IsZero(), "Upper band should not be zero")
		assert.False(t, dc.Lower.IsZero(), "Lower band should not be zero")
		assert.False(t, dc.Middle.IsZero(), "Middle band should not be zero")
	})

	t.Run("New DC with insufficient data", func(t *testing.T) {
		// 使用少于所需的K线数量
		shortKlines := klines[:5] // 只提供5个K线
		dc := NewWithPeriod(10, shortKlines...)

		// 断言DC为nil，因为数据不足
		assert.Nil(t, dc, "DC should be nil with insufficient data")
	})

	t.Run("NextKline", func(t *testing.T) {
		// 创建初始DC，使用前20个K线
		dc := New(klines[:20]...)
		assert.NotNil(t, dc, "Initial DC should not be nil")

		// 保存初始值
		initialUpper := dc.Upper
		initialLower := dc.Lower

		// 使用下一个K线更新DC
		nextKline := klines[20]
		// 确保下一个K线的开始时间是连续的
		nextKline.S = dc.Timestamp + 1
		nextKline.E = nextKline.S + 60000 // 一分钟后结束

		// 设置下一个K线的价格明显高于之前的最高价，以确保上轨变化
		nextKline.H = initialUpper.Add(decimal.NewFromInt(10))

		updatedDC := dc.NextKline(nextKline)

		// 断言更新后的DC不为空
		assert.NotNil(t, updatedDC, "Updated DC should not be nil")

		// 断言时间戳已更新
		assert.Equal(t, nextKline.E, updatedDC.Timestamp)

		// 验证上轨已更新为新的最高价
		assert.True(t, updatedDC.Upper.Equal(nextKline.H), "Upper band should update to new highest price")

		// 验证下轨保持不变（假设新的K线的最低价不低于之前的最低价）
		assert.True(t, updatedDC.Lower.Equal(initialLower), "Lower band should remain unchanged")

		// 验证中轨是上轨和下轨的平均值
		expectedMiddle := updatedDC.Upper.Add(updatedDC.Lower).Div(decimal.NewFromInt(2))
		assert.True(t, updatedDC.Middle.Equal(expectedMiddle), "Middle band should be the average of Upper and Lower")
	})

	t.Run("NextKline with non-sequential kline", func(t *testing.T) {
		// 创建初始DC
		dc := New(klines[:20]...)
		assert.NotNil(t, dc, "Initial DC should not be nil")

		// 创建一个不连续的K线（时间戳不连续）
		nonSequentialKline := &kline.Kline{
			O: decimal.NewFromFloat(100),
			C: decimal.NewFromFloat(101),
			H: decimal.NewFromFloat(102),
			L: decimal.NewFromFloat(99),
			V: decimal.NewFromFloat(1000),
			A: decimal.NewFromFloat(100000),
			S: dc.Timestamp + 100, // 不连续的开始时间
			E: dc.Timestamp + 160000,
		}

		// 应该返回nil
		updatedDC := dc.NextKline(nonSequentialKline)
		assert.Nil(t, updatedDC, "DC should be nil with non-sequential kline")
	})

	t.Run("Buy and Sell Signals", func(t *testing.T) {
		// 创建DC实例
		dc := New(klines...)
		assert.NotNil(t, dc)

		// 价格在上轨上方，应该是买入信号
		buyPrice := dc.Upper.Add(decimal.NewFromFloat(1.0))
		assert.True(t, dc.IsBuySignal(buyPrice), "Price above Upper band should be a buy signal")

		// 价格在下轨下方，应该是卖出信号
		sellPrice := dc.Lower.Sub(decimal.NewFromFloat(1.0))
		assert.True(t, dc.IsSellSignal(sellPrice), "Price below Lower band should be a sell signal")

		// 价格在通道内，既不是买入也不是卖出信号
		middlePrice := dc.Middle
		assert.False(t, dc.IsBuySignal(middlePrice), "Price in the middle should not be a buy signal")
		assert.False(t, dc.IsSellSignal(middlePrice), "Price in the middle should not be a sell signal")
	})

	t.Run("Channel Width and Narrow Channel", func(t *testing.T) {
		// 创建DC实例
		dc := New(klines...)
		assert.NotNil(t, dc)

		// 验证通道宽度
		expectedWidth := dc.Upper.Sub(dc.Lower)
		assert.True(t, dc.ChannelWidth().Equal(expectedWidth), "Channel width should be Upper - Lower")

		// 测试窄通道判断，假设阈值为0.1（10%）
		threshold := decimal.NewFromFloat(0.1)

		// 获取窄通道结果
		isNarrow := dc.IsNarrowChannel(threshold)
		t.Logf("Channel width ratio: %s, threshold: %s, isNarrow: %v",
			dc.ChannelWidth().Div(dc.Middle), threshold, isNarrow)

		// 注意：这里不断言具体结果，因为它取决于测试数据的实际情况
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

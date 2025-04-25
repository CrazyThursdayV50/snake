package donchian_strategy

import (
	"snake/internal/kline"
	"snake/internal/types"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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

// 测试策略初始化
func TestNew(t *testing.T) {
	strategy := New()
	assert.NotNil(t, strategy)
	assert.Equal(t, "Donchian Channel Strategy", strategy.Name())
	assert.Equal(t, 20, strategy.breakoutPeriod)
	assert.Equal(t, 10, strategy.exitPeriod)
	assert.Equal(t, decimal.NewFromFloat(1.0), strategy.riskPercent)
	assert.Equal(t, "none", strategy.position)
}

// 测试策略参数设置
func TestSetParams(t *testing.T) {
	strategy := New()
	strategy.SetParams(30, 15, decimal.NewFromFloat(2.0))

	assert.Equal(t, 30, strategy.breakoutPeriod)
	assert.Equal(t, 15, strategy.exitPeriod)
	assert.Equal(t, decimal.NewFromFloat(2.0), strategy.riskPercent)
	assert.Nil(t, strategy.breakoutChannel)
	assert.Nil(t, strategy.exitChannel)
}

// 测试策略更新-当数据不足时
func TestUpdateInsufficientData(t *testing.T) {
	strategy := New()
	err := strategy.Init(decimal.Zero, decimal.NewFromInt(1000))
	assert.NoError(t, err)

	// 创建少于所需周期的K线数据
	klines := generateTestKlines(time.Now().Unix(), 10)

	// 更新策略（数据不足）
	signal, err := strategy.Update(klines[0])
	assert.NoError(t, err)
	assert.Equal(t, types.SignalTypeHold, signal.Type)
}

// 测试策略更新-当足够数据时应该生成信号
func TestUpdateWithSufficientData(t *testing.T) {
	strategy := New()
	err := strategy.Init(decimal.Zero, decimal.NewFromInt(1000))
	assert.NoError(t, err)

	// 创建足够的K线数据
	klines := generateTestKlines(time.Now().Unix(), 30)

	// 先更新前20个K线以积累数据
	for i := 0; i < 20; i++ {
		signal, err := strategy.Update(klines[i])
		assert.NoError(t, err)
		assert.Equal(t, types.SignalTypeHold, signal.Type)
	}

	// 修改第21个K线，使其突破上轨（设置一个非常高的价格）
	klines[20].H = decimal.NewFromFloat(120.0)
	klines[20].C = decimal.NewFromFloat(115.0)

	// 更新策略，应该生成买入信号
	signal, err := strategy.Update(klines[20])
	assert.NoError(t, err)

	// 由于价格设置为突破上轨，预期是买入信号
	assert.Equal(t, types.SignalTypeBuy, signal.Type)
	assert.Equal(t, "long", strategy.position)
}

// 测试多头平仓
func TestLongPositionExit(t *testing.T) {
	strategy := New()
	err := strategy.Init(decimal.NewFromInt(1), decimal.NewFromInt(1000))
	assert.NoError(t, err)

	// 设置初始状态为多头
	strategy.position = "long"

	// 创建足够的K线数据
	klines := generateTestKlines(time.Now().Unix(), 30)

	// 更新前20个K线以初始化指标
	for i := 0; i < 20; i++ {
		strategy.Update(klines[i])
	}

	// 确保退出通道已经计算
	assert.NotNil(t, strategy.exitChannel)

	// 修改下一个K线，使其跌破下轨（设置一个非常低的价格）
	klines[20].L = decimal.NewFromFloat(80.0)
	klines[20].C = decimal.NewFromFloat(85.0)

	// 更新策略，应该生成卖出信号（平多）
	signal, err := strategy.Update(klines[20])
	assert.NoError(t, err)
	assert.Equal(t, types.SignalTypeSell, signal.Type)
	assert.Equal(t, "none", strategy.position)
}

// 测试空头平仓
func TestShortPositionExit(t *testing.T) {
	strategy := New()
	err := strategy.Init(decimal.NewFromInt(1), decimal.NewFromInt(1000))
	assert.NoError(t, err)

	// 设置初始状态为空头
	strategy.position = "short"

	// 创建足够的K线数据
	klines := generateTestKlines(time.Now().Unix(), 30)

	// 更新前20个K线以初始化指标
	for i := 0; i < 20; i++ {
		strategy.Update(klines[i])
	}

	// 确保退出通道已经计算
	assert.NotNil(t, strategy.exitChannel)

	// 修改下一个K线，使其突破上轨（设置一个非常高的价格）
	klines[20].H = decimal.NewFromFloat(120.0)
	klines[20].C = decimal.NewFromFloat(115.0)

	// 更新策略，应该生成买入信号（平空）
	signal, err := strategy.Update(klines[20])
	assert.NoError(t, err)
	assert.Equal(t, types.SignalTypeBuy, signal.Type)
	assert.Equal(t, "none", strategy.position)
}

// 测试计算仓位大小
func TestCalculatePositionSize(t *testing.T) {
	strategy := New()
	strategy.Init(decimal.Zero, decimal.NewFromInt(1000))

	// 创建足够的K线数据并初始化指标
	klines := generateTestKlines(time.Now().Unix(), 30)
	for i := 0; i < 20; i++ {
		strategy.Update(klines[i])
	}

	// 计算仓位大小
	positionSize := strategy.calculatePositionSize(decimal.NewFromInt(100))

	// 验证结果是否合理（不为零且不过大）
	assert.True(t, positionSize.GreaterThan(decimal.Zero))
	assert.True(t, positionSize.LessThan(decimal.NewFromInt(10)), "仓位不应过大")
}

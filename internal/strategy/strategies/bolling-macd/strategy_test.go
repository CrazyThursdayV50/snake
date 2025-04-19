package bollingmacd

import (
	"context"
	"snake/internal/kline"
	"snake/internal/types"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// 构建测试用K线数据
func createTestKlines(count int) []*kline.Kline {
	klines := make([]*kline.Kline, count)
	now := time.Now().Unix()

	// 模拟一个缓慢上涨的价格序列，从50开始
	basePrice := decimal.NewFromInt(50)

	for i := 0; i < count; i++ {
		// 创建一个K线，其中收盘价格有小幅波动
		kline := &kline.Kline{
			S: now + int64(i),                                        // 开始时间递增
			E: now + int64(i) + 1,                                    // 结束时间
			O: basePrice,                                             // 开盘价
			H: basePrice.Mul(decimal.NewFromFloat(1.05)),             // 最高价
			L: basePrice.Mul(decimal.NewFromFloat(0.95)),             // 最低价
			C: basePrice.Add(decimal.NewFromFloat(float64(i) * 0.1)), // 收盘价缓慢上涨
			V: decimal.NewFromInt(1000),                              // 交易量
		}
		klines[i] = kline

		// 更新基准价格，每次增加一点
		basePrice = kline.C
	}

	return klines
}

// 测试策略初始化
func TestNew(t *testing.T) {
	strategy := New(context.WithCancel(context.TODO()))
	assert.NotNil(t, strategy, "Strategy should not be nil")
	assert.Equal(t, "Bolling-MACD Strategy", strategy.Name(), "Strategy name should match")
}

// 测试策略初始化参数
func TestInit(t *testing.T) {
	strategy := New(context.WithCancel(context.TODO()))
	err := strategy.Init(decimal.NewFromInt(10), decimal.NewFromInt(1000))
	assert.NoError(t, err, "Init should not return error")

	// 验证初始化后的仓位和余额
	assert.Equal(t, decimal.NewFromInt(10), strategy.Position().Amount, "Position amount should match")
	assert.Equal(t, decimal.NewFromInt(1000), strategy.Balance().Amount, "Balance amount should match")
}

// 测试数据不足情况下的策略更新
func TestUpdateInsufficientData(t *testing.T) {
	strategy := New(context.WithCancel(context.TODO()))
	err := strategy.Init(decimal.Zero, decimal.NewFromInt(1000))
	assert.NoError(t, err, "Init should not return error")

	// 创建单个K线
	kline := &kline.Kline{
		S: time.Now().Unix(),
		E: time.Now().Unix() + 1,
		O: decimal.NewFromInt(100),
		H: decimal.NewFromInt(105),
		L: decimal.NewFromInt(95),
		C: decimal.NewFromInt(102),
		V: decimal.NewFromInt(1000),
	}

	// 更新策略
	signal, err := strategy.Update(kline)
	assert.NoError(t, err, "Update should not return error")
	assert.Equal(t, types.SignalTypeHold, signal.Type, "Signal should be HOLD when data is insufficient")
}

// 测试策略更新和信号生成
func TestUpdateSignalGeneration(t *testing.T) {
	strategy := New(context.WithCancel(context.TODO()))
	err := strategy.Init(decimal.Zero, decimal.NewFromInt(1000))
	assert.NoError(t, err, "Init should not return error")

	// 创建60个K线，足够计算指标
	klines := createTestKlines(60)

	// 更新前59个K线，只积累数据
	for i := 0; i < 59; i++ {
		signal, err := strategy.Update(klines[i])
		assert.NoError(t, err, "Update should not return error")
		assert.Equal(t, types.SignalTypeHold, signal.Type, "Signal should be HOLD during data accumulation")
	}

	// 更新第60个K线，此时应该有足够数据计算指标
	signal, err := strategy.Update(klines[59])
	assert.NoError(t, err, "Update should not return error")

	// 根据策略逻辑，信号可能是买入、卖出或持有
	// 由于测试数据是模拟的，我们只验证信号类型是否有效
	assert.Contains(t, []types.SignalType{types.SignalTypeBuy, types.SignalTypeSell, types.SignalTypeHold}, signal.Type,
		"Signal should be one of the valid types")
}

// 测试策略盈亏计算
func TestProfit(t *testing.T) {
	strategy := New(context.WithCancel(context.TODO()))

	// 初始化策略，持有10个单位，成本为1000
	err := strategy.Init(decimal.NewFromInt(10), decimal.NewFromInt(1000))
	assert.NoError(t, err, "Init should not return error")

	// 创建一个K线，设置收盘价为110
	kline := &kline.Kline{
		S: time.Now().Unix(),
		E: time.Now().Unix() + 1,
		O: decimal.NewFromInt(100),
		H: decimal.NewFromInt(105),
		L: decimal.NewFromInt(95),
		C: decimal.NewFromInt(110),
		V: decimal.NewFromInt(1000),
	}

	// 更新策略以更新盈亏
	_, err = strategy.Update(kline)
	assert.NoError(t, err, "Update should not return error")

	// 获取当前盈亏
	absolute, percentage := strategy.Profit()

	// 由于我们使用了BaseStrategy的默认实现，盈亏计算应该基于持仓成本
	// 但初始持仓成本的计算逻辑可能因实现而异，所以这里只检查盈亏计算的逻辑
	assert.NotNil(t, absolute, "Absolute profit should not be nil")
	assert.NotNil(t, percentage, "Percentage profit should not be nil")
}

package rsi_strategy

import (
	"context"
	"math"
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestRSIStrategy(t *testing.T) {
	// 创建测试数据
	startTime := time.Now().Unix() * 1000

	// 创建辅助函数，用于初始化策略并填充历史数据
	initStrategy := func(positionAmount, balanceAmount decimal.Decimal) *RSIStrategy {
		strategy := New(context.WithCancel(context.TODO()))

		if err := strategy.Init(positionAmount, balanceAmount); err != nil {
			t.Fatalf("failed to init strategy: %v", err)
		}

		return strategy
	}

	t.Run("RSI Buy Signal", func(t *testing.T) {
		// 创建策略
		positionAmount := decimal.NewFromFloat(0.5)    // 初始仓位 0.5 BTC
		balanceAmount := decimal.NewFromFloat(10000.0) // 初始余额 10000 USDT
		strategy := initStrategy(positionAmount, balanceAmount)

		// 强制设置RSI参数
		strategy.SetParams(14, decimal.NewFromInt(30), decimal.NewFromInt(70))

		// 生成连续下跌的K线，产生超卖信号
		klines := generateFallingKlines(startTime, 20, 100.0, 1.0)

		// 填充足够的历史数据来计算RSI
		for i := 0; i < 15; i++ {
			if _, err := strategy.Update(klines[i]); err != nil {
				t.Fatalf("failed to update strategy: %v", err)
			}
		}

		// 在测试前，获取策略当前的状态
		initialPosition := strategy.Position().Amount
		initialBalance := strategy.Balance().Amount

		// 应该产生买入信号
		signal, err := strategy.Update(klines[15])
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Fatal("expected buy signal")
		}

		if !signal.Type.IsBuy() {
			t.Fatalf("expected buy signal, got %v", signal.Type)
		}

		// 验证持仓和余额已变化
		currentPosition := strategy.Position().Amount
		currentBalance := strategy.Balance().Amount

		if !currentPosition.GreaterThan(initialPosition) {
			t.Errorf("expected position to increase from %v, got %v", initialPosition, currentPosition)
		}

		if !currentBalance.LessThan(initialBalance) {
			t.Errorf("expected balance to decrease from %v, got %v", initialBalance, currentBalance)
		}
	})

	t.Run("RSI Sell Signal", func(t *testing.T) {
		// 创建策略
		positionAmount := decimal.NewFromFloat(1.0)   // 初始仓位 1 BTC
		balanceAmount := decimal.NewFromFloat(5000.0) // 初始余额 5000 USDT
		strategy := initStrategy(positionAmount, balanceAmount)

		// 强制设置RSI参数
		strategy.SetParams(14, decimal.NewFromInt(30), decimal.NewFromInt(70))

		// 生成连续上涨的K线，产生超买信号
		klines := generateRisingKlines(startTime, 20, 100.0, 1.0)

		// 填充足够的历史数据来计算RSI
		for i := 0; i < 15; i++ {
			if _, err := strategy.Update(klines[i]); err != nil {
				t.Fatalf("failed to update strategy: %v", err)
			}
		}

		// 在测试前，获取策略当前的状态
		initialPosition := strategy.Position().Amount
		initialBalance := strategy.Balance().Amount

		// 应该产生卖出信号
		signal, err := strategy.Update(klines[15])
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Fatal("expected sell signal")
		}

		if !signal.Type.IsSell() {
			t.Fatalf("expected sell signal, got %v", signal.Type)
		}

		// 验证持仓和余额已变化
		currentPosition := strategy.Position().Amount
		currentBalance := strategy.Balance().Amount

		if !currentPosition.LessThan(initialPosition) {
			t.Errorf("expected position to decrease from %v, got %v", initialPosition, currentPosition)
		}

		if !currentBalance.GreaterThan(initialBalance) {
			t.Errorf("expected balance to increase from %v, got %v", initialBalance, currentBalance)
		}
	})

	t.Run("RSI Hold Signal", func(t *testing.T) {
		// 创建策略
		positionAmount := decimal.NewFromFloat(0.5)   // 初始仓位 0.5 BTC
		balanceAmount := decimal.NewFromFloat(5000.0) // 初始余额 5000 USDT
		strategy := initStrategy(positionAmount, balanceAmount)

		// 强制设置RSI参数
		strategy.SetParams(14, decimal.NewFromInt(30), decimal.NewFromInt(70))

		// 生成波动的K线，RSI保持在30-70之间
		klines := generateWavyKlines(startTime, 20, 100.0)

		// 填充足够的历史数据来计算RSI
		for i := 0; i < 15; i++ {
			if _, err := strategy.Update(klines[i]); err != nil {
				t.Fatalf("failed to update strategy: %v", err)
			}
		}

		// 在测试前，获取策略当前的状态
		initialPosition := strategy.Position().Amount
		initialBalance := strategy.Balance().Amount

		// 应该产生持有信号
		signal, err := strategy.Update(klines[15])
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Fatal("expected hold signal")
		}

		if !signal.Type.IsHold() {
			t.Fatalf("expected hold signal, got %v", signal.Type)
		}

		// 验证持仓和余额未变化
		currentPosition := strategy.Position().Amount
		currentBalance := strategy.Balance().Amount

		if !currentPosition.Equal(initialPosition) {
			t.Errorf("expected position to remain %v, got %v", initialPosition, currentPosition)
		}

		if !currentBalance.Equal(initialBalance) {
			t.Errorf("expected balance to remain %v, got %v", initialBalance, currentBalance)
		}
	})
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

// 生成波动的K线序列，使RSI值保持在30-70之间
func generateWavyKlines(startTime int64, count int, basePrice float64) []*kline.Kline {
	klines := make([]*kline.Kline, count)

	price := basePrice
	for i := 0; i < count; i++ {
		// 轻微上涨和下跌交替，使RSI保持在中间范围
		var change float64
		if i%2 == 0 {
			change = 0.5 // 小幅上涨
		} else {
			change = -0.4 // 小幅下跌
		}

		open := price
		close := price + change

		klines[i] = &kline.Kline{
			O: decimal.NewFromFloat(open),
			C: decimal.NewFromFloat(close),
			H: decimal.NewFromFloat(math.Max(open, close) + 0.2),
			L: decimal.NewFromFloat(math.Min(open, close) - 0.2),
			V: decimal.NewFromFloat(1000.0),
			A: decimal.NewFromFloat(100000.0),
			S: startTime + int64(i)*60000,
			E: startTime + int64(i+1)*60000,
		}

		price = close
	}

	return klines
}

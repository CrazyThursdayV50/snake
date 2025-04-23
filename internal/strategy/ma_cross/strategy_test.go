package ma_cross

import (
	"snake/internal/indicates/ma"
	"snake/internal/models"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMACrossStrategy(t *testing.T) {
	// 创建测试数据
	startTime := time.Now().Unix() * 1000
	klines := []*models.Kline{
		{
			O: decimal.NewFromFloat(100.0),
			C: decimal.NewFromFloat(101.0),
			H: decimal.NewFromFloat(102.0),
			L: decimal.NewFromFloat(99.0),
			V: decimal.NewFromFloat(1000.0),
			A: decimal.NewFromFloat(100000.0),
			S: startTime,
			E: startTime + 59999,
		},
		{
			O: decimal.NewFromFloat(101.0),
			C: decimal.NewFromFloat(103.0),
			H: decimal.NewFromFloat(104.0),
			L: decimal.NewFromFloat(100.0),
			V: decimal.NewFromFloat(2000.0),
			A: decimal.NewFromFloat(200000.0),
			S: startTime + 60000,
			E: startTime + 119999,
		},
		{
			O: decimal.NewFromFloat(103.0),
			C: decimal.NewFromFloat(105.0),
			H: decimal.NewFromFloat(106.0),
			L: decimal.NewFromFloat(102.0),
			V: decimal.NewFromFloat(3000.0),
			A: decimal.NewFromFloat(300000.0),
			S: startTime + 120000,
			E: startTime + 179999,
		},
	}

	// 创建 MA20 和 MA60
	ma20 := ma.New(klines[0], klines[1], klines[2])
	ma60 := ma.New(klines[0], klines[1], klines[2])

	// 测试用例 1: 价格高于 MA20，应该卖出
	t.Run("price above MA20", func(t *testing.T) {
		// 创建策略
		strategy := New()
		positionAmount := decimal.NewFromFloat(1.0)  // 初始仓位 1 BTC
		balanceAmount := decimal.NewFromFloat(1000.0) // 初始余额 1000 USDT
		if err := strategy.Init(positionAmount, balanceAmount); err != nil {
			t.Fatalf("failed to init strategy: %v", err)
		}

		// 设置当前价格高于 MA20
		currentKline := &models.Kline{
			O: decimal.NewFromFloat(106.0),
			C: decimal.NewFromFloat(107.0),
			H: decimal.NewFromFloat(108.0),
			L: decimal.NewFromFloat(105.0),
			V: decimal.NewFromFloat(4000.0),
			A: decimal.NewFromFloat(400000.0),
			S: startTime + 180000,
			E: startTime + 239999,
		}

		signal, err := strategy.Update(currentKline, ma20, ma60)
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Error("expected sell signal")
			return
		}

		if !signal.Type.IsSell() {
			t.Errorf("expected sell signal, got %v", signal.Type)
		}

		expectedAmount := positionAmount.Mul(decimal.NewFromFloat(0.05))
		if !signal.Amount.Equal(expectedAmount) {
			t.Errorf("expected amount %v, got %v", expectedAmount, signal.Amount)
		}

		// 检查余额和仓位是否正确更新
		expectedPosition := positionAmount.Sub(expectedAmount)
		if !strategy.Position().Amount.Equal(expectedPosition) {
			t.Errorf("expected position %v, got %v", expectedPosition, strategy.Position().Amount)
		}

		expectedBalance := balanceAmount.Add(expectedAmount.Mul(currentKline.C))
		if !strategy.Balance().Amount.Equal(expectedBalance) {
			t.Errorf("expected balance %v, got %v", expectedBalance, strategy.Balance().Amount)
		}
	})

	// 测试用例 2: 价格低于 MA60，应该买入
	t.Run("price below MA60", func(t *testing.T) {
		// 创建策略
		strategy := New()
		positionAmount := decimal.NewFromFloat(0.95)  // 初始仓位 0.95 BTC
		balanceAmount := decimal.NewFromFloat(1053.0) // 初始余额 1053 USDT
		if err := strategy.Init(positionAmount, balanceAmount); err != nil {
			t.Fatalf("failed to init strategy: %v", err)
		}

		// 设置当前价格低于 MA60
		currentKline := &models.Kline{
			O: decimal.NewFromFloat(95.0),
			C: decimal.NewFromFloat(96.0),
			H: decimal.NewFromFloat(97.0),
			L: decimal.NewFromFloat(94.0),
			V: decimal.NewFromFloat(3000.0),
			A: decimal.NewFromFloat(300000.0),
			S: startTime + 240000,
			E: startTime + 299999,
		}

		signal, err := strategy.Update(currentKline, ma20, ma60)
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Error("expected buy signal")
			return
		}

		if !signal.Type.IsBuy() {
			t.Errorf("expected buy signal, got %v", signal.Type)
		}

		expectedAmount := positionAmount.Mul(decimal.NewFromFloat(0.05))
		if !signal.Amount.Equal(expectedAmount) {
			t.Errorf("expected amount %v, got %v", expectedAmount, signal.Amount)
		}

		// 检查余额和仓位是否正确更新
		expectedPosition := positionAmount.Add(expectedAmount)
		if !strategy.Position().Amount.Equal(expectedPosition) {
			t.Errorf("expected position %v, got %v", expectedPosition, strategy.Position().Amount)
		}

		expectedBalance := balanceAmount.Sub(expectedAmount.Mul(currentKline.C))
		if !strategy.Balance().Amount.Equal(expectedBalance) {
			t.Errorf("expected balance %v, got %v", expectedBalance, strategy.Balance().Amount)
		}
	})

	// 测试用例 3: 价格在 MA20 和 MA60 之间，应该持有
	t.Run("price between MA20 and MA60", func(t *testing.T) {
		// 创建策略
		strategy := New()
		positionAmount := decimal.NewFromFloat(1.0)  // 初始仓位 1 BTC
		balanceAmount := decimal.NewFromFloat(1000.0) // 初始余额 1000 USDT
		if err := strategy.Init(positionAmount, balanceAmount); err != nil {
			t.Fatalf("failed to init strategy: %v", err)
		}

		// 设置当前价格在 MA20 和 MA60 之间
		currentKline := &models.Kline{
			O: decimal.NewFromFloat(102.0),
			C: decimal.NewFromFloat(103.0),
			H: decimal.NewFromFloat(104.0),
			L: decimal.NewFromFloat(101.0),
			V: decimal.NewFromFloat(3000.0),
			A: decimal.NewFromFloat(300000.0),
			S: startTime + 300000,
			E: startTime + 359999,
		}

		signal, err := strategy.Update(currentKline, ma20, ma60)
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Error("expected hold signal")
			return
		}

		if !signal.Type.IsHold() {
			t.Errorf("expected hold signal, got %v", signal.Type)
		}
	})
} 
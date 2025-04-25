package ma_cross

import (
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMACrossStrategy(t *testing.T) {
	// 创建测试数据
	startTime := time.Now().Unix() * 1000

	// 生成足够多的K线数据用于测试
	klines := generateTestKlines(startTime, 60)

	// 创建辅助函数，用于初始化策略并填充历史数据
	initStrategy := func(positionAmount, balanceAmount decimal.Decimal) *MACrossStrategy {
		strategy := New()
		// 使用较小的MA周期值，以便测试能够正确运行
		strategy.ma20Period = 20
		strategy.ma60Period = 60

		if err := strategy.Init(positionAmount, balanceAmount); err != nil {
			t.Fatalf("failed to init strategy: %v", err)
		}

		// 填充历史数据
		for _, k := range klines {
			strategy.Update(k)
		}

		return strategy
	}

	// 测试用例 1: 价格高于 MA20，应该卖出
	t.Run("price above MA20", func(t *testing.T) {
		// 创建策略
		positionAmount := decimal.NewFromFloat(1.0)   // 初始仓位 1 BTC
		balanceAmount := decimal.NewFromFloat(1000.0) // 初始余额 1000 USDT
		strategy := initStrategy(positionAmount, balanceAmount)

		// 获取当前的MA20值，设置当前价格高于MA20
		ma20Value := strategy.ma20.Price.InexactFloat64()
		currentPrice := ma20Value + 1.0 // 确保价格高于MA20

		currentKline := &kline.Kline{
			O: decimal.NewFromFloat(currentPrice - 1.0),
			C: decimal.NewFromFloat(currentPrice),
			H: decimal.NewFromFloat(currentPrice + 1.0),
			L: decimal.NewFromFloat(currentPrice - 2.0),
			V: decimal.NewFromFloat(4000.0),
			A: decimal.NewFromFloat(400000.0),
			S: startTime + int64(len(klines))*60000,
			E: startTime + int64(len(klines)+1)*60000,
		}

		// 在测试前，获取策略当前的状态
		initialPosition := strategy.Position().Amount
		initialBalance := strategy.Balance().Amount

		// 预期的交易数量是当前持仓的5%
		expectedTradeAmount := initialPosition.Mul(decimal.NewFromFloat(0.05))

		signal, err := strategy.Update(currentKline)
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Fatal("expected sell signal")
		}

		if !signal.Type.IsSell() {
			t.Fatalf("expected sell signal, got %v", signal.Type)
		}

		// 验证信号的数量
		if !signal.Amount.Equal(expectedTradeAmount) {
			t.Errorf("expected trade amount %v, got %v", expectedTradeAmount, signal.Amount)
		}

		// 预期的持仓和余额
		expectedPosition := initialPosition.Sub(expectedTradeAmount)
		expectedBalance := initialBalance.Add(expectedTradeAmount.Mul(currentKline.C))

		// 验证持仓和余额
		currentPosition := strategy.Position().Amount
		currentBalance := strategy.Balance().Amount

		if !currentPosition.Equal(expectedPosition) {
			t.Errorf("expected position %v, got %v", expectedPosition, currentPosition)
		}

		if !currentBalance.Equal(expectedBalance) {
			t.Errorf("expected balance %v, got %v", expectedBalance, currentBalance)
		}
	})

	// 测试用例 2: 价格低于 MA60，应该买入
	t.Run("price below MA60", func(t *testing.T) {
		// 创建策略
		positionAmount := decimal.NewFromFloat(0.95)  // 初始仓位 0.95 BTC
		balanceAmount := decimal.NewFromFloat(1053.0) // 初始余额 1053 USDT
		strategy := initStrategy(positionAmount, balanceAmount)

		// 获取当前的MA60值，设置当前价格低于MA60
		ma60Value := strategy.ma60.Price.InexactFloat64()
		currentPrice := ma60Value - 1.0 // 确保价格低于MA60

		currentKline := &kline.Kline{
			O: decimal.NewFromFloat(currentPrice + 1.0),
			C: decimal.NewFromFloat(currentPrice),
			H: decimal.NewFromFloat(currentPrice + 2.0),
			L: decimal.NewFromFloat(currentPrice - 1.0),
			V: decimal.NewFromFloat(3000.0),
			A: decimal.NewFromFloat(300000.0),
			S: startTime + int64(len(klines)+1)*60000,
			E: startTime + int64(len(klines)+2)*60000,
		}

		// 在测试前，获取策略当前的状态
		initialPosition := strategy.Position().Amount
		initialBalance := strategy.Balance().Amount

		// 预期的交易数量是当前持仓的5%
		expectedTradeAmount := initialPosition.Mul(decimal.NewFromFloat(0.05))

		signal, err := strategy.Update(currentKline)
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Fatal("expected buy signal")
		}

		if !signal.Type.IsBuy() {
			t.Fatalf("expected buy signal, got %v", signal.Type)
		}

		// 验证信号的数量
		if !signal.Amount.Equal(expectedTradeAmount) {
			t.Errorf("expected trade amount %v, got %v", expectedTradeAmount, signal.Amount)
		}

		// 预期的持仓和余额
		expectedPosition := initialPosition.Add(expectedTradeAmount)
		expectedBalance := initialBalance.Sub(expectedTradeAmount.Mul(currentKline.C))

		// 验证持仓和余额
		currentPosition := strategy.Position().Amount
		currentBalance := strategy.Balance().Amount

		if !currentPosition.Equal(expectedPosition) {
			t.Errorf("expected position %v, got %v", expectedPosition, currentPosition)
		}

		if !currentBalance.Equal(expectedBalance) {
			t.Errorf("expected balance %v, got %v", expectedBalance, currentBalance)
		}
	})

	// 测试用例 3: 价格在 MA20 和 MA60 之间，应该持有
	t.Run("price between MA20 and MA60", func(t *testing.T) {
		// 创建策略
		positionAmount := decimal.NewFromFloat(1.0)   // 初始仓位 1 BTC
		balanceAmount := decimal.NewFromFloat(1000.0) // 初始余额 1000 USDT
		strategy := initStrategy(positionAmount, balanceAmount)

		// 获取当前的MA20和MA60值，设置当前价格在两者之间
		ma20Value := strategy.ma20.Price.InexactFloat64()
		ma60Value := strategy.ma60.Price.InexactFloat64()

		// 如果MA20小于MA60，则交换它们的值
		if ma20Value < ma60Value {
			ma20Value, ma60Value = ma60Value, ma20Value
		}

		// 确保价格严格在MA20和MA60之间，而不是等于其中之一
		currentPrice := ma60Value + (ma20Value-ma60Value)*0.5

		// 打印调试信息
		t.Logf("MA20: %v, MA60: %v, 计算得到的中间价格: %v", ma20Value, ma60Value, currentPrice)

		currentKline := &kline.Kline{
			O: decimal.NewFromFloat(currentPrice - 1.0),
			C: decimal.NewFromFloat(currentPrice),
			H: decimal.NewFromFloat(currentPrice + 1.0),
			L: decimal.NewFromFloat(currentPrice - 1.0),
			V: decimal.NewFromFloat(3000.0),
			A: decimal.NewFromFloat(300000.0),
			S: startTime + int64(len(klines)+2)*60000,
			E: startTime + int64(len(klines)+3)*60000,
		}

		// 在测试前，获取策略当前的状态
		initialPosition := strategy.Position().Amount
		initialBalance := strategy.Balance().Amount

		signal, err := strategy.Update(currentKline)
		if err != nil {
			t.Fatalf("failed to update strategy: %v", err)
		}

		if signal == nil {
			t.Fatal("expected hold signal")
		}

		if !signal.Type.IsHold() {
			t.Fatalf("expected hold signal, got %v", signal.Type)
		}

		// 验证持仓和余额没有变化
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

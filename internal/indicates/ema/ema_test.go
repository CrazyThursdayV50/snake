package ema

import (
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNew(t *testing.T) {
	builder := New(5)
	if builder == nil {
		t.Fatal("failed to create EMA builder")
	}
	if builder.count != 5 {
		t.Errorf("expected count 5, got %d", builder.count)
	}
}

func TestBuild(t *testing.T) {
	t.Run("build ema with valid klines", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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

		ema := New(3).Klines(klines).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		if ema.count != 3 {
			t.Errorf("expected count 3, got %d", ema.count)
		}

		expectedAlpha := decimal.NewFromFloat(2.0).Div(decimal.NewFromInt(4))
		if !ema.alpha.Equal(expectedAlpha) {
			t.Errorf("expected alpha %v, got %v", expectedAlpha, ema.alpha)
		}

		expectedAlpha1 := decimal.NewFromFloat(0.5)
		if !ema.alpha1.Equal(expectedAlpha1) {
			t.Errorf("expected alpha1 %v, got %v", expectedAlpha1, ema.alpha1)
		}

		if ema.Timestamp != klines[2].S {
			t.Errorf("expected timestamp %d, got %d", klines[2].S, ema.Timestamp)
		}

		if !ema.CurrentPrice.Equal(klines[2].C) {
			t.Errorf("expected current price %v, got %v", klines[2].C, ema.CurrentPrice)
		}
	})

	t.Run("build ema with insufficient klines", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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
		}

		// 尝试用 2 个 kline 构建 count=3 的 EMA
		// 由于 Klines 方法会检查长度，所以不会设置 klines
		builder := New(3).Klines(klines)
		ema := builder.Build()
		if ema != nil {
			t.Error("expected nil for insufficient klines")
		}
	})
}

func TestEMACalculation(t *testing.T) {
	t.Run("ema calculation accuracy", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
			{
				O: decimal.NewFromFloat(100.0),
				C: decimal.NewFromFloat(100.0),
				H: decimal.NewFromFloat(100.0),
				L: decimal.NewFromFloat(100.0),
				V: decimal.NewFromFloat(1000.0),
				A: decimal.NewFromFloat(100000.0),
				S: startTime,
				E: startTime + 59999,
			},
			{
				O: decimal.NewFromFloat(100.0),
				C: decimal.NewFromFloat(110.0),
				H: decimal.NewFromFloat(110.0),
				L: decimal.NewFromFloat(100.0),
				V: decimal.NewFromFloat(2000.0),
				A: decimal.NewFromFloat(200000.0),
				S: startTime + 60000,
				E: startTime + 119999,
			},
			{
				O: decimal.NewFromFloat(110.0),
				C: decimal.NewFromFloat(120.0),
				H: decimal.NewFromFloat(120.0),
				L: decimal.NewFromFloat(110.0),
				V: decimal.NewFromFloat(3000.0),
				A: decimal.NewFromFloat(300000.0),
				S: startTime + 120000,
				E: startTime + 179999,
			},
		}

		// 测试 count=2 的 EMA
		// alpha = 2/(2+1) = 0.6667, alpha1 = 1-0.6667 = 0.3333
		// 第一个 EMA = 100.0
		// 第二个 EMA = 110.0 * 0.6667 + 100.0 * 0.3333 ≈ 106.67
		// 第三个 EMA = 120.0 * 0.6667 + 106.67 * 0.3333 ≈ 115.56
		ema := New(2).Klines([]*kline.Kline{klines[0], klines[1]}).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		// 验证 EMA 值在合理范围内
		if ema.Value.LessThan(decimal.NewFromFloat(106.0)) || ema.Value.GreaterThan(decimal.NewFromFloat(107.0)) {
			t.Errorf("EMA value %v is not in expected range [106.0, 107.0]", ema.Value)
		}

		// 测试下一个 EMA
		nextEMA := ema.Next(klines[2])
		if nextEMA == nil {
			t.Fatal("failed to calculate next EMA")
		}

		// 验证下一个 EMA 值在合理范围内
		if nextEMA.Value.LessThan(decimal.NewFromFloat(115.0)) || nextEMA.Value.GreaterThan(decimal.NewFromFloat(116.0)) {
			t.Errorf("Next EMA value %v is not in expected range [115.0, 116.0]", nextEMA.Value)
		}

		// 验证 EMA 值随着价格上涨而上涨
		if !nextEMA.Value.GreaterThan(ema.Value) {
			t.Errorf("Next EMA value %v should be greater than current EMA value %v", nextEMA.Value, ema.Value)
		}
	})
}

func TestNextKline(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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

		ema := New(2).Klines([]*kline.Kline{klines[0], klines[1]}).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		nextEMA := ema.Next(klines[2])
		if nextEMA == nil {
			t.Fatal("failed to calculate next EMA")
		}

		if nextEMA.Timestamp != klines[2].S {
			t.Errorf("expected timestamp %d, got %d", klines[2].S, nextEMA.Timestamp)
		}

		if nextEMA.count != 2 {
			t.Errorf("expected count 2, got %d", nextEMA.count)
		}

		if !nextEMA.CurrentPrice.Equal(klines[2].C) {
			t.Errorf("expected current price %v, got %v", klines[2].C, nextEMA.CurrentPrice)
		}
	})

	t.Run("earlier kline", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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
		}

		ema := New(2).Klines([]*kline.Kline{klines[0], klines[1]}).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		// 尝试计算下一个 EMA（使用第一个 Kline，时间早于当前 EMA）
		nextEMA := ema.Next(klines[0])
		if nextEMA != nil {
			t.Error("expected nil for earlier kline")
		}
	})

	t.Run("non-continuous kline", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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
				S: startTime + 180000, // 跳过了一个时间间隔
				E: startTime + 239999,
			},
		}

		ema := New(2).Klines([]*kline.Kline{klines[0], klines[1]}).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		// 尝试计算下一个 EMA（使用跳过时间间隔的 Kline）
		nextEMA := ema.Next(klines[2])
		if nextEMA != nil {
			t.Error("expected nil for non-continuous kline")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update current kline", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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
		}

		ema := New(2).Klines([]*kline.Kline{klines[0], klines[1]}).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		updatedKline := &kline.Kline{
			O: decimal.NewFromFloat(101.0),
			C: decimal.NewFromFloat(104.0),
			H: decimal.NewFromFloat(105.0),
			L: decimal.NewFromFloat(100.0),
			V: decimal.NewFromFloat(2500.0),
			A: decimal.NewFromFloat(250000.0),
			S: startTime + 60000,
			E: startTime + 119999,
		}

		updated := ema.Update(updatedKline)
		if !updated {
			t.Fatal("failed to update EMA")
		}

		if !ema.CurrentPrice.Equal(updatedKline.C) {
			t.Errorf("expected updated current price %v, got %v", updatedKline.C, ema.CurrentPrice)
		}

		if ema.Timestamp != updatedKline.S {
			t.Errorf("expected updated timestamp %d, got %d", updatedKline.S, ema.Timestamp)
		}
	})

	t.Run("update with different kline", func(t *testing.T) {
		startTime := time.Now().Unix() * 1000
		klines := []*kline.Kline{
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
		}

		ema := New(2).Klines([]*kline.Kline{klines[0], klines[1]}).Build()
		if ema == nil {
			t.Fatal("failed to create EMA")
		}

		differentKline := &kline.Kline{
			O: decimal.NewFromFloat(105.0),
			C: decimal.NewFromFloat(106.0),
			H: decimal.NewFromFloat(107.0),
			L: decimal.NewFromFloat(104.0),
			V: decimal.NewFromFloat(3000.0),
			A: decimal.NewFromFloat(300000.0),
			S: startTime + 120000, // 不同的时间戳
			E: startTime + 179999,
		}

		updated := ema.Update(differentKline)
		if updated {
			t.Error("expected false for different kline")
		}
	})
}

func TestNextEMA(t *testing.T) {
	currentEMA := decimal.NewFromFloat(100.0)
	alpha := decimal.NewFromFloat(0.5)
	alpha1 := decimal.NewFromFloat(0.5)
	price := decimal.NewFromFloat(110.0)

	result := nextEMA(currentEMA, alpha, alpha1, price)
	expected := decimal.NewFromFloat(105.0)

	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

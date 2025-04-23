package bollingband

import (
	"snake/internal/models"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNextKline(t *testing.T) {
	// 测试用例 1: 正常情况下的下一个 Kline
	t.Run("normal case", func(t *testing.T) {
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

		// 创建布林带
		bb := New(klines[0], klines[1])
		if bb == nil {
			t.Fatal("failed to create Bollinger Bands")
		}

		// 计算下一个布林带
		nextBB := bb.NextKline(klines[2])
		if nextBB == nil {
			t.Fatal("failed to calculate next Bollinger Bands")
		}

		// 验证结果
		expectedMA := decimal.NewFromFloat(104.0) // (103.0 + 105.0) / 2
		if !nextBB.MA.Equal(expectedMA) {
			t.Errorf("expected MA %v, got %v", expectedMA, nextBB.MA)
		}

		// 验证上轨和下轨
		// 标准差 = sqrt(((103-104)^2 + (105-104)^2)/2) = 1
		// 上轨 = 104 + 2*1 = 106
		// 下轨 = 104 - 2*1 = 102
		expectedUpper := decimal.NewFromFloat(106.0)
		expectedLower := decimal.NewFromFloat(102.0)
		if !nextBB.Upper.Equal(expectedUpper) {
			t.Errorf("expected upper band %v, got %v", expectedUpper, nextBB.Upper)
		}
		if !nextBB.Lower.Equal(expectedLower) {
			t.Errorf("expected lower band %v, got %v", expectedLower, nextBB.Lower)
		}

		if nextBB.Timestamp != klines[2].E {
			t.Errorf("expected timestamp %d, got %d", klines[2].E, nextBB.Timestamp)
		}
		if nextBB.count != 2 {
			t.Errorf("expected count 2, got %d", nextBB.count)
		}
	})

	// 测试用例 2: 传入时间早于当前布林带的 Kline
	t.Run("earlier kline", func(t *testing.T) {
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
		}

		// 创建布林带
		bb := New(klines[0], klines[1])
		if bb == nil {
			t.Fatal("failed to create Bollinger Bands")
		}

		// 尝试计算下一个布林带（使用第一个 Kline，时间早于当前布林带）
		nextBB := bb.NextKline(klines[0])
		if nextBB != nil {
			t.Error("expected nil for earlier kline")
		}
	})

	// 测试用例 3: 传入时间间隔不连续的 Kline
	t.Run("non-continuous kline", func(t *testing.T) {
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
				S: startTime + 180000, // 跳过了一个时间间隔
				E: startTime + 239999,
			},
		}

		// 创建布林带
		bb := New(klines[0], klines[1])
		if bb == nil {
			t.Fatal("failed to create Bollinger Bands")
		}

		// 尝试计算下一个布林带（使用跳过时间间隔的 Kline）
		nextBB := bb.NextKline(klines[2])
		if nextBB != nil {
			t.Error("expected nil for non-continuous kline")
		}
	})
} 
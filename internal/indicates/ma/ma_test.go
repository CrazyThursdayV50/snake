package ma

import (
	"snake/internal/kline"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNextKline(t *testing.T) {
	// 测试用例 1: 正常情况下的下一个 Kline
	t.Run("normal case", func(t *testing.T) {
		// 创建测试数据
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

		// 创建 MA
		ma := New(klines[0], klines[1])
		if ma == nil {
			t.Fatal("failed to create MA")
		}

		// 计算下一个 MA
		nextMA := ma.NextKline(klines[2])
		if nextMA == nil {
			t.Fatal("failed to calculate next MA")
		}

		// 验证结果
		expectedPrice := decimal.NewFromFloat(104.0) // (103.0 + 105.0) / 2
		if !nextMA.Price.Equal(expectedPrice) {
			t.Errorf("expected MA price %v, got %v", expectedPrice, nextMA.Price)
		}
		if nextMA.Timestamp != klines[2].E {
			t.Errorf("expected timestamp %d, got %d", klines[2].E, nextMA.Timestamp)
		}
		if nextMA.count != 2 {
			t.Errorf("expected count 2, got %d", nextMA.count)
		}
	})

	// 测试用例 2: 传入时间早于当前 MA 的 Kline
	t.Run("earlier kline", func(t *testing.T) {
		// 创建测试数据
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

		// 创建 MA
		ma := New(klines[0], klines[1])
		if ma == nil {
			t.Fatal("failed to create MA")
		}

		// 尝试计算下一个 MA（使用第一个 Kline，时间早于当前 MA）
		nextMA := ma.NextKline(klines[0])
		if nextMA != nil {
			t.Error("expected nil for earlier kline")
		}
	})

	// 测试用例 3: 传入时间间隔不连续的 Kline
	t.Run("non-continuous kline", func(t *testing.T) {
		// 创建测试数据
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

		// 创建 MA
		ma := New(klines[0], klines[1])
		if ma == nil {
			t.Fatal("failed to create MA")
		}

		// 尝试计算下一个 MA（使用跳过时间间隔的 Kline）
		nextMA := ma.NextKline(klines[2])
		if nextMA != nil {
			t.Error("expected nil for non-continuous kline")
		}
	})
}

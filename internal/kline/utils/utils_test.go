package utils

import (
	"snake/internal/kline"
	"snake/internal/kline/interval"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestFillKlines(t *testing.T) {
	// 创建一个测试用的 interval
	testInterval := interval.Min1()

	// 测试用例 1: 正常情况下的填充
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
				S: startTime + 120000, // 跳过了一个时间间隔
				E: startTime + 179999,
			},
		}

		// 计算结束时间（包含两个完整的时间间隔）
		to := startTime + 180000

		// 执行填充
		result := FillKlines(klines, testInterval, to)

		// 验证结果
		if len(result) != 4 {
			t.Errorf("expected 4 klines, got %d", len(result))
		}

		// 验证第一个 Kline（原始数据）
		if !result[0].O.Equal(decimal.NewFromFloat(100.0)) {
			t.Errorf("first kline open price mismatch")
		}
		if result[0].S != startTime {
			t.Errorf("first kline start time mismatch, expected %d, got %d", startTime, result[0].S)
		}
		if result[0].E != startTime+59999 {
			t.Errorf("first kline end time mismatch, expected %d, got %d", startTime+59999, result[0].E)
		}

		// 验证第二个 Kline（填充的数据）
		if !result[1].O.Equal(decimal.NewFromFloat(101.0)) {
			t.Errorf("second kline open price mismatch")
		}
		if !result[1].V.Equal(decimal.Zero) {
			t.Errorf("second kline volume should be zero")
		}
		if result[1].S != startTime+60000 {
			t.Errorf("second kline start time mismatch, expected %d, got %d", startTime+60000, result[1].S)
		}
		if result[1].E != startTime+119999 {
			t.Errorf("second kline end time mismatch, expected %d, got %d", startTime+119999, result[1].E)
		}

		// 验证第三个 Kline（原始数据）
		if !result[2].O.Equal(decimal.NewFromFloat(101.0)) {
			t.Errorf("third kline open price mismatch")
		}
		if result[2].S != startTime+120000 {
			t.Errorf("third kline start time mismatch, expected %d, got %d", startTime+120000, result[2].S)
		}
		if result[2].E != startTime+179999 {
			t.Errorf("third kline end time mismatch, expected %d, got %d", startTime+179999, result[2].E)
		}

		// 验证第四个 Kline（填充的数据）
		if !result[3].O.Equal(decimal.NewFromFloat(103.0)) {
			t.Errorf("fourth kline open price mismatch")
		}
		if !result[3].V.Equal(decimal.Zero) {
			t.Errorf("fourth kline volume should be zero")
		}
		if result[3].S != startTime+180000 {
			t.Errorf("fourth kline start time mismatch, expected %d, got %d", startTime+180000, result[3].S)
		}
		if result[3].E != startTime+239999 {
			t.Errorf("fourth kline end time mismatch, expected %d, got %d", startTime+239999, result[3].E)
		}
	})

	// 测试用例 2: 空输入
	t.Run("empty input", func(t *testing.T) {
		result := FillKlines([]*kline.Kline{}, testInterval, time.Now().Unix()*1000)
		if len(result) != 0 {
			t.Errorf("expected empty result, got %d klines", len(result))
		}
	})

	// 测试用例 3: 连续 Kline（不需要填充中间，但需要填充最后）
	t.Run("continuous klines", func(t *testing.T) {
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

		to := startTime + 120000
		result := FillKlines(klines, testInterval, to)

		if len(result) != 3 {
			t.Errorf("expected 3 klines, got %d", len(result))
		}

		// 验证第一个 Kline
		if result[0].S != startTime {
			t.Errorf("first kline start time mismatch, expected %d, got %d", startTime, result[0].S)
		}
		if result[0].E != startTime+59999 {
			t.Errorf("first kline end time mismatch, expected %d, got %d", startTime+59999, result[0].E)
		}

		// 验证第二个 Kline
		if result[1].S != startTime+60000 {
			t.Errorf("second kline start time mismatch, expected %d, got %d", startTime+60000, result[1].S)
		}
		if result[1].E != startTime+119999 {
			t.Errorf("second kline end time mismatch, expected %d, got %d", startTime+119999, result[1].E)
		}

		// 验证第三个 Kline（填充的数据）
		if !result[2].O.Equal(decimal.NewFromFloat(103.0)) {
			t.Errorf("third kline open price mismatch")
		}
		if !result[2].V.Equal(decimal.Zero) {
			t.Errorf("third kline volume should be zero")
		}
		if result[2].S != startTime+120000 {
			t.Errorf("third kline start time mismatch, expected %d, got %d", startTime+120000, result[2].S)
		}
		if result[2].E != startTime+179999 {
			t.Errorf("third kline end time mismatch, expected %d, got %d", startTime+179999, result[2].E)
		}
	})
}

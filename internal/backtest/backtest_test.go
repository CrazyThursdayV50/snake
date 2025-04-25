package backtest

import (
	"context"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	"snake/internal/strategy/strategies/ma_cross"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// mockKlineRepository 模拟 K 线数据仓库
type mockKlineRepository struct {
	klines []*models.Kline
}

func (r *mockKlineRepository) ListAll(ctx context.Context, interval interval.Interval) ([]*models.Kline, error) {
	return r.klines, nil
}

func (r *mockKlineRepository) Insert(ctx context.Context, interval interval.Interval, klines []*models.Kline) error {
	return nil
}

func (r *mockKlineRepository) First(ctx context.Context, interval interval.Interval) (*models.Kline, error) {
	if len(r.klines) == 0 {
		return nil, nil
	}
	return r.klines[0], nil
}

func (r *mockKlineRepository) Last(ctx context.Context, interval interval.Interval) (*models.Kline, error) {
	if len(r.klines) == 0 {
		return nil, nil
	}
	return r.klines[len(r.klines)-1], nil
}

func (r *mockKlineRepository) List(ctx context.Context, interval interval.Interval, from, to int64) ([]*models.Kline, error) {
	var result []*models.Kline
	for _, k := range r.klines {
		if k.OpenTs >= from && k.CloseTs <= to {
			result = append(result, k)
		}
	}
	return result, nil
}

func (r *mockKlineRepository) CheckMissing(ctx context.Context, interval interval.Interval, openTs []int64) ([]uint64, error) {
	return nil, nil
}

func TestBacktest(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	klines := []*models.Kline{
		{
			OpenTs:  now.Unix() * 1000,
			CloseTs: now.Add(time.Minute).Unix() * 1000,
			Open:    "100.0",
			Close:   "101.0",
			High:    "102.0",
			Low:     "99.0",
			Volume:  "1000.0",
			Amount:  "100000.0",
		},
		{
			OpenTs:  now.Add(time.Minute).Unix() * 1000,
			CloseTs: now.Add(2*time.Minute).Unix() * 1000,
			Open:    "101.0",
			Close:   "98.0", // 价格下跌，应该产生回撤
			High:    "101.0",
			Low:     "97.0",
			Volume:  "2000.0",
			Amount:  "200000.0",
		},
		{
			OpenTs:  now.Add(2*time.Minute).Unix() * 1000,
			CloseTs: now.Add(3*time.Minute).Unix() * 1000,
			Open:    "98.0",
			Close:   "95.0", // 继续下跌，回撤更大
			High:    "98.0",
			Low:     "94.0",
			Volume:  "3000.0",
			Amount:  "300000.0",
		},
		{
			OpenTs:  now.Add(3*time.Minute).Unix() * 1000,
			CloseTs: now.Add(4*time.Minute).Unix() * 1000,
			Open:    "95.0",
			Close:   "97.0", // 价格回升
			High:    "97.0",
			Low:     "95.0",
			Volume:  "4000.0",
			Amount:  "400000.0",
		},
	}

	// 创建回测配置
	config := &Config{
		InitialBalance:  decimal.NewFromFloat(1000.0),
		InitialPosition: decimal.NewFromFloat(1.0),
		Interval:        interval.Interval1m,
	}

	// 创建策略
	strategy := ma_cross.New()

	// 创建回测器
	backtest := New(config, &mockKlineRepository{klines: klines}, strategy)

	// 运行回测
	result, err := backtest.Run(context.Background())
	if err != nil {
		t.Fatalf("回测失败: %v", err)
	}

	// 验证结果
	if len(result) == 0 {
		t.Fatal("预期有回测结果，但实际没有数据")
	}

	// 获取最后一个K线对应的回测结果
	lastResult := result[len(result)-1]

	// 检查最终持仓和余额
	if lastResult.Balance.IsZero() && lastResult.PositionAmount.IsZero() {
		t.Error("预期有余额或持仓，但实际都为零")
	}

	// 验证最大回撤
	expectedMaxDrawdown := decimal.NewFromFloat(0.35) // 预期最大回撤约为 0.35%
	if lastResult.Drawdown.LessThan(expectedMaxDrawdown) {
		t.Errorf("预期最大回撤大于 %.2f%%，实际为 %.2f%%",
			expectedMaxDrawdown.InexactFloat64(), lastResult.Drawdown.InexactFloat64())
	}
}

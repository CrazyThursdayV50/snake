package ma_cross

import (
	"context"
	"snake/internal/indicates/ma"
	"snake/internal/kline"
	"snake/internal/strategy"

	"math"

	"github.com/shopspring/decimal"
)

// MACrossStrategy MA 交叉策略
type MACrossStrategy struct {
	*strategy.BaseStrategy
	// 保存历史K线用于计算MA
	historicalKlines []*kline.Kline
	// MA参数
	ma20Period int
	ma60Period int
	// MA指标
	ma20 *ma.MA
	ma60 *ma.MA
}

// New 创建 MA 交叉策略
func New(ctx context.Context, cancel context.CancelFunc) *MACrossStrategy {
	return &MACrossStrategy{
		BaseStrategy:     strategy.NewBaseStrategy(ctx,cancel, "MA Cross Strategy"),
		historicalKlines: make([]*kline.Kline, 0, 60), // 预分配足够容量
		ma20Period:       20,
		ma60Period:       60,
	}
}

// Update 更新策略状态
func (s *MACrossStrategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	// 添加新的K线到历史数据
	s.historicalKlines = append(s.historicalKlines, kline)

	// 如果历史数据不足以计算MA，则持有
	if len(s.historicalKlines) < s.ma60Period {
		return s.Hold(), nil
	}

	// 保持历史数据长度不超过需要的最大长度
	if len(s.historicalKlines) > s.ma60Period {
		s.historicalKlines = s.historicalKlines[len(s.historicalKlines)-s.ma60Period:]
	}

	// 计算MA20
	if len(s.historicalKlines) >= s.ma20Period {
		ma20Klines := s.historicalKlines[len(s.historicalKlines)-s.ma20Period:]
		s.ma20 = ma.New(ma20Klines...)
	}

	// 计算MA60
	s.ma60 = ma.New(s.historicalKlines...)

	// 如果无法计算MA，则持有
	if s.ma20 == nil || s.ma60 == nil {
		return s.Hold(), nil
	}

	// 计算交易数量（当前仓位的 5%）
	tradeAmount := s.Position().Amount.Mul(decimal.NewFromFloat(0.05))
	// 如果仓位为0，则使用余额的5%
	if s.Position().Amount.IsZero() {
		tradeAmount = s.Balance().Amount.Mul(decimal.NewFromFloat(0.05))
	}

	// 计算当前盈亏
	absolute, percentage := s.BaseStrategy.Profit()

	// 如果当前有持仓，打印盈亏信息
	if !s.Position().Amount.IsZero() {
		println("当前持仓盈亏：", absolute.String(), "USDT (", percentage.String(), "%)")
	}

	// 获取浮点数值用于比较
	currentPrice := kline.C.InexactFloat64()
	ma20Value := s.ma20.Price.InexactFloat64()
	ma60Value := s.ma60.Price.InexactFloat64()

	// 如果价格接近两个MA之间的中间值，返回持有信号
	if ma20Value > ma60Value {
		midpoint := ma60Value + (ma20Value-ma60Value)*0.5
		// 允许一些误差范围
		if math.Abs(currentPrice-midpoint) < 0.1 {
			return s.Hold(), nil
		}
	} else {
		midpoint := ma20Value + (ma60Value-ma20Value)*0.5
		// 允许一些误差范围
		if math.Abs(currentPrice-midpoint) < 0.1 {
			return s.Hold(), nil
		}
	}

	// 价格大于MA20时卖出
	if kline.C.GreaterThan(s.ma20.Price) {
		signal := s.Sell(tradeAmount, kline.C)
		if signal != nil {
			return signal, nil
		}
		return s.Hold(), nil
	}

	// 价格小于MA60时买入
	if kline.C.LessThan(s.ma60.Price) {
		signal := s.Buy(tradeAmount, kline.C)
		if signal != nil {
			return signal, nil
		}
		return s.Hold(), nil
	}

	// 其他情况，持有
	return s.Hold(), nil
}

// Profit 返回当前盈亏
func (s *MACrossStrategy) Profit() (absolute, percentage decimal.Decimal) {
	return s.BaseStrategy.Profit()
}

package ma_cross

import (
	"snake/internal/indicates/ma"
	"snake/internal/models"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// MACrossStrategy MA 交叉策略
type MACrossStrategy struct {
	*strategy.BaseStrategy
}

// New 创建 MA 交叉策略
func New() *MACrossStrategy {
	return &MACrossStrategy{
		BaseStrategy: strategy.NewBaseStrategy("MA Cross Strategy"),
	}
}

// Update 更新策略状态
func (s *MACrossStrategy) Update(kline *models.Kline, ma20 *ma.MA, ma60 *ma.MA) (*strategy.Signal, error) {
	if ma20 == nil || ma60 == nil {
		return s.Hold(), nil
	}

	// 计算交易数量（当前仓位的 5%）
	tradeAmount := s.Position().Amount.Mul(decimal.NewFromFloat(0.05))

	// 当前价格高于 MA20，卖出
	if kline.C.GreaterThan(ma20.Price) {
		signal := s.Sell(tradeAmount, kline.C)
		if signal != nil {
			return signal, nil
		}
		return s.Hold(), nil
	}

	// 当前价格低于 MA60，买入
	if kline.C.LessThan(ma60.Price) {
		signal := s.Buy(tradeAmount, kline.C)
		if signal != nil {
			return signal, nil
		}
		return s.Hold(), nil
	}

	// 价格在 MA20 和 MA60 之间，持有
	return s.Hold(), nil
} 
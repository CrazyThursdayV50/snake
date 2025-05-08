package macd

import (
	"context"
	"snake/internal/indicates/macd"
	"snake/internal/kline"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// MACDStrategy MACD策略
type MACDStrategy struct {
	*strategy.BaseStrategy
	lastMACD *macd.MACD
	// 保存历史K线用于计算指标
	historicalKlines []*kline.Kline
}

// New 创建新的MACD策略实例
func New(ctx context.Context, cancel context.CancelFunc) strategy.Strategy {
	return &MACDStrategy{
		BaseStrategy:     strategy.NewBaseStrategy(ctx, cancel, "MACD Strategy"),
		historicalKlines: make([]*kline.Kline, 0, 60), // 预分配足够容量
	}
}

// Update 更新策略状态并生成交易信号
func (s *MACDStrategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	// 计算当前盈亏
	absolute, percentage := s.BaseStrategy.Profit()

	// 如果当前有持仓，打印盈亏信息
	if !s.Position().Amount.IsZero() {
		println("当前持仓盈亏：", absolute.String(), "USDT (", percentage.String(), "%)")
	}

	// 添加新的K线到历史数据
	s.historicalKlines = append(s.historicalKlines, kline)

	// 如果历史数据不足以计算MACD，则返回持有信号
	if len(s.historicalKlines) < 26 { // 至少需要26根K线才能计算MACD
		return s.Hold(), nil
	}

	// 保持历史数据长度不超过需要的最大长度
	if len(s.historicalKlines) > 60 {
		s.historicalKlines = s.historicalKlines[len(s.historicalKlines)-60:]
	}

	// 计算MACD指标
	var m *macd.MACD
	if s.lastMACD == nil {
		m = macd.New(s.historicalKlines...)
	} else {
		m = s.lastMACD.NextKline(kline)
	}
	if m == nil {
		return s.Hold(), nil
	}
	s.lastMACD = m

	// 生成交易信号
	signal := s.generateSignal(m)
	if signal == nil {
		return s.Hold(), nil
	}

	return signal, nil
}

// generateSignal 根据MACD指标生成交易信号
func (s *MACDStrategy) generateSignal(m *macd.MACD) *strategy.Signal {
	// 当前没有持仓
	if s.Position().Amount.IsZero() {
		// MACD金叉（MACD线上穿信号线）
		if m.MACD.GreaterThan(m.Signal) && s.lastMACD.MACD.LessThanOrEqual(s.lastMACD.Signal) {
			// 生成买入信号，使用95%的余额，按市价买入
			return s.Buy(s.Balance().Amount.Mul(decimal.NewFromFloat(0.95)), m.LastPrice)
		}
	} else {
		// 当前有持仓
		// MACD死叉（MACD线下穿信号线）
		if m.MACD.LessThan(m.Signal) && s.lastMACD.MACD.GreaterThanOrEqual(s.lastMACD.Signal) {
			// 生成卖出信号，卖出全部持仓，按市价卖出
			return s.Sell(s.Position().Amount, m.LastPrice)
		}
	}

	return s.Hold()
}

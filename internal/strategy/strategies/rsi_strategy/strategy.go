package rsi_strategy

import (
	"snake/internal/indicates/rsi"
	"snake/internal/kline"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// RSIStrategy 基于RSI指标的交易策略
type RSIStrategy struct {
	*strategy.BaseStrategy
	// 保存历史K线用于计算RSI
	historicalKlines []*kline.Kline
	// RSI参数
	rsiPeriod       int
	oversoldLevel   decimal.Decimal // 超卖水平，默认30
	overboughtLevel decimal.Decimal // 超买水平，默认70
	// RSI指标
	rsiIndicator *rsi.RSI
}

// New 创建RSI策略
func New() *RSIStrategy {
	return &RSIStrategy{
		BaseStrategy:     strategy.NewBaseStrategy("RSI Strategy"),
		historicalKlines: make([]*kline.Kline, 0, 30), // 预分配足够容量
		rsiPeriod:        14,
		oversoldLevel:    decimal.NewFromInt(30),
		overboughtLevel:  decimal.NewFromInt(70),
	}
}

// Update 更新策略状态
func (s *RSIStrategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	// 添加新的K线到历史数据
	s.historicalKlines = append(s.historicalKlines, kline)

	// 如果历史数据不足以计算RSI，则持有
	if len(s.historicalKlines) < s.rsiPeriod+1 {
		return s.Hold(), nil
	}

	// 初始化或更新RSI指标
	if s.rsiIndicator == nil {
		s.rsiIndicator = rsi.NewWithPeriod(s.rsiPeriod, s.historicalKlines...)
	} else {
		newRSI := s.rsiIndicator.NextKline(kline)
		if newRSI != nil {
			s.rsiIndicator = newRSI
		}
	}

	// 如果无法计算RSI，则持有
	if s.rsiIndicator == nil {
		return s.Hold(), nil
	}

	// 计算交易数量（仓位大小）
	// 对于买入：使用可用余额的10%
	// 对于卖出：使用当前持仓的10%
	var buyAmount, sellAmount decimal.Decimal

	// 计算可用的买入数量（余额的10%除以当前价格）
	buyAmount = s.Balance().Amount.Mul(decimal.NewFromFloat(0.1)).Div(kline.C)
	// 计算可用的卖出数量（持仓的10%）
	sellAmount = s.Position().Amount.Mul(decimal.NewFromFloat(0.1))

	// 计算当前盈亏
	absolute, percentage := s.BaseStrategy.Profit(kline.C)

	// 如果当前有持仓，打印盈亏信息
	if !s.Position().Amount.IsZero() {
		println("当前持仓盈亏：", absolute.String(), "USDT (", percentage.String(), "%)")
	}

	// 打印当前RSI值
	println("当前RSI值：", s.rsiIndicator.Value.String())

	// RSI策略
	// 1. 如果RSI低于超卖水平(30)，买入信号
	if s.rsiIndicator.Value.LessThanOrEqual(s.oversoldLevel) {
		// 只有当有足够余额时才买入
		if !s.Balance().Amount.IsZero() && buyAmount.GreaterThan(decimal.Zero) {
			signal := s.Buy(buyAmount, kline.C)
			if signal != nil {
				return signal, nil
			}
		}
	}

	// 2. 如果RSI高于超买水平(70)，卖出信号
	if s.rsiIndicator.Value.GreaterThanOrEqual(s.overboughtLevel) {
		// 只有当有持仓时才卖出
		if !s.Position().Amount.IsZero() && sellAmount.GreaterThan(decimal.Zero) {
			signal := s.Sell(sellAmount, kline.C)
			if signal != nil {
				return signal, nil
			}
		}
	}

	// 其他情况，持有
	return s.Hold(), nil
}

// SetParams 设置RSI策略参数
func (s *RSIStrategy) SetParams(period int, oversold, overbought decimal.Decimal) {
	s.rsiPeriod = period
	s.oversoldLevel = oversold
	s.overboughtLevel = overbought
	// 重置指标
	s.rsiIndicator = nil
}

// Profit 计算盈亏
func (s *RSIStrategy) Profit(currentPrice decimal.Decimal) (absolute, percentage decimal.Decimal) {
	return s.BaseStrategy.Profit(currentPrice)
}

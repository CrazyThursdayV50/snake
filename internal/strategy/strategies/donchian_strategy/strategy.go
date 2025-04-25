package donchian_strategy

import (
	donchianchannel "snake/internal/indicates/donchian-channel"
	"snake/internal/kline"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// DonchianStrategy 基于唐奇安通道的交易策略
type DonchianStrategy struct {
	*strategy.BaseStrategy
	// 保存历史K线用于计算指标
	historicalKlines []*kline.Kline
	// 唐奇安通道参数
	breakoutPeriod int             // 突破周期（默认20）
	exitPeriod     int             // 退出周期（默认10）
	riskPercent    decimal.Decimal // 风险百分比（每笔交易的风险）
	// 唐奇安通道指标
	breakoutChannel *donchianchannel.DC // 用于入场信号的通道
	exitChannel     *donchianchannel.DC // 用于出场信号的通道
	// 策略状态
	position string // "long", "short", "none"
}

// New 创建唐奇安通道策略
func New() *DonchianStrategy {
	return &DonchianStrategy{
		BaseStrategy:     strategy.NewBaseStrategy("Donchian Channel Strategy"),
		historicalKlines: make([]*kline.Kline, 0, 50), // 预分配足够容量
		breakoutPeriod:   20,
		exitPeriod:       10,
		riskPercent:      decimal.NewFromFloat(1.0), // 1%风险
		position:         "none",
	}
}

// Update 更新策略状态
func (s *DonchianStrategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	// 添加新的K线到历史数据
	s.historicalKlines = append(s.historicalKlines, kline)

	// 如果历史数据不足以计算指标，则持有
	requiredBars := s.breakoutPeriod
	if s.exitPeriod > s.breakoutPeriod {
		requiredBars = s.exitPeriod
	}

	if len(s.historicalKlines) <= requiredBars {
		return s.Hold(), nil
	}

	// 计算唐奇安通道指标
	s.calculateIndicators()

	// 获取当前价格（收盘价）
	currentPrice := kline.C

	// 计算当前盈亏
	absolute, percentage := s.BaseStrategy.Profit(currentPrice)

	// 如果当前有持仓，打印盈亏信息
	if !s.Position().Amount.IsZero() {
		println("当前持仓盈亏：", absolute.String(), "USDT (", percentage.String(), "%)")
	}

	// 打印当前通道值
	println("突破通道上轨: ", s.breakoutChannel.Upper.String())
	println("突破通道下轨: ", s.breakoutChannel.Lower.String())
	println("退出通道上轨: ", s.exitChannel.Upper.String())
	println("退出通道下轨: ", s.exitChannel.Lower.String())

	// 计算交易数量（基于风险管理）
	tradeAmount := s.calculatePositionSize(currentPrice)

	// 执行交易策略
	switch s.position {
	case "none":
		// 无持仓状态，检查是否应该入场
		if s.breakoutChannel.IsBuySignal(currentPrice) {
			// 价格突破上轨，买入做多
			s.position = "long"
			signal := s.Buy(tradeAmount, currentPrice)
			if signal != nil {
				return signal, nil
			}
		} else if currentPrice.LessThanOrEqual(s.breakoutChannel.Lower) {
			// 价格突破下轨，卖出做空
			s.position = "short"
			signal := s.Sell(tradeAmount, currentPrice)
			if signal != nil {
				return signal, nil
			}
		}
	case "long":
		// 做多状态，检查是否应该退出
		if currentPrice.LessThanOrEqual(s.exitChannel.Lower) {
			// 价格跌破退出通道下轨，平多
			totalPosition := s.Position().Amount
			if !totalPosition.IsZero() {
				s.position = "none"
				signal := s.Sell(totalPosition, currentPrice)
				if signal != nil {
					return signal, nil
				}
			}
		}
	case "short":
		// 做空状态，检查是否应该退出
		if currentPrice.GreaterThanOrEqual(s.exitChannel.Upper) {
			// 价格突破退出通道上轨，平空
			totalPosition := s.Position().Amount
			if !totalPosition.IsZero() {
				s.position = "none"
				signal := s.Buy(totalPosition, currentPrice)
				if signal != nil {
					return signal, nil
				}
			}
		}
	}

	// 如果没有交易信号，则持有
	return s.Hold(), nil
}

// calculateIndicators 计算唐奇安通道指标
func (s *DonchianStrategy) calculateIndicators() {
	// 确保有足够的历史数据
	if len(s.historicalKlines) <= s.breakoutPeriod {
		return
	}

	// 计算突破通道
	s.breakoutChannel = donchianchannel.NewWithPeriod(s.breakoutPeriod, s.historicalKlines...)

	// 计算退出通道
	s.exitChannel = donchianchannel.NewWithPeriod(s.exitPeriod, s.historicalKlines...)
}

// calculatePositionSize 计算仓位大小
func (s *DonchianStrategy) calculatePositionSize(currentPrice decimal.Decimal) decimal.Decimal {
	// 获取账户总价值
	balance := s.Balance().Amount
	positionValue := s.Position().Amount.Mul(currentPrice)
	accountValue := balance.Add(positionValue)

	// 根据风险百分比计算能够承受的风险金额
	riskAmount := accountValue.Mul(s.riskPercent.Div(decimal.NewFromInt(100)))

	// 计算止损距离
	var stopDistance decimal.Decimal
	if s.exitChannel != nil {
		if s.breakoutChannel.Upper.LessThan(currentPrice) {
			// 做多，止损是退出通道的下轨
			stopDistance = currentPrice.Sub(s.exitChannel.Lower)
		} else if s.breakoutChannel.Lower.GreaterThan(currentPrice) {
			// 做空，止损是退出通道的上轨
			stopDistance = s.exitChannel.Upper.Sub(currentPrice)
		} else {
			// 默认使用1%的止损距离
			stopDistance = currentPrice.Mul(decimal.NewFromFloat(0.01))
		}
	} else {
		// 默认使用1%的止损距离
		stopDistance = currentPrice.Mul(decimal.NewFromFloat(0.01))
	}

	// 防止除以零
	if stopDistance.IsZero() {
		stopDistance = currentPrice.Mul(decimal.NewFromFloat(0.01))
	}

	// 计算仓位大小 = 风险金额 / 止损距离
	positionSize := riskAmount.Div(stopDistance)

	// 将资金金额转换为资产数量
	units := positionSize.Div(currentPrice)

	// 设置最小交易单位
	minUnit := decimal.NewFromFloat(0.01)
	if units.LessThan(minUnit) {
		return minUnit
	}

	return units
}

// SetParams 设置唐奇安通道策略参数
func (s *DonchianStrategy) SetParams(breakoutPeriod, exitPeriod int, riskPercent decimal.Decimal) {
	s.breakoutPeriod = breakoutPeriod
	s.exitPeriod = exitPeriod
	s.riskPercent = riskPercent
	// 重置指标
	s.breakoutChannel = nil
	s.exitChannel = nil
}

// Profit 计算盈亏
func (s *DonchianStrategy) Profit(currentPrice decimal.Decimal) (absolute, percentage decimal.Decimal) {
	return s.BaseStrategy.Profit(currentPrice)
}

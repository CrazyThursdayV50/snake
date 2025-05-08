package bollingmacd

import (
	"context"
	bollingband "snake/internal/indicates/bolling-band"
	"snake/internal/indicates/macd"
	"snake/internal/kline"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// BollingMACDStrategy 布林带-MACD联合策略
type BollingMACDStrategy struct {
	*strategy.BaseStrategy
	// 保存历史K线用于计算指标
	historicalKlines []*kline.Kline
	// 布林带参数
	bbPeriod int
	// MACD参数
	fastEMAPeriod int
	slowEMAPeriod int
	signalPeriod  int
	// 当前指标
	bb   *bollingband.BB
	macd *macd.MACD
}

// New 创建布林带-MACD联合策略
func New(ctx context.Context, cancel context.CancelFunc) *BollingMACDStrategy {
	return &BollingMACDStrategy{
		BaseStrategy:     strategy.NewBaseStrategy(ctx, cancel, "Bolling-MACD Strategy"),
		historicalKlines: make([]*kline.Kline, 0, 60), // 预分配足够容量
		bbPeriod:         20,
		fastEMAPeriod:    12,
		slowEMAPeriod:    26,
		signalPeriod:     9,
	}
}

// Update 更新策略状态
func (s *BollingMACDStrategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	// 添加新的K线到历史数据
	s.historicalKlines = append(s.historicalKlines, kline)

	// 如果历史数据不足以计算指标，则持有
	requiredPeriod := s.bbPeriod
	if s.slowEMAPeriod+s.signalPeriod > requiredPeriod {
		requiredPeriod = s.slowEMAPeriod + s.signalPeriod
	}

	if len(s.historicalKlines) < requiredPeriod {
		return s.Hold(), nil
	}

	// 保持历史数据长度不超过需要的最大长度
	if len(s.historicalKlines) > requiredPeriod {
		s.historicalKlines = s.historicalKlines[len(s.historicalKlines)-requiredPeriod:]
	}

	// 计算布林带
	bbKlines := s.historicalKlines[len(s.historicalKlines)-s.bbPeriod:]
	s.bb = bollingband.New(bbKlines...)

	// 计算MACD
	s.macd = macd.NewWithParams(s.fastEMAPeriod, s.slowEMAPeriod, s.signalPeriod, s.historicalKlines...)

	// 如果无法计算指标，则持有
	if s.bb == nil || s.macd == nil {
		return s.Hold(), nil
	}

	// 计算当前盈亏
	absolute, percentage := s.BaseStrategy.Profit()

	// 如果当前有持仓，打印盈亏信息
	if !s.Position().Amount.IsZero() {
		println("当前持仓盈亏：", absolute.String(), "USDT (", percentage.String(), "%)")
	}

	// 获取交易信号
	signal := s.getSignal(kline)
	return signal, nil
}

// getSignal 根据布林带和MACD指标生成交易信号
func (s *BollingMACDStrategy) getSignal(kline *kline.Kline) *strategy.Signal {
	// 默认交易量为当前仓位的5%
	tradeAmount := s.Position().Amount.Mul(decimal.NewFromFloat(0.05))
	// 如果仓位为0，则使用余额的5%
	if s.Position().Amount.IsZero() {
		tradeAmount = s.Balance().Amount.Mul(decimal.NewFromFloat(0.05))
	}

	// 买入条件：
	// 1. 价格低于布林带下轨
	// 2. MACD直方图由负变正（MACD金叉死叉）
	if kline.C.LessThan(s.bb.Lower) &&
		s.macd.Histogram.IsPositive() &&
		s.macd.MACD.GreaterThan(s.macd.Signal) {
		signal := s.Buy(tradeAmount, kline.C)
		if signal != nil {
			return signal
		}
	}

	// 卖出条件：
	// 1. 价格高于布林带上轨
	// 2. MACD直方图由正变负（MACD死叉金叉）
	if kline.C.GreaterThan(s.bb.Upper) &&
		s.macd.Histogram.IsNegative() &&
		s.macd.MACD.LessThan(s.macd.Signal) {
		signal := s.Sell(tradeAmount, kline.C)
		if signal != nil {
			return signal
		}
	}

	// 其他情况，持有
	return s.Hold()
}

// Profit 返回当前盈亏
func (s *BollingMACDStrategy) Profit() (absolute, percentage decimal.Decimal) {
	return s.BaseStrategy.Profit()
}

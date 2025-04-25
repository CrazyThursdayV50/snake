package turtle

import (
	donchianchannel "snake/internal/indicates/donchian-channel"
	"snake/internal/kline"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// TurtleStrategy 海龟交易法则策略
type TurtleStrategy struct {
	*strategy.BaseStrategy
	// 保存历史K线用于计算指标
	historicalKlines []*kline.Kline
	// 策略参数
	donchianPeriod int     // 唐奇安通道周期（一般为20）
	atrPeriod      int     // ATR计算周期（一般为14）
	riskPercent    float64 // 风险比例（每次交易风险占总资产的百分比，一般为1-2%）
	entryUnits     int     // 入场单元数（一般为1-4）
	currentUnits   int     // 当前持有单元数
	// 策略状态
	position       string          // 当前仓位：long, short, none
	lastEntryPrice decimal.Decimal // 上次入场价格
	lastExitPrice  decimal.Decimal // 上次出场价格
	// 唐奇安通道指标
	donchianChannel *donchianchannel.DC // 唐奇安通道指标
	atr             decimal.Decimal     // 当前ATR值
	stopLoss        decimal.Decimal     // 止损价
}

// New 创建海龟交易法则策略
func New() *TurtleStrategy {
	return &TurtleStrategy{
		BaseStrategy:     strategy.NewBaseStrategy("Turtle Trading Strategy"),
		historicalKlines: make([]*kline.Kline, 0, 50), // 预分配足够容量
		donchianPeriod:   20,
		atrPeriod:        14,
		riskPercent:      2.0, // 2%
		entryUnits:       4,
		currentUnits:     0,
		position:         "none",
		lastEntryPrice:   decimal.Zero,
		lastExitPrice:    decimal.Zero,
		donchianChannel:  nil,
		atr:              decimal.Zero,
		stopLoss:         decimal.Zero,
	}
}

// Update 更新策略状态
func (s *TurtleStrategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	println("Update - 添加K线:", kline.S, "->", kline.E, "收盘价:", kline.C.String())

	// 添加新的K线到历史数据
	s.historicalKlines = append(s.historicalKlines, kline)

	// 需要足够的历史数据来生成信号
	if len(s.historicalKlines) < s.donchianPeriod {
		println("历史数据不足:", len(s.historicalKlines), "<", s.donchianPeriod)
		return s.Hold(), nil
	}

	// 限制历史数据长度
	if len(s.historicalKlines) > s.donchianPeriod*2 {
		s.historicalKlines = s.historicalKlines[len(s.historicalKlines)-s.donchianPeriod*2:]
	}

	// 计算指标
	s.calculateIndicators()
	println("计算指标 - 最高高点:", s.donchianChannel.Upper.String(), "最低低点:", s.donchianChannel.Lower.String(), "ATR:", s.atr.String())

	// 计算交易数量
	tradeAmount := s.calculatePositionSize(kline.C)
	println("计算交易数量:", tradeAmount.String(), "基于价格:", kline.C.String())

	// 当前持仓盈亏
	if !s.Position().Amount.IsZero() {
		absolute, percentage := s.BaseStrategy.Profit(kline.C)
		println("当前持仓盈亏：", absolute.String(), "USDT (", percentage.String(), "%)")
	}

	// 根据当前仓位执行不同的交易逻辑
	var signal *strategy.Signal
	var err error

	switch s.position {
	case "none":
		// 根据突破判断入场（系统1：20日突破）
		signal, err = s.evaluateEntry(kline, tradeAmount)
	case "long":
		// 检查是否应该追加仓位或离场
		signal, err = s.evaluateLongPosition(kline, tradeAmount)
	case "short":
		// 检查是否应该追加仓位或离场
		signal, err = s.evaluateShortPosition(kline, tradeAmount)
	}

	if err != nil {
		return nil, err
	}

	if signal != nil {
		return signal, nil
	}

	return s.Hold(), nil
}

// calculateIndicators 计算策略所需的技术指标
func (s *TurtleStrategy) calculateIndicators() {
	// 计算唐奇安通道
	s.calculateDonchianChannel()

	// 计算ATR
	s.calculateATR()
}

// calculateDonchianChannel 计算唐奇安通道
func (s *TurtleStrategy) calculateDonchianChannel() {
	// 确保有足够的历史数据
	if len(s.historicalKlines) < s.donchianPeriod {
		return
	}

	// 使用donchian-channel包计算唐奇安通道
	s.donchianChannel = donchianchannel.NewWithPeriod(s.donchianPeriod, s.historicalKlines...)
}

// calculateATR 计算平均真实范围（Average True Range）
func (s *TurtleStrategy) calculateATR() {
	// 确保有足够的历史数据
	if len(s.historicalKlines) < s.atrPeriod+1 {
		return
	}

	// 计算真实范围（True Range）序列
	trValues := make([]decimal.Decimal, s.atrPeriod)

	// 获取计算范围内的K线
	dataRange := s.historicalKlines[len(s.historicalKlines)-s.atrPeriod-1:]

	// 计算每个周期的TR值
	for i := 1; i <= s.atrPeriod; i++ {
		current := dataRange[i]
		previous := dataRange[i-1]

		// TR = max(high-low, |high-prev_close|, |low-prev_close|)
		hl := current.H.Sub(current.L)
		hpc := current.H.Sub(previous.C).Abs()
		lpc := current.L.Sub(previous.C).Abs()

		// 找出最大值
		tr := hl
		if hpc.GreaterThan(tr) {
			tr = hpc
		}
		if lpc.GreaterThan(tr) {
			tr = lpc
		}

		trValues[i-1] = tr
	}

	// 计算ATR（简单平均）
	sum := decimal.Zero
	for _, tr := range trValues {
		sum = sum.Add(tr)
	}

	s.atr = sum.Div(decimal.NewFromInt(int64(s.atrPeriod)))
}

// calculatePositionSize 计算交易头寸大小
func (s *TurtleStrategy) calculatePositionSize(price decimal.Decimal) decimal.Decimal {
	// 获取账户总值
	accountValue := s.Balance().Amount
	positionValue := decimal.Zero
	if !s.Position().Amount.IsZero() {
		positionValue = s.Position().Amount.Mul(price)
	}
	totalValue := accountValue.Add(positionValue)

	// 计算单位资金
	riskAmount := totalValue.Mul(decimal.NewFromFloat(s.riskPercent / 100.0))

	// 计算可以承担的单位
	if s.atr.IsZero() {
		return decimal.NewFromFloat(0.01) // 默认最小交易量
	}

	// 每个单位的美元价值为：N * 美元乘数
	// 这里假设美元乘数为1（即1个单位代表1美元）
	dollarPerUnit := s.atr.Mul(decimal.NewFromInt(1))

	// 计算单位数量
	units := riskAmount.Div(dollarPerUnit)

	// 将单位转换为实际的交易量
	tradeAmount := units.Div(price)

	// 设置最小交易量（为了测试方便）
	minTradeAmount := decimal.NewFromFloat(0.01)
	if tradeAmount.LessThan(minTradeAmount) {
		return minTradeAmount
	}

	return tradeAmount
}

// evaluateEntry 评估是否入场
func (s *TurtleStrategy) evaluateEntry(kline *kline.Kline, tradeAmount decimal.Decimal) (*strategy.Signal, error) {
	// 调试信息
	println("evaluateEntry - 收盘价:", kline.C.String(), "最高高点:", s.donchianChannel.Upper.String(), "最低低点:", s.donchianChannel.Lower.String())
	println("突破条件:", kline.C.GreaterThanOrEqual(s.donchianChannel.Upper), s.position != "long")

	// 系统1：价格突破20日高点，做多入场
	if kline.C.GreaterThanOrEqual(s.donchianChannel.Upper) && s.position != "long" {
		println("满足多头入场条件")
		// 设置止损价（通常为入场价减去2个ATR）
		s.stopLoss = kline.C.Sub(s.atr.Mul(decimal.NewFromInt(2)))

		// 更新状态
		s.position = "long"
		s.currentUnits = 1
		s.lastEntryPrice = kline.C

		// 生成买入信号
		return s.Buy(tradeAmount, kline.C), nil
	}

	println("突破条件:", kline.C.LessThanOrEqual(s.donchianChannel.Lower), s.position != "short")

	// 系统1：价格突破20日低点，做空入场
	if kline.C.LessThanOrEqual(s.donchianChannel.Lower) && s.position != "short" {
		println("满足空头入场条件")
		// 设置止损价（通常为入场价加上2个ATR）
		s.stopLoss = kline.C.Add(s.atr.Mul(decimal.NewFromInt(2)))

		// 更新状态
		s.position = "short"
		s.currentUnits = 1
		s.lastEntryPrice = kline.C

		// 生成卖出信号
		return s.Sell(tradeAmount, kline.C), nil
	}

	return nil, nil
}

// evaluateLongPosition 评估多头持仓
func (s *TurtleStrategy) evaluateLongPosition(kline *kline.Kline, tradeAmount decimal.Decimal) (*strategy.Signal, error) {
	// 调试信息
	println("evaluateLongPosition - 收盘价:", kline.C.String(), "止损价:", s.stopLoss.String())

	// 检查是否触发止损
	if kline.C.LessThanOrEqual(s.stopLoss) {
		println("触发止损出场")
		// 平掉所有仓位
		totalPosition := s.Position().Amount
		if !totalPosition.IsZero() {
			s.position = "none"
			s.currentUnits = 0
			s.lastExitPrice = kline.C
			return s.Sell(totalPosition, kline.C), nil
		}
	}

	// 检查是否触发利润保护（如价格跌破10日低点）
	exitPeriod := 10 // 退出使用较短周期
	if len(s.historicalKlines) > exitPeriod {
		// 创建退出用的短周期唐奇安通道
		exitDC := donchianchannel.NewWithPeriod(exitPeriod, s.historicalKlines...)

		if exitDC != nil && kline.C.LessThanOrEqual(exitDC.Lower) {
			println("触发利润保护出场")
			// 平掉所有仓位
			totalPosition := s.Position().Amount
			if !totalPosition.IsZero() {
				s.position = "none"
				s.currentUnits = 0
				s.lastExitPrice = kline.C
				return s.Sell(totalPosition, kline.C), nil
			}
		}
	}

	// 检查是否可以加仓（海龟法则的分批入场机制）
	if s.currentUnits < s.entryUnits {
		// 只有价格比上次入场价高0.5个ATR，才允许加仓
		atrJump := s.atr.Mul(decimal.NewFromFloat(0.5))
		nextEntryPrice := s.lastEntryPrice.Add(atrJump)

		if kline.C.GreaterThanOrEqual(nextEntryPrice) {
			println("触发加仓条件")
			// 执行加仓
			s.currentUnits++
			s.lastEntryPrice = kline.C
			// 更新止损
			s.stopLoss = kline.C.Sub(s.atr.Mul(decimal.NewFromFloat(2)))
			return s.Buy(tradeAmount, kline.C), nil
		}
	}

	return nil, nil
}

// evaluateShortPosition 评估空头持仓
func (s *TurtleStrategy) evaluateShortPosition(kline *kline.Kline, tradeAmount decimal.Decimal) (*strategy.Signal, error) {
	// 调试信息
	println("evaluateShortPosition - 收盘价:", kline.C.String(), "止损价:", s.stopLoss.String())

	// 检查是否触发止损
	if kline.C.GreaterThanOrEqual(s.stopLoss) {
		println("触发止损出场")
		// 平掉所有仓位
		totalPosition := s.Position().Amount
		if !totalPosition.IsZero() {
			s.position = "none"
			s.currentUnits = 0
			s.lastExitPrice = kline.C
			return s.Buy(totalPosition, kline.C), nil
		}
	}

	// 检查是否触发利润保护（如价格突破10日高点）
	exitPeriod := 10 // 退出使用较短周期
	if len(s.historicalKlines) > exitPeriod {
		// 创建退出用的短周期唐奇安通道
		exitDC := donchianchannel.NewWithPeriod(exitPeriod, s.historicalKlines...)

		if exitDC != nil && kline.C.GreaterThanOrEqual(exitDC.Upper) {
			println("触发利润保护出场")
			// 平掉所有仓位
			totalPosition := s.Position().Amount
			if !totalPosition.IsZero() {
				s.position = "none"
				s.currentUnits = 0
				s.lastExitPrice = kline.C
				return s.Buy(totalPosition, kline.C), nil
			}
		}
	}

	// 检查是否可以加仓（海龟法则的分批入场机制）
	if s.currentUnits < s.entryUnits {
		// 只有价格比上次入场价低0.5个ATR，才允许加仓
		atrJump := s.atr.Mul(decimal.NewFromFloat(0.5))
		nextEntryPrice := s.lastEntryPrice.Sub(atrJump)

		if kline.C.LessThanOrEqual(nextEntryPrice) {
			println("触发加仓条件")
			// 执行加仓
			s.currentUnits++
			s.lastEntryPrice = kline.C
			// 更新止损
			s.stopLoss = kline.C.Add(s.atr.Mul(decimal.NewFromFloat(2)))
			return s.Sell(tradeAmount, kline.C), nil
		}
	}

	return nil, nil
}

// Profit 计算盈亏
func (s *TurtleStrategy) Profit(currentPrice decimal.Decimal) (absolute, percentage decimal.Decimal) {
	return s.BaseStrategy.Profit(currentPrice)
}

// 辅助方法，转发到BaseStrategy
func (s *TurtleStrategy) Buy(amount, price decimal.Decimal) *strategy.Signal {
	println("发出买入信号, 数量:", amount.String(), "价格:", price.String())
	return s.BaseStrategy.Buy(amount, price)
}

func (s *TurtleStrategy) Sell(amount, price decimal.Decimal) *strategy.Signal {
	println("发出卖出信号, 数量:", amount.String(), "价格:", price.String())
	return s.BaseStrategy.Sell(amount, price)
}

func (s *TurtleStrategy) Hold() *strategy.Signal {
	return s.BaseStrategy.Hold()
}

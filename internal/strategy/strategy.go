package strategy

import (
	"snake/internal/kline"
	"snake/internal/types"
	"time"

	"github.com/shopspring/decimal"
)

// Position 表示当前持仓
type Position struct {
	// 当前持仓数量（例如 BTC 数量）
	Amount decimal.Decimal
	// 持仓成本（USDT）
	Cost decimal.Decimal
	// 持仓时间
	Time time.Time
}

// Balance 表示当前余额
type Balance struct {
	// 余额数量（例如 USDT 数量）
	Amount decimal.Decimal
	// 余额时间
	Time time.Time
}

// Signal 表示交易信号
type Signal struct {
	// 信号类型：买入、卖出或持有
	Type types.SignalType
	// 交易数量
	Amount decimal.Decimal
	// 交易价格
	Price decimal.Decimal
	// 信号时间
	Time time.Time
}

// Strategy 策略接口
type Strategy interface {
	// Name 返回策略名称
	Name() string
	// Init 初始化策略
	Init(positionAmount, balanceAmount decimal.Decimal) error
	// Update 更新策略状态
	Update(kline *kline.Kline) (*Signal, error)
	// Position 返回当前持仓
	Position() *Position
	// Balance 返回当前余额
	Balance() *Balance
	// Profit 计算盈亏
	Profit(currentPrice decimal.Decimal) (absolute, percentage decimal.Decimal)
}

// BaseStrategy 基础策略结构体
type BaseStrategy struct {
	// 策略名称
	name string
	// 当前持仓
	position *Position
	// 当前余额
	balance *Balance
}

// NewBaseStrategy 创建基础策略
func NewBaseStrategy(name string) *BaseStrategy {
	return &BaseStrategy{
		name:     name,
		position: &Position{Amount: decimal.Zero, Cost: decimal.Zero},
		balance:  &Balance{Amount: decimal.Zero},
	}
}

// Name 返回策略名称
func (s *BaseStrategy) Name() string {
	return s.name
}

// Init 初始化策略
func (s *BaseStrategy) Init(positionAmount, balanceAmount decimal.Decimal) error {
	s.position.Amount = positionAmount
	s.position.Time = time.Now()
	s.balance.Amount = balanceAmount
	s.balance.Time = time.Now()

	// 设置初始持仓成本（使用初始价格 100 USDT）
	if !positionAmount.IsZero() {
		initialPrice := decimal.NewFromFloat(100.0)
		s.position.Cost = positionAmount.Mul(initialPrice)
	}

	return nil
}

// Position 返回当前持仓
func (s *BaseStrategy) Position() *Position {
	return s.position
}

// Balance 返回当前余额
func (s *BaseStrategy) Balance() *Balance {
	return s.balance
}

// Buy 执行买入操作
func (s *BaseStrategy) Buy(amount, price decimal.Decimal) *Signal {
	// 计算需要的 USDT 数量
	usdtAmount := amount.Mul(price)

	// 检查余额是否足够
	if s.balance.Amount.LessThan(usdtAmount) {
		return nil
	}

	// 更新余额和仓位
	s.balance.Amount = s.balance.Amount.Sub(usdtAmount)
	s.position.Amount = s.position.Amount.Add(amount)
	// 更新持仓成本：新成本 = 旧成本 + 新买入成本
	s.position.Cost = s.position.Cost.Add(usdtAmount)
	s.position.Time = time.Now()
	s.balance.Time = time.Now()

	return &Signal{
		Type:   types.SignalTypeBuy,
		Amount: amount,
		Price:  price,
		Time:   time.Now(),
	}
}

// Sell 执行卖出操作
func (s *BaseStrategy) Sell(amount, price decimal.Decimal) *Signal {
	// 检查仓位是否足够
	if s.position.Amount.LessThan(amount) {
		return nil
	}

	// 计算获得的 USDT 数量
	usdtAmount := amount.Mul(price)

	// 更新余额和仓位
	s.balance.Amount = s.balance.Amount.Add(usdtAmount)
	s.position.Amount = s.position.Amount.Sub(amount)
	// 更新持仓成本：新成本 = 旧成本 * (1 - 卖出比例)
	sellRatio := amount.Div(s.position.Amount.Add(amount))
	s.position.Cost = s.position.Cost.Mul(decimal.NewFromInt(1).Sub(sellRatio))
	s.position.Time = time.Now()
	s.balance.Time = time.Now()

	return &Signal{
		Type:   types.SignalTypeSell,
		Amount: amount,
		Price:  price,
		Time:   time.Now(),
	}
}

// Hold 返回持有信号
func (s *BaseStrategy) Hold() *Signal {
	return &Signal{
		Type: types.SignalTypeHold,
		Time: time.Now(),
	}
}

// Profit 计算盈亏
func (s *BaseStrategy) Profit(currentPrice decimal.Decimal) (absolute, percentage decimal.Decimal) {
	if s.position.Amount.IsZero() {
		return decimal.Zero, decimal.Zero
	}

	// 计算当前持仓市值
	currentValue := s.position.Amount.Mul(currentPrice)

	// 计算盈亏绝对数量
	absolute = currentValue.Sub(s.position.Cost)

	// 计算盈亏百分比
	if s.position.Cost.IsZero() {
		percentage = decimal.Zero
	} else {
		// 计算盈亏百分比：(当前市值 - 持仓成本) / 持仓成本 * 100
		percentage = absolute.Div(s.position.Cost).Mul(decimal.NewFromInt(100))
	}

	return absolute, percentage
}

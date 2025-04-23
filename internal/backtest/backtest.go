package backtest

import (
	"context"
	"fmt"
	"snake/internal/indicates/ma"
	"snake/internal/kline/interval"
	"snake/internal/models"
	"snake/internal/repository"
	"snake/internal/strategy"

	"github.com/shopspring/decimal"
)

// Config 回测配置
type Config struct {
	// 初始资金（USDT）
	InitialBalance decimal.Decimal
	// 初始持仓（BTC）
	InitialPosition decimal.Decimal
	// K线时间间隔
	Interval interval.Interval
	// MA20 周期
	MA20Period int
	// MA60 周期
	MA60Period int
}

// Result 回测结果
type Result struct {
	// 最终资金（USDT）
	FinalBalance decimal.Decimal
	// 最终持仓（BTC）
	FinalPosition decimal.Decimal
	// 总交易次数
	TotalTrades int
	// 盈利次数
	WinningTrades int
	// 亏损次数
	LosingTrades int
	// 最大回撤
	MaxDrawdown decimal.Decimal
	// 收益率
	ROI decimal.Decimal
	// 交易记录
	Trades []*Trade
}

// Trade 交易记录
type Trade struct {
	// 交易时间
	Time int64
	// 交易类型
	Type strategy.SignalType
	// 交易数量
	Amount decimal.Decimal
	// 交易价格
	Price decimal.Decimal
	// 交易后的余额
	Balance decimal.Decimal
	// 交易后的持仓
	Position decimal.Decimal
	// 当前盈亏
	ProfitLoss decimal.Decimal
	// 盈亏百分比
	ProfitLossPercentage decimal.Decimal
}

// Backtest 回测器
type Backtest struct {
	config     *Config
	repository repository.KlineRepository
	strategy   strategy.Strategy
}

// New 创建回测器
func New(config *Config, repository repository.KlineRepository, strategy strategy.Strategy) *Backtest {
	return &Backtest{
		config:     config,
		repository: repository,
		strategy:   strategy,
	}
}

// Run 运行回测
func (b *Backtest) Run(ctx context.Context) (*Result, error) {
	// 初始化策略
	if err := b.strategy.Init(b.config.InitialPosition, b.config.InitialBalance); err != nil {
		return nil, fmt.Errorf("初始化策略失败: %v", err)
	}

	// 获取所有 K 线数据
	mysqlKlines, err := b.repository.ListAll(ctx, b.config.Interval)
	if err != nil {
		return nil, fmt.Errorf("获取 K 线数据失败: %v", err)
	}

	// 转换 K 线数据
	klines := convertKlines(mysqlKlines)

	// 初始化结果
	result := &Result{
		Trades: make([]*Trade, 0),
	}

	// 初始化 MA 指标
	var ma20Klines, ma60Klines []*models.Kline
	ma20 := (*ma.MA)(nil)
	ma60 := (*ma.MA)(nil)

	// 记录每个时间点的资产总值
	peakValue := b.config.InitialBalance.Add(b.config.InitialPosition.Mul(klines[0].C))
	maxDrawdown := decimal.Zero

	// 遍历每个 K 线
	for i, kline := range klines {
		// 更新 MA20
		if i >= b.config.MA20Period-1 {
			if ma20 == nil {
				ma20Klines = klines[i-b.config.MA20Period+1 : i+1]
				ma20 = ma.New(ma20Klines...)
			} else {
				ma20 = ma20.NextKline(kline)
			}
		}

		// 更新 MA60
		if i >= b.config.MA60Period-1 {
			if ma60 == nil {
				ma60Klines = klines[i-b.config.MA60Period+1 : i+1]
				ma60 = ma.New(ma60Klines...)
			} else {
				ma60 = ma60.NextKline(kline)
			}
		}

		// 如果两个 MA 都准备好了，执行策略
		if ma20 != nil && ma60 != nil {
			signal, err := b.strategy.Update(kline, ma20, ma60)
			if err != nil {
				return nil, fmt.Errorf("更新策略失败: %v", err)
			}

			// 计算当前资产总值（使用最高价）
			currentValue := b.strategy.Balance().Amount.Add(b.strategy.Position().Amount.Mul(kline.H))
			
			// 更新峰值
			if currentValue.GreaterThan(peakValue) {
				peakValue = currentValue
			}

			// 计算回撤（使用最低价）
			if !peakValue.IsZero() {
				minValue := b.strategy.Balance().Amount.Add(b.strategy.Position().Amount.Mul(kline.L))
				drawdown := peakValue.Sub(minValue).Div(peakValue).Mul(decimal.NewFromInt(100))
				if drawdown.GreaterThan(maxDrawdown) {
					maxDrawdown = drawdown
				}
			}

			// 记录交易
			if !signal.Type.IsHold() {
				// 计算当前盈亏
				profitLoss, profitLossPercentage := b.strategy.Profit(kline.C)

				trade := &Trade{
					Time:                signal.Time.Unix(),
					Type:                signal.Type,
					Amount:              signal.Amount,
					Price:               signal.Price,
					Balance:             b.strategy.Balance().Amount,
					Position:            b.strategy.Position().Amount,
					ProfitLoss:         profitLoss,
					ProfitLossPercentage: profitLossPercentage,
				}
				result.Trades = append(result.Trades, trade)

				// 更新交易统计
				result.TotalTrades++
				if profitLoss.IsPositive() {
					result.WinningTrades++
				} else if profitLoss.IsNegative() {
					result.LosingTrades++
				}
			}
		}
	}

	// 计算最终结果
	result.FinalBalance = b.strategy.Balance().Amount
	result.FinalPosition = b.strategy.Position().Amount

	// 计算收益率
	initialValue := b.config.InitialBalance.Add(b.config.InitialPosition.Mul(klines[0].C))
	finalValue := result.FinalBalance.Add(result.FinalPosition.Mul(klines[len(klines)-1].C))
	result.ROI = finalValue.Sub(initialValue).Div(initialValue).Mul(decimal.NewFromInt(100))

	// 设置最大回撤
	result.MaxDrawdown = maxDrawdown

	return result, nil
}

// calculateMaxDrawdown 计算最大回撤
func (b *Backtest) calculateMaxDrawdown(trades []*Trade) decimal.Decimal {
	if len(trades) == 0 {
		return decimal.Zero
	}

	maxDrawdown := decimal.Zero
	peak := decimal.Zero

	for _, trade := range trades {
		// 计算当前资产总值
		currentValue := trade.Balance.Add(trade.Position.Mul(trade.Price))
		
		// 更新峰值
		if currentValue.GreaterThan(peak) {
			peak = currentValue
		}

		// 计算回撤
		if !peak.IsZero() {
			drawdown := peak.Sub(currentValue).Div(peak).Mul(decimal.NewFromInt(100))
			if drawdown.GreaterThan(maxDrawdown) {
				maxDrawdown = drawdown
			}
		}
	}

	return maxDrawdown
} 
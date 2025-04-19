package backtest

import (
	"context"
	"fmt"
	"snake/internal/kline"
	"snake/internal/kline/interval"
	"snake/internal/strategy"
	"snake/internal/types"
	"time"

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

// Result 回测结果（每个K线的回测结果）
type Result []*kline.PositionKline

// Trade 交易记录
type Trade struct {
	// 交易时间
	Time time.Time
	// 交易类型
	Type types.SignalType
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

// TradeType 交易类型
type TradeType int

const (
	// TradeTypeBuy 买入
	TradeTypeBuy TradeType = iota
	// TradeTypeSell 卖出
	TradeTypeSell
)

// Backtest 回测结构体
type Backtest struct {
	config     *Config
	repository kline.Repository
	strategy   strategy.Strategy
	// 当前最高资产值
	peakValue decimal.Decimal
	// 交易记录
	trades []*Trade
}

// New 创建回测实例
func New(config *Config, repository kline.Repository, strategy strategy.Strategy) *Backtest {
	return &Backtest{
		config:     config,
		repository: repository,
		strategy:   strategy,
		peakValue:  decimal.Zero,
		trades:     make([]*Trade, 0),
	}
}

// Run 执行回测
func (b *Backtest) Run(ctx context.Context) (Result, error) {
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

	// 记录交易
	b.trades = make([]*Trade, 0)

	// 初始化最高资产值
	b.peakValue = b.config.InitialBalance.Add(b.config.InitialPosition.Mul(decimal.NewFromFloat(100.0)))

	// 初始化回测结果
	result := make(Result, 0)

	// 遍历 K 线
	for _, k := range klines {
		// 更新策略
		signal, err := b.strategy.Update(k)
		if err != nil {
			return nil, fmt.Errorf("更新策略失败: %v", err)
		}

		// 记录当前资产情况
		positionKline := &kline.PositionKline{
			Time:           time.Unix(k.S/1000, 0),
			Kline:          k,
			PositionAmount: b.strategy.Position().Amount,
			PositionCost:   b.strategy.Position().Cost,
			Balance:        b.strategy.Balance().Amount,
		}

		// 计算当前持仓市值
		positionValue := positionKline.PositionAmount.Mul(k.C)
		// 计算总资产
		positionKline.TotalValue = positionValue.Add(positionKline.Balance)

		// 更新最高资产值
		if positionKline.TotalValue.GreaterThan(b.peakValue) {
			b.peakValue = positionKline.TotalValue
		}
		positionKline.PeakValue = b.peakValue

		// 计算回撤
		if !b.peakValue.IsZero() {
			positionKline.Drawdown = b.peakValue.Sub(positionKline.TotalValue).Div(b.peakValue).Mul(decimal.NewFromInt(100))
		}

		// 计算盈亏
		profitAbsolute, profitPercentage := b.strategy.Profit()
		positionKline.ProfitAbsolute = profitAbsolute
		positionKline.ProfitPercentage = profitPercentage

		// 记录当前 K 线的资产情况
		result = append(result, positionKline)

		// 处理交易信号
		if signal != nil {
			if signal.Type.IsBuy() {
				b.trades = append(b.trades, &Trade{
					Time:   signal.Time,
					Type:   signal.Type,
					Amount: signal.Amount,
					Price:  signal.Price,
				})
			} else if signal.Type.IsSell() {
				b.trades = append(b.trades, &Trade{
					Time:   signal.Time,
					Type:   signal.Type,
					Amount: signal.Amount,
					Price:  signal.Price,
				})
			}
		}
	}

	return result, nil
}

// DisplaySummary 显示回测结果摘要
func (b *Backtest) DisplaySummary(result Result) {
	if len(result) == 0 {
		fmt.Println("回测结果为空")
		return
	}

	// 获取初始和最终状态
	initialKline := result[0]
	finalKline := result[len(result)-1]

	// 计算交易次数
	buyCount := 0
	sellCount := 0
	for _, trade := range b.trades {
		if trade.Type.IsBuy() {
			buyCount++
		} else if trade.Type.IsSell() {
			sellCount++
		}
	}
	totalTrades := buyCount + sellCount

	// 计算收益率
	initialValue := b.config.InitialBalance.Add(b.config.InitialPosition.Mul(initialKline.Kline.C))
	finalValue := finalKline.Balance.Add(finalKline.PositionAmount.Mul(finalKline.Kline.C))
	roi := decimal.Zero
	if !initialValue.IsZero() {
		roi = finalValue.Sub(initialValue).Div(initialValue).Mul(decimal.NewFromInt(100))
	}

	// 计算最大回撤
	maxDrawdown := decimal.Zero
	for _, pk := range result {
		if pk.Drawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = pk.Drawdown
		}
	}

	// 输出摘要
	fmt.Println("\n======================== 回测结果摘要 ========================")
	fmt.Printf("策略名称: %s\n", b.strategy.Name())
	fmt.Printf("回测周期: %s\n", b.config.Interval.String())
	fmt.Printf("回测K线数量: %d\n", len(result))
	fmt.Printf("开始日期: %s\n", initialKline.Time.Format("2006-01-02 15:04:05"))
	fmt.Printf("结束日期: %s\n", finalKline.Time.Format("2006-01-02 15:04:05"))
	fmt.Println("\n---------------------- 资金情况 ------------------------")
	fmt.Printf("初始资金: %.4f USDT\n", b.config.InitialBalance.InexactFloat64())
	fmt.Printf("初始持仓: %.8f BTC\n", b.config.InitialPosition.InexactFloat64())
	fmt.Printf("初始总资产: %.4f USDT\n", initialValue.InexactFloat64())
	fmt.Printf("最终资金: %.4f USDT\n", finalKline.Balance.InexactFloat64())
	fmt.Printf("最终持仓: %.8f BTC\n", finalKline.PositionAmount.InexactFloat64())
	fmt.Printf("最终总资产: %.4f USDT\n", finalValue.InexactFloat64())
	fmt.Printf("收益率: %.2f%%\n", roi.InexactFloat64())
	fmt.Println("\n---------------------- 交易统计 ------------------------")
	fmt.Printf("总交易次数: %d\n", totalTrades)
	fmt.Printf("买入次数: %d\n", buyCount)
	fmt.Printf("卖出次数: %d\n", sellCount)
	fmt.Printf("最大回撤: %.2f%%\n", maxDrawdown.InexactFloat64())

	// 显示部分交易记录
	if len(b.trades) > 0 {
		fmt.Println("\n---------------------- 交易记录示例 ------------------------")
		fmt.Printf("%-20s %-6s %-12s %-12s\n", "时间", "类型", "数量", "价格")

		// 显示最多10条交易记录
		displayCount := 10
		if len(b.trades) < displayCount {
			displayCount = len(b.trades)
		}

		for i := 0; i < displayCount; i++ {
			trade := b.trades[i]
			tradeType := "买入"
			if trade.Type.IsSell() {
				tradeType = "卖出"
			}
			fmt.Printf("%-20s %-6s %-12.8f %-12.2f\n",
				trade.Time.Format("2006-01-02 15:04:05"),
				tradeType,
				trade.Amount.InexactFloat64(),
				trade.Price.InexactFloat64())
		}

		if len(b.trades) > displayCount {
			fmt.Printf("... 还有 %d 条交易记录未显示 ...\n", len(b.trades)-displayCount)
		}
	}

	fmt.Println("\n=============================================================")
}

// DisplayKlineResults 显示每个K线的回测结果
func (b *Backtest) DisplayKlineResults(result Result, limit int) {
	if len(result) == 0 {
		fmt.Println("回测结果为空")
		return
	}

	// 确定显示的K线数量
	displayCount := limit
	if displayCount <= 0 || displayCount > len(result) {
		displayCount = len(result)
	}

	// 打印表头
	fmt.Println("\n========================= K线回测结果 =========================")
	fmt.Printf("%-20s %-8s %-8s %-12s %-12s %-12s %-12s %-8s %-8s\n",
		"时间", "开盘价", "收盘价", "持仓数量", "持仓成本", "余额", "总资产", "回撤(%)", "盈亏(%)")

	// 根据显示数量决定显示方式
	if displayCount == len(result) {
		// 显示全部结果
		for _, pk := range result {
			fmt.Printf("%-20s %-8.2f %-8.2f %-12.8f %-12.4f %-12.4f %-12.4f %-8.2f %-8.2f\n",
				pk.Time.Format("2006-01-02 15:04:05"),
				pk.Kline.O.InexactFloat64(),
				pk.Kline.C.InexactFloat64(),
				pk.PositionAmount.InexactFloat64(),
				pk.PositionCost.InexactFloat64(),
				pk.Balance.InexactFloat64(),
				pk.TotalValue.InexactFloat64(),
				pk.Drawdown.InexactFloat64(),
				pk.ProfitPercentage.InexactFloat64())
		}
	} else {
		// 只显示部分结果：前3个、最后3个、以及其他特别的结果点（比如最大回撤点）
		// 找到最大回撤的K线
		maxDrawdownIndex := 0
		maxDrawdown := decimal.Zero
		for i, pk := range result {
			if pk.Drawdown.GreaterThan(maxDrawdown) {
				maxDrawdown = pk.Drawdown
				maxDrawdownIndex = i
			}
		}

		// 显示前3个K线
		showCount := 3
		if len(result) < showCount {
			showCount = len(result)
		}
		for i := 0; i < showCount; i++ {
			pk := result[i]
			fmt.Printf("%-20s %-8.2f %-8.2f %-12.8f %-12.4f %-12.4f %-12.4f %-8.2f %-8.2f\n",
				pk.Time.Format("2006-01-02 15:04:05"),
				pk.Kline.O.InexactFloat64(),
				pk.Kline.C.InexactFloat64(),
				pk.PositionAmount.InexactFloat64(),
				pk.PositionCost.InexactFloat64(),
				pk.Balance.InexactFloat64(),
				pk.TotalValue.InexactFloat64(),
				pk.Drawdown.InexactFloat64(),
				pk.ProfitPercentage.InexactFloat64())
		}

		// 如果结果数量大于6，显示省略号
		if len(result) > 6 {
			fmt.Println("                    ............")
		}

		// 显示最大回撤点（如果不在前3个或后3个K线中）
		if maxDrawdownIndex >= showCount && maxDrawdownIndex < len(result)-showCount {
			pk := result[maxDrawdownIndex]
			fmt.Printf("%-20s %-8.2f %-8.2f %-12.8f %-12.4f %-12.4f %-12.4f %-8.2f %-8.2f (最大回撤点)\n",
				pk.Time.Format("2006-01-02 15:04:05"),
				pk.Kline.O.InexactFloat64(),
				pk.Kline.C.InexactFloat64(),
				pk.PositionAmount.InexactFloat64(),
				pk.PositionCost.InexactFloat64(),
				pk.Balance.InexactFloat64(),
				pk.TotalValue.InexactFloat64(),
				pk.Drawdown.InexactFloat64(),
				pk.ProfitPercentage.InexactFloat64())
		}

		// 如果结果数量大于6，显示省略号
		if len(result) > 6 {
			fmt.Println("                    ............")
		}

		// 显示最后3个K线
		for i := len(result) - showCount; i < len(result); i++ {
			if i < 0 {
				continue
			}
			pk := result[i]
			fmt.Printf("%-20s %-8.2f %-8.2f %-12.8f %-12.4f %-12.4f %-12.4f %-8.2f %-8.2f\n",
				pk.Time.Format("2006-01-02 15:04:05"),
				pk.Kline.O.InexactFloat64(),
				pk.Kline.C.InexactFloat64(),
				pk.PositionAmount.InexactFloat64(),
				pk.PositionCost.InexactFloat64(),
				pk.Balance.InexactFloat64(),
				pk.TotalValue.InexactFloat64(),
				pk.Drawdown.InexactFloat64(),
				pk.ProfitPercentage.InexactFloat64())
		}
	}

	fmt.Println("\n=============================================================")
}

// DisplayTrades 显示详细的交易记录
func (b *Backtest) DisplayTrades(limit int) {
	if len(b.trades) == 0 {
		fmt.Println("没有交易记录")
		return
	}

	// 确定显示的交易数量
	displayCount := limit
	if displayCount <= 0 || displayCount > len(b.trades) {
		displayCount = len(b.trades)
	}

	// 打印表头
	fmt.Println("\n========================= 交易记录 =========================")
	fmt.Printf("%-20s %-6s %-12s %-12s %-12s\n",
		"时间", "类型", "数量", "价格", "交易额(USDT)")

	// 显示交易记录
	for i := 0; i < displayCount; i++ {
		trade := b.trades[i]
		tradeType := "买入"
		if trade.Type.IsSell() {
			tradeType = "卖出"
		}

		tradeValue := trade.Amount.Mul(trade.Price)
		fmt.Printf("%-20s %-6s %-12.8f %-12.2f %-12.2f\n",
			trade.Time.Format("2006-01-02 15:04:05"),
			tradeType,
			trade.Amount.InexactFloat64(),
			trade.Price.InexactFloat64(),
			tradeValue.InexactFloat64())
	}

	// 如果有更多交易未显示
	if displayCount < len(b.trades) {
		fmt.Printf("\n... 还有 %d 条交易记录未显示 ...\n", len(b.trades)-displayCount)
	}

	fmt.Println("\n=============================================================")
}

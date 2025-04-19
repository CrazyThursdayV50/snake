# MACD (Moving Average Convergence Divergence) 指标

## 概述

MACD（Moving Average Convergence Divergence，移动平均线收敛与发散）是一种趋势跟踪动量指标，通过计算两条指数移动平均线（EMA）之间的差值来显示短期趋势相对于长期趋势的关系。MACD是交易中常用的技术指标之一，用于辨别市场的超买超卖状态以及潜在的趋势反转点。

## 计算方法

MACD由三个主要组成部分：

1. **MACD线**：快速EMA与慢速EMA的差值
   - 标准设置：MACD = 12日EMA - 26日EMA

2. **信号线**：MACD线的EMA
   - 标准设置：信号线 = 9日MACD的EMA

3. **柱状图**：MACD线与信号线的差值
   - 柱状图 = MACD - 信号线

## 使用方法

### 基本用法

```go
import (
    "snake/internal/indicates/macd"
    "snake/internal/models"
)

// 使用默认参数(12,26,9)创建MACD指标
macdIndicator := macd.New(klines...)

// 使用自定义参数创建MACD指标
customMacd := macd.NewWithParams(5, 10, 3, klines...)

// 使用新的K线更新MACD
updatedMacd := macdIndicator.NextKline(newKline)

// 访问MACD指标的各个值
macdValue := macdIndicator.MACD       // MACD线值
signalValue := macdIndicator.Signal    // 信号线值
histogramValue := macdIndicator.Histogram  // 柱状图值
```

### 交易信号解读

1. **MACD线穿越信号线**
   - MACD线从下方穿越信号线：买入信号
   - MACD线从上方穿越信号线：卖出信号

2. **零线交叉**
   - MACD线从负值交叉到正值：看涨信号
   - MACD线从正值交叉到负值：看跌信号

3. **背离**
   - 价格创新高，但MACD未创新高：潜在顶部背离，看跌信号
   - 价格创新低，但MACD未创新低：潜在底部背离，看涨信号

## 注意事项

- MACD是滞后指标，因为它基于历史价格数据
- 在横盘市场中，MACD可能产生误导性信号
- 最好与其他技术指标和分析方法结合使用
- 参数可以根据不同的市场和资产类型进行调整 
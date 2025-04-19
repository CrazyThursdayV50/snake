# 布林带-MACD联合策略

## 概述

布林带-MACD联合策略是一种结合了布林带和MACD两个技术指标的交易策略。该策略利用布林带识别价格波动范围和潜在的超买超卖点，同时使用MACD识别趋势方向和动量变化，从而生成更为可靠的交易信号。

## 策略原理

### 布林带 (Bollinger Bands)

布林带由三条线组成：
- 中轨：通常是20期的简单移动平均线(SMA)
- 上轨：中轨加上2倍标准差
- 下轨：中轨减去2倍标准差

布林带可以帮助识别价格波动范围，当价格接近或突破上轨时可能表示超买，接近或突破下轨时可能表示超卖。

### MACD (Moving Average Convergence Divergence)

MACD由三个组成部分：
- MACD线：快速EMA(通常是12日)减去慢速EMA(通常是26日)
- 信号线：MACD的9日EMA
- 柱状图：MACD线减去信号线

MACD可以帮助识别趋势方向和动量变化，当MACD线从下方穿过信号线时形成金叉(看涨信号)，从上方穿过信号线时形成死叉(看跌信号)。

## 策略逻辑

### 买入条件

同时满足以下条件时产生买入信号：
1. 价格低于布林带下轨（表示可能超卖）
2. MACD直方图由负变正（MACD形成金叉）

### 卖出条件

同时满足以下条件时产生卖出信号：
1. 价格高于布林带上轨（表示可能超买）
2. MACD直方图由正变负（MACD形成死叉）

### 默认参数设置

- 布林带周期：20
- MACD快速EMA周期：12
- MACD慢速EMA周期：26
- MACD信号线周期：9
- 每次交易量：当前仓位或可用余额的5%

## 使用方法

```go
import (
    "snake/internal/strategy/bolling-macd"
    "snake/internal/models"
)

// 创建策略实例
strategy := bollingmacd.New()

// 初始化策略
strategy.Init(initialPosition, initialBalance)

// 接收新的K线数据并更新策略
signal, err := strategy.Update(newKline)

// 根据返回的信号执行交易
switch signal.Type {
case types.SignalTypeBuy:
    // 执行买入
case types.SignalTypeSell:
    // 执行卖出
case types.SignalTypeHold:
    // 持有不动
}
```

## 优势与注意事项

### 优势

1. 结合两个强大的技术指标，提高信号可靠性
2. 布林带识别价格波动范围，MACD确认趋势方向
3. 有效避免假突破和虚假信号
4. 适用于各种市场环境

### 注意事项

1. 在震荡市场中可能产生较多交易信号，增加交易成本
2. 两个指标都有一定的滞后性
3. 在强趋势市场中，布林带的超买超卖信号可能不够可靠
4. 建议结合其他分析方法，如基本面分析、交易量分析等

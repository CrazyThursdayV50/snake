# RSI (相对强弱指标) 指标

## 概述

RSI（Relative Strength Index，相对强弱指标）是一种动量振荡器，由J. Welles Wilder Jr.在1978年开发。RSI测量价格变动的速度和变化，帮助识别超买或超卖条件。RSI值在0到100之间，传统上认为RSI高于70表示可能处于超买状态，低于30表示可能处于超卖状态。

## 计算方法

RSI通过比较上涨和下跌的平均幅度来计算，公式如下：

1. **计算价格变化**：
   - 对于每个周期，计算当前价格与前一周期价格的差值

2. **分离上涨和下跌**：
   - 如果价格上涨，则将变化值记为增益(gain)，下跌为0
   - 如果价格下跌，则将变化的绝对值记为损失(loss)，上涨为0

3. **计算平均上涨和平均下跌**：
   - 首次计算：对于前n个周期的上涨和下跌取平均值
   - 后续计算：使用平滑算法 
     - `Avg Gain = ((n-1) * 前一个 Avg Gain + 当前 Gain) / n`
     - `Avg Loss = ((n-1) * 前一个 Avg Loss + 当前 Loss) / n`

4. **计算相对强度(RS)**：
   - `RS = 平均上涨 / 平均下跌`

5. **计算RSI**：
   - `RSI = 100 - (100 / (1 + RS))`
   或者等价于：
   - `RSI = RS / (1 + RS) * 100`

## 使用方法

### 基本用法

```go
import (
    "snake/internal/indicates/rsi"
    "snake/internal/models"
)

// 创建RSI指标，使用默认14天周期
rsiIndicator := rsi.New(klines...)

// 使用自定义周期创建RSI指标
customRSI := rsi.NewWithPeriod(7, klines...)

// 使用新的K线更新RSI
updatedRSI := rsiIndicator.NextKline(newKline)

// 获取RSI值
rsiValue := rsiIndicator.Value

// 判断是否为买入信号（RSI < 30）
isBuy := rsiIndicator.IsBuy()

// 判断是否为卖出信号（RSI > 70）
isSell := rsiIndicator.IsSell()
```

### 交易信号解读

1. **超买和超卖**
   - RSI > 70：可能处于超买状态，考虑卖出
   - RSI < 30：可能处于超卖状态，考虑买入

2. **背离**
   - 价格创新高，但RSI未创新高：看跌背离
   - 价格创新低，但RSI未创新低：看涨背离

3. **中线交叉**
   - RSI从下向上穿过50线：看涨信号
   - RSI从上向下穿过50线：看跌信号

4. **区域突破**
   - RSI从超卖区（<30）向上突破：强烈看涨信号
   - RSI从超买区（>70）向下突破：强烈看跌信号

## 注意事项

- RSI在盘整市场中效果最佳，在强趋势市场中可能产生误导信号
- 超买不一定立即导致下跌，超卖不一定立即导致上涨
- 建议与其他技术指标和分析方法结合使用
- 不同资产和不同时间周期可能需要调整RSI的超买超卖阈值
- 默认周期为14，但可根据需要进行调整：短周期（如7）对价格变化更敏感，长周期（如21）产生更少但更可靠的信号 
package kline

import (
	"errors"
	"time"

	"github.com/CrazyThursdayV50/pkgo/json"
	"github.com/shopspring/decimal"
)

type Kline struct {
	// open
	O decimal.Decimal
	// close
	C decimal.Decimal
	// high
	H decimal.Decimal
	// low
	L decimal.Decimal
	// volume
	V decimal.Decimal
	// amount
	A decimal.Decimal
	// start: milisecond
	S int64
	// end: milisecond
	E int64
}

func (k *Kline) IsCurrent(kline *Kline) bool {
	return k.S == kline.S && k.E == kline.E
}

// is k before kline
func (k *Kline) IsBefore(kline *Kline) bool {
	diffK := k.E - k.S
	diffKline := kline.E - kline.S
	return k.E+1 == kline.S && diffK == diffKline
}

// is k after kline
func (s *Kline) IsAfter(kline *Kline) bool {
	return kline.IsBefore(s)
}

func (k *Kline) MarshalBinary() ([]byte, error) {
	if k == nil {
		return nil, errors.New("invalid receiver")
	}

	return json.JSON().Marshal(k)
}

func (k *Kline) UnmarshalBinary(data []byte) error {
	if k == nil {
		return errors.New("invalid receiver")
	}

	return json.JSON().Unmarshal(data, k)
}

// PositionKline 记录每个 K 线出现时的资产情况
type PositionKline struct {
	// K 线时间
	Time time.Time
	// K 线数据
	Kline *Kline
	// 当前持仓数量
	PositionAmount decimal.Decimal
	// 当前持仓成本
	PositionCost decimal.Decimal
	// 当前余额
	Balance decimal.Decimal
	// 当前资产总值（持仓市值 + 余额）
	TotalValue decimal.Decimal
	// 当前持仓盈亏（绝对值）
	ProfitAbsolute decimal.Decimal
	// 当前持仓盈亏（百分比）
	ProfitPercentage decimal.Decimal
	// 当前回撤（相对于最高资产值的百分比）
	Drawdown decimal.Decimal
	// 最高资产值
	PeakValue decimal.Decimal
}

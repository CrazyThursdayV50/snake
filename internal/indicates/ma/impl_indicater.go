package ma

import (
	"snake/internal/indicates"
)

// Next 计算下一个 Kline 对应的 MA
// 如果传入的 Kline 时间戳小于当前 MA 的时间戳，返回 nil
func (m MA) Next(kline *indicates.Kline) *MA {
	last := m.klines[len(m.klines)-1]
	if !last.IsBefore(kline) {
		return nil
	}

	klines := append(m.klines[1:], kline)
	return New(m.count).Klines(klines).Build()
}

func (m *MA) Update(kline *indicates.Kline) bool {
	last := m.klines[len(m.klines)-1]
	if !last.IsCurrent(kline) {
		return false
	}

	m.klines[len(m.klines)-1] = kline
	m.prices[len(m.prices)-1] = kline.C
	m.calculate()
	return true
}

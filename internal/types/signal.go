package types

// SignalType 表示信号类型
type SignalType int

const (
	// SignalTypeHold 持有信号
	SignalTypeHold SignalType = iota
	// SignalTypeBuy 买入信号
	SignalTypeBuy
	// SignalTypeSell 卖出信号
	SignalTypeSell
)

// IsBuy 判断是否为买入信号
func (s SignalType) IsBuy() bool {
	return s == SignalTypeBuy
}

// IsSell 判断是否为卖出信号
func (s SignalType) IsSell() bool {
	return s == SignalTypeSell
}

// IsHold 判断是否为持有信号
func (s SignalType) IsHold() bool {
	return s == SignalTypeHold
} 
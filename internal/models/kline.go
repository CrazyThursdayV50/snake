package models

import "github.com/shopspring/decimal"

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
	// start: second
	S int64
	// start: second
	E int64
}

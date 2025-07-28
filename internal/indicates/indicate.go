package indicates

import "snake/internal/kline"

type Kline = kline.Kline

type Indicater[T any] interface {
	Next(*Kline) T
	Update(*Kline) T
}

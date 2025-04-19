package workers

import (
	"snake/internal/kline/interval"
)

type IntervalTrigger[T any] struct {
	t map[interval.Interval]func(T)
}

func NewIntervalTrigger[T any]() *IntervalTrigger[T] {
	return &IntervalTrigger[T]{
		t: make(map[interval.Interval]func(T)),
	}
}

func (t *IntervalTrigger[T]) Add(i interval.Interval, f func(T)) {
	t.t[i] = f
}

func (t *IntervalTrigger[T]) Trigger(i interval.Interval, v T) bool {
	if f, ok := t.t[i]; ok {
		f(v)
		return true
	}
	return false
}

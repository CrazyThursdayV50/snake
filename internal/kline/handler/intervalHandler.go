package handler

type IntervalHandler[T any] struct {
	t map[string]T
}

func NewIntervalHandler[T any]() *IntervalHandler[T] {
	return &IntervalHandler[T]{
		t: make(map[string]T),
	}
}

func (t *IntervalHandler[T]) Add(i string, h T) {
	t.t[i] = h
}

func (t *IntervalHandler[T]) Get(i string) (T, bool) {
	if f, ok := t.t[i]; ok {
		return f, true
	}

	var h T
	return h, false
}

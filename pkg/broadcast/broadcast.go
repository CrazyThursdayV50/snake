package broadcast

import (
	"sync"
)

// Broadcast 广播器
type Broadcast[T any] struct {
	mu       sync.RWMutex
	channels map[chan T]struct{}
	closed   bool
}

// New 创建一个新的广播器
func New[T any]() *Broadcast[T] {
	return &Broadcast[T]{
		channels: make(map[chan T]struct{}),
	}
}

// Subscribe 订阅广播
func (b *Broadcast[T]) Subscribe() chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan T, 1)
	b.channels[ch] = struct{}{}
	return ch
}

// Unsubscribe 取消订阅
func (b *Broadcast[T]) Unsubscribe(ch chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.channels[ch]; ok {
		delete(b.channels, ch)
		close(ch)
	}
}

// Broadcast 广播消息
func (b *Broadcast[T]) Broadcast(msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return
	}

	for ch := range b.channels {
		select {
		case ch <- msg:
		default:
			// 如果channel已满，跳过该订阅者
		}
	}
}

// Close 关闭广播器
func (b *Broadcast[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	for ch := range b.channels {
		close(ch)
	}
	b.channels = nil
} 
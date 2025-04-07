package deepl

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Callback[T any] func(ctx context.Context, result T, err error)

type CMD[T any] struct {
	ctx    context.Context
	closed int32
	fn     func() (T, error)
}

func NewCMD[T any](ctx context.Context, fn func() (T, error)) *CMD[T] {
	return &CMD[T]{
		ctx: ctx,
		fn:  fn,
	}
}

func (self *CMD[T]) Closed() bool {
	return atomic.LoadInt32(&self.closed) >= 1
}

func (self *CMD[T]) Sync() (T, error) {
	if atomic.LoadInt32(&self.closed) >= 1 {
		var zero T
		return zero, fmt.Errorf("cmd is closed")
	}
	atomic.StoreInt32(&self.closed, 1)
	return self.fn()
}

func (self *CMD[T]) Async(handler Callback[T]) {
	if atomic.LoadInt32(&self.closed) >= 1 {
		var zero T
		handler(self.ctx, zero, fmt.Errorf("cmd is closed"))
		return
	}
	atomic.StoreInt32(&self.closed, 1)
	go func() {
		result, err := self.fn()
		handler(self.ctx, result, err)
	}()
}

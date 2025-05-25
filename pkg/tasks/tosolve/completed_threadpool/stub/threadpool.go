package stub

import (
	"context"
	completedthreadpool "threadpool_example/pkg/tasks/tosolve/completed_threadpool"
)

func NewThreadPool[T, E any](ctx context.Context, workers int, applier completedthreadpool.ApplierFunc[T, E]) *ThreadPool[T, E] {
	panic("Not implemented yet")
}

type ThreadPool[T, E any] struct{}

func (t *ThreadPool[T, E]) Incoming() chan<- T {
	panic("Not implemented yet")
}
func (t *ThreadPool[T, E]) Outgoing() <-chan E {
	panic("Not implemented yet")
}

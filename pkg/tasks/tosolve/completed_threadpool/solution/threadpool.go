package solution

import (
	"context"
	"sync"
	completedthreadpool "threadpool_example/pkg/tasks/tosolve/completed_threadpool"
)

func NewThreadPool[T, E any](ctx context.Context, workers int, applier completedthreadpool.ApplierFunc[T, E]) *ThreadPool[T, E] {
	p := &ThreadPool[T, E]{
		ctx:     ctx,
		in:      make(chan T, workers),
		out:     make(chan E, workers),
		workers: workers,
		applier: applier,
		wg:      &sync.WaitGroup{},
	}
	p.start()
	return p
}

type ThreadPool[T, E any] struct {
	ctx     context.Context
	in      chan T
	out     chan E
	workers int
	applier completedthreadpool.ApplierFunc[T, E]
	wg      *sync.WaitGroup
}

func (p *ThreadPool[T, E]) Incoming() chan<- T {
	return p.in
}

func (p *ThreadPool[T, E]) Outgoing() <-chan E {
	return p.out
}

func (p *ThreadPool[T, E]) worker() {
	defer p.wg.Done()
	for value := range p.in {
		p.out <- p.applier(value)
	}
}

func (p *ThreadPool[T, E]) start() {
	p.wg.Add(p.workers)
	go func(in chan T) {
		<-p.ctx.Done()
		close(p.in)
	}(p.in)
	go func(out chan<- E, wg *sync.WaitGroup) {
		wg.Wait()
		close(out)
	}(p.out, p.wg)
	for range p.workers {
		go p.worker()
	}
}

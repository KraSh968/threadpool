package good

import (
	"context"
	"sync"
	"threadpool_example/threadpool"
)

const (
	chanSizeMult = 2
)

func NewThreadPool[T, E any](ctx context.Context, workers int, applier threadpool.ApplierFunc[T, E]) *GoodThreadPool[T, E] {
	p := &GoodThreadPool[T, E]{
		ctx:     ctx,
		in:      make(chan T, workers*chanSizeMult),
		out:     make(chan E, workers*chanSizeMult),
		workers: workers,
		applier: applier,
		wg:      &sync.WaitGroup{},
	}
	p.start()
	return p
}

type GoodThreadPool[T, E any] struct {
	ctx     context.Context
	in      chan T
	out     chan E
	workers int
	applier threadpool.ApplierFunc[T, E]
	wg      *sync.WaitGroup
}

func (p *GoodThreadPool[T, E]) Incoming() chan<- T {
	return p.in
}

func (p *GoodThreadPool[T, E]) Outgoing() <-chan E {
	return p.out
}

func (p *GoodThreadPool[T, E]) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case value, ok := <-p.in:
			if !ok {
				return
			}
			p.out <- p.applier(value)
		}
	}
}

func (p *GoodThreadPool[T, E]) start() {
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

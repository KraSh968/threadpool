package bad

import (
	"context"
	"sync"
	"threadpool_example/threadpool"
)

func NewThreadPool[T, E any](ctx context.Context, chanSize int, applier threadpool.ApplierFunc[T, E]) *BadThreadPool[T, E] {
	p := &BadThreadPool[T, E]{
		ctx:     ctx,
		in:      make(chan T, chanSize),
		out:     make(chan E, chanSize),
		applier: applier,
		wg:      &sync.WaitGroup{},
	}
	p.start()
	return p
}

type BadThreadPool[T, E any] struct {
	ctx     context.Context
	in      chan T
	out     chan E
	applier threadpool.ApplierFunc[T, E]
	wg      *sync.WaitGroup
}

func (p *BadThreadPool[T, E]) Incoming() chan<- T {
	return p.in
}

func (p *BadThreadPool[T, E]) Outgoing() <-chan E {
	return p.out
}

func (p *BadThreadPool[T, E]) worker(value T) {
	defer p.wg.Done()
	p.out <- p.applier(value)
}

func (p *BadThreadPool[T, E]) start() {
	go func() {
		<-p.ctx.Done()
		close(p.in)
	}()
	go func() {
		for value := range p.in {
			p.wg.Add(1)
			go p.worker(value)
		}
		p.wg.Wait()
		close(p.out)
	}()
}

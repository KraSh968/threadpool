package good

import (
	"sync"
	"threadpool_example/threadpool"
)

const (
	chanSizeMult = 2
)

func NewThreadPool[T, E any](workers int, applier threadpool.ApplierFunc[T, E]) *GoodThreadPool[T, E] {
	p := &GoodThreadPool[T, E]{
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

func (p *GoodThreadPool[T, E]) Finalize() {
	close(p.in)
}

func (p *GoodThreadPool[T, E]) worker() {
	for value := range p.in {
		p.out <- p.applier(value)
	}
	p.wg.Done()
}

func (p *GoodThreadPool[T, E]) start() {
	p.wg.Add(p.workers)
	go func(out chan<- E, wg *sync.WaitGroup) {
		wg.Wait()
		close(out)
	}(p.out, p.wg)
	for range p.workers {
		go p.worker()
	}
}

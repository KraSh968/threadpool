package bad

import (
	"sync"
	"threadpool_example/threadpool"
)

func NewThreadPool[T, E any](chanSize int, applier threadpool.ApplierFunc[T, E]) *BadThreadPool[T, E] {
	p := &BadThreadPool[T, E]{
		in:      make(chan T, chanSize),
		out:     make(chan E, chanSize),
		applier: applier,
		wg:      &sync.WaitGroup{},
	}
	p.start()
	return p
}

type BadThreadPool[T, E any] struct {
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

func (p *BadThreadPool[T, E]) Finalize() {
	close(p.in)
}

func (p *BadThreadPool[T, E]) worker(value T) {
	defer p.wg.Done()
	p.out <- p.applier(value)
}

func (p *BadThreadPool[T, E]) start() {
	go func() {
		for value := range p.in {
			p.wg.Add(1)
			go p.worker(value)
		}
		p.wg.Wait()
		close(p.out)
	}()
}

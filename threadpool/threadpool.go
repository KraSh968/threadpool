package threadpool

type ApplierFunc[T any, E any] func(T) E

type ThreadPool[T any, E any] interface {
	Incoming() chan<- T
	Outgoing() <-chan E
	Finalize()
}

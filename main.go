package main

import (
	"fmt"
	threadPool "threadpool_example/threadpool/good"
)

func main() {
	threadPool := threadPool.NewThreadPool(10, func(a int) int {
		return 2 * a
	})

	go func() {
		defer threadPool.Finalize()
		for i := range 10 {
			threadPool.Incoming() <- i
		}
	}()

	lock := make(chan struct{})

	go func() {
		defer close(lock)
		for value := range threadPool.Outgoing() {
			fmt.Println(value)
		}
	}()

	<-lock
}

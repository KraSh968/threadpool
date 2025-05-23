package main

import (
	"context"
	"fmt"
	threadPool "threadpool_example/threadpool/good"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	threadPool := threadPool.NewThreadPool(ctx, 10, func(a int) int {
		return 2 * a
	})

	go func() {
		defer cancel()
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

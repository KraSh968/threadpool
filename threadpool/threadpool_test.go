package threadpool_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"threadpool_example/threadpool"
	"threadpool_example/threadpool/bad"
	"threadpool_example/threadpool/good"
	"time"
)

type PoolConstructor[T any, E any] func(context.Context, int) threadpool.ThreadPool[T, E]

type ErrTooManyGorotinesInUse struct {
	current int
	allowed int
}

func (e ErrTooManyGorotinesInUse) Error() string {
	return fmt.Sprintf("using %d gorotines, but allowed only %d", e.current, e.allowed)
}

type ErrGoroutinesStillAliveAfterShutdown struct {
	alive int
}

func (e ErrGoroutinesStillAliveAfterShutdown) Error() string {
	return fmt.Sprintf("%d goroutines still active after shutdown, too slow resources releasing or potential memory leak", e.alive)
}

type ErrTooManyTimeForProcessingSmallTask struct {
	time time.Duration
}

func (e ErrTooManyTimeForProcessingSmallTask) Error() string {
	return fmt.Sprintf("%v for processing small task, it is too much", e.time)
}

const allowedGoroCount = 20

var applier = func(value int) int { return value * 2 }

func badImplementation(ctx context.Context, size int) threadpool.ThreadPool[int, int] {
	return bad.NewThreadPool(ctx, size, applier)
}

func goodImplementation(ctx context.Context, size int) threadpool.ThreadPool[int, int] {
	return good.NewThreadPool(ctx, size, applier)
}

func LoadThreadPool(createFunc PoolConstructor[int, int], size int, topValue int) error {
	ctx, cancel := context.WithCancel(context.Background())
	pool := createFunc(ctx, size)
	errChan := make(chan error, 1)
	defer close(errChan)
	go func() {
		defer cancel()
		for i := 0; i < topValue; i++ {
			pool.Incoming() <- i

			currentGoroCount := runtime.NumGoroutine()
			if runtime.NumGoroutine() > allowedGoroCount {
				errChan <- ErrTooManyGorotinesInUse{current: currentGoroCount, allowed: allowedGoroCount}
				return
			}
		}
	}()

	for {
		select {
		case value, ok := <-pool.Outgoing():
			if !ok {
				return nil
			}
			_ = value
		case err := <-errChan:
			return err
		}
	}
}

func TestGoroutinesUsageInThreadPool(t *testing.T) {
	if err := LoadThreadPool(goodImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing good implementation: %s", err.Error())
	}

	if err := LoadThreadPool(badImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing bad implementation: %s", err.Error())
	}
}

func LoadAndShutDownThreadPool(createFunc PoolConstructor[int, int], size int, topValue int) error {
	ctx, cancel := context.WithCancel(context.Background())
	pool := createFunc(ctx, size)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				errChan <- fmt.Errorf("panic while processing test: %s", err)
			}
		}()
		beforeStart := runtime.NumGoroutine()
		for i := 0; i < topValue; i++ {
			pool.Incoming() <- i

			if i == topValue/2 {
				cancel()
				currentAmount := runtime.NumGoroutine()
				if currentAmount > beforeStart {
					errChan <- ErrGoroutinesStillAliveAfterShutdown{alive: currentAmount - beforeStart}
				}
				return
			}
		}
	}()

	for {
		select {
		case value, ok := <-pool.Outgoing():
			if !ok {
				return nil
			}
			_ = value
		case err := <-errChan:
			return err
		}
	}
}

func TestThreadPoolShutdown(t *testing.T) {

	if err := LoadAndShutDownThreadPool(goodImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing good implementation: %s", err.Error())
	}

	if err := LoadAndShutDownThreadPool(badImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing bad implementation: %s", err.Error())
	}
}

type variant struct {
	queueSize int
	values    int
}

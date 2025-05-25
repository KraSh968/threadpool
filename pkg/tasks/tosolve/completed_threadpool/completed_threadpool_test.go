package completedthreadpool_test

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"testing"
	"threadpool_example/pkg/data"
	completedthreadpool "threadpool_example/pkg/tasks/tosolve/completed_threadpool"
	"threadpool_example/pkg/tasks/tosolve/completed_threadpool/correct"
	"threadpool_example/pkg/tasks/tosolve/completed_threadpool/stub"
)

type PoolConstructor[T any, E any] func(context.Context, int, completedthreadpool.ApplierFunc[T, E]) completedthreadpool.ThreadPool[T, E]

const (
	allowedGoroCount = 20
)

var (
	multApplier = func(value int) int { return value * 2 }
)

func stubImplementation[T, E any](ctx context.Context, size int, applier completedthreadpool.ApplierFunc[T, E]) completedthreadpool.ThreadPool[T, E] {
	return stub.NewThreadPool(ctx, size, applier)
}

func solutionImplementation[T, E any](ctx context.Context, size int, applier completedthreadpool.ApplierFunc[T, E]) completedthreadpool.ThreadPool[T, E] {
	return correct.NewThreadPool(ctx, size, applier)
}

type ErrTooManyGorotinesInUse struct {
	current int
	allowed int
}

func (e ErrTooManyGorotinesInUse) Error() string {
	return fmt.Sprintf("using %d gorotines, but allowed only %d", e.current, e.allowed)
}
func LoadThreadPool(createFunc PoolConstructor[int, int], size int, topValue int) error {
	ctx, cancel := context.WithCancel(context.Background())
	pool := createFunc(ctx, size, multApplier)
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
	if err := LoadThreadPool(stubImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing stub implementation: %s", err.Error())
	}

	if err := LoadThreadPool(solutionImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing solution implementation: %s", err.Error())
	}
}

type ErrGoroutinesStillAliveAfterShutdown struct {
	alive int
}

func (e ErrGoroutinesStillAliveAfterShutdown) Error() string {
	return fmt.Sprintf("%d goroutines still active after shutdown, too slow resources releasing or potential memory leak", e.alive)
}
func LoadAndShutDownThreadPool(createFunc PoolConstructor[int, int], size int, topValue int) error {
	ctx, cancel := context.WithCancel(context.Background())
	pool := createFunc(ctx, size, multApplier)
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

	if err := LoadAndShutDownThreadPool(stubImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing stub implementation: %s", err.Error())
	}

	if err := LoadAndShutDownThreadPool(solutionImplementation, 10, 100000); err != nil {
		t.Errorf("error while testing solution implementation: %s", err.Error())
	}
}

func LoadAndVerifyValues(createFunc PoolConstructor[int, int], size int, topValue int) error {
	ctx, cancel := context.WithCancel(context.Background())
	pool := createFunc(ctx, size, multApplier)
	inputData := *data.GenerateInts(topValue)
	go func() {
		defer cancel()
		for _, value := range inputData {
			pool.Incoming() <- value
		}
	}()
	results := make([]int, 0, topValue)
	for value := range pool.Outgoing() {
		results = append(results, value)
	}
	if len(results) != len(inputData) {
		return fmt.Errorf("sent into %d values, got back %d", len(results), len(inputData))
	}

	slices.Sort(results)

	for i, item := range results {
		expected := multApplier(inputData[i])
		if item != multApplier(inputData[i]) {
			return fmt.Errorf("expected %d, got %d", expected, item)
		}
	}
	return nil
}

func TestThreadPoolValuesCorrect(t *testing.T) {
	if err := LoadAndShutDownThreadPool(stubImplementation, 4, 1000); err != nil {
		t.Errorf("error while testing stub implementation: %s", err.Error())
	}

	if err := LoadAndShutDownThreadPool(solutionImplementation, 4, 1000); err != nil {
		t.Errorf("error while testing solution implementation: %s", err.Error())
	}
}

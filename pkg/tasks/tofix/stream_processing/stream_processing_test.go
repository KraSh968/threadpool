package streamprocessing_test

import (
	"runtime"
	"slices"
	"testing"
	"threadpool_example/pkg/data"
	"threadpool_example/pkg/tasks"
	streamprocessing "threadpool_example/pkg/tasks/tofix/stream_processing"
	"time"
)

func TestProcessData(t *testing.T) {
	tests := []struct {
		name        string
		workers     int
		inputData   []int
		awaitedData []int
	}{
		{
			name:        "single worker",
			workers:     1,
			inputData:   *data.GenerateRange(1000),
			awaitedData: *data.MultiplyInts(data.GenerateRange(1000), 2),
		},
		{
			name:        "multiple workers",
			workers:     4,
			inputData:   *data.GenerateRange(5000),
			awaitedData: *data.MultiplyInts(data.GenerateRange(5000), 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			input := make(chan int, tt.workers)
			output := make(chan int, tt.workers)
			errChan := make(chan error, 1)

			go func() {
				defer close(input)
				startCount := runtime.NumGoroutine()
				for _, value := range tt.inputData {
					input <- value
					current := runtime.NumGoroutine()
					if current > startCount+tt.workers {
						errChan <- tasks.ErrTooManyGorotinesInUse{Current: current, Allowed: startCount + tt.workers}
						break
					}
				}
			}()

			streamprocessing.ProcessData(input, output, tt.workers)

			results := make([]int, 0, len(tt.inputData))

			timeout := time.NewTimer(time.Second * 5)
			defer timeout.Stop()

		loop:
			for {
				select {
				case value, ok := <-output:
					if !ok {
						break loop
					}
					results = append(results, value)
				case <-timeout.C:
					t.Fatalf("timeout exeeded while running test")
				case err := <-errChan:
					t.Fatalf("to many goroutines in use: %s", err)
				}
			}

			if len(results) != len(tt.awaitedData) {
				t.Fatalf("non-equal lengths: %d of inputData and %d of results", len(tt.inputData), len(results))
			}

			slices.Sort(results)

			for i, value := range results {
				if value != tt.awaitedData[i] {
					t.Errorf("non-equal values at index %d: awaited %d, got %d", i, tt.awaitedData[i], value)
				}
			}
		})
	}
}

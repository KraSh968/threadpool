package contextusing_test

import (
	"context"
	"testing"
	"threadpool_example/pkg/data"
	contextusing "threadpool_example/pkg/tasks/tofix/context_using"
	"time"
)

func TestProcessDataWithContext(t *testing.T) {
	tests := []struct {
		name      string
		workers   int
		inputData []int
	}{
		{
			name:      "little amount of workers",
			workers:   4,
			inputData: *data.GenerateInts(10000),
		},
		{
			name:      "big amount of workers",
			workers:   12,
			inputData: *data.GenerateInts(100000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			input := make(chan int, tt.workers)
			output := make(chan int, tt.workers)
			ctx, cancel := context.WithCancel(context.Background())

			dropIndex := len(tt.inputData) / 2

			go func() {
				for i, value := range tt.inputData {
					input <- value
					if i == dropIndex {
						cancel()
					}
				}
			}()

			contextusing.ProcessDataWithContext(ctx, input, output, tt.workers)

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
					t.Fatalf("timeout exeeded while running test: probably wrong usage of context.Context")
				}
			}
		})
	}
}

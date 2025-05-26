package contextusing

import (
	"context"
	"sync"
)

// Некорректная реализация обработки данных с использованим контекста.
// Функция должна использовать context.Context для возможности отмены обработки,
// при этом корректно завершая работу, выполняя graceful shutdown
func ProcessDataWithContext(ctx context.Context, input <-chan int, output chan<- int, workers int) {
	wg := &sync.WaitGroup{}

	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for value := range input {
				output <- value * 2
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()
}

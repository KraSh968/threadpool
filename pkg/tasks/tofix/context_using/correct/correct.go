package correct

import (
	"context"
	"sync"
)

// Исправленная версия функции для обработки данных с использованием context.Context
func ProcessDataWithContext(ctx context.Context, input <-chan int, output chan<- int, workers int) {
	wg := &sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Обработка значений из двух каналов
				select {
				// Канал контекста, который закрывается после отмены контекста и инициирует выход из цикла
				case <-ctx.Done():
					return

				// Входной канал, из которого читаются значения
				case value, ok := <-input:
					if !ok {
						return
					}
					output <- value * 2
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()
}

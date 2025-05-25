package correct

import "sync"

// Исправленная версия функции для обработки данных с использованием каналов
func ProcessData(input <-chan int, output chan<- int, workers int) {
	wg := &sync.WaitGroup{}

	// Создается фиксированное количество каналов
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for value := range input {
				output <- value * 2
			}
		}()
	}

	go func() {
		// Производиться освобождение ресурсов после обработки всех входящих значений
		wg.Wait()
		close(output)
	}()
}

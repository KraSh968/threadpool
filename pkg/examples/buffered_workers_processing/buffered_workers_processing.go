package bufferedworkersprocessing

import (
	"fmt"
	"sync"
	"threadpool_example/pkg/data"
)

// Более эффективная обработка данных использующая воркеры и буферизированные каналы. Улучшения:
//   - Использование буферизированных каналов для получения и отдачи данных обеспечивает более стабильную нагрузку
//     при неравномерном поступлении и отдаче данных а также потенциально более высокую производительность за счет
//     отсутствия жесткой двусторонней синхронизацией при чтении или записи в поток
//
// Проблемы:
//   - Отсутствие возможности завершить обработку механизмами, встроенными в стандартную библиотеку языка,
//     соблюдая при этом принцип graceful shutdown
func RunBaseExample() {
	fmt.Println("Запуск примера с обработкой воркерами и буферизированными каналами...")

	const inputSize = 1_000_000
	const workersCount = 5
	const bufferSize = 20

	fmt.Printf("Кол-во значений: %d, кол-во воркеров: %d, буфер каналов: %d\n", inputSize, workersCount, bufferSize)

	// Создание буферизированных каналов
	inputChan := make(chan int, bufferSize)
	resChan := make(chan int, bufferSize)
	wg := &sync.WaitGroup{}

	// Запуск отдельной горутины для записи значений во входящий канал
	go func() {
		for _, value := range *data.GenerateInts(inputSize) {
			inputChan <- value
		}
		// Закрытие канала после окончания записи
		close(inputChan)
		fmt.Println("Канал для записи закрыт")
	}()

	// Запуск фиксированного количества горутин-воркеров
	for range workersCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Чтение из входного канала и запись в выходной
			for value := range inputChan {
				resChan <- value * 2
			}
		}()
	}

	go func() {
		// Ожидание завершения работы воркеров для закрытия выходного канала
		wg.Wait()
		close(resChan)
		fmt.Println("Канал для чтения закрыт")
	}()

	extracted := 0
	// Чтение значений из канала
	for value := range resChan {
		_ = value
		extracted++
	}
	fmt.Printf("Обработано %d значений\n", extracted)
}

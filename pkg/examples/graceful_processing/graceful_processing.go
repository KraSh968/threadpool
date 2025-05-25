package gracefulprocessing

import (
	"context"
	"fmt"
	"sync"
	"threadpool_example/pkg/data"
)

// Более эффективная обработка данных использующая воркеры, буферизированные каналы и возможность остановки. Улучшения:
//   - Использование context.Context для получения сигнала о завершении обработки вместо прямого вызова close() для канала,
//     что позволяет корректно завершить работу пула, освободив все ресурсы и не потеряв значения, которые уже были в него
//     помещены

// Данный пример гарантирует обработку всех значений, успевших попасть в пул до момента завершения его контекста.
func RunBaseExample() {
	fmt.Println("Запуск примера с обработкой воркерами, буферизированными каналами и использованием context.Context, гарантирующего обработку всех значений...")

	const inputSize = 1_000_000
	const workersCount = 5
	const bufferSize = 20

	fmt.Printf("Кол-во значений: %d, кол-во воркеров: %d, буфер каналов: %d\n", inputSize, workersCount, bufferSize)

	// Создание буферизированных каналов
	inputChan := make(chan int, bufferSize)
	resChan := make(chan int, bufferSize)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Запуск отдельной горутины для записи значений во входящий канал
	go func() {
		for _, value := range *data.GenerateInts(inputSize) {
			inputChan <- value
		}
		// Отмена контекста
		cancel()
		fmt.Println("Вызвана отмена контекста")
	}()

	// Отдельная горутина для отслеживания состояния контекста и закрытия канала по завершению контекста
	go func() {
		<-ctx.Done()
		close(inputChan)
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

// Альтернативная версия с немедленным завершением обработчиков при закрытии контекста (не гарантирует обработку всех поступивших значение
// и допускает потерю данных)
func RunAlternativeExample() {
	fmt.Println("Запуск примера с обработкой воркерами, буферизированными каналами и использованием context.Context, не гарантирующего обработку всех значений...")

	const inputSize = 1_000_000
	const workersCount = 5
	const bufferSize = 20

	fmt.Printf("Кол-во значений: %d, кол-во воркеров: %d, буфер каналов: %d\n", inputSize, workersCount, bufferSize)

	// Создание буферизированных каналов
	inputChan := make(chan int, bufferSize)
	resChan := make(chan int, bufferSize)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Запуск отдельной горутины для записи значений во входящий канал
	go func() {
		for _, value := range *data.GenerateInts(inputSize) {
			inputChan <- value
		}
		// Отмена контекста
		cancel()
		fmt.Println("Вызвана отмена контекста")
	}()

	// Отдельная горутина для отслеживания состояния контекста и закрытия канала по завершению контекста
	go func() {
		<-ctx.Done()
		close(inputChan)
	}()

	// Запуск фиксированного количества горутин-воркеров
	for range workersCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Чтение из входного канала и, одновременно, ожидание сигнала завершения от контекста
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-inputChan:
					if !ok {
						return
					}
					resChan <- value * 2
				}
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

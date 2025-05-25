package streamprocessing

// Некорректная реализация многопоточной обработки данных.
// Функция должна запустить определенное количество воркеров,
// фиксируемое входным параметром, а также освобождать все ресурсы
// и не терять значения
func ProcessData(input <-chan int, output chan<- int, workers int) {

	for value := range input {
		go func() {
			output <- value * 2
		}()
	}
}

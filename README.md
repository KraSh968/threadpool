# Параллельное программирование: Threadpool

*Проект содержит примеры кода, демонстрирующие реализацию и использование пулов потоков в языке Golang*

## Примеры
1. Синхронная обработка
2. Неэффективная обработка с помощью горутин
3. Обработка обработка с использованием фиксированного числа горутин
4. Обработка с использованием фиксированного числа горутин и буферизированных каналов канала
5. Обработка с использованием фиксированного числа горутин, буферизированных каналов и контекста для поддержки graceful shutdown

## Структура
| Директория          | Назначение                                                              |
|---------------------|-------------------------------------------------------------------------|
| `cmd/main`          | Входная точка приложения                                                |
| `pkg/examples`      | Примеры, отражающие постепенное совершенствование реализации threadpool |
| `pkg/tasks/tosolve` | Задачи на корректную реализацию интерфейса                              |
| `pkg/tasks/tofix`   | Задачи на исправление существующей реализации                           |

## Запуск примеров
Запуск без параметров, демонстрирующий использование:
```go
go run cmd/main/entry.go
```
Запуск конкретного примера:
```go
go run cmd/main/entry.go -example [имя]
```
Запуск всех примеров:
```go
go run cmd/main/entry.go -example all
```

## Запуск тестов

Все тесты
```go
go test -count=1 ./...`
```

Все тесты со всеми бенчмарками
```go
go test ./... -count=1 -bench=./... -benchmem`
```

Тесты конкретного пакета (в примере `completed_threadpool`)
```go
go test -count=1 threadpool_example/pkg/tasks/tosolve/completed_threadpool`
```
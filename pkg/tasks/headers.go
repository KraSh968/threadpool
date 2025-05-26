package tasks

import "fmt"

type ErrTooManyGorotinesInUse struct {
	Current int
	Allowed int
}

func (e ErrTooManyGorotinesInUse) Error() string {
	return fmt.Sprintf("using %d gorotines, but allowed only %d", e.Current, e.Allowed)
}

type ErrGoroutinesStillAliveAfterShutdown struct {
	Alive int
}

func (e ErrGoroutinesStillAliveAfterShutdown) Error() string {
	return fmt.Sprintf("%d goroutines still active after shutdown, too slow resources releasing or potential memory leak", e.Alive)
}

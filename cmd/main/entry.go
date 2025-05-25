package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	bufferedworkersprocessing "threadpool_example/pkg/examples/buffered_workers_processing"
	gracefulprocessing "threadpool_example/pkg/examples/graceful_processing"
	sequentialprocessing "threadpool_example/pkg/examples/sequential_processing"
	uneffectiveparallelprocessing "threadpool_example/pkg/examples/uneffective_parallel_processing"
	workersprocessing "threadpool_example/pkg/examples/workers_processing"
)

type Example struct {
	name   string
	runner func()
}

type Examples []Example

func (ms Examples) Run(name string) error {
	for _, mode := range ms {
		if mode.name == name {
			mode.runner()
			return nil
		}
	}
	return fmt.Errorf("mode with name %s was not found", name)
}

func (ms Examples) List() []string {
	result := make([]string, 0, len(ms))
	for _, value := range ms {
		result = append(result, value.name)
	}
	return result
}

var examples Examples = Examples{
	Example{"sequential", sequentialprocessing.RunBaseExample},
	Example{"uneffective", uneffectiveparallelprocessing.RunBaseExample},
	Example{"with-workers", workersprocessing.RunBaseExample},
	Example{"buffered-with-workers", bufferedworkersprocessing.RunBaseExample},
	Example{"graceful", gracefulprocessing.RunBaseExample},
	Example{"graceful-alt", gracefulprocessing.RunAlternativeExample},
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func allModes() []string {
	return append([]string{"all"}, examples.List()...)

}

func runAll() {
	for _, exampleName := range examples {
		exampleName.runner()
	}
}

func runMode(exampleName string) {
	if exampleName == "all" {
		runAll()
	} else {
		if err := examples.Run(exampleName); err != nil {
			usage()
		}
	}
}

func main() {
	exampleName := flag.String("example", "", fmt.Sprintf("Доступные значения: %s", strings.Join(allModes(), ", ")))
	flag.Parse()
	if len(os.Args) < 2 {
		usage()
	}
	runMode(*exampleName)
}

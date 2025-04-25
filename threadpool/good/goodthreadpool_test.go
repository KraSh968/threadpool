package good

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"testing"
)

type variant struct {
	workers int
	values  int
}

var trackGoroCount bool

func init() {
	flag.BoolVar(&trackGoroCount, "trackGoroCount", false, "enable test fail when too many gorotines in use")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func BenchmarkGoodThreadpool(b *testing.B) {

	b.ReportAllocs()
	var variants []variant = []variant{{
		workers: 2,
		values:  5000,
	}, {
		workers: 10,
		values:  100000,
	}, {
		workers: runtime.NumCPU(),
		values:  1000000,
	}}

	for _, variant := range variants {
		b.Run(fmt.Sprintf("%d workers %d values", variant.workers, variant.values), func(b *testing.B) {
			p := NewThreadPool(variant.workers, func(i int) int { return i * 2 })
			go func() {
				defer p.Finalize()
				for i := 0; i < variant.values; i++ {
					p.Incoming() <- i

					if trackGoroCount {
						gos := runtime.NumGoroutine()
						if gos > 5+variant.workers {
							b.Fatalf("to many goroutines in use: %d goroutines with %d workers\n", gos, variant.workers)
						}
					}
				}
			}()
			acc := 0
			for value := range p.Outgoing() {
				acc += value
			}
		})
	}

}

package bad

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"testing"
)

type variant struct {
	queueSize int
	values    int
}

var trackGoroCount bool

func init() {
	flag.BoolVar(&trackGoroCount, "trackGoroCount", false, "enable test fail when too many gorotines in use")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func BenchmarkBadThreadPool(b *testing.B) {

	b.ReportAllocs()
	runtime.GC()

	var variants []variant = []variant{{
		queueSize: 2,
		values:    5000,
	}, {
		queueSize: 10,
		values:    100000,
	}, {
		queueSize: runtime.NumCPU(),
		values:    1000000,
	}}

	for _, variant := range variants {
		b.Run(fmt.Sprintf("%d queue size %d values", variant.queueSize, variant.values), func(b *testing.B) {
			p := NewThreadPool(variant.queueSize, func(i int) int { return i * 2 })
			go func() {
				defer p.Finalize()
				for i := 0; i < variant.values; i++ {
					p.Incoming() <- i

					if trackGoroCount {
						gos := runtime.NumGoroutine()
						if gos > 5+variant.queueSize {
							b.Fatalf("to many goroutines in use: %d goroutines with %d queue size\n", gos, variant.queueSize)
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

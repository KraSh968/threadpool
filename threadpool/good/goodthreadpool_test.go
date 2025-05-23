package good

import (
	"context"
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
			ctx, cancel := context.WithCancel(context.Background())
			p := NewThreadPool(ctx, variant.workers, func(i int) int { return i * 2 })
			go func() {
				defer cancel()
				for i := 0; i < variant.values; i++ {
					p.Incoming() <- i
				}
			}()
			for value := range p.Outgoing() {
				_ = value
			}
		})
	}

}

package bad

import (
	"context"
	"fmt"
	"runtime"
	"testing"
)

type variant struct {
	queueSize int
	values    int
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
			ctx, cancel := context.WithCancel(context.Background())
			p := NewThreadPool(ctx, variant.queueSize, func(i int) int { return i * 2 })
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

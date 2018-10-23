package meow_test

import (
	"fmt"
	"testing"

	"github.com/mmcloughlin/meow"
)

// buffer is pre-allocated to avoid allocations in benchmarking functions themselves.
var buffer = make([]byte, 1<<20)

// sink is a running sum of checksum output, to ensure the benchmarked function isn't optimized away.
var sink byte

func BenchmarkChecksum(b *testing.B) {
	var sizes []int
	for n := uint(0); n <= 20; n++ {
		sizes = append(sizes, 1<<n)
	}

	for _, size := range sizes {
		name := fmt.Sprintf("size=%d", size)
		b.Run(name, func(b *testing.B) {
			data := buffer[:size]
			b.SetBytes(int64(size))
			for i := 0; i < b.N; i++ {
				s := meow.Checksum(0, data)
				sink += s[0]
			}
		})
	}
}

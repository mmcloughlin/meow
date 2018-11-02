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

// benchmarkTarget is a hash function to be benchmarked.
type benchmarkTarget func([]byte) byte

func benchmarkChecksum(b *testing.B, f benchmarkTarget) {
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
				sink += f(data)
			}
		})
	}
}

func BenchmarkChecksum(b *testing.B) {
	benchmarkChecksum(b, func(data []byte) byte {
		h := meow.Checksum(0, data)
		return h[0]
	})
}

func BenchmarkHash(b *testing.B) {
	benchmarkChecksum(b, func(data []byte) byte {
		h := meow.New(0)
		h.Write(data)
		return h.Sum(nil)[0]
	})
}

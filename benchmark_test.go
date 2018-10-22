package meow_test

import (
	"fmt"
	"testing"

	"github.com/mmcloughlin/meow"
)

func BenchmarkChecksum(b *testing.B) {
	var sizes []int
	for n := uint(1); n <= 20; n++ {
		sizes = append(sizes, 1<<n)
	}

	for _, size := range sizes {
		name := fmt.Sprintf("size=%d", size)
		b.Run(name, func(b *testing.B) {
			data := make([]byte, size)
			for i := 0; i < b.N; i++ {
				_ = meow.Checksum(0, data)
			}
			b.SetBytes(int64(size))
		})
	}
}

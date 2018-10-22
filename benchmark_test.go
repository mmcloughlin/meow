package meow_test

import (
	"fmt"
	"testing"

	"github.com/mmcloughlin/meow"
)

func BenchmarkChecksum(b *testing.B) {
	sizes := []int{1}

	for _, size := range sizes {
		name := fmt.Sprintf("size=%d", size)
		b.Run(name, func(b *testing.B) {
			data := make([]byte, size)
			for i := 0; i < b.N; i++ {
				_ = meow.Sum(0, data)
			}
			b.SetBytes(int64(size))
		})
	}
}

package meow_test

import (
	"fmt"

	"github.com/mmcloughlin/meow"
)

func ExampleChecksum() {
	checksum := meow.Checksum(0, []byte("Hello, World!"))
	fmt.Printf("%x\n", checksum)
	// Output: a8cfb4aad7eada8ef007aafe27135386
}

package meow_test

import (
	"fmt"
	"io"

	"github.com/mmcloughlin/meow"
)

func ExampleChecksum() {
	checksum := meow.Checksum(0, []byte("Hello, World!"))
	fmt.Printf("%x\n", checksum)
	// Output: a8cfb4aad7eada8ef007aafe27135386
}

func ExampleNew() {
	h := meow.New(0)
	io.WriteString(h, "Hello, ")
	io.WriteString(h, "World!")
	fmt.Printf("%x\n", h.Sum(nil))
	// Output: a8cfb4aad7eada8ef007aafe27135386
}

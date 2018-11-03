// +build gofuzz

package meow

import (
	"bytes"
	"encoding/binary"
)

func Fuzz(data []byte) int {
	if len(data) < 8 {
		return 0
	}
	seed := binary.BigEndian.Uint64(data)
	data = data[8:]

	// Compute checksum via Checksum function.
	c0 := Checksum(seed, data)

	// Compute with the hash.Hash interface.
	h := New(seed)
	h.Write(data)
	c1 := h.Sum(nil)

	if !bytes.Equal(c0[:], c1) {
		panic("Checksum != hash.Hash")
	}

	// Compute checksum via pure Go.
	checksumgo(seed, c1, data)

	if !bytes.Equal(c0[:], c1) {
		panic("Checksum != checksumgo")
	}

	return 0
}

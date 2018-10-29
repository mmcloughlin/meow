package meow

import (
	"crypto/aes"
	"hash"
)

//go:generate go run make_block.go

// Meow hash version implemented by this package.
const (
	Version     = 2
	VersionName = "0.2/Ragdoll"
)

// BlockSize is the underlying block size of Meow hash in bytes.
const BlockSize = 256

// Size of a Meow checksum in bytes.
const Size = 16

// Variables capturing the implementation. Default to the pure go fallback.
var (
	implementation = "go"
	checksum       = fallback
)

// Checksum returns the Meow checksum of data.
func Checksum(seed uint64, data []byte) [Size]byte {
	var dst [Size]byte
	checksum(seed, dst[:], data)
	return dst
}

// New returns a 128-bit Meow hash.
func New(seed uint64) hash.Hash {
	return &digest{seed: seed}
}

// digest computes Meow hash in a streaming fashion.
type digest struct {
	seed   uint64
	s      [BlockSize]byte // streams
	b      [BlockSize]byte // pending block
	n      int             // number of (initial) bytes populated in b
	t      []byte          // the trailing block of data written to the hash
	length uint64          // total length written
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Reset() {
	for i := 0; i < BlockSize; i++ {
		d.s[i] = 0
	}
	d.n = 0
	d.length = 0
	d.t = nil
}

func (d *digest) Write(p []byte) (int, error) {
	N := len(p)
	d.length += uint64(N)

	// Update trailing block.
	if len(p) >= aes.BlockSize {
		d.t = p[N-aes.BlockSize:]
	} else {
		d.t = append(d.t, p...)
	}
	if len(d.t) > aes.BlockSize {
		d.t = d.t[len(d.t)-aes.BlockSize:]
	}

	// Combine with any pending data.
	if d.n > 0 {
		n := copy(d.b[d.n:], p)
		d.n += n
		if d.n == BlockSize {
			blocks(d.s[:], d.b[:])
			d.n = 0
		}
		p = p[n:]
	}

	// Hash any entire blocks.
	if len(p) >= BlockSize {
		n := len(p) &^ (BlockSize - 1)
		blocks(d.s[:], p[:n])
		p = p[n:]
	}

	// Keep any remaining data.
	if len(p) > 0 {
		d.n = copy(d.b[:], p)
	}

	return N, nil
}

func (d *digest) Sum(b []byte) []byte {
	var dst [Size]byte
	dt := *d
	finish(dt.seed, dt.s[:], dst[:], dt.b[:d.n], dt.t, dt.length)
	return append(b, dst[:]...)
}

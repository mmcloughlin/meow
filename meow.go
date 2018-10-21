package meow

// BlockSize is the underlying block size of Meow hash in bytes.
const BlockSize = 256

// Size of a Meow checksum in bytes.
const Size = 64

// block hashes one Meow block.
func block(state *byte, src []byte)

// sum computes the Meow checksum of data and writes to dst.
//go:noescape
func sum(seed uint64, dst *byte, src []byte)

// Sum returns the Meow checksum of data.
func Sum(seed uint64, data []byte) [Size]byte {
	var dst [Size]byte
	sum(seed, &dst[0], data)
	return dst
}

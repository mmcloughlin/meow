package meow

// BlockSize is the underlying block size of Meow hash in bytes.
const BlockSize = 256

// Size of a Meow checksum in bytes.
const Size = 64

// checksum computes the Meow checksum of data and writes to dst.
func checksum(seed uint64, dst *byte, src []byte)

// Checksum returns the Meow checksum of data.
func Checksum(seed uint64, data []byte) [Size]byte {
	var dst [Size]byte
	checksum(seed, &dst[0], data)
	return dst
}

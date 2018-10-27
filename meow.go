package meow

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

// checksum is the underlying implementation. Default to pure go fallback.
var checksum = fallback

// Checksum returns the Meow checksum of data.
func Checksum(seed uint64, data []byte) [Size]byte {
	var dst [Size]byte
	checksum(seed, dst[:], data)
	return dst
}

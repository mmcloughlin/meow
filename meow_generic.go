package meow

import (
	"crypto/aes"
	"encoding/binary"
)

// fallback is a pure go implementation of Meow checksum.
func fallback(seed uint64, dst, src []byte) {
	// Initialize streams.
	var s [BlockSize]byte

	// Handle full 256-byte blocks.
	cur := src
	for len(cur) >= BlockSize {
		for i := 0; i < BlockSize; i += aes.BlockSize {
			aesdec(cur[i:], s[i:], s[i:])
		}
		cur = cur[BlockSize:]
	}

	// Handle full 16-byte blocks.
	i := 0
	for len(cur) >= aes.BlockSize {
		aesdec(cur, s[i:], s[i:])
		cur = cur[aes.BlockSize:]
		i += aes.BlockSize
	}

	// Partial block.
	if len(cur) > 0 {
		var partial []byte
		if len(src) >= aes.BlockSize {
			partial = src[len(src)-aes.BlockSize:]
		} else {
			partial = make([]byte, aes.BlockSize)
			copy(partial[:], cur)
		}

		aesdec(partial, s[15*aes.BlockSize:], s[15*aes.BlockSize:])
	}

	// Combine.
	m0 := s[7*aes.BlockSize : 8*aes.BlockSize]
	ordering := []int{10, 4, 5, 12, 8, 0, 1, 9, 13, 2, 6, 14, 3, 11, 15}
	for _, i := range ordering {
		aesdec(s[i*aes.BlockSize:], m0, m0)
	}

	// Mixer.
	var mixer [aes.BlockSize]byte
	binary.LittleEndian.PutUint64(mixer[:], seed-uint64(len(src)))
	binary.LittleEndian.PutUint64(mixer[8:], seed+uint64(len(src))+1)

	for i := 0; i < 3; i++ {
		aesdec(mixer[:], m0, m0)
	}

	copy(dst, m0)
}

// aesdec performs one round of AES decryption.
func aesdec(key, dst, src []byte) {
	s0 := binary.BigEndian.Uint32(src[0:4])
	s1 := binary.BigEndian.Uint32(src[4:8])
	s2 := binary.BigEndian.Uint32(src[8:12])
	s3 := binary.BigEndian.Uint32(src[12:16])

	k0 := binary.BigEndian.Uint32(key[0:4])
	k1 := binary.BigEndian.Uint32(key[4:8])
	k2 := binary.BigEndian.Uint32(key[8:12])
	k3 := binary.BigEndian.Uint32(key[12:16])

	t0 := k0 ^ td0[uint8(s0>>24)] ^ td1[uint8(s3>>16)] ^ td2[uint8(s2>>8)] ^ td3[uint8(s1)]
	t1 := k1 ^ td0[uint8(s1>>24)] ^ td1[uint8(s0>>16)] ^ td2[uint8(s3>>8)] ^ td3[uint8(s2)]
	t2 := k2 ^ td0[uint8(s2>>24)] ^ td1[uint8(s1>>16)] ^ td2[uint8(s0>>8)] ^ td3[uint8(s3)]
	t3 := k3 ^ td0[uint8(s3>>24)] ^ td1[uint8(s2>>16)] ^ td2[uint8(s1>>8)] ^ td3[uint8(s0)]

	binary.BigEndian.PutUint32(dst[0:4], t0)
	binary.BigEndian.PutUint32(dst[4:8], t1)
	binary.BigEndian.PutUint32(dst[8:12], t2)
	binary.BigEndian.PutUint32(dst[12:16], t3)
}

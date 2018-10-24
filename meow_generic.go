package meow

import (
	"crypto/aes"
	"encoding/binary"
)

// fallback is a pure go implementation of Meow checksum.
func fallback(seed uint64, dst, src []byte) {
	// Build IV.
	var iv [64]byte
	for i := 0; i < 64; i += 16 {
		binary.LittleEndian.PutUint64(iv[i+0:], seed)
		binary.LittleEndian.PutUint64(iv[i+8:], seed+uint64(len(src))+1)
	}

	// Initialize streams with IV.
	var s [BlockSize]byte
	copy(s[0:], iv[:])
	copy(s[64:], iv[:])
	copy(s[128:], iv[:])
	copy(s[192:], iv[:])

	// Handle full blocks.
	for len(src) >= BlockSize {
		for i := 0; i < BlockSize; i += aes.BlockSize {
			aesdec(src[i:], s[i:], s[i:])
		}
		src = src[BlockSize:]
	}

	// Partial block.
	if len(src) > 0 {
		var p [BlockSize]byte
		copy(p[0:], iv[:])
		copy(p[64:], iv[:])
		copy(p[128:], iv[:])
		copy(p[192:], iv[:])
		copy(p[0:], src)

		for i := 0; i < BlockSize; i += aes.BlockSize {
			aesdec(p[i:], s[i:], s[i:])
		}
	}

	// Rotations.
	var r [Size]byte
	copy(r[:], iv[:])
	for t := 0; t < 4; t++ {
		for j := 0; j < 4; j++ {
			for i := 0; i < 4; i++ {
				idx := 4*j + (i+t)%4
				aesdec(s[idx*aes.BlockSize:], r[i*aes.BlockSize:], r[i*aes.BlockSize:])
			}
		}
	}

	// Final merge.
	for t := 0; t < 5; t++ {
		for i := 0; i < Size; i += aes.BlockSize {
			aesdec(iv[i:], r[i:], r[i:])
		}
	}

	copy(dst, r[:])
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

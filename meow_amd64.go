// +build !noasm

package meow

import "golang.org/x/sys/cpu"

func init() {
	// AVX required for VEX-encoded AES instruction, which allows non-aligned memory addresses.
	if cpu.X86.HasAES && cpu.X86.HasAVX {
		checksum = checksumAsm
	}
}

// checksumAsm implements Meow checksum in assembly.
func checksumAsm(seed uint64, dst, src []byte)

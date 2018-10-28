// +build !noasm

package meow

func init() {
	//// AVX required for VEX-encoded AES instruction, which allows non-aligned memory addresses.
	//if cpu.X86.HasAES && cpu.X86.HasAVX {
	//	checksum = checksum128
	//}
	checksum = checksum512
}

// checksum128 implements Meow checksum with AES-NI.
func checksum128(seed uint64, dst, src []byte)

// checksum512 implements Meow checksum with VAES-512.
func checksum512(seed uint64, dst, src []byte)

// +build !noasm

package meow

var cpu struct {
	HasOSXSAVE    bool
	HasAES        bool
	HasVAES       bool
	HasAVX        bool
	HasAVX512F    bool
	HasAVX512VL   bool
	HasAVX512DQ   bool
	EnabledAVX    bool
	EnabledAVX512 bool
}

func init() {
	determineCPUFeatures()

	switch {
	case cpu.HasVAES && cpu.HasAVX512F && cpu.EnabledAVX512:
		checksum = checksum512
	case cpu.HasAES && cpu.HasAVX && cpu.EnabledAVX:
		// AVX required for VEX-encoded AES instruction, which allows non-aligned memory addresses.
		checksum = checksum128
	}
}

// checksum128 implements Meow checksum with AES-NI.
func checksum128(seed uint64, dst, src []byte)

// checksum256 implements Meow checksum with VAES-256.
func checksum256(seed uint64, dst, src []byte)

// checksum512 implements Meow checksum with VAES-512.
func checksum512(seed uint64, dst, src []byte)

// determineCPUFeatures populates flags in global cpu variable by querying CPUID.
func determineCPUFeatures() {
	maxID, _, _, _ := cpuid(0, 0)
	if maxID < 1 {
		return
	}

	_, _, ecx1, _ := cpuid(1, 0)
	cpu.HasOSXSAVE = isSet(ecx1, 27)
	cpu.HasAES = isSet(ecx1, 25)
	cpu.HasAVX = isSet(ecx1, 28)

	if cpu.HasOSXSAVE {
		eax, _ := xgetbv()
		cpu.EnabledAVX = (eax & 0x6) == 0x6
		cpu.EnabledAVX512 = (eax & 0xe) == 0xe
	}

	if maxID < 7 {
		return
	}
	_, ebx7, ecx7, _ := cpuid(7, 0)
	cpu.HasVAES = isSet(ecx7, 9)
	cpu.HasAVX512F = isSet(ebx7, 16)
	cpu.HasAVX512VL = isSet(ebx7, 31)
	cpu.HasAVX512DQ = isSet(ebx7, 17)
}

// cpuid executes the CPUID instruction with the given EAX, ECX inputs.
func cpuid(eaxArg, ecxArg uint32) (eax, ebx, ecx, edx uint32)

// xgetbv executes the XGETBV instruction.
func xgetbv() (eax, edx uint32)

// isSet determines if bit i of x is set.
func isSet(x uint32, i uint) bool {
	return (x>>i)&1 == 1
}

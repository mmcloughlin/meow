// +build ignore

package main

import (
	"crypto/aes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const BlockSize = 256

var output = flag.String("out", "block_amd64.s", "output filename")

func main() {
	flag.Parse()
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m := NewMeow(f)
	if err := m.Generate(); err != nil {
		log.Fatal(err)
	}
}

// Array represents a byte array based at Offset.
type Array struct {
	Base   string
	Offset int
}

// Addr returns a reference to the given byte index into the array.
func (a Array) Addr(idx int) string {
	return fmt.Sprintf("%d(%s)", a.Offset+idx, a.Base)
}

// Slice returns an array based at position idx into the array.
func (a Array) Slice(idx int) Array {
	return Array{Base: a.Base, Offset: a.Offset + idx}
}

// StackFrame represents the stack frame of a function.
type StackFrame struct {
	Size int
}

// Alloc allocates an Array on the stack frame.
func (s *StackFrame) Alloc(size int) Array {
	a := Array{Base: "SP", Offset: s.Size}
	s.Size += size
	return a
}

// Backend encapsulates the instruction set we are targetting.
type Backend interface {
	// Width is the register width in bits.
	Width() int
}

// AESNI implements the AES-NI backend.
type AESNI struct{}

func (a AESNI) Width() int { return 128 }

// VAES256 implements the VAES 256-bit backend.
type VAES256 struct{}

func (v VAES256) Width() int { return 256 }

// VAES512 implements the VAES 512-bit backend.
type VAES512 struct{}

func (v VAES512) Width() int { return 512 }

// Meow writes an assembly implementation of Meow hash components.
type Meow struct {
	w       io.Writer // where to write assembly output
	defines []string  // names of defined macros
	err     error     // saved error from writing
}

// NewMeow builds a new assembly builder writing to w.
func NewMeow(w io.Writer) *Meow {
	return &Meow{
		w: w,
	}
}

// Generate triggers assembly generation.
func (m *Meow) Generate() error {
	m.header()

	m.checksum(AESNI{})
	m.checksum(VAES256{})
	m.checksum(VAES512{})

	return m.err
}

// checksum outputs the entire checksum function.
func (m *Meow) checksum(b Backend) {
	f := StackFrame{}
	iv := f.Alloc(16)
	partial := f.Alloc(BlockSize)
	spill := f.Alloc(64)

	name := fmt.Sprintf("checksum%d", b.Width())
	m.text(name, f.Size, 56)

	m.arg("seed", 0, "R8")
	m.arg("dst_ptr", 8, "DI")
	m.arg("src_ptr", 32, "SI")
	m.arg("src_len", 40, "AX")

	m.section("Prepare IV.")
	m.alloc("IV", "R9")
	m.inst("MOVQ", "SEED, IV")
	m.inst("MOVQ", "IV, %s", iv.Addr(0))
	m.inst("ADDQ", "SRC_LEN, IV")
	m.inst("INCQ", "IV")
	m.inst("MOVQ", "IV, %s", iv.Addr(8))

	m.section("Load IV.")
	for i := 0; i < 16; i++ {
		m.inst("MOVOU", "%s, X%d", iv.Addr(0), i)
	}

	m.blockloop("residual")

	m.label("residual")
	m.inst("CMPQ", "SRC_LEN, $0")
	m.inst("JE", "finish")

	m.section("Duplicate IV.")
	m.inst("MOVQ", "%s, R10", iv.Addr(0))
	m.inst("MOVQ", "%s, R11", iv.Addr(8))
	for i := 0; i < BlockSize; i += 16 {
		m.inst("MOVQ", "R10, %s", partial.Addr(i))
		m.inst("MOVQ", "R11, %s", partial.Addr(i+8))
	}

	m.alloc("BLOCK_PTR", "BX")
	m.inst("LEAQ", "%s, BLOCK_PTR", partial.Addr(0))
	m.label("byteloop")
	m.inst("MOVB", "(SRC_PTR), R10")
	m.inst("MOVB", "R10, (BX)")
	m.inst("INCQ", "SRC_PTR")
	m.inst("INCQ", "BLOCK_PTR")
	m.inst("DECQ", "SRC_LEN")
	m.inst("JNE", "byteloop")

	for i := 0; i < BlockSize; i += aes.BlockSize {
		m.inst("VAESDEC", "%s, X%d, X%d", partial.Addr(i), i/aes.BlockSize, i/aes.BlockSize)
	}

	m.label("finish")

	s := make([]string, 16)
	for i := 0; i < 16; i++ {
		s[i] = fmt.Sprintf("X%d", i)
	}

	for i := 0; i < 4; i++ {
		addr := spill.Addr(16 * i)
		m.inst("MOVOU", "X%d, %s", i, addr)
		s[i] = addr
	}

	for i := 0; i < 4; i++ {
		m.inst("MOVOU", "%s, X%d", iv.Addr(0), i)
	}

	for r := 0; r < 4; r++ {
		m.section(fmt.Sprintf("Rotation block %d.", r))
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				idx := 4*i + (j+r)%4
				m.inst("VAESDEC", "%s, X%d, X%d", s[idx], j, j)
			}
		}
	}

	m.section("Final merge.")
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			m.inst("VAESDEC", "%s, X%d, X%d", iv.Addr(0), j, j)
		}
	}

	m.section("Store hash.")
	for i := 0; i < 4; i++ {
		m.inst("MOVOU", "X%d, %d(DST_PTR)", i, 16*i)
	}

	m.inst("RET", "")
	m.undefall()
}

// blockloop outputs a loop to encrypt entire blocks, exiting to the provided label.
func (m *Meow) blockloop(exit string) {
	m.label("loop")
	m.inst("CMPQ", "SRC_LEN, $%d", BlockSize)
	m.inst("JL", exit)

	m.section("Hash block.")
	for i := 0; i < BlockSize; i += aes.BlockSize {
		m.inst("VAESDEC", "%d(SRC_PTR), X%d, X%d", i, i/aes.BlockSize, i/aes.BlockSize)
	}

	m.section("Update source pointer.")
	m.inst("ADDQ", "$%d, SRC_PTR", BlockSize)
	m.inst("SUBQ", "$%d, SRC_LEN", BlockSize)
	m.inst("JMP", "loop")
}

// header outputs the file header with code generation warning and standard header includes.
func (m *Meow) header() {
	_, self, _, _ := runtime.Caller(0)
	m.printf("// Code generated by go run %s. DO NOT EDIT.\n\n", filepath.Base(self))
	m.printf("// +build !noasm\n\n")
	m.printf("#include \"textflag.h\"\n")
}

// section marks a section of the code with a comment.
func (m *Meow) section(description string) {
	m.printf("\n\t// %s\n", description)
}

// label defines a label.
func (m *Meow) label(name string) {
	m.printf("\n%s:\n", name)
}

// text defines a function header.
func (m *Meow) text(name string, frame, args int) {
	m.printf("\nTEXT %s,0,$%d-%d\n", local(name), frame, args)
}

// alloc informally "allocates" a register with a #define statement.
func (m *Meow) alloc(name, reg string) string {
	macro := strings.ToUpper(name)
	m.define(macro, reg)
	return macro
}

// define a macro.
func (m *Meow) define(name, value string) {
	m.printf("#define %s %s\n", name, value)
	m.defines = append(m.defines, name)
}

// undefall undefs all defined macros.
func (m *Meow) undefall() {
	for _, name := range m.defines {
		m.printf("#undef %s\n", name)
	}
	m.defines = nil
}

// arg reads an argument, and allocates a register for it.
func (m *Meow) arg(name string, offset int, reg string) {
	macro := m.alloc(name, reg)
	m.inst("MOVQ", "%s+%d(FP), %s", name, offset, macro)
}

// inst writes an instruction.
func (m *Meow) inst(name, format string, args ...interface{}) {
	args = append([]interface{}{name}, args...)
	m.printf("\t%-8s "+format+"\n", args...)
}

func (m *Meow) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(m.w, format, args...); err != nil {
		m.err = err
	}
}

// local returns a reference to a local symbol (primarily useful for the unicode dot).
func local(name string) string {
	return fmt.Sprintf("\u00b7%s(SB)", name)
}

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

const (
	LaneSize  = 4 * aes.BlockSize
	BlockSize = 256
)

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

// Generator is an interface for assembly generation.
type Generator interface {
	inst(name, format string, args ...interface{})
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

	// StackAlloc allocates necessary space on the stack for the given number of lanes.
	StackAlloc(*StackFrame)

	// LoadLane loads a 64-byte lane i from the given array.
	LoadLane(g Generator, m Array, i int)

	// StoreLane writes lane i to the given array.
	StoreLane(g Generator, i int, m Array)

	// AESLoad AES decrypts with key loaded from an array.
	AESLoad(g Generator, i int, m Array)

	// AESMerge AES decrypts lane i with key from j.
	AESMerge(g Generator, i, j int)

	// R0 returns a lane ID for the final return lane.
	R0(g Generator) int

	// Rotate lane i.
	Rotate(g Generator, i int)
}

// AESNI implements the AES-NI backend.
type AESNI struct {
	spill  Array
	stream []string
}

func NewAESNI() *AESNI {
	stream := make([]string, 16)
	for i := 0; i < 16; i++ {
		stream[i] = fmt.Sprintf("X%d", i)
	}
	return &AESNI{
		stream: stream,
	}
}

func (a AESNI) Width() int { return 128 }

func (a *AESNI) StackAlloc(f *StackFrame) {
	a.spill = f.Alloc(LaneSize)
}

func (a *AESNI) R0(g Generator) int {
	a.StoreLane(g, 0, a.spill)
	for i := 0; i < 4; i++ {
		a.stream[i] = a.spill.Addr(16 * i)
		a.stream = append(a.stream, fmt.Sprintf("X%d", i))
	}
	return 4
}

func (a AESNI) LoadLane(g Generator, m Array, i int) {
	for j := 0; j < 4; j++ {
		g.inst("MOVOU", "%s, %s", m.Addr(j*aes.BlockSize), a.stream[4*i+j])
	}
}

func (a AESNI) StoreLane(g Generator, i int, m Array) {
	for j := 0; j < 4; j++ {
		g.inst("MOVOU", "%s, %s", a.stream[4*i+j], m.Addr(j*aes.BlockSize))
	}
}

func (a AESNI) AESLoad(g Generator, i int, m Array) {
	for j := 0; j < 4; j++ {
		ref := a.stream[4*i+j]
		g.inst("VAESDEC", "%s, %s, %s", m.Addr(j*aes.BlockSize), ref, ref)
	}
}

func (a AESNI) AESMerge(g Generator, i, j int) {
	for k := 0; k < 4; k++ {
		g.inst("VAESDEC", "%s, %s, %s", a.stream[4*j+k], a.stream[4*i+k], a.stream[4*i+k])
	}
}

func (a AESNI) Rotate(g Generator, i int) {
	l := a.stream[4*i : 4*i+4]
	t := l[0]
	l[0] = l[1]
	l[1] = l[2]
	l[2] = l[3]
	l[3] = t
}

// VAES256 implements the VAES 256-bit backend.
type VAES256 struct{}

func (v VAES256) Width() int               { return 256 }
func (v VAES256) StackAlloc(f *StackFrame) {}

// VAES512 implements the VAES 512-bit backend.
type VAES512 struct{}

func (v VAES512) Width() int               { return 512 }
func (v VAES512) StackAlloc(f *StackFrame) {}

func (v VAES512) LoadLane(g Generator, m Array, i int) {
	g.inst("VMOVDQU64", "%s, Z%d", m.Addr(0), i)
}

func (v VAES512) StoreLane(g Generator, i int, m Array) {
	g.inst("VMOVDQU64", "Z%d, %s", i, m.Addr(0))
}

func (v VAES512) AESLoad(g Generator, i int, m Array) {
	g.inst("VAESDEC", "%s, Z%d, Z%d", m.Addr(0), i, i)
}

func (v VAES512) AESMerge(g Generator, i, j int) {
	g.inst("VAESDEC", "Z%d, Z%d, Z%d", j, i, i)
}

func (v VAES512) R0(g Generator) int { return 4 }

func (v VAES512) Rotate(g Generator, i int) {}

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

	m.checksum(NewAESNI())
	//m.checksum(&VAES256{})
	m.checksum(&VAES512{})

	return m.err
}

// checksum outputs the entire checksum function.
func (m *Meow) checksum(b Backend) {
	f := &StackFrame{}
	iv := f.Alloc(LaneSize)
	partial := f.Alloc(BlockSize)
	b.StackAlloc(f)

	name := fmt.Sprintf("checksum%d", b.Width())
	m.text(name, f.Size, 56)

	m.arg("seed", 0, "R8")
	m.arg("dst_ptr", 8, "DI")
	m.arg("src_ptr", 32, "SI")
	m.arg("src_len", 40, "AX")

	m.section("Prepare IV.")
	m.alloc("IV0", "R9")
	m.alloc("IV1", "R10")
	m.inst("MOVQ", "SEED, IV0")
	m.inst("MOVQ", "SEED, IV1")
	m.inst("ADDQ", "SRC_LEN, IV1")
	m.inst("INCQ", "IV1")
	for i := 0; i < LaneSize; i += 8 {
		m.inst("MOVQ", "IV%d, %s", (i%16)/8, iv.Addr(i))
	}

	m.section("Load IV.")
	for l := 0; l < 4; l++ {
		b.LoadLane(m, iv, l)
	}

	m.blockloop(b, "residual")

	m.label("residual")
	m.inst("CMPQ", "SRC_LEN, $0")
	m.inst("JE", "finish")

	m.section("Duplicate IV.")
	m.inst("MOVQ", "%s, R11", iv.Addr(0))
	m.inst("MOVQ", "%s, R12", iv.Addr(8))
	for i := 0; i < BlockSize; i += 16 {
		m.inst("MOVQ", "R11, %s", partial.Addr(i))
		m.inst("MOVQ", "R12, %s", partial.Addr(i+8))
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

	for l := 0; l < 4; l++ {
		b.AESLoad(m, l, partial.Slice(l*LaneSize))
	}

	m.label("finish")

	r0 := b.R0(m)
	b.LoadLane(m, iv, r0)

	for r := 0; r < 4; r++ {
		m.section(fmt.Sprintf("Rotation block %d.", r))
		for l := 0; l < 4; l++ {
			b.AESMerge(m, r0, l)
			b.Rotate(m, l)
		}
	}

	m.section("Final merge.")
	for i := 0; i < 5; i++ {
		b.AESLoad(m, r0, iv)
	}

	m.section("Store hash.")
	b.StoreLane(m, r0, Array{Base: "DST_PTR"})

	m.inst("RET", "")
	m.undefall()
}

// blockloop outputs a loop to encrypt entire blocks, exiting to the provided label.
func (m *Meow) blockloop(b Backend, exit string) {
	m.label("loop")
	m.inst("CMPQ", "SRC_LEN, $%d", BlockSize)
	m.inst("JL", exit)

	m.section("Hash block.")
	src := Array{Base: "SRC_PTR"}
	for l := 0; l < 4; l++ {
		b.AESLoad(m, l, src.Slice(l*LaneSize))
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

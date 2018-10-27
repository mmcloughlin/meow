// Code generated by go run make_block.go. DO NOT EDIT.

// +build !noasm

#include "textflag.h"

TEXT ·checksum128(SB),0,$32-56
#define SEED R8
	MOVQ     seed+0(FP), SEED
#define DST_PTR DI
	MOVQ     dst_ptr+8(FP), DST_PTR
#define SRC_PTR SI
	MOVQ     src_ptr+32(FP), SRC_PTR
#define SRC_LEN AX
	MOVQ     src_len+40(FP), SRC_LEN

	// Allocate general purpose registers.
#define TOTAL_LEN R9
#define MIX0 R10
#define MIX1 R11
#define PARTIAL_PTR R12
#define TMP R13
#define ZERO R15

	// Prepare a zero register.
	XORQ     ZERO, ZERO

	// Backup total input length.
	MOVQ     SRC_LEN, TOTAL_LEN

	// Prepare Mixer.
	MOVQ     SEED, MIX0
	SUBQ     SRC_LEN, MIX0
	MOVQ     SEED, MIX1
	ADDQ     SRC_LEN, MIX1
	INCQ     MIX1
	MOVQ     MIX0, 0(SP)
	MOVQ     MIX1, 8(SP)

	// Load zero "IV".
	PXOR     X0, X0
	PXOR     X1, X1
	PXOR     X2, X2
	PXOR     X3, X3
	PXOR     X4, X4
	PXOR     X5, X5
	PXOR     X6, X6
	PXOR     X7, X7
	PXOR     X8, X8
	PXOR     X9, X9
	PXOR     X10, X10
	PXOR     X11, X11
	PXOR     X12, X12
	PXOR     X13, X13
	PXOR     X14, X14
	PXOR     X15, X15

	// Handle full 256-byte blocks.

loop:
	CMPQ     SRC_LEN, $256
	JB       sub256

	// Hash block.
	VAESDEC  0(SRC_PTR), X0, X0
	VAESDEC  16(SRC_PTR), X1, X1
	VAESDEC  32(SRC_PTR), X2, X2
	VAESDEC  48(SRC_PTR), X3, X3
	VAESDEC  64(SRC_PTR), X4, X4
	VAESDEC  80(SRC_PTR), X5, X5
	VAESDEC  96(SRC_PTR), X6, X6
	VAESDEC  112(SRC_PTR), X7, X7
	VAESDEC  128(SRC_PTR), X8, X8
	VAESDEC  144(SRC_PTR), X9, X9
	VAESDEC  160(SRC_PTR), X10, X10
	VAESDEC  176(SRC_PTR), X11, X11
	VAESDEC  192(SRC_PTR), X12, X12
	VAESDEC  208(SRC_PTR), X13, X13
	VAESDEC  224(SRC_PTR), X14, X14
	VAESDEC  240(SRC_PTR), X15, X15

	// Update source pointer.
	ADDQ     $256, SRC_PTR
	SUBQ     $256, SRC_LEN
	JMP      loop

	// Handle final sub 256-byte block.

sub256:
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X0, X0
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X1, X1
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X2, X2
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X3, X3
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X4, X4
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X5, X5
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X6, X6
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X7, X7
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X8, X8
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X9, X9
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X10, X10
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X11, X11
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X12, X12
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X13, X13
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN
	CMPQ     SRC_LEN, $16
	JB       sub16
	VAESDEC  0(SRC_PTR), X14, X14
	ADDQ     $16, SRC_PTR
	SUBQ     $16, SRC_LEN

	// Handle final sub 16-byte block.

sub16:
	CMPQ     SRC_LEN, $0
	JE       combine
	MOVQ     ZERO, 16(SP)
	MOVQ     ZERO, 24(SP)
	LEAQ     16(SP), PARTIAL_PTR
	CMPQ     TOTAL_LEN, $16
	JB       byteloop
	LEAQ     -16(SRC_PTR)(SRC_LEN*1), SRC_PTR
	MOVQ     $16, SRC_LEN

byteloop:
	MOVB     (SRC_PTR), TMP
	MOVB     TMP, (PARTIAL_PTR)
	INCQ     SRC_PTR
	INCQ     PARTIAL_PTR
	DECQ     SRC_LEN
	JNE      byteloop
	VAESDEC  16(SP), X15, X15

	// Combine.

combine:
	VAESDEC  X10, X7, X7
	VAESDEC  X4, X7, X7
	VAESDEC  X5, X7, X7
	VAESDEC  X12, X7, X7
	VAESDEC  X8, X7, X7
	VAESDEC  X0, X7, X7
	VAESDEC  X1, X7, X7
	VAESDEC  X9, X7, X7
	VAESDEC  X13, X7, X7
	VAESDEC  X2, X7, X7
	VAESDEC  X6, X7, X7
	VAESDEC  X14, X7, X7
	VAESDEC  X3, X7, X7
	VAESDEC  X11, X7, X7
	VAESDEC  X15, X7, X7

	// Mixing.
	VAESDEC  0(SP), X7, X7
	VAESDEC  0(SP), X7, X7
	VAESDEC  0(SP), X7, X7

	// Store hash.
	MOVOU    X7, 0(DST_PTR)
	RET      
#undef SEED
#undef DST_PTR
#undef SRC_PTR
#undef SRC_LEN
#undef TOTAL_LEN
#undef MIX0
#undef MIX1
#undef PARTIAL_PTR
#undef TMP
#undef ZERO

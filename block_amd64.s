// Code generated by go run make_block.go. DO NOT EDIT.

// +build !noasm

#include "textflag.h"

TEXT ·checksum128(SB),0,$384-56
#define SEED R8
	MOVQ     seed+0(FP), SEED
#define DST_PTR DI
	MOVQ     dst_ptr+8(FP), DST_PTR
#define SRC_PTR SI
	MOVQ     src_ptr+32(FP), SRC_PTR
#define SRC_LEN AX
	MOVQ     src_len+40(FP), SRC_LEN

	// Prepare IV.
#define IV0 R9
#define IV1 R10
	MOVQ     SEED, IV0
	MOVQ     SEED, IV1
	ADDQ     SRC_LEN, IV1
	INCQ     IV1
	MOVQ     IV0, 0(SP)
	MOVQ     IV1, 8(SP)
	MOVQ     IV0, 16(SP)
	MOVQ     IV1, 24(SP)
	MOVQ     IV0, 32(SP)
	MOVQ     IV1, 40(SP)
	MOVQ     IV0, 48(SP)
	MOVQ     IV1, 56(SP)

	// Load IV.
	MOVOU    0(SP), X0
	MOVOU    16(SP), X1
	MOVOU    32(SP), X2
	MOVOU    48(SP), X3
	MOVOU    0(SP), X4
	MOVOU    16(SP), X5
	MOVOU    32(SP), X6
	MOVOU    48(SP), X7
	MOVOU    0(SP), X8
	MOVOU    16(SP), X9
	MOVOU    32(SP), X10
	MOVOU    48(SP), X11
	MOVOU    0(SP), X12
	MOVOU    16(SP), X13
	MOVOU    32(SP), X14
	MOVOU    48(SP), X15

loop:
	CMPQ     SRC_LEN, $256
	JL       residual

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

residual:
	CMPQ     SRC_LEN, $0
	JE       finish

	// Duplicate IV.
	MOVQ     0(SP), R11
	MOVQ     8(SP), R12
	MOVQ     R11, 64(SP)
	MOVQ     R12, 72(SP)
	MOVQ     R11, 80(SP)
	MOVQ     R12, 88(SP)
	MOVQ     R11, 96(SP)
	MOVQ     R12, 104(SP)
	MOVQ     R11, 112(SP)
	MOVQ     R12, 120(SP)
	MOVQ     R11, 128(SP)
	MOVQ     R12, 136(SP)
	MOVQ     R11, 144(SP)
	MOVQ     R12, 152(SP)
	MOVQ     R11, 160(SP)
	MOVQ     R12, 168(SP)
	MOVQ     R11, 176(SP)
	MOVQ     R12, 184(SP)
	MOVQ     R11, 192(SP)
	MOVQ     R12, 200(SP)
	MOVQ     R11, 208(SP)
	MOVQ     R12, 216(SP)
	MOVQ     R11, 224(SP)
	MOVQ     R12, 232(SP)
	MOVQ     R11, 240(SP)
	MOVQ     R12, 248(SP)
	MOVQ     R11, 256(SP)
	MOVQ     R12, 264(SP)
	MOVQ     R11, 272(SP)
	MOVQ     R12, 280(SP)
	MOVQ     R11, 288(SP)
	MOVQ     R12, 296(SP)
	MOVQ     R11, 304(SP)
	MOVQ     R12, 312(SP)
#define BLOCK_PTR BX
	LEAQ     64(SP), BLOCK_PTR

byteloop:
	MOVB     (SRC_PTR), R10
	MOVB     R10, (BX)
	INCQ     SRC_PTR
	INCQ     BLOCK_PTR
	DECQ     SRC_LEN
	JNE      byteloop
	VAESDEC  64(SP), X0, X0
	VAESDEC  80(SP), X1, X1
	VAESDEC  96(SP), X2, X2
	VAESDEC  112(SP), X3, X3
	VAESDEC  128(SP), X4, X4
	VAESDEC  144(SP), X5, X5
	VAESDEC  160(SP), X6, X6
	VAESDEC  176(SP), X7, X7
	VAESDEC  192(SP), X8, X8
	VAESDEC  208(SP), X9, X9
	VAESDEC  224(SP), X10, X10
	VAESDEC  240(SP), X11, X11
	VAESDEC  256(SP), X12, X12
	VAESDEC  272(SP), X13, X13
	VAESDEC  288(SP), X14, X14
	VAESDEC  304(SP), X15, X15

finish:
	MOVOU    X0, 320(SP)
	MOVOU    X1, 336(SP)
	MOVOU    X2, 352(SP)
	MOVOU    X3, 368(SP)
	MOVOU    0(SP), X0
	MOVOU    16(SP), X1
	MOVOU    32(SP), X2
	MOVOU    48(SP), X3

	// Rotation block 0.
	VAESDEC  320(SP), X0, X0
	VAESDEC  336(SP), X1, X1
	VAESDEC  352(SP), X2, X2
	VAESDEC  368(SP), X3, X3
	VAESDEC  X4, X0, X0
	VAESDEC  X5, X1, X1
	VAESDEC  X6, X2, X2
	VAESDEC  X7, X3, X3
	VAESDEC  X8, X0, X0
	VAESDEC  X9, X1, X1
	VAESDEC  X10, X2, X2
	VAESDEC  X11, X3, X3
	VAESDEC  X12, X0, X0
	VAESDEC  X13, X1, X1
	VAESDEC  X14, X2, X2
	VAESDEC  X15, X3, X3

	// Rotation block 1.
	VAESDEC  336(SP), X0, X0
	VAESDEC  352(SP), X1, X1
	VAESDEC  368(SP), X2, X2
	VAESDEC  320(SP), X3, X3
	VAESDEC  X5, X0, X0
	VAESDEC  X6, X1, X1
	VAESDEC  X7, X2, X2
	VAESDEC  X4, X3, X3
	VAESDEC  X9, X0, X0
	VAESDEC  X10, X1, X1
	VAESDEC  X11, X2, X2
	VAESDEC  X8, X3, X3
	VAESDEC  X13, X0, X0
	VAESDEC  X14, X1, X1
	VAESDEC  X15, X2, X2
	VAESDEC  X12, X3, X3

	// Rotation block 2.
	VAESDEC  352(SP), X0, X0
	VAESDEC  368(SP), X1, X1
	VAESDEC  320(SP), X2, X2
	VAESDEC  336(SP), X3, X3
	VAESDEC  X6, X0, X0
	VAESDEC  X7, X1, X1
	VAESDEC  X4, X2, X2
	VAESDEC  X5, X3, X3
	VAESDEC  X10, X0, X0
	VAESDEC  X11, X1, X1
	VAESDEC  X8, X2, X2
	VAESDEC  X9, X3, X3
	VAESDEC  X14, X0, X0
	VAESDEC  X15, X1, X1
	VAESDEC  X12, X2, X2
	VAESDEC  X13, X3, X3

	// Rotation block 3.
	VAESDEC  368(SP), X0, X0
	VAESDEC  320(SP), X1, X1
	VAESDEC  336(SP), X2, X2
	VAESDEC  352(SP), X3, X3
	VAESDEC  X7, X0, X0
	VAESDEC  X4, X1, X1
	VAESDEC  X5, X2, X2
	VAESDEC  X6, X3, X3
	VAESDEC  X11, X0, X0
	VAESDEC  X8, X1, X1
	VAESDEC  X9, X2, X2
	VAESDEC  X10, X3, X3
	VAESDEC  X15, X0, X0
	VAESDEC  X12, X1, X1
	VAESDEC  X13, X2, X2
	VAESDEC  X14, X3, X3

	// Final merge.
	VAESDEC  0(SP), X0, X0
	VAESDEC  16(SP), X1, X1
	VAESDEC  32(SP), X2, X2
	VAESDEC  48(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  16(SP), X1, X1
	VAESDEC  32(SP), X2, X2
	VAESDEC  48(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  16(SP), X1, X1
	VAESDEC  32(SP), X2, X2
	VAESDEC  48(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  16(SP), X1, X1
	VAESDEC  32(SP), X2, X2
	VAESDEC  48(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  16(SP), X1, X1
	VAESDEC  32(SP), X2, X2
	VAESDEC  48(SP), X3, X3

	// Store hash.
	MOVOU    X0, 0(DST_PTR)
	MOVOU    X1, 16(DST_PTR)
	MOVOU    X2, 32(DST_PTR)
	MOVOU    X3, 48(DST_PTR)
	RET      
#undef SEED
#undef DST_PTR
#undef SRC_PTR
#undef SRC_LEN
#undef IV0
#undef IV1
#undef BLOCK_PTR

DATA rotate512<>+0x00(SB)/8, $0x0000000000000002
DATA rotate512<>+0x08(SB)/8, $0x0000000000000003
DATA rotate512<>+0x10(SB)/8, $0x0000000000000004
DATA rotate512<>+0x18(SB)/8, $0x0000000000000005
DATA rotate512<>+0x20(SB)/8, $0x0000000000000006
DATA rotate512<>+0x28(SB)/8, $0x0000000000000007
DATA rotate512<>+0x30(SB)/8, $0x0000000000000000
DATA rotate512<>+0x38(SB)/8, $0x0000000000000001
GLOBL rotate512<>(SB), (NOPTR+RODATA), $64

TEXT ·checksum512(SB),0,$320-56
#define SEED R8
	MOVQ     seed+0(FP), SEED
#define DST_PTR DI
	MOVQ     dst_ptr+8(FP), DST_PTR
#define SRC_PTR SI
	MOVQ     src_ptr+32(FP), SRC_PTR
#define SRC_LEN AX
	MOVQ     src_len+40(FP), SRC_LEN

	// Prepare IV.
#define IV0 R9
#define IV1 R10
	MOVQ     SEED, IV0
	MOVQ     SEED, IV1
	ADDQ     SRC_LEN, IV1
	INCQ     IV1
	MOVQ     IV0, 0(SP)
	MOVQ     IV1, 8(SP)
	MOVQ     IV0, 16(SP)
	MOVQ     IV1, 24(SP)
	MOVQ     IV0, 32(SP)
	MOVQ     IV1, 40(SP)
	MOVQ     IV0, 48(SP)
	MOVQ     IV1, 56(SP)

	// Load IV.
	VMOVDQU64 0(SP), Z0
	VMOVDQU64 0(SP), Z1
	VMOVDQU64 0(SP), Z2
	VMOVDQU64 0(SP), Z3

loop:
	CMPQ     SRC_LEN, $256
	JL       residual

	// Hash block.
	VAESDEC  0(SRC_PTR), Z0, Z0
	VAESDEC  64(SRC_PTR), Z1, Z1
	VAESDEC  128(SRC_PTR), Z2, Z2
	VAESDEC  192(SRC_PTR), Z3, Z3

	// Update source pointer.
	ADDQ     $256, SRC_PTR
	SUBQ     $256, SRC_LEN
	JMP      loop

residual:
	CMPQ     SRC_LEN, $0
	JE       finish

	// Duplicate IV.
	MOVQ     0(SP), R11
	MOVQ     8(SP), R12
	MOVQ     R11, 64(SP)
	MOVQ     R12, 72(SP)
	MOVQ     R11, 80(SP)
	MOVQ     R12, 88(SP)
	MOVQ     R11, 96(SP)
	MOVQ     R12, 104(SP)
	MOVQ     R11, 112(SP)
	MOVQ     R12, 120(SP)
	MOVQ     R11, 128(SP)
	MOVQ     R12, 136(SP)
	MOVQ     R11, 144(SP)
	MOVQ     R12, 152(SP)
	MOVQ     R11, 160(SP)
	MOVQ     R12, 168(SP)
	MOVQ     R11, 176(SP)
	MOVQ     R12, 184(SP)
	MOVQ     R11, 192(SP)
	MOVQ     R12, 200(SP)
	MOVQ     R11, 208(SP)
	MOVQ     R12, 216(SP)
	MOVQ     R11, 224(SP)
	MOVQ     R12, 232(SP)
	MOVQ     R11, 240(SP)
	MOVQ     R12, 248(SP)
	MOVQ     R11, 256(SP)
	MOVQ     R12, 264(SP)
	MOVQ     R11, 272(SP)
	MOVQ     R12, 280(SP)
	MOVQ     R11, 288(SP)
	MOVQ     R12, 296(SP)
	MOVQ     R11, 304(SP)
	MOVQ     R12, 312(SP)
#define BLOCK_PTR BX
	LEAQ     64(SP), BLOCK_PTR

byteloop:
	MOVB     (SRC_PTR), R10
	MOVB     R10, (BX)
	INCQ     SRC_PTR
	INCQ     BLOCK_PTR
	DECQ     SRC_LEN
	JNE      byteloop
	VAESDEC  64(SP), Z0, Z0
	VAESDEC  128(SP), Z1, Z1
	VAESDEC  192(SP), Z2, Z2
	VAESDEC  256(SP), Z3, Z3

finish:
	VMOVDQU64 0(SP), Z4

	// Rotation block 0.
	VAESDEC  Z0, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z0, Z15, Z0
	VAESDEC  Z1, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z1, Z15, Z1
	VAESDEC  Z2, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z2, Z15, Z2
	VAESDEC  Z3, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z3, Z15, Z3

	// Rotation block 1.
	VAESDEC  Z0, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z0, Z15, Z0
	VAESDEC  Z1, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z1, Z15, Z1
	VAESDEC  Z2, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z2, Z15, Z2
	VAESDEC  Z3, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z3, Z15, Z3

	// Rotation block 2.
	VAESDEC  Z0, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z0, Z15, Z0
	VAESDEC  Z1, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z1, Z15, Z1
	VAESDEC  Z2, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z2, Z15, Z2
	VAESDEC  Z3, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z3, Z15, Z3

	// Rotation block 3.
	VAESDEC  Z0, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z0, Z15, Z0
	VAESDEC  Z1, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z1, Z15, Z1
	VAESDEC  Z2, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z2, Z15, Z2
	VAESDEC  Z3, Z4, Z4
	LEAQ     rotate512<>(SB), R15
	VMOVDQU64 (R15), Z15
	VPERMQ   Z3, Z15, Z3

	// Final merge.
	VAESDEC  0(SP), Z4, Z4
	VAESDEC  0(SP), Z4, Z4
	VAESDEC  0(SP), Z4, Z4
	VAESDEC  0(SP), Z4, Z4
	VAESDEC  0(SP), Z4, Z4

	// Store hash.
	VMOVDQU64 Z4, 0(DST_PTR)
	RET      
#undef SEED
#undef DST_PTR
#undef SRC_PTR
#undef SRC_LEN
#undef IV0
#undef IV1
#undef BLOCK_PTR

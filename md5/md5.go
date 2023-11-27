package md5

import (
	"encoding/binary"
	"math"
)

type myMd5 struct{}

func New() *myMd5 {
	return &myMd5{}
}

func (m *myMd5) leftRotate(n, b uint32) uint32 {
	return (n << b) | (n >> (32 - b))
}

func (m *myMd5) Md5Hash(message []byte) []byte {
	// Constants
	T := [64]uint32{}
	for i := 0; i < 64; i++ {
		T[i] = uint32(math.Abs(math.Sin(float64(i+1))) * math.Pow(2, 32))
	}

	s := [4][4]uint32{
		{7, 12, 17, 22},
		{5, 9, 14, 20},
		{4, 11, 16, 23},
		{6, 10, 15, 21},
	}

	// Initialize variables
	A, B, C, D := uint32(0x67452301), uint32(0xEFCDAB89), uint32(0x98BADCFE), uint32(0x10325476)

	// Pad the message
	message = append(message, 0x80)
	for len(message)%64 != 56 {
		message = append(message, 0)
	}

	length := uint64(len(message)) * 8
	lengthBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBytes, length)
	message = append(message, lengthBytes...)

	// Process message in 16-word blocks
	for i := 0; i < len(message); i += 64 {
		X := make([]uint32, 16)
		for j := 0; j < 16; j++ {
			X[j] = binary.LittleEndian.Uint32(message[i+(j*4) : i+((j+1)*4)])
		}

		A_, B_, C_, D_ := A, B, C, D

		// Main loop
		for j := 0; j < 64; j++ {
			var F, F_index uint32
			if j < 16 {
				F = (B & C) | ((^B) & D)
				F_index = uint32(j)
			} else if j < 32 {
				F = (D & B) | ((^D) & C)
				F_index = uint32((5*j + 1) % 16)
			} else if j < 48 {
				F = B ^ C ^ D
				F_index = uint32((3*j + 5) % 16)
			} else {
				F = C ^ (B | (^D))
				F_index = uint32((7 * j) % 16)
			}

			dTemp := D
			D = C
			C = B
			B = B + m.leftRotate((A+F+T[j]+X[F_index])&0xFFFFFFFF, s[j%4][j%4])
			A = dTemp
		}

		// Update state
		A = (A + A_) & 0xFFFFFFFF
		B = (B + B_) & 0xFFFFFFFF
		C = (C + C_) & 0xFFFFFFFF
		D = (D + D_) & 0xFFFFFFFF
	}

	// Output
	result := make([]byte, 16)
	binary.LittleEndian.PutUint32(result, A)
	binary.LittleEndian.PutUint32(result[4:], B)
	binary.LittleEndian.PutUint32(result[8:], C)
	binary.LittleEndian.PutUint32(result[12:], D)

	return result
}

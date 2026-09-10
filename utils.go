package br

import (
	"math/bits"
	"math/rand/v2"
)

func isSpace(b byte) bool {
	return b == ' '
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphaUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func isAlphaNumericalUpper(b byte) bool {
	return isDigit(b) || isAlphaUpper(b)
}

func asciiLowerToUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		b -= 'a' - 'A'
	}
	return b
}

// load64 returns the first 8 bytes of s as a little-endian word. The shift
// pattern is recognized by the compiler and becomes a single wide load on
// little-endian platforms.
func load64(s string) uint64 {
	_ = s[7]
	return uint64(s[0]) | uint64(s[1])<<8 | uint64(s[2])<<16 | uint64(s[3])<<24 |
		uint64(s[4])<<32 | uint64(s[5])<<40 | uint64(s[6])<<48 | uint64(s[7])<<56
}

// load32 returns the first 4 bytes of s as a little-endian word.
func load32(s string) uint32 {
	_ = s[3]
	return uint32(s[0]) | uint32(s[1])<<8 | uint32(s[2])<<16 | uint32(s[3])<<24
}

// nonDigits64 returns 0 if and only if all 8 bytes of v are ASCII digits.
//
// A byte is a digit if its high nibble is 3 and adding 6 to it does not carry
// into the high nibble (that is, its low nibble is at most 9).
func nonDigits64(v uint64) uint64 {
	return ((v & 0xF0F0F0F0F0F0F0F0) |
		(((v + 0x0606060606060606) & 0xF0F0F0F0F0F0F0F0) >> 4)) ^ 0x3333333333333333
}

// nonDigits32 is nonDigits64 for 4 byte words.
func nonDigits32(v uint32) uint32 {
	return ((v & 0xF0F0F0F0) | (((v + 0x06060606) & 0xF0F0F0F0) >> 4)) ^ 0x33333333
}

var pcg = rand.NewPCG(rand.Uint64(), rand.Uint64())

var cnsFirstDigits = []byte{'1', '2', '7', '8', '9'}

func randomCNSFirstDigit() byte {
	n := uint64(len(cnsFirstDigits))

	// This code here is taken from the stdlib.
	// You can check it at the math/rand/v2 package under func '(r *Rand) uint64n(n uint64) uint64'.
	hi, lo := bits.Mul64(pcg.Uint64(), n)
	if lo < n {
		thresh := -n % n
		for lo < thresh {
			hi, lo = bits.Mul64(pcg.Uint64(), n)
		}
	}

	return cnsFirstDigits[int(hi)]
}

func randomDigit() byte {
	var n uint64 = ('9' + 1) - '0'

	// This code here is taken from the stdlib.
	// You can check it at the math/rand/v2 package under func '(r *Rand) uint64n(n uint64) uint64'.
	hi, lo := bits.Mul64(pcg.Uint64(), n)
	if lo < n {
		thresh := -n % n
		for lo < thresh {
			hi, lo = bits.Mul64(pcg.Uint64(), n)
		}
	}

	return byte(hi) + '0'
}

func randomAlphaUpper() byte {
	var n uint64 = ('Z' + 1) - 'A'

	// This code here is taken from the stdlib.
	// You can check it at the math/rand/v2 package under func '(r *Rand) uint64n(n uint64) uint64'.
	hi, lo := bits.Mul64(pcg.Uint64(), n)
	if lo < n {
		thresh := -n % n
		for lo < thresh {
			hi, lo = bits.Mul64(pcg.Uint64(), n)
		}
	}

	return byte(hi) + 'A'
}

var alphaNumericals = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

func randomAlphaNumericalUpper() byte {
	n := uint64(len(alphaNumericals))

	// This code here is taken from the stdlib.
	// You can check it at the math/rand/v2 package under func '(r *Rand) uint64n(n uint64) uint64'.
	hi, lo := bits.Mul64(pcg.Uint64(), n)
	if lo < n {
		thresh := -n % n
		for lo < thresh {
			hi, lo = bits.Mul64(pcg.Uint64(), n)
		}
	}

	return alphaNumericals[int(hi)]
}

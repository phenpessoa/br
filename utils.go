package br

import (
	"math/bits"
	"math/rand/v2"
)

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

var pcg = rand.NewPCG(rand.Uint64(), rand.Uint64())

func randomDigit() byte {
	var n uint64 = '9' - '0'

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
	var n uint64 = 'Z' - 'A'

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

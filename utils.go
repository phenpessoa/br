package br

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

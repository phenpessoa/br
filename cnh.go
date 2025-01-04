package br

import (
	"errors"
)

// CNH represents a Brazilian driver's license number.
type CNH string

// NewCNH creates a new CNH instance from a string representation.
//
// It verifies the CNH's validity using checksum digits.
func NewCNH(s string) (CNH, error) {
	cnh := CNH(s)
	if !cnh.IsValid() {
		return "", ErrInvalidCNH
	}
	return cnh, nil
}

// ErrInvalidCNH is an error returned when an invalid CNH is encountered.
var ErrInvalidCNH = errors.New("br: invalid cnh")

var cnhFirstTable = []int{2, 3, 4, 5, 6, 7, 8, 9, 10}

// IsValid checks whether the provided CNH is valid based on its checksum digits.
func (cnh CNH) IsValid() bool {
	switch len(cnh) {
	case 11:
		dByte, cacheSum, ok := cnhIterFirst(cnh)
		if !ok {
			return false
		}

		if cnh[len(cnh)-2] != dByte {
			return false
		}

		dByte, ok = cnhIterSecond(cnh, cacheSum)
		if !ok {
			return false
		}

		return cnh[len(cnh)-1] == dByte
	default:
		return false
	}
}

func cnhIterFirst[T string | CNH | []byte](cnh T) (byte, int, bool) {
	if len(cnh) != 11 || len(cnhFirstTable) != 9 {
		panic("not 11 or 9 - cnh")
	}

	var sum, cacheSum, rest, out int

	for i, d := range cnhFirstTable {
		cur := cnh[i]
		if !isDigit(cur) {
			return 0, 0, false
		}
		parsed := int(cur - '0')
		sum += d * parsed
		cacheSum += parsed
	}

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', cacheSum + sum, true
}

func cnhIterSecond[T string | CNH | []byte](cnh T, sum int) (byte, bool) {
	if len(cnh) != 11 {
		panic("not 11 - cnh")
	}

	var rest, out int

	last := cnh[9]
	if !isDigit(last) {
		return 0, false
	}
	sum += 2 * int(last-'0')

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', true
}

// String returns the string representation of CNH.
func (cnh CNH) String() string {
	if !cnh.IsValid() {
		return ""
	}
	return string(cnh)
}

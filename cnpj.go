package br

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/phenpessoa/gutils/unsafex"
)

var (
	// ErrInvalidCNPJ is an error returned when an invalid CNPJ is encountered.
	ErrInvalidCNPJ = errors.New("br: invalid cnpj passed")

	cnpjFirstTable  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondTable = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

// NewCNPJ creates a new CNPJ instance from a string representation.
//
// It verifies the CNPJ's validity using checksum digits.
func NewCNPJ(s string) (CNPJ, error) {
	cnpj := CNPJ(s)
	if !cnpj.IsValid() {
		return "", ErrInvalidCNPJ
	}
	return cnpj, nil
}

// CNPJ represents a Brazilian CNPJ.
type CNPJ string

// IsValid checks whether the provided CNPJ is valid based on its checksum
// digits.
func (cnpj CNPJ) IsValid() bool {
	l := len(cnpj)
	if l != 14 && l != 18 {
		return false
	}

	dByte, ok := iterCNPJTable(cnpj, cnpjFirstTable)
	if !ok {
		return false
	}
	if cnpj[l-2] != dByte {
		return false
	}

	dByte, ok = iterCNPJTable(cnpj, cnpjSecondTable)
	if !ok {
		return false
	}

	return cnpj[l-1] == dByte
}

func iterCNPJTable(cnpj CNPJ, table []int) (byte, bool) {
	var sum, pad, rest, d int

	for i, d := range table {
		cur := cnpj[i+pad]
		switch cur {
		case '.':
			if i != 2 && i != 5 {
				return 0, false
			}
			pad++
			cur = cnpj[i+pad]
		case '/':
			if i != 8 {
				return 0, false
			}
			pad++
			cur = cnpj[i+pad]
		case '-':
			if i != 12 {
				return 0, false
			}
			pad++
			cur = cnpj[i+pad]
		}

		cur = asciiLowerToUpper(cur)

		if cur < '0' || (cur > '9' && cur < 'A') || cur > 'Z' {
			return 0, false
		}

		sum += d * int(cur-'0')
	}

	rest = sum % 11

	if rest >= 2 {
		d = 11 - rest
	}

	return byte(d) + '0', true
}

// String returns the formatted CNPJ string with punctuation as
// XX.XXX.XXX/XXXX-XX.
func (cnpj CNPJ) String() string {
	if !cnpj.IsValid() {
		return ""
	}

	if len(cnpj) == 18 {
		return strings.ToUpper(string(cnpj))
	}

	out := make([]byte, 18)

	for i := range cnpj {
		cur := asciiLowerToUpper(cnpj[i])

		switch {
		case i < 2:
			out[i] = cur
		case i == 2:
			out[i] = '.'
			out[i+1] = cur
		case i < 5:
			out[i+1] = cur
		case i == 5:
			out[i+1] = '.'
			out[i+2] = cur
		case i < 8:
			out[i+2] = cur
		case i == 8:
			out[i+2] = '/'
			out[i+3] = cur
		case i < 12:
			out[i+3] = cur
		case i == 12:
			out[i+3] = '-'
			out[i+4] = cur
		default:
			out[i+4] = cur
		}
	}

	return unsafex.String(out)
}

func asciiLowerToUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		b -= 'a' - 'A'
	}
	return b
}

func (cnpj CNPJ) Value() (driver.Value, error) {
	return cnpj.String(), nil
}

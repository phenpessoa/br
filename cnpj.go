package br

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/phenpessoa/gutils/unsafex"
)

// CNPJ represents a Brazilian CNPJ.
type CNPJ string

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

// GenerateCNPJ generates a pseudo-random valid CNPJ.
func GenerateCNPJ() CNPJ {
	data := make([]byte, 18)
	data[2] = '.'
	data[6] = '.'
	data[10] = '/'
	data[15] = '-'

	for i := range 2 {
		data[i] = randomAlphaNumericalUpper()
	}

	for i := 3; i < 6; i++ {
		data[i] = randomAlphaNumericalUpper()
	}

	for i := 7; i < 10; i++ {
		data[i] = randomAlphaNumericalUpper()
	}

	for i := 11; i < 15; i++ {
		data[i] = randomAlphaNumericalUpper()
	}

	data[16], _ = cnpjIterFirst18(data)
	data[17], _ = cnpjIterSecond18(data)

	return CNPJ(string(data))
}

// ErrInvalidCNPJ is an error returned when an invalid CNPJ is encountered.
var ErrInvalidCNPJ = errors.New("br: invalid cnpj")

// cnpj tables
var (
	cnpjFirstTable  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondTable = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

// IsValid checks whether the provided CNPJ is valid based on its checksum digits.
func (cnpj CNPJ) IsValid() bool {
	switch len(cnpj) {
	case 14:
		dByte, ok := cnpjIterFirst14(cnpj)
		if !ok {
			return false
		}

		if cnpj[len(cnpj)-2] != dByte {
			return false
		}

		dByte, ok = cnpjIterSecond14(cnpj)
		if !ok {
			return false
		}

		return cnpj[len(cnpj)-1] == dByte
	case 18:
		if cnpj[2] != '.' || cnpj[6] != '.' || cnpj[10] != '/' || cnpj[15] != '-' {
			return false
		}

		dByte, ok := cnpjIterFirst18(cnpj)
		if !ok {
			return false
		}

		if cnpj[len(cnpj)-2] != dByte {
			return false
		}

		dByte, ok = cnpjIterSecond18(cnpj)
		if !ok {
			return false
		}

		return cnpj[len(cnpj)-1] == dByte
	default:
		return false
	}
}

func cnpjIterFirst18[T string | CNPJ | []byte](cnpj T) (byte, bool) {
	if len(cnpj) != 18 || len(cnpjFirstTable) != 12 {
		panic("not 18 or 12")
	}

	var sum, rest, out int
	for i, d := range cnpjFirstTable[:2] {
		cur := asciiLowerToUpper(cnpj[i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cnpjFirstTable[2:5] {
		cur := asciiLowerToUpper(cnpj[3:6][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cnpjFirstTable[5:8] {
		cur := asciiLowerToUpper(cnpj[7:10][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cnpjFirstTable[8:12] {
		cur := asciiLowerToUpper(cnpj[11:15][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', true
}

func cnpjIterSecond18[T string | CNPJ | []byte](cnpj T) (byte, bool) {
	if len(cnpj) != 18 || len(cnpjSecondTable) != 13 {
		panic("not 18 or 12")
	}

	var sum, rest, out int
	for i, d := range cnpjSecondTable[:2] {
		cur := asciiLowerToUpper(cnpj[i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cnpjSecondTable[2:5] {
		cur := asciiLowerToUpper(cnpj[3:6][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cnpjSecondTable[5:8] {
		cur := asciiLowerToUpper(cnpj[7:10][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cnpjSecondTable[8:12] {
		cur := asciiLowerToUpper(cnpj[11:15][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	last := cnpj[16]
	if !isDigit(last) {
		return 0, false
	}
	sum += cnpjSecondTable[12] * int(last-'0')

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', true
}

func cnpjIterFirst14[T string | CNPJ | []byte](cnpj T) (byte, bool) {
	if len(cnpj) != 14 || len(cnpjFirstTable) != 12 {
		panic("not 14 or 12")
	}

	var sum, rest, out int

	for i, d := range cnpjFirstTable {
		cur := asciiLowerToUpper(cnpj[i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', true
}

func cnpjIterSecond14[T string | CNPJ | []byte](cnpj T) (byte, bool) {
	if len(cnpj) != 14 || len(cnpjSecondTable) != 13 {
		panic("not 14 or 12")
	}

	var sum, rest, out int

	for i, d := range cnpjSecondTable[:12] {
		cur := asciiLowerToUpper(cnpj[i])
		if !isAlphaNumericalUpper(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	last := cnpj[12]
	if !isDigit(last) {
		return 0, false
	}
	sum += cnpjSecondTable[12] * int(last-'0')

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', true
}

// String returns the formatted CNPJ string with punctuation as XX.XXX.XXX/XXXX-XX.
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

// Value implements the driver.Valuer interface for CNPJ.
func (cnpj CNPJ) Value() (driver.Value, error) {
	return cnpj.String(), nil
}

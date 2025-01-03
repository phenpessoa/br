package br

import (
	"database/sql/driver"
	"errors"
	"strings"
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

	data[4] = '0'

	var cacheSum int
	data[16], cacheSum, _ = cnpjIterFirst18(data)
	data[17], _ = cnpjIterSecond18(data, cacheSum)

	return CNPJ(string(data))
}

// ErrInvalidCNPJ is an error returned when an invalid CNPJ is encountered.
var ErrInvalidCNPJ = errors.New("br: invalid cnpj")

var cnpjFirstTable = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

// IsValid checks whether the provided CNPJ is valid based on its checksum digits.
func (cnpj CNPJ) IsValid() bool {
	switch len(cnpj) {
	case 14:
		dByte, cacheSum, ok := cnpjIterFirst14(cnpj)
		if !ok {
			return false
		}

		if cnpj[len(cnpj)-2] != dByte {
			return false
		}

		dByte, ok = cnpjIterSecond14(cnpj, cacheSum)
		if !ok {
			return false
		}

		return cnpj[len(cnpj)-1] == dByte
	case 18:
		if cnpj[2] != '.' || cnpj[6] != '.' || cnpj[10] != '/' || cnpj[15] != '-' {
			return false
		}

		dByte, cacheSum, ok := cnpjIterFirst18(cnpj)
		if !ok {
			return false
		}

		if cnpj[len(cnpj)-2] != dByte {
			return false
		}

		dByte, ok = cnpjIterSecond18(cnpj, cacheSum)
		if !ok {
			return false
		}

		return cnpj[len(cnpj)-1] == dByte
	default:
		return false
	}
}

func cnpjIterFirst18[T string | CNPJ | []byte](cnpj T) (byte, int, bool) {
	if len(cnpj) != 18 || len(cnpjFirstTable) != 12 {
		panic("not 18 or 12")
	}

	var sum, cacheSum, rest, out int

	for i, d := range cnpjFirstTable[:2] {
		cur := asciiLowerToUpper(cnpj[i])
		if !isAlphaNumericalUpper(cur) {
			return 0, 0, false
		}
		parsed := int(cur - '0')
		sum += d * parsed
		cacheSum += parsed
	}

	for i, d := range cnpjFirstTable[2:5] {
		cur := asciiLowerToUpper(cnpj[3:6][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, 0, false
		}
		parsed := int(cur - '0')
		sum += d * parsed

		if i == 2 {
			// The delta of all digits in the 2 CNPJ tables is 1,
			// the only exception is on index 4.
			// That's why this is the only index that needs special treatment.
			cacheSum += parsed * (2 - 9) // cnpjSecondTable[4] - cnpjFirstTable[4]
		} else {
			cacheSum += parsed
		}
	}

	for i, d := range cnpjFirstTable[5:8] {
		cur := asciiLowerToUpper(cnpj[7:10][i])
		if !isAlphaNumericalUpper(cur) {
			return 0, 0, false
		}
		parsed := int(cur - '0')
		sum += d * parsed
		cacheSum += parsed
	}

	for i, d := range cnpjFirstTable[8:12] {
		cur := asciiLowerToUpper(cnpj[11:15][i])
		if !isAlphaNumericalUpper(cur) {
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

func cnpjIterSecond18[T string | CNPJ | []byte](cnpj T, sum int) (byte, bool) {
	if len(cnpj) != 18 {
		panic("not 18 or 12")
	}

	var rest, out int

	last := cnpj[16]
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

func cnpjIterFirst14[T string | CNPJ | []byte](cnpj T) (byte, int, bool) {
	if len(cnpj) != 14 || len(cnpjFirstTable) != 12 {
		panic("not 14 or 12")
	}

	var sum, cacheSum, rest, out int

	for i, d := range cnpjFirstTable {
		cur := asciiLowerToUpper(cnpj[i])
		if !isAlphaNumericalUpper(cur) {
			return 0, 0, false
		}
		parsed := int(cur - '0')
		sum += d * parsed
		if i == 4 {
			// The delta of all digits in the two CNPJ tables is 1,
			// the only exception is on index 4.
			// That's why this is the only index that needs special treatment.
			cacheSum += parsed * (2 - 9) // cnpjSecondTable[4] - cnpjFirstTable[4]
		} else {
			cacheSum += parsed
		}
	}

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', cacheSum + sum, true
}

func cnpjIterSecond14[T string | CNPJ | []byte](cnpj T, sum int) (byte, bool) {
	if len(cnpj) != 14 {
		panic("not 14 or 12")
	}

	var rest, out int

	last := cnpj[12]
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

// String returns the formatted CNPJ string with punctuation as XX.XXX.XXX/XXXX-XX.
func (cnpj CNPJ) String() string {
	if !cnpj.IsValid() {
		return ""
	}

	if len(cnpj) == 18 {
		return strings.ToUpper(string(cnpj))
	}

	if len(cnpj) != 14 {
		return ""
	}

	out := make([]byte, 18)
	out[2] = '.'
	out[6] = '.'
	out[10] = '/'
	out[15] = '-'

	copy(out[0:2], cnpj[0:2])
	copy(out[3:6], cnpj[2:5])
	copy(out[7:10], cnpj[5:8])
	copy(out[11:15], cnpj[8:12])
	copy(out[16:18], cnpj[12:14])

	for i, d := range out {
		out[i] = asciiLowerToUpper(d)
	}

	return string(out)
}

// Value implements the driver.Valuer interface for CNPJ.
func (cnpj CNPJ) Value() (driver.Value, error) {
	return cnpj.String(), nil
}

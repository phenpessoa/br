package br

import (
	"database/sql/driver"
	"errors"
)

// CPF represents a Brazilian CPF.
type CPF string

// NewCPF creates a new CPF instance from a string representation.
//
// It verifies the CPF's validity using checksum digits.
func NewCPF(s string) (CPF, error) {
	cpf := CPF(s)
	if !cpf.IsValid() {
		return "", ErrInvalidCPF
	}
	return cpf, nil
}

func GenerateCPF() CPF {
	data := make([]byte, 14)
	data[3] = '.'
	data[7] = '.'
	data[11] = '-'

	for i := range 3 {
		data[i] = randomDigit()
	}

	for i := 4; i < 7; i++ {
		data[i] = randomDigit()
	}

	for i := 8; i < 11; i++ {
		data[i] = randomDigit()
	}

	data[12], _ = cpfIterFirst14(data)
	data[13], _ = cpfIterSecond14(data)

	return CPF(string(data))
}

// ErrInvalidCPF is an error returned when an invalid CPF is encountered.
var ErrInvalidCPF = errors.New("br: invalid cpf")

// cpf tables
var (
	cpfFirstTable  = []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	cpfSecondTable = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
)

// IsValid checks whether the provided CPF is valid based on its checksum
// digits.
func (cpf CPF) IsValid() bool {
	switch len(cpf) {
	case 11:
		dByte, ok := cpfIterFirst11(cpf)
		if !ok {
			return false
		}

		if cpf[len(cpf)-2] != dByte {
			return false
		}

		dByte, ok = cpfIterSecond11(cpf)
		if !ok {
			return false
		}

		return cpf[len(cpf)-1] == dByte
	case 14:
		if cpf[3] != '.' || cpf[7] != '.' || cpf[11] != '-' {
			return false
		}

		dByte, ok := cpfIterFirst14(cpf)
		if !ok {
			return false
		}

		if cpf[len(cpf)-2] != dByte {
			return false
		}

		dByte, ok = cpfIterSecond14(cpf)
		if !ok {
			return false
		}

		return cpf[len(cpf)-1] == dByte
	default:
		return false
	}
}

func cpfIterFirst14[T string | CPF | []byte](cpf T) (byte, bool) {
	if len(cpf) != 14 || len(cpfFirstTable) != 9 {
		panic("not 14 or 9")
	}

	var sum, rest, out int

	for i, d := range cpfFirstTable[:3] {
		cur := cpf[i]
		if !isDigit(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cpfFirstTable[3:6] {
		cur := cpf[4:7][i]
		if !isDigit(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cpfFirstTable[6:9] {
		cur := cpf[8:11][i]
		if !isDigit(cur) {
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

func cpfIterSecond14[T string | CPF | []byte](cpf T) (byte, bool) {
	if len(cpf) != 14 || len(cpfSecondTable) != 10 {
		panic("not 14 or 10")
	}

	var sum, rest, out int

	for i, d := range cpfSecondTable[:3] {
		cur := cpf[i]
		if !isDigit(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cpfSecondTable[3:6] {
		cur := cpf[4:7][i]
		if !isDigit(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	for i, d := range cpfSecondTable[6:9] {
		cur := cpf[8:11][i]
		if !isDigit(cur) {
			return 0, false
		}
		sum += d * int(cur-'0')
	}

	last := cpf[12]
	if !isDigit(last) {
		return 0, false
	}
	sum += cpfSecondTable[9] * int(last-'0')

	rest = sum % 11

	if rest >= 2 {
		out = 11 - rest
	}

	return byte(out) + '0', true
}

func cpfIterFirst11[T string | CPF | []byte](cpf T) (byte, bool) {
	if len(cpf) != 11 || len(cpfFirstTable) != 9 {
		panic("not 11 or 9")
	}

	var sum, rest, out int

	for i, d := range cpfFirstTable {
		cur := cpf[i]
		if !isDigit(cur) {
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

func cpfIterSecond11[T string | CPF | []byte](cpf T) (byte, bool) {
	if len(cpf) != 11 || len(cpfSecondTable) != 10 {
		panic("not 11 or 10")
	}

	var sum, rest, out int

	for i, d := range cpfSecondTable {
		cur := cpf[i]
		if !isDigit(cur) {
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

// String returns the formatted CPF string with punctuation as XXX.XXX.XXX-XX.
func (cpf CPF) String() string {
	if !cpf.IsValid() {
		return ""
	}

	if len(cpf) == 14 {
		return string(cpf)
	}

	if len(cpf) != 11 {
		panic("not 11")
	}

	out := make([]byte, 14)

	out[3] = '.'
	out[7] = '.'
	out[11] = '-'

	copy(out[0:3], cpf[0:3])
	copy(out[4:7], cpf[3:6])
	copy(out[8:11], cpf[6:9])
	copy(out[12:14], cpf[9:11])

	return string(out)
}

func (c CPF) Value() (driver.Value, error) {
	return c.String(), nil
}

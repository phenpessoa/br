package br

import (
	"database/sql/driver"
	"errors"

	"github.com/phenpessoa/gutils/unsafex"
)

var (
	// ErrInvalidCPF is an error returned when an invalid CPF is encountered.
	ErrInvalidCPF = errors.New("br: invalid cpf passed")

	cpfFirstTable  = []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	cpfSecondTable = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
)

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

// CPF represents a Brazilian CPF.
type CPF string

// IsValid checks whether the provided CPF is valid based on its checksum
// digits.
func (cpf CPF) IsValid() bool {
	l := len(cpf)
	if l != 11 && l != 14 {
		return false
	}

	dByte, ok := iterCPFTable(cpf, cpfFirstTable)
	if !ok {
		return false
	}

	if cpf[l-2] != dByte {
		return false
	}

	dByte, ok = iterCPFTable(cpf, cpfSecondTable)
	if !ok {
		return false
	}

	return cpf[l-1] == dByte
}

func iterCPFTable(cpf CPF, table []int) (byte, bool) {
	var sum, pad, rest, d int

	for i, d := range table {
		cur := cpf[i+pad]
		switch cur {
		case '.':
			if i != 3 && i != 6 {
				return 0, false
			}
			pad++
			cur = cpf[i+pad]
		case '-':
			if i != 9 {
				return 0, false
			}
			pad++
			cur = cpf[i+pad]
		}

		if cur < '0' || cur > '9' {
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

// String returns the formatted CPF string with punctuation as XXX.XXX.XXX-XX.
func (cpf CPF) String() string {
	if !cpf.IsValid() {
		return ""
	}

	if len(cpf) == 14 {
		return string(cpf)
	}

	out := make([]byte, 14)
	for i := range cpf {
		switch {
		case i < 3:
			out[i] = cpf[i]
		case i == 3:
			out[i] = '.'
			out[i+1] = cpf[i]
		case i < 6:
			out[i+1] = cpf[i]
		case i == 6:
			out[i+1] = '.'
			out[i+2] = cpf[i]
		case i < 9:
			out[i+2] = cpf[i]
		case i == 9:
			out[i+2] = '-'
			out[i+3] = cpf[i]
		default:
			out[i+3] = cpf[i]
		}
	}

	return unsafex.String(out)
}

func (c CPF) Value() (driver.Value, error) {
	return c.String(), nil
}

package br

import (
	"bytes"
	"errors"
	"strings"

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
// It removes any punctuation characters and verifies the CPF's validity using
// checksum digits.
func NewCPF(s string) (CPF, error) {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "-", "")

	cpf := CPF(s)
	if !cpf.IsValid() {
		return "", ErrInvalidCPF
	}

	return cpf, nil
}

// CPF represents a Brazilian individual CPF.
type CPF string

// IsValid checks whether the provided CPF is valid based on its checksum
// digits.
func (cpf CPF) IsValid() bool {
	if len(cpf) != 11 {
		return false
	}

	buf := make([]byte, 11)
	for i := range cpf {
		buf[i] = cpf[i]
	}

	var sum1 int
	for i, d := range cpfFirstTable {
		c := int(buf[i] - '0')
		if c < 0 || c > 9 {
			return false
		}
		sum1 += d * c
	}
	rest1 := sum1 % 11
	d1 := 0

	if rest1 >= 2 {
		d1 = 11 - rest1
	}

	if d1 == 0 {
		buf[9] = '0'
	} else {
		buf[9] = byte(d1) + '0'
	}

	var sum2 int
	for i, d := range cpfSecondTable {
		c := int(buf[i] - '0')
		if c < 0 || c > 9 {
			return false
		}
		sum2 += d * c
	}
	rest2 := sum2 % 11
	d2 := 0

	if rest2 >= 2 {
		d2 = 11 - rest2
	}

	if d2 == 0 {
		buf[10] = '0'
	} else {
		buf[10] = byte(d2) + '0'
	}

	return bytes.Equal(buf, unsafex.ByteSlice(string(cpf)))
}

// String returns the formatted CPF string with punctuation as XXX.XXX.XXX-XX.
func (cpf CPF) String() string {
	if !cpf.IsValid() {
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

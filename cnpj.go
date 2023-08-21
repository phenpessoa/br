package br

import (
	"bytes"
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
// It removes any punctuation characters and verifies the CNPJ's validity using
// checksum digits.
func NewCNPJ(s string) (CNPJ, error) {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "/", "")
	s = strings.ReplaceAll(s, "-", "")

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
	s := string(cnpj)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "/", "")
	s = strings.ReplaceAll(s, "-", "")

	if len(s) != 14 {
		return false
	}

	buf := make([]byte, 14)
	for i := range s {
		buf[i] = s[i]
	}

	var sum1 int
	for i, d := range cnpjFirstTable {
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
		buf[12] = '0'
	} else {
		buf[12] = byte(d1) + '0'
	}

	var sum2 int
	for i, d := range cnpjSecondTable {
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
		buf[13] = '0'
	} else {
		buf[13] = byte(d2) + '0'
	}

	return bytes.Equal(buf, unsafex.ByteSlice(s))
}

// String returns the formatted CNPJ string with punctuation as
// XX.XXX.XXX/XXXX-XX.
func (cnpj CNPJ) String() string {
	if !cnpj.IsValid() {
		return string(cnpj)
	}

	out := make([]byte, 18)

	for i := range cnpj {
		switch {
		case i < 2:
			out[i] = cnpj[i]
		case i == 2:
			out[i] = '.'
			out[i+1] = cnpj[i]
		case i < 5:
			out[i+1] = cnpj[i]
		case i == 5:
			out[i+1] = '.'
			out[i+2] = cnpj[i]
		case i < 8:
			out[i+2] = cnpj[i]
		case i == 8:
			out[i+2] = '/'
			out[i+3] = cnpj[i]
		case i < 12:
			out[i+3] = cnpj[i]
		case i == 12:
			out[i+3] = '-'
			out[i+4] = cnpj[i]
		default:
			out[i+4] = cnpj[i]
		}
	}

	return unsafex.String(out)
}

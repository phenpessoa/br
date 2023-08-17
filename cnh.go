package br

import (
	"bytes"
	"errors"

	"github.com/phenpessoa/gutils/unsafex"
)

var (
	// ErrInvalidCNH is an error returned when an invalid CNH is encountered.
	ErrInvalidCNH = errors.New("br: invalid cnh passed")

	cnhFirstTable  = []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	cnhSecondTable = []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 2}
)

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

// CNH represents a Brazilian driver's license number.
type CNH string

// IsValid checks whether the provided CNH is valid based on its checksum digits.
func (cnh CNH) IsValid() bool {
	if len(cnh) != 11 {
		return false
	}

	buf := make([]byte, 11)
	for i := range cnh {
		buf[i] = cnh[i]
	}

	var sum1 int
	for i, d := range cnhFirstTable {
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
	for i, d := range cnhSecondTable {
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

	return bytes.Equal(buf, unsafex.ByteSlice(string(cnh)))
}

// String returns the string representation of CNH.
func (cnh CNH) String() string {
	return string(cnh)
}

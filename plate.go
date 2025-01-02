package br

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/phenpessoa/gutils/unsafex"
)

// ErrInvalidPlate is an error returned when an invalid license plate is
// encountered.
var ErrInvalidPlate = errors.New("br: invalid license plate")

// NewPlate creates a new Plate instance from a string representation.
func NewPlate(s string) (Plate, error) {
	plate := Plate(s)
	if !plate.IsValid() {
		return "", ErrInvalidPlate
	}
	return plate, nil
}

// Plate represents a Brazilian vehicle license plate.
type Plate string

// IsValid checks whether the provided license plate is valid based on specific
// formatting rules.
//
// IsValid will return true if the plate if either a MercoSul or a Brazilian
// type plate.
//
// The formats accepted are: XXXXXXX, XXX-XXXX, XXX.XXXX
func (p Plate) IsValid() bool {
	l := len(p)
	if l != 7 && l != 8 {
		return false
	}

	if l == 8 && p[3] != '.' && p[3] != '-' {
		return false
	}

	var pad int
	for i := 0; i < 7; i++ {
		switch cur := p[i+pad]; {
		case i < 3:
			if !isAlphaUpper(cur) {
				return false
			}
		case i == 3:
			if cur == '.' || cur == '-' {
				pad++
				cur = p[i+pad]
			}

			if !isDigit(cur) {
				return false
			}
		case i == 4:
			if !isAlphaNumericalUpper(cur) {
				return false
			}
		default:
			if !isDigit(cur) {
				return false
			}
		}
	}

	return true
}

// String returns the license plate as an uppercase formatted string.
//
// String always returns the plate in the XXX-XXXX format and in uppercase.
func (p Plate) String() string {
	if !p.IsValid() {
		return ""
	}

	if len(p) == 8 {
		if strings.ContainsRune(string(p), '-') {
			return strings.ToUpper(string(p))
		}

		out := make([]byte, 8)
		copy(out, p)
		out[3] = '-'
		return strings.ToUpper(unsafex.String(out))
	}

	out := make([]byte, 8)
	for i := range out {
		switch {
		case i < 3:
			out[i] = p[i]
		case i == 3:
			out[i] = '-'
		default:
			out[i] = p[i-1]
		}
	}

	return strings.ToUpper(unsafex.String(out))
}

func (p Plate) Value() (driver.Value, error) {
	return p.String(), nil
}

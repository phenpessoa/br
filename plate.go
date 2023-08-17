package br

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidPlate is an error returned when an invalid license plate is
	// encountered.
	ErrInvalidPlate = errors.New("br: invalid license plate")
)

// NewPlate creates a new Plate instance from a string representation.
//
// It removes any punctuation characters and converts the string to uppercase.
func NewPlate(s string) (Plate, error) {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ToUpper(s)

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
func (p Plate) IsValid() bool {
	if len(p) != 7 {
		return false
	}

	for i := range p {
		cur := p[i]
		switch {
		case i < 3:
			if cur < 'A' || cur > 'Z' {
				return false
			}
		case i == 3:
			if cur < '0' || cur > '9' {
				return false
			}
		case i == 4:
			if (cur < 'A' || cur > 'Z') &&
				(cur < '0' || cur > '9') {
				return false
			}
		default:
			if cur < '0' || cur > '9' {
				return false
			}
		}
	}

	return true
}

// String returns the license plate as an uppercase formatted string.
func (p Plate) String() string {
	return strings.ToUpper(string(p))
}

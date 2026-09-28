package br

import (
	"database/sql/driver"
	"errors"
)

// Plate represents a Brazilian vehicle license plate.
type Plate string

// NewPlate creates a new Plate instance from a string representation.
func NewPlate(s string) (Plate, error) {
	plate := Plate(s)
	if !plate.IsValid() {
		return "", ErrInvalidPlate
	}
	return plate, nil
}

// GeneratePlate generates a pseudo-random valid Plate.
func GeneratePlate() Plate {
	data := make([]byte, 8)
	data[3] = '-'

	for i := range 3 {
		data[i] = randomAlphaUpper()
	}

	data[4] = randomDigit()
	data[5] = randomAlphaNumericalUpper()
	data[6], data[7] = randomDigit(), randomDigit()

	return Plate(string(data))
}

// ErrInvalidPlate is an error returned when an invalid license plate is encountered.
var ErrInvalidPlate = errors.New("br: invalid license plate")

// IsValid checks whether the provided license plate is valid based on specific formatting rules.
//
// IsValid will return true if the plate if either a MercoSul or a Brazilian type plate.
//
// The formats accepted are: XXXXXXX, XXX-XXXX, XXX.XXXX
func (p Plate) IsValid() bool {
	switch len(p) {
	case 7:
		for i := range 3 {
			cur := asciiLowerToUpper(p[i])
			if !isAlphaUpper(cur) {
				return false
			}
		}

		if !isDigit(p[3]) {
			return false
		}

		if !isAlphaNumericalUpper(asciiLowerToUpper(p[4])) {
			return false
		}

		return isDigit(p[5]) && isDigit(p[6])
	case 8:
		if p[3] != '.' && p[3] != '-' {
			return false
		}

		for i := range 3 {
			cur := asciiLowerToUpper(p[i])
			if !isAlphaUpper(cur) {
				return false
			}
		}

		if !isDigit(p[4]) {
			return false
		}

		if !isAlphaNumericalUpper(asciiLowerToUpper(p[5])) {
			return false
		}

		return isDigit(p[6]) && isDigit(p[7])
	default:
		return false
	}
}

// String returns the license plate as an uppercase formatted string.
//
// String always returns the plate in the XXX-XXXX format and in uppercase.
func (p Plate) String() string {
	if !p.IsValid() {
		return ""
	}

	if len(p) == 8 {
		if p[3] == '-' && !hasLowerPlate(p) {
			return string(p)
		}

		out := make([]byte, 8)
		copy(out, p)
		out[3] = '-'
		upperPlate(out)
		return string(out)
	}

	if len(p) != 7 {
		return ""
	}

	out := make([]byte, 8)
	copy(out[:3], p[:3])
	out[3] = '-'
	copy(out[4:8], p[3:7])
	upperPlate(out)

	return string(out)
}

// hasLowerPlate reports whether a valid plate has a lowercase letter.
func hasLowerPlate(p Plate) bool {
	return p[0] >= 'a' || p[1] >= 'a' || p[2] >= 'a' || p[5] >= 'a'
}

// upperPlate changes the letters of a valid plate in the XXX-XXXX format to
// uppercase, in place.
func upperPlate(b []byte) {
	_ = b[7]
	b[0] &^= 0x20
	b[1] &^= 0x20
	b[2] &^= 0x20
	b[5] = asciiLowerToUpper(b[5])
}

// Value implements the driver.Valuer interface for Plate.
func (p Plate) Value() (driver.Value, error) {
	return p.String(), nil
}

package br

import (
	"database/sql/driver"
	"errors"
)

// CNS represents a Brazilian CNS.
type CNS string

// NewCNS creates a new CNS instance from a string representation.
//
// It verifies the CNS's validity using checksum digits.
func NewCNS(s string) (CNS, error) {
	cns := CNS(s)
	if !cns.IsValid() {
		return "", ErrInvalidCNS
	}
	return cns, nil
}

// ErrInvalidCNS is an error returned when an invalid CNS is encountered.
var ErrInvalidCNS = errors.New("br: invalid cns")

// IsValid checks whether the provided CNS is valid based on its checksum digits.
func (cns CNS) IsValid() bool {
	switch len(cns) {
	case 15:
		switch cns[0] {
		case '1', '2', '7', '8', '9':
			var sum int
			for i := range 15 {
				cur := cns[i]
				if !isDigit(cur) {
					return false
				}
				sum += int(cur-'0') * (15 - i)
			}
			return sum%11 == 0
		default:
			return false
		}
	case 18:
		switch cns[0] {
		case '1', '2', '7', '8', '9':
			if !isSpace(cns[3]) || !isSpace(cns[8]) || !isSpace(cns[13]) {
				return false
			}

			var sum int
			for i := range 3 {
				cur := cns[i]
				if !isDigit(cur) {
					return false
				}
				sum += int(cur-'0') * (15 - i)
			}

			_cns := cns[4:8]
			for i := range 4 {
				cur := _cns[i]
				if !isDigit(cur) {
					return false
				}
				sum += int(cur-'0') * (12 - i)
			}

			_cns = cns[9:13]
			for i := range 4 {
				cur := _cns[i]
				if !isDigit(cur) {
					return false
				}
				sum += int(cur-'0') * (8 - i)
			}

			_cns = cns[14:18]
			for i := range 4 {
				cur := _cns[i]
				if !isDigit(cur) {
					return false
				}
				sum += int(cur-'0') * (4 - i)
			}

			return sum%11 == 0
		default:
			return false
		}
	default:
		return false
	}
}

// String returns the CNS formatted as XXX XXXX XXXX XXXX.
func (cns CNS) String() string {
	if !cns.IsValid() {
		return ""
	}

	if len(cns) == 18 {
		return string(cns)
	}

	if len(cns) != 15 {
		return ""
	}

	out := make([]byte, 18)
	out[3] = ' '
	out[8] = ' '
	out[13] = ' '

	copy(out[:3], cns[:3])
	copy(out[4:8], cns[3:7])
	copy(out[9:13], cns[7:11])
	copy(out[14:18], cns[11:15])

	return string(out)
}

// Value implements the driver.Valuer interface for CNS.
func (cns CNS) Value() (driver.Value, error) {
	return cns.String(), nil
}

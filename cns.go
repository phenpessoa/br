package br

import (
	"errors"
	"unicode"

	"github.com/phenpessoa/gutils/unsafex"
)

var (
	ErrInvalidCNS = errors.New("br: invalid cns passed")
)

func NewCNS(s string) (CNS, error) {
	cns := CNS(s)
	if !cns.IsValid() {
		return "", ErrInvalidCNS
	}
	return cns, nil
}

type CNS string

func (cns CNS) IsValid() bool {
	if len(cns) != 15 && len(cns) != 18 {
		return false
	}

	if !isFirstCNSDigitValid(cns[0]) {
		return false
	}

	var sum, pad int
	for i := 0; i < 15; i++ {
		cur := cns[i+pad]
		if unicode.IsSpace(rune(cur)) {
			if i != 3 && i != 7 && i != 11 {
				return false
			}
			pad++
			cur = cns[i+pad]
		}

		if cur < '0' || cur > '9' {
			return false
		}

		sum += int(cur-'0') * (15 - i)
	}

	return sum%11 == 0
}

func (cns CNS) String() string {
	if !cns.IsValid() {
		return ""
	}

	if len(cns) == 18 {
		return string(cns)
	}

	out := make([]byte, 18)

	var pad int
	for i := range out {
		switch i {
		case 3, 8, 13:
			out[i] = ' '
			pad++
			continue
		}
		out[i] = cns[i-pad]
	}

	return unsafex.String(out)
}

func isFirstCNSDigitValid(d byte) bool {
	switch d {
	case '1', '2', '7', '8', '9':
		return true
	default:
		return false
	}
}

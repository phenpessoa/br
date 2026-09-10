package br

import (
	"database/sql/driver"
	"errors"
	"math/bits"
)

// CPF represents a Brazilian CPF.
type CPF string

// NewCPF creates a new CPF instance from a string representation.
//
// It verifies the CPF's validity using checksum digits and rejects CPFs whose
// digits are all the same.
func NewCPF(s string) (CPF, error) {
	cpf := CPF(s)
	if !cpf.IsValid() {
		return "", ErrInvalidCPF
	}
	return cpf, nil
}

// GenerateCPF generates a pseudo-random valid CPF.
func GenerateCPF() CPF {
	// Taken from stdlib, take a look at the method:
	// func (r *Rand) uint64n(n uint64) uint64
	//
	// See also:
	// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction
	// https://lemire.me/blog/2016/06/30/fast-random-shuffling

	const (
		n      = 1_000_000_000
		thresh = (1 << 64) % n
	)

	var v uint32
	for {
		hi, lo := bits.Mul64(pcg.Uint64(), n)
		for lo < thresh {
			hi, lo = bits.Mul64(pcg.Uint64(), n)
		}
		v = uint32(hi)

		if v%111_111_111 != 0 {
			break
		}
	}

	return cpfFromBase(v)
}

func cpfFromBase(v uint32) CPF {
	q := v / 100_000 // first four digits
	r := v % 100_000 // last five digits

	d0 := q / 1000
	d1 := q / 100 % 10
	d2 := q / 10 % 10
	d3 := q % 10
	d4 := r / 10_000
	d5 := r / 1000 % 10
	d6 := r / 100 % 10
	d7 := r / 10 % 10
	d8 := r % 10

	sum1 := 10*d0 + 9*d1 + 8*d2 + 7*d3 + 6*d4 + 5*d5 + 4*d6 + 3*d7 + 2*d8
	var e1 uint32
	if rest := sum1 % 11; rest >= 2 {
		e1 = 11 - rest
	}

	sum2 := sum1 + d0 + d1 + d2 + d3 + d4 + d5 + d6 + d7 + d8 + 2*e1
	var e2 uint32
	if rest := sum2 % 11; rest >= 2 {
		e2 = 11 - rest
	}

	return CPF([]byte{
		byte(d0) + '0', byte(d1) + '0', byte(d2) + '0', '.',
		byte(d3) + '0', byte(d4) + '0', byte(d5) + '0', '.',
		byte(d6) + '0', byte(d7) + '0', byte(d8) + '0', '-',
		byte(e1) + '0', byte(e2) + '0',
	})
}

// ErrInvalidCPF is an error returned when an invalid CPF is encountered.
var ErrInvalidCPF = errors.New("br: invalid cpf")

// IsValid checks whether the provided CPF is valid based on its checksum
// digits.
//
// CPFs whose digits are all the same (e.g. 111.111.111-11) satisfy the
// checksum but are never issued, and are rejected.
func (cpf CPF) IsValid() bool {
	switch len(cpf) {
	case 11:
		return validCPF11(string(cpf))
	case 14:
		return validCPF14(string(cpf))
	default:
		return false
	}
}

func validCPF11(s string) bool {
	if len(s) != 11 {
		return false
	}

	lo := load64(s)
	hi := load32(s[7:])

	if nonDigits64(lo)|uint64(nonDigits32(hi)) != 0 {
		return false
	}

	// A word rotated by one byte equals itself only if all its bytes are equal.
	if (lo^bits.RotateLeft64(lo, 8))|uint64(hi^bits.RotateLeft32(hi, 8)) == 0 {
		return false
	}

	v := lo & 0x0F0F0F0F0F0F0F0F // ASCII digits to their values
	d8 := uint64(hi>>8) & 0x0F
	d9 := uint64(hi>>16) & 0x0F
	d10 := uint64(hi>>24) & 0x0F

	// Spreading every other digit into 16-bit lanes and multiplying by the
	// weights packed in reverse makes the top lane of each product accumulate
	// a 4-term weighted sum. No partial sum can exceed 16 bits, so lanes
	// never carry into each other.
	e := v & 0x00FF00FF00FF00FF        // d0 d2 d4 d6
	o := (v >> 8) & 0x00FF00FF00FF00FF // d1 d3 d5 d7
	const we = 10<<48 | 8<<32 | 6<<16 | 4
	const wo = 9<<48 | 7<<32 | 5<<16 | 3
	sum1 := (e*we)>>48 + (o*wo)>>48 + 2*d8

	var e1 uint64
	if rest := sum1 % 11; rest >= 2 {
		e1 = 11 - rest
	}

	// The second checksum weights each base digit one higher than the first:
	// sum2 = sum1 + d0+..+d8 + 2*d9. Multiplying by 0x01 repeated makes the
	// top byte of the product the sum of all eight bytes (at most 72, so no
	// carry).
	sum2 := sum1 + (v*0x0101010101010101)>>56 + d8 + 2*d9
	var e2 uint64
	if rest := sum2 % 11; rest >= 2 {
		e2 = 11 - rest
	}

	return (d9^e1)|(d10^e2) == 0
}

func validCPF14(s string) bool {
	if len(s) != 14 {
		return false
	}

	// NOTICE: d5 is the last digit of lo and the first digit of hi.
	lo := load64(s)
	hi := load64(s[6:])

	const (
		pMaskLo = 0xFF000000FF000000 // punctuation lanes of lo (bytes 3 and 7)
		pWantLo = uint64('.')<<24 | uint64('.')<<56
		pMaskHi = 0x0000FF000000FF00 // punctuation lanes of hi (bytes 1 and 5)
		pWantHi = uint64('.')<<8 | uint64('-')<<40
		blendLo = uint64('0')<<24 | uint64('0')<<56
		blendHi = uint64('0')<<8 | uint64('0')<<40
	)

	pd := ((lo ^ pWantLo) & pMaskLo) | ((hi ^ pWantHi) & pMaskHi)

	lod := (lo &^ pMaskLo) | blendLo
	hid := (hi &^ pMaskHi) | blendHi
	if pd|nonDigits64(lod)|nonDigits64(hid) != 0 {
		return false
	}

	// Repeated-digit rejection: broadcast the first digit to every lane and
	// compare against both words, with the punctuation lanes blended to '0'
	// on both sides.
	bb := (lo & 0xFF) * 0x0101010101010101
	rep := (lod ^ ((bb &^ pMaskLo) | blendLo)) | (hid ^ ((bb &^ pMaskHi) | blendHi))
	if rep == 0 {
		return false
	}

	v := lod & 0x0F0F0F0F0F0F0F0F // d0 d1 d2 0 d3 d4 d5 0
	h := hid & 0x0F0F0F0F0F0F0F0F // d5 0 d6 d7 d8 0 dv1 dv2
	d6 := (h >> 16) & 0x0F
	d7 := (h >> 24) & 0x0F
	d8 := (h >> 32) & 0x0F
	dv1 := (h >> 48) & 0x0F
	dv2 := h >> 56

	// Same lane dot products as validCPF11, with the weights placed to match
	// this layout.
	e := v & 0x00FF00FF00FF00FF        // d0 d2 d3 d5
	o := (v >> 8) & 0x00FF00FF00FF00FF // d1 0 d4 0
	const we = 10<<48 | 8<<32 | 7<<16 | 5
	const wo = 9<<48 | 6<<16
	sum1 := (e*we)>>48 + (o*wo)>>48 + 4*d6 + 3*d7 + 2*d8

	var e1 uint64
	if rest := sum1 % 11; rest >= 2 {
		e1 = 11 - rest
	}

	sum2 := sum1 + (v*0x0101010101010101)>>56 + d6 + d7 + d8 + 2*dv1
	var e2 uint64
	if rest := sum2 % 11; rest >= 2 {
		e2 = 11 - rest
	}

	return (dv1^e1)|(dv2^e2) == 0
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
		return ""
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

// Value implements the driver.Valuer interface for CPF.
func (cpf CPF) Value() (driver.Value, error) {
	return cpf.String(), nil
}

package br

import (
	"math/rand/v2"
	"strings"
	"testing"
)

var cpfSink CPF

func BenchmarkGenerateCPF(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		cpfSink = GenerateCPF()
	}
}

func TestGenerateCPF(t *testing.T) {
	for range 1_000_000 {
		if cpf := GenerateCPF(); !cpf.IsValid() {
			t.Errorf("invalid CPF generated: %s", string(cpf))
		}
	}
}

func BenchmarkCPF_IsValid14(b *testing.B) {
	const cpf = CPF("453.178.287-91")
	if !cpf.IsValid() {
		b.Error("invalid cpf on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cpf.IsValid()
	}
}

func BenchmarkCPF_IsValid11(b *testing.B) {
	const cpf = CPF("45317828791")
	if !cpf.IsValid() {
		b.Error("invalid cpf on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cpf.IsValid()
	}
}

func BenchmarkCPF_IsValid14Invalid(b *testing.B) {
	const cpf = CPF("453.178.287-92")
	if cpf.IsValid() {
		b.Error("valid cpf on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cpf.IsValid()
	}
}

func BenchmarkCPF_IsValid11Invalid(b *testing.B) {
	const cpf = CPF("45317828792")
	if cpf.IsValid() {
		b.Error("valid cpf on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cpf.IsValid()
	}
}

func BenchmarkCPF_String14(b *testing.B) {
	const cpf = CPF("453.178.287-91")
	if !cpf.IsValid() {
		b.Error("invalid cpf on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		stringSink = cpf.String()
	}
}

func BenchmarkCPF_String11(b *testing.B) {
	const cpf = CPF("45317828791")
	if !cpf.IsValid() {
		b.Error("invalid cpf on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		stringSink = cpf.String()
	}
}

func TestCPF_IsValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cpf   CPF
		valid bool
	}{
		{
			name:  "formatted CPF",
			cpf:   CPF("453.178.287-91"),
			valid: true,
		},
		{
			name:  "raw CPF",
			cpf:   CPF("45317828791"),
			valid: true,
		},
		{
			name:  "invalid first digit formatted CPF",
			cpf:   CPF("453.178.287-81"),
			valid: false,
		},
		{
			name:  "invalid first digit raw CPF",
			cpf:   CPF("45317828781"),
			valid: false,
		},
		{
			name:  "invalid second digit formatted CPF",
			cpf:   CPF("453.178.287-92"),
			valid: false,
		},
		{
			name:  "invalid second digit raw CPF",
			cpf:   CPF("45317828792"),
			valid: false,
		},
		{
			name:  "empty cpf",
			cpf:   CPF(""),
			valid: false,
		},
		{
			name:  "incorrect length cpf",
			cpf:   CPF("123"),
			valid: false,
		},
		{
			name:  "invalid characters",
			cpf:   CPF("abc.def.ghi-jk"),
			valid: false,
		},
		{
			name:  "invalid separators",
			cpf:   CPF("453-178-287.92"),
			valid: false,
		},
		{
			name:  "repeated digit formatted CPF",
			cpf:   CPF("111.111.111-11"),
			valid: false,
		},
		{
			name:  "repeated digit raw CPF",
			cpf:   CPF("11111111111"),
			valid: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cpf.IsValid() != tc.valid {
				t.Errorf(
					"\ncpf: %s\nshould be valid: %v\nis valid: %v",
					tc.cpf, tc.valid, tc.cpf.IsValid(),
				)
			}
		})
	}
}

func TestCPF_String(t *testing.T) {
	for _, tc := range []struct {
		name string
		cpf  CPF
		want string
	}{
		{
			name: "formatted CPF",
			cpf:  CPF("453.178.287-91"),
			want: "453.178.287-91",
		},
		{
			name: "raw CPF",
			cpf:  CPF("45317828791"),
			want: "453.178.287-91",
		},
		{
			name: "empty CPF",
			cpf:  CPF(""),
			want: "",
		},
		{
			name: "invalid",
			cpf:  CPF("123"),
			want: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cpf.String() != tc.want {
				t.Errorf(
					"\ncpf: %s\nshould be formatted like: %s\nis formatted like: %s",
					tc.cpf, tc.want, tc.cpf.String(),
				)
			}
		})
	}
}

func TestCPF_RepeatedDigits(t *testing.T) {
	for d := byte('0'); d <= '9'; d++ {
		raw := strings.Repeat(string(rune(d)), 11)
		formatted := raw[:3] + "." + raw[3:6] + "." + raw[6:9] + "-" + raw[9:]

		if CPF(raw).IsValid() {
			t.Errorf("repeated-digit CPF should be invalid: %s", raw)
		}
		if CPF(formatted).IsValid() {
			t.Errorf("repeated-digit CPF should be invalid: %s", formatted)
		}
		if got := CPF(raw).String(); got != "" {
			t.Errorf("String of repeated-digit CPF should be empty, got: %s", got)
		}
		if _, err := NewCPF(formatted); err == nil {
			t.Errorf("NewCPF should reject repeated-digit CPF: %s", formatted)
		}
	}
}

// cpfRefValid is a deliberately straightforward reference implementation of
// the CPF validation rules, used to cross-check the optimized IsValid.
func cpfRefValid(s string) bool {
	var digits [11]byte
	switch len(s) {
	case 11:
		copy(digits[:], s)
	case 14:
		if s[3] != '.' || s[7] != '.' || s[11] != '-' {
			return false
		}
		copy(digits[0:3], s[0:3])
		copy(digits[3:6], s[4:7])
		copy(digits[6:9], s[8:11])
		copy(digits[9:11], s[12:14])
	default:
		return false
	}

	same := true
	for _, c := range digits {
		if c < '0' || c > '9' {
			return false
		}
		if c != digits[0] {
			same = false
		}
	}
	if same {
		return false
	}

	var sum1, sum2 int
	for i, c := range digits[:9] {
		d := int(c - '0')
		sum1 += (10 - i) * d
		sum2 += (11 - i) * d
	}
	sum2 += 2 * int(digits[9]-'0')

	dv1 := 11 - sum1%11
	if dv1 >= 10 {
		dv1 = 0
	}
	dv2 := 11 - sum2%11
	if dv2 >= 10 {
		dv2 = 0
	}
	return int(digits[9]-'0') == dv1 && int(digits[10]-'0') == dv2
}

func cpfToRaw(c string) string {
	return c[0:3] + c[4:7] + c[8:11] + c[12:14]
}

func TestCPF_IsValidDifferential(t *testing.T) {
	check := func(s string) {
		t.Helper()
		if got, want := CPF(s).IsValid(), cpfRefValid(s); got != want {
			t.Fatalf("IsValid(%q) = %v, reference implementation says %v", s, got, want)
		}
	}

	for range 50_000 {
		c := string(GenerateCPF())
		check(c)
		check(cpfToRaw(c))
	}

	for range 100_000 {
		s := string(GenerateCPF())
		if rand.IntN(2) == 0 {
			s = cpfToRaw(s)
		}
		b := []byte(s)
		b[rand.IntN(len(b))] = byte(rand.IntN(256))
		check(string(b))
	}

	const alphabet = "0123456789.-/: \x00"
	for range 100_000 {
		n := 11
		if rand.IntN(2) == 0 {
			n = 14
		}
		b := make([]byte, n)
		for i := range b {
			b[i] = alphabet[rand.IntN(len(alphabet))]
		}
		check(string(b))
	}
}

func TestCPF_NonDigitBytes(t *testing.T) {
	bad := []byte{0x00, ' ', '*', '+', ',', '-', '.', '/', ':', ';', '@', 'a', 0x7F, 0x80, 0xFF}

	c := string(GenerateCPF())
	raw := cpfToRaw(c)

	for i := range raw {
		for _, b := range bad {
			mut := []byte(raw)
			mut[i] = b
			if CPF(mut).IsValid() {
				t.Errorf("CPF with byte %#x at position %d should be invalid: %q", b, i, mut)
			}
		}
	}

	for _, i := range []int{0, 1, 2, 4, 5, 6, 8, 9, 10, 12, 13} {
		for _, b := range bad {
			mut := []byte(c)
			mut[i] = b
			if CPF(mut).IsValid() {
				t.Errorf("CPF with byte %#x at position %d should be invalid: %q", b, i, mut)
			}
		}
	}

	for _, i := range []int{3, 7, 11} {
		for _, b := range []byte{'0', ' ', '/', ':', '-', '.'} {
			if (i == 11 && b == '-') || (i != 11 && b == '.') {
				continue
			}
			mut := []byte(c)
			mut[i] = b
			if CPF(mut).IsValid() {
				t.Errorf("CPF with separator %q at position %d should be invalid: %q", b, i, mut)
			}
		}
	}
}

func TestCPFFromBase(t *testing.T) {
	if got := cpfFromBase(453_178_287); got != CPF("453.178.287-91") {
		t.Errorf("cpfFromBase(453178287) = %s, want 453.178.287-91", got)
	}
	if got := cpfFromBase(0); got != CPF("000.000.000-00") {
		t.Errorf("cpfFromBase(0) = %s, want 000.000.000-00", got)
	}
	for _, v := range []uint32{1, 12, 123_456_789, 453_178_287, 999_999_998} {
		c := string(cpfFromBase(v))
		if !cpfRefValid(c) {
			t.Errorf("cpfFromBase(%d) produced invalid CPF: %s", v, c)
		}
	}
}

func TestGenerateCPF_Distribution(t *testing.T) {
	positions := []int{0, 1, 2, 4, 5, 6, 8, 9, 10}
	var counts [9][10]int
	const rounds = 1_000_000
	for range rounds {
		cpf := GenerateCPF()
		for i, p := range positions {
			counts[i][cpf[p]-'0']++
		}
	}
	for i, c := range counts {
		for d, n := range c {
			ratio := float64(n) / (rounds / 10)
			if ratio < 0.95 || ratio > 1.05 {
				t.Errorf("digit %d at position %d: count %d (ratio %.3f)", d, positions[i], n, ratio)
			}
		}
	}
}

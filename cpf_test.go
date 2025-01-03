package br

import "testing"

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

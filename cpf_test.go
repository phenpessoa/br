package br

import "testing"

func BenchmarkCPF_IsValid(b *testing.B) {
	const cpfBolsonaro = CPF("453.178.287-91")
	if !cpfBolsonaro.IsValid() {
		b.Error("invalid cpfBolsonaro on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = cpfBolsonaro.IsValid()
	}
}

func TestCPF_IsValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cpf   CPF
		valid bool
	}{
		{
			name:  "formatted CPF Bolsonaro",
			cpf:   CPF("453.178.287-91"),
			valid: true,
		},
		{
			name:  "raw CPF Bolsonaro",
			cpf:   CPF("45317828791"),
			valid: true,
		},
		{
			name:  "invalid first digit formatted CPF Bolsonaro",
			cpf:   CPF("453.178.287-81"),
			valid: false,
		},
		{
			name:  "invalid first digit raw CPF Bolsonaro",
			cpf:   CPF("453.178.287-81"),
			valid: false,
		},
		{
			name:  "invalid second digit formatted CPF Bolsonaro",
			cpf:   CPF("453.178.287-92"),
			valid: false,
		},
		{
			name:  "invalid second digit raw CPF Bolsonaro",
			cpf:   CPF("453.178.287-92"),
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
			name: "formatted CPF Bolsonaro",
			cpf:  CPF("453.178.287-91"),
			want: "453.178.287-91",
		},
		{
			name: "raw CPF Bolsonaro",
			cpf:  CPF("45317828791"),
			want: "453.178.287-91",
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

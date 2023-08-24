package br

import "testing"

var sink bool

func BenchmarkCNPJ_IsValid(b *testing.B) {
	const cnpjPetrobras = CNPJ("33.000.167/1002-46")
	if !cnpjPetrobras.IsValid() {
		b.Error("invalid cnpjPetrobras on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = cnpjPetrobras.IsValid()
	}
}

func TestCNPJ_IsValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cnpj  CNPJ
		valid bool
	}{
		{
			name:  "formatted CNPJ Petrobras",
			cnpj:  CNPJ("33.000.167/1002-46"),
			valid: true,
		},
		{
			name:  "raw CNPJ Petrobras",
			cnpj:  CNPJ("33000167100246"),
			valid: true,
		},
		{
			name:  "invalid first digit formatted CNPJ Petrobras",
			cnpj:  CNPJ("33.000.167/1002-56"),
			valid: false,
		},
		{
			name:  "invalid first digit raw CNPJ Petrobras",
			cnpj:  CNPJ("33000167100256"),
			valid: false,
		},
		{
			name:  "invalid second digit formatted CNPJ Petrobras",
			cnpj:  CNPJ("33.000.167/1002-45"),
			valid: false,
		},
		{
			name:  "invalid second digit raw CNPJ Petrobras",
			cnpj:  CNPJ("33000167100245"),
			valid: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cnpj.IsValid() != tc.valid {
				t.Errorf(
					"\ncnpj: %s\nshould be valid: %v\nis valid: %v",
					tc.cnpj, tc.valid, tc.cnpj.IsValid(),
				)
			}
		})
	}
}

func TestCNPJ_String(t *testing.T) {
	for _, tc := range []struct {
		name string
		cnpj CNPJ
		want string
	}{
		{
			name: "formatted CNPJ Petrobras",
			cnpj: CNPJ("33.000.167/1002-46"),
			want: "33.000.167/1002-46",
		},
		{
			name: "raw CNPJ Petrobras",
			cnpj: CNPJ("33000167100246"),
			want: "33.000.167/1002-46",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cnpj.String() != tc.want {
				t.Errorf(
					"\ncnpj: %s\nshould be formatted like: %s\nis formatted like: %s",
					tc.cnpj, tc.want, tc.cnpj.String(),
				)
			}
		})
	}
}

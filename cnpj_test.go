package br

import "testing"

var cnpjSink CNPJ

func BenchmarkGenerateCNPJ(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		cnpjSink = GenerateCNPJ()
	}
}

func TestGenerateCNPJ(t *testing.T) {
	for range 1_000_000 {
		if cnpj := GenerateCNPJ(); !cnpj.IsValid() {
			t.Errorf("invalid CNPJ generated: %s", string(cnpj))
		}
	}
}

func BenchmarkCNPJ_IsValid18(b *testing.B) {
	const cnpj = CNPJ("33.000.167/1002-46")
	if !cnpj.IsValid() {
		b.Error("invalid cnpj on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cnpj.IsValid()
	}
}

func BenchmarkCNPJ_IsValid14(b *testing.B) {
	const cnpj = CNPJ("33000167100246")
	if !cnpj.IsValid() {
		b.Error("invalid cnpj on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cnpj.IsValid()
	}
}

func BenchmarkCNPJ_IsValid18Invalid(b *testing.B) {
	const cnpj = CNPJ("33.000.167/1002-47")
	if cnpj.IsValid() {
		b.Error("valid cnpj on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cnpj.IsValid()
	}
}

func BenchmarkCNPJ_IsValid14Invalid(b *testing.B) {
	const cnpj = CNPJ("33000167100247")
	if cnpj.IsValid() {
		b.Error("valid cnpj on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cnpj.IsValid()
	}
}

func BenchmarkCNPJ_String18(b *testing.B) {
	const cnpj = CNPJ("33.000.167/1002-46")
	if !cnpj.IsValid() {
		b.Error("invalid cnpj on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		stringSink = cnpj.String()
	}
}

func BenchmarkCNPJ_String14(b *testing.B) {
	const cnpj = CNPJ("33000167100246")
	if !cnpj.IsValid() {
		b.Error("invalid cnpj on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		stringSink = cnpj.String()
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
		{
			name:  "formatted CNPJ BB",
			cnpj:  CNPJ("00.000.000/0001-91"),
			valid: true,
		},
		{
			name:  "raw CNPJ BB",
			cnpj:  CNPJ("00000000000191"),
			valid: true,
		},
		{
			name:  "formatted CNPJ alfanumerico",
			cnpj:  CNPJ("AA.AAA.AAA/AAAA-45"),
			valid: true,
		},
		{
			name:  "raw CNPJ alfanumerico",
			cnpj:  CNPJ("AAAAAAAAAAAA45"),
			valid: true,
		},
		{
			name:  "formatted CNPJ alfanumerico lower",
			cnpj:  CNPJ("aa.aaa.aaa/aaaa-45"),
			valid: true,
		},
		{
			name:  "raw CNPJ alfanumerico lower",
			cnpj:  CNPJ("aaaaaaaaaaaa45"),
			valid: true,
		},
		{
			name:  "formatted CNPJ random 1",
			cnpj:  CNPJ("34.588.324/0001-04"),
			valid: true,
		},
		{
			name:  "formatted CNPJ random 2",
			cnpj:  CNPJ("72.285.712/0001-05"),
			valid: true,
		},
		{
			name:  "formatted CNPJ alfanumerico 2",
			cnpj:  CNPJ("AB.CDE.FGI/HIJK-56"),
			valid: true,
		},
		{
			name:  "raw alfanumerico 3",
			cnpj:  CNPJ("12ABC34501DE35"),
			valid: true,
		},
		{
			name:  "formatted alfanumerico 3",
			cnpj:  CNPJ("12.ABC.345/01DE-35"),
			valid: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.cnpj.IsValid()
			if isValid != tc.valid {
				t.Errorf(
					"\ncnpj: %s\nshould be valid: %v\nis valid: %v",
					tc.cnpj, tc.valid, isValid,
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
		{
			name: "formatted CNPJ alfanumerico lower",
			cnpj: CNPJ("aa.aaa.aaa/aaaa-45"),
			want: "AA.AAA.AAA/AAAA-45",
		},
		{
			name: "raw CNPJ alfanumerico lower",
			cnpj: CNPJ("aaaaaaaaaaaa45"),
			want: "AA.AAA.AAA/AAAA-45",
		},
		{
			name: "formatted CNPJ alfanumerico",
			cnpj: CNPJ("AA.AAA.AAA/AAAA-45"),
			want: "AA.AAA.AAA/AAAA-45",
		},
		{
			name: "raw CNPJ alfanumerico",
			cnpj: CNPJ("AAAAAAAAAAAA45"),
			want: "AA.AAA.AAA/AAAA-45",
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

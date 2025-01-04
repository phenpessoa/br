package br

import "testing"

var cnhSink CNH

func BenchmarkGenerateCNH(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		cnhSink = GenerateCNH()
	}
}

func TestGenerateCNH(t *testing.T) {
	for range 1_000_000 {
		if cnh := GenerateCNH(); !cnh.IsValid() {
			t.Errorf("invalid CNH generated: %s", string(cnh))
		}
	}
}

func BenchmarkCNH_IsValid(b *testing.B) {
	const cnh = CNH("96300689842")
	if !cnh.IsValid() {
		b.Error("invalid cnh on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cnh.IsValid()
	}
}

func BenchmarkCNH_IsValidInvalid(b *testing.B) {
	const cnh = CNH("96300689843")
	if cnh.IsValid() {
		b.Error("valid cnh on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		boolSink = cnh.IsValid()
	}
}

func BenchmarkCNH_String(b *testing.B) {
	const cnh = CNH("96300689842")
	if !cnh.IsValid() {
		b.Error("invalid cnh on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for range b.N {
		stringSink = cnh.String()
	}
}

func TestCNH_IsValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cnh   CNH
		valid bool
	}{
		{
			name:  "raw CNH 1",
			cnh:   CNH("96300689842"),
			valid: true,
		},
		{
			name:  "raw CNH 2",
			cnh:   CNH("74510118051"),
			valid: true,
		},
		{
			name:  "raw CNH 3",
			cnh:   CNH("93826104830"),
			valid: true,
		},
		{
			name:  "invalid first CNH digit",
			cnh:   CNH("96300689852"),
			valid: false,
		},
		{
			name:  "invalid first CNH digit",
			cnh:   CNH("96300689843"),
			valid: false,
		},
		{
			name:  "empty cnh",
			cnh:   CNH(""),
			valid: false,
		},
		{
			name:  "incorrect length cnh",
			cnh:   CNH("123"),
			valid: false,
		},
		{
			name:  "invalid characters",
			cnh:   CNH("aaaaaaaaaaa"),
			valid: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cnh.IsValid() != tc.valid {
				t.Errorf(
					"\ncnh: %s\nshould be valid: %v\nis valid: %v",
					tc.cnh, tc.valid, tc.cnh.IsValid(),
				)
			}
		})
	}
}

func TestCNH_String(t *testing.T) {
	for _, tc := range []struct {
		name string
		cnh  CNH
		want string
	}{
		{
			name: "raw CNH",
			cnh:  CNH("96300689842"),
			want: "96300689842",
		},
		{
			name: "empty CNH",
			cnh:  CNH(""),
			want: "",
		},
		{
			name: "invalid",
			cnh:  CNH("123"),
			want: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cnh.String() != tc.want {
				t.Errorf(
					"\ncnh: %s\nshould be formatted like: %s\nis formatted like: %s",
					tc.cnh, tc.want, tc.cnh.String(),
				)
			}
		})
	}
}

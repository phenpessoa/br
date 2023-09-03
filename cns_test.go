package br

import (
	"strings"
	"testing"
)

func BenchmarkCNS_IsValid(b *testing.B) {
	const randomCNS = CNS("708521331850008")
	if !randomCNS.IsValid() {
		b.Error("invalid randomCNS on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = randomCNS.IsValid()
	}
}

func TestCNS_IsValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cns   CNS
		valid bool
	}{
		{
			name:  "formatted random CNS",
			cns:   CNS("708 5213 3185 0008"),
			valid: true,
		},
		{
			name:  "raw random CNS",
			cns:   CNS("708521331850008"),
			valid: true,
		},
		{
			name:  "invalid formatted CNS",
			cns:   CNS("708 5213 3185 0001"),
			valid: false,
		},
		{
			name:  "cns with len 14",
			cns:   CNS(strings.Repeat("1", 14)),
			valid: false,
		},
		{
			name:  "cns with len 19",
			cns:   CNS(strings.Repeat("1", 19)),
			valid: false,
		},
		{
			name:  "cns with invalid first digit",
			cns:   CNS("008521331850008"),
			valid: false,
		},
		{
			name:  "empty cns",
			cns:   CNS(""),
			valid: false,
		},
		{
			name:  "cns with invalid separator 1",
			cns:   CNS("708A5213A3185A0008"),
			valid: false,
		},
		{
			name:  "cns with invalid separator 2",
			cns:   CNS("708.5213.3185.0008"),
			valid: false,
		},
		{
			name:  "cns with invalid separator 3",
			cns:   CNS("708-5213-3185-0008"),
			valid: false,
		},
		{
			name:  "cns with invalid format",
			cns:   CNS("7085 213 3185 0008"),
			valid: false,
		},
		{
			name:  "cns with invalid digits",
			cns:   CNS("915 5017 0193 0306"),
			valid: false,
		},
		{
			name:  "cns with invalid digits 2",
			cns:   CNS("915501701930306"),
			valid: false,
		},
		{
			name:  "cns with invalid digits 3",
			cns:   CNS("174 2241 7133 0004"),
			valid: false,
		},
		{
			name:  "cns with invalid digits 4",
			cns:   CNS("174224171330004"),
			valid: false,
		},
		{
			name:  "cns with invalid digits 5",
			cns:   CNS("259 7557 3388 0001"),
			valid: false,
		},
		{
			name:  "cns with invalid digits 6",
			cns:   CNS("259755733880001"),
			valid: false,
		},
		{
			name:  "valid 1",
			cns:   CNS("174 5984 3528 0018"),
			valid: true,
		},
		{
			name:  "valid 2",
			cns:   CNS("174598435280018"),
			valid: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cns.IsValid() != tc.valid {
				t.Errorf(
					"\ncns: %s\nshould be valid: %v\nis valid: %v",
					tc.cns, tc.valid, tc.cns.IsValid(),
				)
			}
		})
	}
}

func TestCNS_String(t *testing.T) {
	for _, tc := range []struct {
		name string
		cns  CNS
		want string
	}{
		{
			name: "formatted random CNS",
			cns:  CNS("708 5213 3185 0008"),
			want: "708 5213 3185 0008",
		},
		{
			name: "raw random CNS",
			cns:  CNS("708521331850008"),
			want: "708 5213 3185 0008",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cns.String() != tc.want {
				t.Errorf(
					"\ncns: %s\nshould be formatted like: %s\nis formatted like: %s",
					tc.cns, tc.want, tc.cns.String(),
				)
			}
		})
	}
}

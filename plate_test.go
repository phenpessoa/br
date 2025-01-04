package br

import "testing"

var plateSink Plate

func BenchmarkGeneratePlate(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		plateSink = GeneratePlate()
	}
}

func TestGeneratePlate(t *testing.T) {
	for range 1_000_000 {
		if plate := GeneratePlate(); !plate.IsValid() {
			t.Errorf("invalid Plate generated: %s", string(plate))
		}
	}
}

func BenchmarkPlate_IsValid8(b *testing.B) {
	const plate = Plate("BRA-2023")
	if !plate.IsValid() {
		b.Error("invalid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		boolSink = plate.IsValid()
	}
}

func BenchmarkPlate_IsValid7(b *testing.B) {
	const plate = Plate("BRA2023")
	if !plate.IsValid() {
		b.Error("invalid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		boolSink = plate.IsValid()
	}
}

func BenchmarkPlate_IsValidInvalid8(b *testing.B) {
	const plate = Plate("BRA-202A")
	if plate.IsValid() {
		b.Error("valid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		boolSink = plate.IsValid()
	}
}

func BenchmarkPlate_IsValidInvalid7(b *testing.B) {
	const plate = Plate("BRA202A")
	if plate.IsValid() {
		b.Error("valid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		boolSink = plate.IsValid()
	}
}

func BenchmarkPlate_String8Dash(b *testing.B) {
	const plate = Plate("BRA-2023")
	if !plate.IsValid() {
		b.Error("invalid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stringSink = plate.String()
	}
}

func BenchmarkPlate_String8Dot(b *testing.B) {
	const plate = Plate("BRA.2023")
	if !plate.IsValid() {
		b.Error("invalid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stringSink = plate.String()
	}
}

func BenchmarkPlate_String7(b *testing.B) {
	const plate = Plate("BRA2023")
	if !plate.IsValid() {
		b.Error("invalid plate on benchmark")
		b.FailNow()
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stringSink = plate.String()
	}
}

func TestPlate_IsValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		plate Plate
		valid bool
	}{
		{
			name:  "formatted BR 1 plate",
			plate: Plate("BRA-2023"),
			valid: true,
		},
		{
			name:  "formatted BR 2 plate",
			plate: Plate("BRA.2023"),
			valid: true,
		},
		{
			name:  "raw BR plate",
			plate: Plate("BRA2023"),
			valid: true,
		},
		{
			name:  "formatted Mercosul 1 plate",
			plate: Plate("BRA-2A23"),
			valid: true,
		},
		{
			name:  "formatted Mercosul 2 plate",
			plate: Plate("BRA.2A23"),
			valid: true,
		},
		{
			name:  "raw Mercosul plate",
			plate: Plate("BRA2A23"),
			valid: true,
		},
		{
			name:  "invalid plate 1",
			plate: Plate("BRA223"),
			valid: false,
		},
		{
			name:  "invalid plate 1",
			plate: Plate("BRA02023"),
			valid: false,
		},
		{
			name:  "invalid plate 1",
			plate: Plate("BRAA2023"),
			valid: false,
		},
		{
			name:  "invalid plate 2",
			plate: Plate("BRA-A023"),
			valid: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.plate.IsValid() != tc.valid {
				t.Errorf(
					"\nplate: %s\nshould be valid: %v\nis valid: %v",
					tc.plate, tc.valid, tc.plate.IsValid(),
				)
			}
		})
	}
}

func TestPlate_String(t *testing.T) {
	for _, tc := range []struct {
		name  string
		plate Plate
		want  string
	}{
		{
			name:  "invalid plate",
			plate: Plate("34fsd"),
			want:  "",
		},
		{
			name:  "raw plate",
			plate: Plate("BRA2023"),
			want:  "BRA-2023",
		},
		{
			name:  "already formatted plate",
			plate: Plate("BRA-2023"),
			want:  "BRA-2023",
		},
		{
			name:  "dot formatted plate",
			plate: Plate("BRA.2023"),
			want:  "BRA-2023",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.plate.String() != tc.want {
				t.Errorf(
					"\nplate: %s\nshould be formatted like: %s\nis formatted like: %s",
					tc.plate, tc.want, tc.plate.String(),
				)
			}
		})
	}
}

package select5_test

import (
	"fmt"
	"github.com/g1eng/select5"
	"strings"
	"testing"
)

var (
	pointedString          = "POINT"
	pInt          int      = 366
	pInt8         int8     = 36
	pInt16        int16    = 366
	pInt32        int32    = 366
	pInt64        int64    = 366
	pF32          float32  = 366.0
	pF64          float64  = 365.0051
	pB                     = true
	pNilInt       *int     = nil
	pNilInt8      *int8    = nil
	pNilInt16     *int16   = nil
	pNilInt32     *int32   = nil
	pNilInt64     *int64   = nil
	pNilF32       *float32 = nil
	pNilF64       *float64 = nil

	//
	pNilString *string = nil
	pNilB      *bool   = nil
)

func TestGetS(t *testing.T) {
	var (
		v      = "neko jiro"
		a  any = v
		a2     = &v
		b  any = 8342
	)

	w, err := select5.GetS(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetB got %v, want %v", w, v)
	}

	w, err = select5.GetS(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetB got %v, want %v", w, v)
	}

	w, err = select5.GetS(b)
	if err == nil {
		t.Errorf("GetB got %v, want error", w)
	}
}

func TestGetB(t *testing.T) {
	var (
		v      = true
		a  any = v
		a2     = &v
		b  any = "neko jealousy"
	)
	w, err := select5.GetB(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetB got %v, want %v", w, v)
	}

	w, err = select5.GetB(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetF32 got %v, want %v", w, v)
	}

	w, err = select5.GetB(b)
	if err == nil {
		t.Errorf("GetB got %v, want error", w)
	}
}

func TestGetF32(t *testing.T) {
	var (
		v  float32 = 3.14
		a  any     = v
		a2         = &v
		b  any     = "neko jealousy"
	)
	w, err := select5.GetF32(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetF32 got %v, want %v", w, v)
	}

	w, err = select5.GetF32(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetF32 got %v, want %v", w, v)
	}

	w, err = select5.GetF32(b)
	if err == nil {
		t.Errorf("GetF32 got %v, want error", w)
	}
}

func TestGetF64(t *testing.T) {
	var (
		v  float64 = 3.14
		a  any     = v
		a2         = &v
		b  any     = "neko jealousy"
	)
	w, err := select5.GetF64(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetF64 got %v, want %v", w, v)
	}

	w, err = select5.GetF64(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetF64 got %v, want %v", w, v)
	}

	w, err = select5.GetF64(b)
	if err == nil {
		t.Errorf("GetF64 got %v, want error", w)
	}
}

func TestGetI64(t *testing.T) {
	var (
		v  int64 = 1945
		a  any   = v
		a2       = &v
		b  any   = "neko jealousy"
	)
	w, err := select5.GetI64(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI64 got %v, want %v", w, v)
	}

	w, err = select5.GetI64(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetF64 got %v, want %v", w, v)
	}

	w, err = select5.GetI64(b)
	if err == nil {
		t.Errorf("GetI64 got %v, want error", w)
	}
}

func TestGetI32(t *testing.T) {
	var (
		v  int32 = 1945
		a  any   = v
		a2       = &v
		b  any   = "neko jealousy"
	)
	w, err := select5.GetI32(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI32 got %v, want %v", w, v)
	}

	w, err = select5.GetI32(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI32 got %v, want %v", w, v)
	}

	w, err = select5.GetI32(b)
	if err == nil {
		t.Errorf("GetI32 got %v, want error", w)
	}
}

func TestGetI16(t *testing.T) {
	var (
		v  int16 = 1945
		a  any   = v
		a2       = &v
		b  any   = "neko jealousy"
	)
	w, err := select5.GetI16(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI16 got %v, want %v", w, v)
	}

	w, err = select5.GetI16(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI16 got %v, want %v", w, v)
	}

	w, err = select5.GetI16(b)
	if err == nil {
		t.Errorf("GetI16 got %v, want error", w)
	}
}

func TestGetI8(t *testing.T) {
	var (
		v  int8 = 11
		a  any  = v
		a2      = &v
		b  any  = "neko jealousy"
	)
	w, err := select5.GetI8(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI8 got %v, want %v", w, v)
	}

	w, err = select5.GetI8(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI8 got %v, want %v", w, v)
	}

	w, err = select5.GetI8(b)
	if err == nil {
		t.Errorf("GetI8 got %v, want error", w)
	}
}

func TestGetI(t *testing.T) {
	var (
		v      = 1984
		a  any = v
		a2     = &v
		b  any = "neko jealousy"
	)
	w, err := select5.GetI(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI got %v, want %v", w, v)
	}

	w, err = select5.GetI(a2)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI got %v, want %v", w, v)
	}

	w, err = select5.GetI(b)
	if err == nil {
		t.Errorf("GetI got %v, want error", w)
	}
}

func TestCheckPrimitive(t *testing.T) {
	type args struct {
		s any
	}

	typesTT := []struct {
		name    string
		args    args
		wantRes byte
	}{
		{"int", args{1}, select5.IsInt},
		{"int8", args{int8(1)}, select5.IsInt8},
		{"int16", args{int16(1)}, select5.IsInt16},
		{"int32", args{int32(1)}, select5.IsInt32},
		{"int64", args{int64(1)}, select5.IsInt64},
		{"float32", args{float32(1.1)}, select5.IsFloat32},
		{"float64", args{float64(1.1)}, select5.IsFloat64},
		{"bool", args{true}, select5.IsBool},
		{"string", args{"a ok "}, select5.IsString},
		{"*string", args{&pointedString}, select5.IsString | select5.IsPointer},
		{"*int", args{&pInt}, select5.IsInt | select5.IsPointer},
		{"*int8", args{&pInt8}, select5.IsInt8 | select5.IsPointer},
		{"*int16", args{&pInt16}, select5.IsInt16 | select5.IsPointer},
		{"*int32", args{&pInt32}, select5.IsInt32 | select5.IsPointer},
		{"*int64", args{&pInt64}, select5.IsInt64 | select5.IsPointer},
		{"*float32", args{&pF32}, select5.IsFloat32 | select5.IsPointer},
		{"*float64", args{&pF64}, select5.IsFloat64 | select5.IsPointer},
		{"*bool", args{&pB}, select5.IsBool | select5.IsPointer},
		{"any", args{struct{}{}}, select5.IsAny},
	}
	for _, tt := range typesTT {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := select5.CheckPrimitive(tt.args.s); gotRes != tt.wantRes {
				t.Errorf("CheckPrimitive() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetV(t *testing.T) {

	type args struct {
		s any
	}
	typesTT := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"int", args{1}, "1"},
		{"int8", args{int8(12)}, "12"},
		{"int16", args{int16(1234)}, "1234"},
		{"int32", args{int32(365)}, "365"},
		{"int64", args{int64(123456789)}, "123456789"},
		{"float32", args{float32(1.1)}, fmt.Sprintf("%f", float32(1.1))},
		{"float64", args{float64(1.1)}, fmt.Sprintf("%f", float64(1.1))},
		{"bool", args{true}, "✓"},
		{"string", args{"a ok "}, "a ok "},
		{"*string", args{&pointedString}, pointedString},
		{"*int", args{&pInt}, fmt.Sprintf("%d", pInt)},
		{"*int8", args{&pInt8}, fmt.Sprintf("%d", pInt8)},
		{"*int16", args{&pInt16}, fmt.Sprintf("%d", pInt16)},
		{"*int32", args{&pInt32}, fmt.Sprintf("%d", pInt32)},
		{"*int64", args{&pInt64}, fmt.Sprintf("%d", pInt64)},
		{"*float32", args{&pF32}, fmt.Sprintf("%f", pF32)},
		{"*float64", args{&pF64}, fmt.Sprintf("%f", pF64)},
		{"*bool", args{&pB}, "✓"},
		{"any", args{struct{}{}}, ""},
	}
	for _, tt := range typesTT {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes, err := select5.GetV(tt.args.s); err != nil {
				if strings.TrimLeft(tt.name, "*") == tt.name && tt.name != "any" {
					t.Errorf("GetV() failed for %s, got %s: %v", tt.name, gotRes, err)
				}
			} else if gotRes != tt.wantRes {
				t.Errorf("GetV() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetVP(t *testing.T) {

	type args struct {
		s any
	}
	typesTT := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"int", args{1}, "1"},
		{"int8", args{int8(12)}, "12"},
		{"int16", args{int16(1234)}, "1234"},
		{"int32", args{int32(365)}, "365"},
		{"int64", args{int64(123456789)}, "123456789"},
		{"float32", args{float32(1.1)}, fmt.Sprintf("%f", float32(1.1))},
		{"float64", args{float64(1.1)}, fmt.Sprintf("%f", float64(1.1))},
		{"bool", args{true}, "✓"},
		{"string", args{"a ok "}, "a ok "},
		{"*string", args{&pointedString}, pointedString},
		{"*int", args{&pInt}, fmt.Sprintf("%d", pInt)},
		{"*int8", args{&pInt8}, fmt.Sprintf("%d", pInt8)},
		{"*int16", args{&pInt16}, fmt.Sprintf("%d", pInt16)},
		{"*int32", args{&pInt32}, fmt.Sprintf("%d", pInt32)},
		{"*int64", args{&pInt64}, fmt.Sprintf("%d", pInt64)},
		{"*float32", args{&pF32}, fmt.Sprintf("%f", pF32)},
		{"*float64", args{&pF64}, fmt.Sprintf("%f", pF64)},
		{"*bool", args{&pB}, "✓"},
		{"any", args{struct{}{}}, ""},
		{"*(nil)int", args{pNilInt}, ""},
		{"*(nil)int8", args{pNilInt8}, ""},
		{"*(nil)int16", args{pNilInt16}, ""},
		{"*(nil)int32", args{pNilInt32}, ""},
		{"*(nil)int64", args{pNilInt64}, ""},
		{"*(nil)float32", args{pNilF32}, ""},
		{"*(nil)float64", args{pNilF64}, ""},
		{"*(nil)bool", args{pNilB}, ""},        //printable
		{"*(nil)string", args{pNilString}, ""}, //printable
	}
	for _, tt := range typesTT {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes, err := select5.GetVP(tt.args.s); err != nil {
				if strings.TrimLeft(tt.name, "*") == tt.name && tt.name != "any" {
					t.Errorf("GetVP() failed for %s, got %s: %v", tt.name, gotRes, err)
				} else if tt.name == "*(nil)bool" || tt.name == "*(nil)string" {
					t.Errorf("GetVP() failed for %s, got %s: %v", tt.name, gotRes, err)
				}
			} else if gotRes != tt.wantRes {
				t.Errorf("GetVP() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

package select5_test

import (
	"github.com/g1eng/select5"
	"testing"
)

func TestGetX(t *testing.T) {
	var (
		v     = true
		a any = v
		b any = "neko jealousy"
	)
	w, err := select5.GetB(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetB got %v, want %v", w, v)
	}
	w, err = select5.GetB(b)
	if err == nil {
		t.Errorf("GetB got %v, want error", w)
	}
}

func TestGetF32(t *testing.T) {
	var (
		v float32 = 3.14
		a any     = v
		b any     = "neko jealousy"
	)
	w, err := select5.GetF32(a)
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
		v float64 = 3.14
		a any     = v
		b any     = "neko jealousy"
	)
	w, err := select5.GetF64(a)
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
		v int64 = 1945
		a any   = v
		b any   = "neko jealousy"
	)
	w, err := select5.GetI64(a)
	if err != nil {
		t.Fatal(err)
	}
	if w != v {
		t.Errorf("GetI64 got %v, want %v", w, v)
	}
	w, err = select5.GetI64(b)
	if err == nil {
		t.Errorf("GetI64 got %v, want error", w)
	}
}

func TestGetI32(t *testing.T) {
	var (
		v int32 = 1945
		a any   = v
		b any   = "neko jealousy"
	)
	w, err := select5.GetI32(a)
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
		v int16 = 1945
		a any   = v
		b any   = "neko jealousy"
	)
	w, err := select5.GetI16(a)
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
		v int8 = 11
		a any  = v
		b any  = "neko jealousy"
	)
	w, err := select5.GetI8(a)
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
		v     = 1945
		a any = v
		b any = "neko jealousy"
	)
	w, err := select5.GetI(a)
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

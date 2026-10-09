package model

import (
	"errors"
	"math"
	"testing"
)

func TestValue(t *testing.T) {
	decimal, err := NewDecimal(1.25)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		value *Value
		want  string
	}{
		{NewText(""), ""}, {NewUnsigned(0), "0"},
		{NewUnsigned(math.MaxUint64), "18446744073709551615"},
		{NewBoolean(false), "false"}, {decimal, "1.25"},
	} {
		got, err := tt.value.Format()
		if err != nil || got != tt.want {
			t.Fatalf("Format() = %q, %v; want %q", got, err, tt.want)
		}
		copy, err := tt.value.Clone()
		if err != nil {
			t.Fatal(err)
		}
		if copy == tt.value {
			t.Fatal("clone aliases value")
		}
		cloned, err := copy.Format()
		if err != nil || cloned != got {
			t.Fatal("clone differs")
		}
	}
}

func TestValueInvalid(t *testing.T) {
	x, y := "x", uint64(0)
	for _, value := range []*Value{nil, {}, {Text: &x, Unsigned: &y}} {
		if !errors.Is(value.Validate(), ErrInvalid) {
			t.Fatal("expected ErrInvalid")
		}
	}
	for _, x := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NewDecimal(x); !errors.Is(err, ErrInvalid) {
			t.Fatal("expected rejection of non-finite decimal")
		}
	}
}

func TestCloneIndependent(t *testing.T) {
	original := NewUnsigned(42)
	copy, err := original.Clone()
	if err != nil {
		t.Fatal(err)
	}
	*copy.Unsigned = 7
	if *original.Unsigned != 42 {
		t.Fatal("clone aliases original storage")
	}
}

package model

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// ErrInvalid marks an invalid model value, observation or snapshot.
var ErrInvalid = errors.New("invalid model")

// NewText constructs a text value, including an empty string.
func NewText(v string) *Value { return &Value{Text: &v} }

// NewUnsigned constructs an unsigned integer value, including zero.
func NewUnsigned(v uint64) *Value { return &Value{Unsigned: &v} }

// NewDecimal constructs a finite decimal value.
func NewDecimal(v float64) (*Value, error) {
	value := &Value{Decimal: &v}
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return value, nil
}

// NewBoolean constructs a boolean value, including false.
func NewBoolean(v bool) *Value { return &Value{Boolean: &v} }

// Validate requires exactly one active variant and rejects NaN and infinities.
func (v *Value) Validate() error {
	if v == nil {
		return fmt.Errorf("%w: missing value", ErrInvalid)
	}
	count := 0
	for _, present := range []bool{v.Text != nil, v.Unsigned != nil, v.Decimal != nil, v.Boolean != nil} {
		if present {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("%w: value has %d variants, want exactly one", ErrInvalid, count)
	}
	if v.Decimal != nil && (math.IsNaN(*v.Decimal) || math.IsInf(*v.Decimal, 0)) {
		return fmt.Errorf("%w: decimal must be finite", ErrInvalid)
	}
	return nil
}

// Format returns a locale-independent representation, not CSV escaping.
func (v *Value) Format() (string, error) {
	if err := v.Validate(); err != nil {
		return "", err
	}
	switch {
	case v.Text != nil:
		return *v.Text, nil
	case v.Unsigned != nil:
		return strconv.FormatUint(*v.Unsigned, 10), nil
	case v.Decimal != nil:
		return strconv.FormatFloat(*v.Decimal, 'g', -1, 64), nil
	default:
		return strconv.FormatBool(*v.Boolean), nil
	}
}

// Clone creates an independent copy of a validated value.
func (v *Value) Clone() (*Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	switch {
	case v.Text != nil:
		return NewText(*v.Text), nil
	case v.Unsigned != nil:
		return NewUnsigned(*v.Unsigned), nil
	case v.Decimal != nil:
		return NewDecimal(*v.Decimal)
	default:
		return NewBoolean(*v.Boolean), nil
	}
}

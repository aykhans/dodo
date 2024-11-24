package utils

import (
	"encoding/json"
	"errors"
)

type NonNilT interface {
	~int | ~float64 | ~string | ~bool
}

type Option[T NonNilT] interface {
	IsNone() bool
	ValueOrErr() (T, error)
	ValueOr(def T) T
	ValueOrPanic() T
	UnmarshalJSON(data []byte) error
}

// Don't call this struct directly, use NewOption[T] or NewNoneOption[T] instead.
type option[T NonNilT] struct {
	// value holds the actual value of the Option if it is not None.
	value T
	// none indicates whether the Option is None (i.e., has no value).
	none bool
}

func (o *option[T]) IsNone() bool {
	return o.none
}

// If the Option is None, it will return zero value of the type and an error.
func (o *option[T]) ValueOrErr() (T, error) {
	if o.IsNone() {
		return o.value, errors.New("Option is None")
	}
	return o.value, nil
}

// If the Option is None, it will return the default value.
func (o *option[T]) ValueOr(def T) T {
	if o.IsNone() {
		return def
	}
	return o.value
}

// If the Option is None, it will panic.
func (o *option[T]) ValueOrPanic() T {
	if o.IsNone() {
		panic("Option is None")
	}
	return o.value
}

func (o *option[T]) SetValue(value T) {
	o.value = value
	o.none = false
}

func (o *option[T]) SetNone() {
	var zeroValue T
	o.value = zeroValue
	o.none = true
}

func (o *option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.none = true
		return nil
	}
	o.none = false
	return json.Unmarshal(data, &o.value)
}

func NewOption[T NonNilT](value T) *option[T] {
	return &option[T]{value: value}
}

func NewNoneOption[T NonNilT]() *option[T] {
	return &option[T]{none: true}
}

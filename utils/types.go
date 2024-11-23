package utils

import (
	"encoding/json"
	"errors"
)

type NonNilConcrete interface {
	~int | ~float64 | ~string | ~bool
}

type Option[T NonNilConcrete] interface {
	IsNone() bool
	ValueOrErr() (*T, error)
	ValueOr(def *T) *T
	ValueOrPanic() *T
	UnmarshalJSON(data []byte) error
}

// Don't call this struct directly, use NewOption[T] or NewNoneOption[T] instead.
type option[T NonNilConcrete] struct {
	// value holds the actual value of the Option if it is not None.
	value T
	// none indicates whether the Option is None (i.e., has no value).
	none bool
}

func (o *option[T]) IsNone() bool {
	return o.none
}

// The returned value can be nil, if the Option is None, it will return nil and an error.
func (o *option[T]) ValueOrErr() (*T, error) {
	if o.IsNone() {
		return nil, errors.New("Option is None")
	}
	return &o.value, nil
}

// The returned value can't be nil, if the Option is None, it will return the default value.
func (o *option[T]) ValueOr(def *T) *T {
	if o.IsNone() {
		return def
	}
	return &o.value
}

// The returned value can't be nil, if the Option is None, it will panic.
func (o *option[T]) ValueOrPanic() *T {
	if o.IsNone() {
		panic("Option is None")
	}
	return &o.value
}

func (o *option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.none = true
		return nil
	}
	o.none = false
	return json.Unmarshal(data, &o.value)
}

func NewOption[T NonNilConcrete](value T) *option[T] {
	return &option[T]{value: value}
}

func NewNoneOption[T NonNilConcrete]() *option[T] {
	return &option[T]{none: true}
}

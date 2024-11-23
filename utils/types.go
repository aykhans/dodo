package utils

import (
	"encoding/json"
	"errors"
)

// Don't call this struct directly, use NewOption[T] or NewNoneOption[T] instead.
type Option[T any] struct {
	// value holds the actual value of the Option if it is not None.
	value T
	// none indicates whether the Option is None (i.e., has no value).
	none bool
}

func (o *Option[T]) IsNone() bool {
	return o.none
}

// The returned value can be nil, if the Option is None, it will return nil and an error.
func (o *Option[T]) ValueOrErr() (*T, error) {
	if o.IsNone() {
		return nil, errors.New("Option is None")
	}
	return &o.value, nil
}

// The returned value can't be nil, if the Option is None, it will return the default value.
func (o *Option[T]) ValueOr(def *T) *T {
	if o.IsNone() {
		return def
	}
	return &o.value
}

// The returned value can't be nil, if the Option is None, it will panic.
func (o *Option[T]) ValueOrPanic() *T {
	if o.IsNone() {
		panic("Option is None")
	}
	return &o.value
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.none = true
		return nil
	}
	o.none = false
	return json.Unmarshal(data, &o.value)
}

func NewOption[T any](value T) Option[T] {
	return Option[T]{value: value}
}

func NewNoneOption[T any]() Option[T] {
	return Option[T]{none: true}
}

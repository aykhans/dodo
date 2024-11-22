package utils

import (
	"encoding/json"
	"errors"
)

type Option[T any] struct {
	value T
	none  bool
}

func (o *Option[T]) IsNone() bool {
	return o.none
}

func (o *Option[T]) ValueOrErr() (*T, error) {
	if o.IsNone() {
		return nil, errors.New("Option is None")
	}
	return &o.value, nil
}

func (o *Option[T]) ValueOr(def *T) *T {
	if o.IsNone() {
		return def
	}
	return &o.value
}

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

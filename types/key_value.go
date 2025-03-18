package types

type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

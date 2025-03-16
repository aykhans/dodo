package types

import (
	"errors"
)

var (
	ErrInterrupt = errors.New("interrupted")
	ErrTimeout   = errors.New("timeout")
)

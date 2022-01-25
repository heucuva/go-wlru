package wlru

import (
	"github.com/pkg/errors"
)

var (
	// ErrNilArgument is returned when an argument is unexpectedly nil.
	ErrNilArgument = errors.New("nil argument")
)

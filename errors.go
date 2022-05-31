package applicator

import (
	"errors"
)

// ErrUnsupported is returned from an applicator that received an incorrect type.
// Like the trim receiving an integer.
var ErrUnsupported = errors.New("unsupported type encountered")

// ErrNotFound is returned when the tag specified an unknown applicator.
var ErrNotFound = errors.New("applicator name not found")

// ErrInvalidSet is returned when the applicator returned a value that doesn't
// match the field type.
var ErrInvalidSet = errors.New("cannot set new value with unmatching type")

// ErrCannotApply is returned when you call Apply on an incompatible type.
var ErrCannotApply = errors.New("cannot apply to invalid type")

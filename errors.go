package applicator

import (
	"errors"
)

var ErrUnsupported = errors.New("unsupported type encountered")

var ErrNotFound = errors.New("applicator name not found")

var ErrInvalidSet = errors.New("cannot set new value with unmatching type")

package helper

import (
	"errors"
)

var ErrUnsupported = errors.New("unsupported type encountered")

var ErrNotFound = errors.New("helper name not found")

var ErrInvalidSet = errors.New("cannot set new value with unmatching type")

package helper

import (
	"reflect"
	. "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrongType(t *T) {
	AddFunc("wrongtype", func(i interface{}, _ string) (interface{}, error) {
		return interface{}(""), nil
	})
	s := &struct {
		A uint `helper:"wrongtype"`
	}{
		A: 1,
	}
	err := Run(s)
	assert.Equal(t, ErrInvalidSet, err)
}

func TestDiffUintBytes(t *T) {
	AddFunc("diffuintbytes", func(i interface{}, _ string) (interface{}, error) {
		v := reflect.ValueOf(i)

		return interface{}(uint64(v.Uint())), nil
	})
	s := &struct {
		A uint `helper:"diffuintbytes"`
	}{
		A: 1,
	}
	err := Run(s)
	assert.Equal(t, ErrInvalidSet, err)
}

func TestMultiple(t *T) {
	s := &struct {
		A string `helper:"trim,lower"`
	}{
		A: " ABC ",
	}
	err := Run(s)
	require.Nil(t, err)
	assert.Equal(t, "abc", s.A)
}

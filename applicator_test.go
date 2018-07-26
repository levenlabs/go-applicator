package applicator

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
		A uint `apply:"wrongtype"`
	}{
		A: 1,
	}
	err := Apply(s)
	assert.Equal(t, ErrInvalidSet, err)
}

func TestDiffUintBytes(t *T) {
	AddFunc("diffuintbytes", func(i interface{}, _ string) (interface{}, error) {
		v := reflect.ValueOf(i)

		return interface{}(uint64(v.Uint())), nil
	})
	s := &struct {
		A uint `apply:"diffuintbytes"`
	}{
		A: 1,
	}
	err := Apply(s)
	assert.Equal(t, ErrInvalidSet, err)
}

func TestMultiple(t *T) {
	s := &struct {
		A string `apply:"trim,lower"`
	}{
		A: " ABC ",
	}
	err := Apply(s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)
}

func TestEmbedded(t *T) {
	type Inner struct {
		A string `apply:"lower"`
	}
	s := &struct {
		Inner
	}{
		Inner{"ABC"},
	}
	err := Apply(s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)

	s2 := &struct {
		*Inner
	}{
		&Inner{"ABC"},
	}
	err = Apply(s2)
	require.NoError(t, err)
	assert.Equal(t, "abc", s2.A)
}

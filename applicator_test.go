package applicator

import (
	"reflect"
	. "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStruct(t *T) {
	type B struct {
		A string `apply:"trim"`
	}
	s := B{" abc"}
	err := Apply(&s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)

	err = Apply(s)
	assert.Equal(t, ErrCannotApply, err)
}

func TestStructInterface(t *T) {
	type B struct {
		A string `apply:"trim"`
	}
	s := B{" abc"}
	var i interface{} = &s
	err := Apply(i)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)

	s = B{" abc"}
	i = s
	err = Apply(i)
	assert.Equal(t, ErrCannotApply, err)
}

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

func TestStructInStruct(t *T) {
	type B struct {
		A string `apply:"trim,lower"`
	}
	s := &struct {
		B *B
	}{
		B: &B{
			A: " ABC ",
		},
	}
	err := Apply(s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.B.A)
}

func TestSliceInStruct(t *T) {
	type B struct {
		A string `apply:"trim,lower"`
	}
	s := &struct {
		B []B
	}{
		B: []B{
			{
				A: " ABC ",
			},
		},
	}
	err := Apply(s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.B[0].A)
}

func TestSlice(t *T) {
	type A struct {
		S string `apply:"trim,lower"`
	}

	{
		s := []A{{" abc "}, {"ABC"}}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, []A{{"abc"}, {"abc"}}, s)
	}

	{
		s := []A{{" abc "}, {"ABC"}}
		err := Apply(&s)
		require.NoError(t, err)
		assert.Equal(t, []A{{"abc"}, {"abc"}}, s)
	}

	{
		s := []*A{{" abc "}, {"ABC"}}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, []*A{{"abc"}, {"abc"}}, s)
	}

	{
		s := []interface{}{&A{" abc "}, &A{"ABC"}}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, []interface{}{&A{"abc"}, &A{"abc"}}, s)
	}
}

func TestArray(t *T) {
	type A struct {
		S string `apply:"trim,lower"`
	}

	{
		s := [2]A{{" abc "}, {"ABC"}}
		err := Apply(s)
		assert.Equal(t, ErrCannotApply, err)
	}

	{
		s := [2]A{{" abc "}, {"ABC"}}
		err := Apply(&s)
		require.NoError(t, err)
		assert.Equal(t, [2]A{{"abc"}, {"abc"}}, s)
	}

	{
		s := [2]*A{{" abc "}, {"ABC"}}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, [2]*A{{"abc"}, {"abc"}}, s)
	}

	{
		s := [2]interface{}{&A{" abc "}, &A{"ABC"}}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, [2]interface{}{&A{"abc"}, &A{"abc"}}, s)
	}
}

func TestMap(t *T) {
	type A struct {
		S string `apply:"trim,lower"`
	}

	{
		s := map[int]A{
			1: {" abc "},
			2: {"ABC"},
		}
		err := Apply(s)
		assert.Equal(t, ErrCannotApply, err)
	}

	{
		s := map[int]A{
			1: {" abc "},
			2: {"ABC"},
		}
		err := Apply(&s)
		assert.Equal(t, ErrCannotApply, err)
	}

	{
		s := map[int]*A{
			1: {" abc "},
			2: {"ABC"},
		}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, map[int]*A{
			1: {"abc"},
			2: {"abc"},
		}, s)
	}

	{
		s := map[int]interface{}{
			1: &A{" abc "},
			2: &A{"ABC"},
		}
		err := Apply(s)
		require.NoError(t, err)
		assert.Equal(t, map[int]interface{}{
			1: &A{"abc"},
			2: &A{"abc"},
		}, s)
	}
}

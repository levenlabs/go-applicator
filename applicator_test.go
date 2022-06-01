package applicator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStruct(t *testing.T) {
	type B struct {
		A string `apply:"trim"`
	}
	s := B{" abc"}
	err := Apply(&s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)

	err = Apply(s)
	assert.Equal(t, ErrCannotApply, err)

	s = B{" abc"}
	ptr := &s
	err = Apply(ptr)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)

	s = B{" abc"}
	ptr = &s
	err = Apply(&ptr)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)
}

func TestStructInterface(t *testing.T) {
	type B struct {
		A string `apply:"trim"`
	}
	s := B{" abc"}
	var i interface{} = &s
	err := Apply(i)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)
	assert.Equal(t, &s, i.(*B))

	s = B{" abc"}
	i = &s
	err = Apply(&i)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)

	s = B{" abc"}
	i = s
	err = Apply(i)
	assert.Equal(t, ErrCannotApply, err)

	s = B{" abc"}
	i = s
	err = Apply(&i)
	require.NoError(t, err)
	assert.Equal(t, "abc", i.(B).A)

	s = B{" abc"}
	i = s
	var ii interface{} = &i
	err = Apply(ii)
	require.NoError(t, err)
	assert.Equal(t, "abc", i.(B).A)

	s = B{" abc"}
	i = s
	ii = &i
	err = Apply(&ii)
	require.NoError(t, err)
	assert.Equal(t, "abc", i.(B).A)
}

func TestWrongType(t *testing.T) {
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

func TestDiffUintBytes(t *testing.T) {
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

func TestMultiple(t *testing.T) {
	s := &struct {
		A string `apply:"trim,lower"`
	}{
		A: " ABC ",
	}
	err := Apply(s)
	require.NoError(t, err)
	assert.Equal(t, "abc", s.A)
}

func TestEmbedded(t *testing.T) {
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

func TestStructInStruct(t *testing.T) {
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

func TestSliceInStruct(t *testing.T) {
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

func TestSlice(t *testing.T) {
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

	{
		s := []interface{}{&A{" abc "}, &A{"ABC"}}
		var i interface{} = s
		err := Apply(i)
		require.NoError(t, err)
		assert.Equal(t, []interface{}{&A{"abc"}, &A{"abc"}}, s)
	}

	{
		s := []interface{}{&A{" abc "}, &A{"ABC"}}
		var i interface{} = s
		err := Apply(&i)
		require.NoError(t, err)
		assert.Equal(t, []interface{}{&A{"abc"}, &A{"abc"}}, s)
	}

	{
		s := []interface{}{&A{" abc "}, &A{"ABC"}}
		var i interface{} = &s
		err := Apply(i)
		require.NoError(t, err)
		assert.Equal(t, []interface{}{&A{"abc"}, &A{"abc"}}, s)
	}
}

func TestArray(t *testing.T) {
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

	{
		s := [2]A{{" abc "}, {"ABC"}}
		var i interface{} = &s
		err := Apply(i)
		require.NoError(t, err)
		assert.Equal(t, [2]A{{"abc"}, {"abc"}}, s)
	}

	{
		s := [2]A{{" abc "}, {"ABC"}}
		var i interface{} = s
		err := Apply(&i)
		require.NoError(t, err)
		assert.Equal(t, [2]A{{"abc"}, {"abc"}}, i.([2]A))
	}
}

func TestMap(t *testing.T) {
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

	{
		s := map[int]interface{}{
			1: &A{" abc "},
			2: &A{"ABC"},
		}
		var i interface{} = s
		err := Apply(i)
		require.NoError(t, err)
		assert.Equal(t, map[int]interface{}{
			1: &A{"abc"},
			2: &A{"abc"},
		}, s)
	}

	{
		s := map[int]interface{}{
			1: &A{" abc "},
			2: &A{"ABC"},
		}
		var i interface{} = s
		err := Apply(&i)
		require.NoError(t, err)
		assert.Equal(t, map[int]interface{}{
			1: &A{"abc"},
			2: &A{"abc"},
		}, s)
	}

	{
		s := map[int]interface{}{
			1: &A{" abc "},
			2: &A{"ABC"},
		}
		var i interface{} = &s
		err := Apply(i)
		require.NoError(t, err)
		assert.Equal(t, map[int]interface{}{
			1: &A{"abc"},
			2: &A{"abc"},
		}, s)
	}
}

func TestString(t *testing.T) {
	s := " abc"
	err := Apply(&s)
	assert.Equal(t, ErrCannotApply, err)

	err = Apply(s)
	assert.Equal(t, ErrCannotApply, err)

	var i interface{} = &s
	err = Apply(i)
	assert.Equal(t, ErrCannotApply, err)

	i = &s
	err = Apply(&i)
	assert.Equal(t, ErrCannotApply, err)

	i = s
	err = Apply(&i)
	assert.Equal(t, ErrCannotApply, err)
}

func TestNilInterface(t *testing.T) {
	var i interface{}
	err := Apply(i)
	assert.Equal(t, ErrCannotApply, err)

	err = Apply(&i)
	assert.Equal(t, ErrCannotApply, err)
}

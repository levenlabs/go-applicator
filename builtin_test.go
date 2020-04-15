package applicator

import (
	. "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrim(t *T) {
	s := &struct {
		A string `apply:"trim"`
	}{
		A: " 123 ",
	}
	err := Apply(s)
	require.Nil(t, err)
	assert.Equal(t, "123", s.A)

	str := " 234 "
	s2 := &struct {
		A *string `apply:"trim"`
	}{
		A: &str,
	}
	err = Apply(s2)
	require.Nil(t, err)
	assert.Equal(t, "234", *s2.A)

	s3 := &struct {
		A *string `apply:"trim"`
	}{}
	err = Apply(s3)
	require.Nil(t, err)
	assert.Nil(t, s3.A)
}

func TestLower(t *T) {
	s := &struct {
		A string `apply:"lower"`
	}{
		A: "AAA",
	}
	err := Apply(s)
	require.Nil(t, err)
	assert.Equal(t, "aaa", s.A)

	str := "BBB"
	s2 := &struct {
		A *string `apply:"lower"`
	}{
		A: &str,
	}
	err = Apply(s2)
	require.Nil(t, err)
	assert.Equal(t, "bbb", *s2.A)

	s3 := &struct {
		A *string `apply:"lower"`
	}{}
	err = Apply(s3)
	require.Nil(t, err)
	assert.Nil(t, s3.A)
}

func TestNonNil(t *T) {
	s := &struct {
		A []string `apply:"fillNil"`
	}{}
	err := Apply(s)
	require.Nil(t, err)
	require.NotNil(t, s.A)
	assert.Equal(t, []string{}, s.A)

	s2 := &struct {
		A *string `apply:"fillNil"`
	}{}
	err = Apply(s2)
	require.Nil(t, err)
	require.NotNil(t, s2.A)
	assert.Equal(t, "", *s2.A)

	str := "A"
	s2.A = &str
	err = Apply(s2)
	require.Nil(t, err)
	require.NotNil(t, s2.A)
	assert.Equal(t, &str, s2.A)
	assert.Equal(t, str, *s2.A)

	s3 := &struct {
		A map[string]string `apply:"fillNil"`
	}{}
	err = Apply(s3)
	require.Nil(t, err)
	require.NotNil(t, s3.A)
	assert.Equal(t, map[string]string{}, s3.A)
}

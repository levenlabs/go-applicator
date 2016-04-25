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
}

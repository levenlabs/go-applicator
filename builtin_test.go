package helper

import (
	. "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrim(t *T) {
	s := &struct {
		A string `helper:"trim"`
	}{
		A: " 123 ",
	}
	err := Run(s)
	require.Nil(t, err)
	assert.Equal(t, "123", s.A)

	str := " 234 "
	s2 := &struct {
		A *string `helper:"trim"`
	}{
		A: &str,
	}
	err = Run(s2)
	require.Nil(t, err)
	assert.Equal(t, "234", *s2.A)
}

func TestLower(t *T) {
	s := &struct {
		A string `helper:"lower"`
	}{
		A: "AAA",
	}
	err := Run(s)
	require.Nil(t, err)
	assert.Equal(t, "aaa", s.A)

	str := "BBB"
	s2 := &struct {
		A *string `helper:"lower"`
	}{
		A: &str,
	}
	err = Run(s2)
	require.Nil(t, err)
	assert.Equal(t, "bbb", *s2.A)
}

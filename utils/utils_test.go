package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test1Average(t *testing.T) {
	got := Average(1, 2)
	want := 1
	require.Equal(t, got, want)
}
func Test2Average(t *testing.T) {
	got := Average(-5, 5)
	want := 0
	require.Equal(t, got, want)
}

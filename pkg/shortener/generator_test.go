package shortener

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShorten(t *testing.T) {
	input := "some_long_url"
	expect := "PaGhXB2p"

	shorten, err := Shorten(input)
	require.NoError(t, err)
	require.Equal(t, expect, shorten)
}

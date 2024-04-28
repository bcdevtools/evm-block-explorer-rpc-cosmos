package utils

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnpackAbiString(t *testing.T) {
	bz, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b68656c6c6f20776f726c64000000000000000000000000000000000000000000")
	require.NoError(t, err)
	got, err := UnpackAbiString(bz, "")
	require.NoError(t, err)
	require.Equal(t, "hello world", got)
}

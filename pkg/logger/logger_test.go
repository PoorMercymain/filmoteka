package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	SetLogFile("abc.txt")

	str := "abc"
	n, err := Logger().Write([]byte(str))
	require.NoError(t, err)
	require.Equal(t, len(str), n)
}
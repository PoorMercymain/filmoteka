package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg := Config{}

	dsn := cfg.DSN()
	require.NotEmpty(t, dsn)
}

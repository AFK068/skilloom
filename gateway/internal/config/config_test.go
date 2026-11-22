package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AFK068/skilloom/gateway/internal/config"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("project config", func(t *testing.T) {
		t.Parallel()

		cfg, err := config.Load("../../config.yaml")
		require.NoError(t, err)
		assert.NotEmpty(t, cfg.Listen)
		assert.NotEmpty(t, cfg.Routes)
	})

	t.Run("partial config", func(t *testing.T) {
		t.Parallel()

		file := filepath.Join(t.TempDir(), "config.yaml")
		require.NoError(t, os.WriteFile(file, []byte("listen: ':8080'\nextra: true\n"), 0o600))

		cfg, err := config.Load(file)
		require.NoError(t, err)
		assert.Equal(t, ":8080", cfg.Listen)
		assert.Empty(t, cfg.Routes)
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
		assert.ErrorContains(t, err, "missing.yaml")
	})
}

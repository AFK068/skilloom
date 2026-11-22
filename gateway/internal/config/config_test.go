package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AFK068/skilloom/gateway/internal/config"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("project config", func(t *testing.T) {
		t.Parallel()

		cfg, err := config.Load("../../config.yaml")
		if err != nil {
			t.Fatal(err)
		}
		wantPrefixes := []string{"/openai/v1", "/anthropic/v1", "/kimi/v1", "/zai/v1", "/deepseek/v1"}
		if cfg.Listen == "" || len(cfg.Routes) != len(wantPrefixes) {
			t.Fatalf("unexpected config: %+v", cfg)
		}
		for i, prefix := range wantPrefixes {
			if route := cfg.Routes[i]; route.Prefix != prefix || route.Upstream == "" {
				t.Errorf("route = %+v, want prefix %q and a nonempty upstream", route, prefix)
			}
		}
	})

	t.Run("partial config", func(t *testing.T) {
		t.Parallel()

		file := filepath.Join(t.TempDir(), "config.yaml")
		if err := os.WriteFile(file, []byte("listen: ':8080'\nextra: true\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(file)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Listen != ":8080" || len(cfg.Routes) != 0 {
			t.Fatalf("unexpected config: %+v", cfg)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
		if err == nil || !strings.Contains(err.Error(), "missing.yaml") {
			t.Fatalf("error = %v, want missing file name", err)
		}
	})
}

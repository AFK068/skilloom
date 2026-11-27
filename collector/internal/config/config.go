package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}

func Load() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("get user config dir: %w", err)
	}

	var cfg Config
	file := filepath.Join(dir, "skilloom", "config.json")
	if err := cleanenv.ReadConfig(file, &cfg); err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", file, err)
	}

	cfg.normalize()

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config %q: %w", file, err)
	}

	return cfg, nil
}

func (c *Config) normalize() {
	c.Endpoint = strings.TrimSpace(c.Endpoint)
	c.Token = strings.TrimSpace(c.Token)
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Endpoint,
			validation.Required,
			is.RequestURL,
			validation.NewStringRule(func(s string) bool {
				u, err := url.Parse(s)
				if err != nil || u.Hostname() == "" {
					return false
				}
				return u.Scheme == "https" ||
					(u.Scheme == "http" && slices.Contains([]string{"localhost", "127.0.0.1", "::1"}, strings.ToLower(u.Hostname())))
			}, "must use HTTPS outside localhost"),
		),
		validation.Field(&c.Token, validation.Required),
	)
}

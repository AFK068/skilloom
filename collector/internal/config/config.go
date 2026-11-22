package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Token    string `mapstructure:"token"`
}

func Load() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("get user config dir: %w", err)
	}

	file := filepath.Join(dir, "skilloom", "config.json")

	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", file, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config %q: %w", file, err)
	}

	return cfg, cfg.validate()
}

func (cfg Config) validate() error {
	endpoint, err := url.Parse(cfg.Endpoint)
	if err != nil || endpoint.Hostname() == "" || strings.TrimSpace(cfg.Token) == "" {
		return errors.New("a valid Skilloom endpoint and token are required")
	}

	host := endpoint.Hostname()
	local := strings.EqualFold(host, "localhost") || net.ParseIP(host).IsLoopback()

	if endpoint.Scheme != "https" && (endpoint.Scheme != "http" || !local) {
		return errors.New("use HTTPS for Skilloom endpoints outside localhost")
	}

	return nil
}

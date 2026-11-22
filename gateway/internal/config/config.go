package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Listen string  `mapstructure:"listen"`
	Routes []Route `mapstructure:"routes"`
}

type Route struct {
	Prefix    string     `mapstructure:"prefix"`
	Upstream  string     `mapstructure:"upstream"`
	Endpoints []Endpoint `mapstructure:"endpoints"`
}

type Endpoint struct {
	Method string `mapstructure:"method"`
	Path   string `mapstructure:"path"`
}

func Load(file string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", file, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config %q: %w", file, err)
	}

	return cfg, nil
}

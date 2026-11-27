package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultConfigPath = "config.yaml"

type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	NATS     NATSConfig     `yaml:"nats"`
	Redis    RedisConfig    `yaml:"redis"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr" env:"HTTP_ADDR" env-default:":8080"`
}

type NATSConfig struct {
	URL string `yaml:"url" env:"NATS_URL"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"     env:"REDIS_ADDR"     env-default:"127.0.0.1:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db"       env:"REDIS_DB"       env-default:"0"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn" env:"POSTGRES_DSN"`
}

func Load() (Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultConfigPath
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	return cfg, nil
}

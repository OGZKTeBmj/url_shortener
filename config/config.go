package config

import (
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/adapter/postgres"
	"github.com/OGZKTeBmj/url_shortener/internal/adapter/redis"
	"github.com/OGZKTeBmj/url_shortener/pkg/httpserver"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Name    string
		Version string
		Env     string `default:"dev"`
		GuestTL time.Duration
	}
	Log struct {
		Level string `default:"info"`
	}
	HTTP     httpserver.Config
	Postgres postgres.Config
	Redis    redis.Config
}

func InitConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

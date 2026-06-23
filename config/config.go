package config

import (
	"os"

	"github.com/OGZKTeBmj/url_shortener/internal/adapter/postgres"
	"github.com/OGZKTeBmj/url_shortener/internal/adapter/redis"
	"github.com/OGZKTeBmj/url_shortener/internal/service/auth"
	"github.com/OGZKTeBmj/url_shortener/internal/service/shortener"
	"github.com/OGZKTeBmj/url_shortener/pkg/httpserver"
	"github.com/goccy/go-yaml"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App struct {
		Name    string
		Version string
		Env     string `default:"dev"`
	}
	Log struct {
		Level string `default:"info"`
	}
	HTTP     httpserver.Config
	Postgres postgres.Config
	Redis    redis.Config

	Auth      auth.AuthConfig
	Shortener shortener.ShortenerConfig
}

type YAMLConfig struct {
	Auth      auth.AuthConfig           `yaml:"auth"`
	Shortener shortener.ShortenerConfig `yaml:"shortener"`
}

func InitConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	var yamlCfg YAMLConfig
	if err := loadYAML("./config.yaml", &yamlCfg); err != nil {
		return Config{}, err
	}

	cfg.Auth.AccessTTL = yamlCfg.Auth.AccessTTL
	cfg.Auth.RefreshTTL = yamlCfg.Auth.RefreshTTL

	cfg.Shortener = yamlCfg.Shortener
	cfg.Auth.AppName = cfg.App.Name

	return cfg, nil
}

func loadYAML(path string, cfg any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

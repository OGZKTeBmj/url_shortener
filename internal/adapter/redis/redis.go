package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `envconfig:"HOST"     required:"true"`
	Port     string `envconfig:"PORT"     required:"true"`
	Password string `envconfig:"PASSWORD"`
	DB       int    `envconfig:"DB"       default:"0"`
}

type Client struct {
	client *redis.Client
	cfg    Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Connect(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", c.cfg.Host, c.cfg.Port)

	var rdb *redis.Client
	var err error

	for i := 1; i <= 5; i++ {
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: c.cfg.Password,
			DB:       c.cfg.DB,
		})
		attemptCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

		err = rdb.Ping(attemptCtx).Err()
		cancel()

		if err == nil {
			c.client = rdb
			return nil
		}
		_ = rdb.Close()

		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("redis connect failed after retries: %w", err)
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

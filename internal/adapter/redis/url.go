package redis

import (
	"context"
	"errors"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func (c *Client) Save(ctx context.Context, short string, url string) error {
	const op = "redis.url.Save"

	if err := c.client.Set(ctx, shortKey(short), url, time.Hour*24*30).Err(); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func shortKey(short string) string {
	return "short:" + short
}

func (c *Client) Get(ctx context.Context, short string) (string, error) {
	const op = "redis.url.Get"

	url, err := c.client.Get(ctx, short).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", utils.ErrWrap(op, domain.ErrEntityNotFound)
		}
		return "", utils.ErrWrap(op, err)
	}
	return url, nil
}

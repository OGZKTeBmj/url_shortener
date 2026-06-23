package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type UrlStorage struct {
	client *redis.Client
}

func NewUrlStorage(client *redis.Client) *UrlStorage {
	return &UrlStorage{client: client}
}

func (u *UrlStorage) Save(ctx context.Context, short string, url string, tl time.Duration) error {
	const op = "redis.url.Save"

	if err := u.client.Set(ctx, shortKey(short), url, tl).Err(); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (u *UrlStorage) Get(ctx context.Context, short string) (string, error) {
	const op = "redis.url.Get"

	url, err := u.client.Get(ctx, shortKey(short)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", utils.ErrWrap(op, domain.ErrEntityNotFound)
		}
		return "", utils.ErrWrap(op, err)
	}
	return url, nil
}

func (u *UrlStorage) Allow(ctx context.Context, key string, limit int64) error {
	const op = "redis.url.Allow"

	count, err := u.client.Get(ctx, rateLimitKey(key)).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return utils.ErrWrap(op, err)
	}

	if count >= limit {
		return utils.ErrWrap(op, domain.ErrRateLimitExceeded)
	}

	return nil
}

func (u *UrlStorage) Consume(ctx context.Context, key string, window time.Duration) error {
	const op = "redis.url.Consume"

	count, err := u.client.Incr(ctx, rateLimitKey(key)).Result()
	if err != nil {
		return utils.ErrWrap(op, err)
	}

	if count == 1 {
		if err := u.client.Expire(ctx, rateLimitKey(key), window).Err(); err != nil {
			return utils.ErrWrap(op, err)
		}
	}

	return nil
}

func (u *UrlStorage) VisitIncr(ctx context.Context, short string) error {
	const op = "redis.url.VisitIncr"

	_, err := u.client.Incr(ctx, visitKey(short)).Result()
	if err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (u *UrlStorage) VisitsDrain(ctx context.Context) (map[string]int64, error) {
	const op = "redis.url.VisitsDrain"

	keys, err := u.client.Keys(ctx, visitKey("*")).Result()
	if err != nil {
		return nil, utils.ErrWrap(op, err)
	}
	visits := make(map[string]int64)

	for _, key := range keys {
		val, err := u.client.GetDel(ctx, key).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return nil, utils.ErrWrap(op, err)
		}
		key = strings.TrimPrefix(key, visitKey(""))
		visits[key] = val
	}

	return visits, nil
}

func visitKey(key string) string {
	return "visits:" + key
}

func shortKey(short string) string {
	return "short:" + short
}

func rateLimitKey(key string) string {
	return "ratelimit:" + key
}

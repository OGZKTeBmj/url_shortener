package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenStorage struct {
	client *redis.Client
}

func NewRefreshTokenStorage(client *redis.Client) *RefreshTokenStorage {
	return &RefreshTokenStorage{client: client}
}

func (r *RefreshTokenStorage) Save(ctx context.Context, hash string, session domain.RefreshSession, tl time.Duration) error {
	const op = "redis.refreshToken.Save"

	data, err := json.Marshal(session)
	if err != nil {
		return utils.ErrWrap(op, err)
	}

	if err := r.client.Set(ctx, r.sessionKey(hash), data, tl).Err(); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (r *RefreshTokenStorage) Update(ctx context.Context, hash string, newHash string) error {
	const op = "redis.refreshToken.Update"

	err := r.client.Rename(ctx, r.sessionKey(hash), r.sessionKey(newHash)).Err()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return utils.ErrWrap(op, domain.ErrEntityNotFound)
		}
		return utils.ErrWrap(op, err)
	}

	return nil
}

func (r *RefreshTokenStorage) Delete(ctx context.Context, hash string) error {
	const op = "redis.refreshToken.Delete"

	val, err := r.client.Del(ctx, r.sessionKey(hash)).Result()
	if err != nil {
		return utils.ErrWrap(op, err)
	}
	if val == 0 {
		return utils.ErrWrap(op, domain.ErrEntityNotFound)
	}

	return nil
}

func (r *RefreshTokenStorage) GetSession(ctx context.Context, hash string) (*domain.RefreshSession, error) {
	const op = "redis.refreshToken.GetSession"

	val, err := r.client.Get(ctx, r.sessionKey(hash)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &domain.RefreshSession{}, utils.ErrWrap(op, domain.ErrEntityNotFound)
		}
		return &domain.RefreshSession{}, utils.ErrWrap(op, err)
	}

	var session domain.RefreshSession
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return &domain.RefreshSession{}, utils.ErrWrap(op, err)
	}

	return &session, nil
}

func (r *RefreshTokenStorage) sessionKey(hash string) string {
	return "seesion:" + hash
}

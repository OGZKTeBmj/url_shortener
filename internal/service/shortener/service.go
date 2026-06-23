package shortener

import (
	"context"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
)

type RateLimiter interface {
	Allow(сtx context.Context, key string, limit int64) error
	Consume(ctx context.Context, key string, window time.Duration) error
}

type StatisticProvider interface {
	VisitIncr(ctx context.Context, short string) error
}

type СacheURLProvider interface {
	Save(ctx context.Context, short string, url string, tl time.Duration) error
	Get(ctx context.Context, short string) (string, error)
}

type URLProvider interface {
	Save(ctx context.Context, short string, url string, userId domain.UUID, tl time.Duration) error
	Get(ctx context.Context, short string) (string, error)
	GetURLsByUserID(ctx context.Context, userId domain.UUID) ([]domain.UserURL, error)
}

type UserURLProvider interface {
}

type Shortener struct {
	cfg ShortenerConfig
	log logger.Logger

	mainURLProvider  URLProvider
	cacheURLProvider СacheURLProvider

	statisticProvider StatisticProvider
	rateLimiter       RateLimiter
}

type ShortenerConfig struct {
	GuestTTL    time.Duration `yaml:"guest_ttl"`
	GuestLimit  int64         `yaml:"guest_limit"`
	GuestWindow time.Duration `yaml:"guest_window"`
}

func NewShortener(
	cfg ShortenerConfig,
	log logger.Logger,
	mainURLProvider URLProvider,
	cacheURLProvider СacheURLProvider,
	statisticCounter StatisticProvider,
	rateLimiter RateLimiter,
) *Shortener {
	return &Shortener{
		cfg:               cfg,
		log:               log,
		mainURLProvider:   mainURLProvider,
		cacheURLProvider:  cacheURLProvider,
		statisticProvider: statisticCounter,
		rateLimiter:       rateLimiter,
	}
}

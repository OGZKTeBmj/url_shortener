package shortener

import (
	"context"
	"errors"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/internal/dto"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
)

type URLProvider interface {
	Save(ctx context.Context, short string, url string, tl time.Duration) error
	Get(ctx context.Context, short string) (string, error)
}

type Shortener struct {
	guestShortTL     time.Duration
	log              logger.Logger
	mainURLProvider  URLProvider
	cacheURLProvider URLProvider
}

func NewShortener(
	log logger.Logger,
	mainURLProvider URLProvider,
	cacheURLProvider URLProvider,
	guestShortTL time.Duration,
) *Shortener {
	return &Shortener{
		log:              log,
		mainURLProvider:  mainURLProvider,
		cacheURLProvider: cacheURLProvider,
		guestShortTL:     guestShortTL,
	}
}

func (s *Shortener) Short(ctx context.Context, input dto.ShortInput) (string, error) {
	const op = "shortenerService.Short"

	log := s.log.With("op", op, "url", input.URL)

	short, err := utils.GenerateShortCode()
	if err != nil {
		log.Error("generate short", "error", err)
		return "", utils.ErrWrap(op, err)
	}

	if err := s.mainURLProvider.Save(ctx, short, input.URL, s.guestShortTL); err != nil {
		log.Error("save url", "provider", "main", "error", err)
		return "", utils.ErrWrap(op, err)
	}

	if err := s.cacheURLProvider.Save(ctx, short, input.URL, s.guestShortTL); err != nil {
		log.Error("save url", "provider", "cache", "error", err)
	} else {
		log.Debug("save url", "provider", "cache")
	}

	return short, nil
}

func (s *Shortener) GetUrl(ctx context.Context, short string) (string, error) {
	const op = "shortenerService.GetUrl"

	log := s.log.With("op", op, "short", short)

	url, err := s.cacheURLProvider.Get(ctx, short)
	if err == nil {
		log.Debug("get url", "provider", "cache")
		return url, nil
	}

	if !errors.Is(err, domain.ErrEntityNotFound) {
		log.Error("get url", "provider", "cache")
	}

	url, err = s.mainURLProvider.Get(ctx, short)
	if err != nil {
		if !errors.Is(err, domain.ErrEntityNotFound) {
			log.Error("get url", "provider", "main", "error", err)
		}
		log.Debug("get url", "provider", "main", "info", "url not found")
		return "", utils.ErrWrap(op, err)
	}
	return url, nil
}

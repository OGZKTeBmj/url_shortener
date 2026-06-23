package shortener

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/internal/dto"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
)

func (s *Shortener) Short(ctx context.Context, input dto.ShortInput) (string, error) {
	const op = "shortenerService.Short"

	log := s.log.With("op", op, "url", input.URL)

	if input.IsGuest {
		if err := s.rateLimiter.Allow(ctx, input.IP, s.cfg.GuestLimit); err != nil {
			if errors.Is(err, domain.ErrRateLimitExceeded) {
				log.Debug("rate limit allow", "message", "limit exceeded")
				return "", utils.ErrWrap(op, err)
			}
			log.Error("rate limit allow", "error", err)
			return "", utils.ErrWrap(op, err)
		}
	}

	short, err := utils.GenerateShortCode()
	if err != nil {
		log.Error("generate short", "error", err)
		return "", utils.ErrWrap(op, err)
	}

	var ttl time.Duration
	if input.IsGuest {
		ttl = s.cfg.GuestTTL
	} else {
		ttl = 0
	}

	if err := s.mainURLProvider.Save(ctx, short, input.URL, input.UserID, ttl); err != nil {
		log.Error("save url", "provider", "main", "error", err)
		return "", utils.ErrWrap(op, err)
	}

	if err := s.cacheURLProvider.Save(ctx, short, input.URL, ttl); err != nil {
		log.Error("save url", "provider", "cache", "error", err)
	} else {
		log.Debug("save url", "provider", "cache")
	}

	if err := s.rateLimiter.Consume(ctx, input.IP, s.cfg.GuestWindow); err != nil {
		log.Error("rate limit consume", "error", err)
	}

	return short, nil
}

func (s *Shortener) VisitUrl(ctx context.Context, short string) (string, error) {
	const op = "shortenerService.VisitUrl"

	log := s.log.With("op", op, "short", short)

	url, err := s.cacheURLProvider.Get(ctx, short)
	if err == nil {
		log.Debug("get url", "provider", "cache")
	} else {
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
	}
	if err := s.statisticProvider.VisitIncr(ctx, short); err != nil {
		log.Error("add visit", "error", err)
	}

	return url, nil
}

func (s *Shortener) Urls(ctx context.Context, userId domain.UUID) ([]domain.UserURL, error) {
	const op = "shortenerService.Urls"

	log := s.log.With("op", op)

	urls, err := s.mainURLProvider.GetURLsByUserID(ctx, userId)
	if err != nil {
		log.Error("get urls", "error", err)
		return nil, utils.ErrWrap(op, err)
	}
	fmt.Println(urls)
	return urls, nil
}

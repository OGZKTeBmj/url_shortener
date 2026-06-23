package job

import (
	"context"
	"time"

	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
)

type Cleaner interface {
	DeleteExpired(context.Context) error
}

type CleanupJob struct {
	cleanUpInt time.Duration
	cleaner    Cleaner
	log        logger.Logger
}

func NewCleanupJob(log logger.Logger, cleanUpInt time.Duration, cleaner Cleaner) *CleanupJob {
	return &CleanupJob{
		log:        log,
		cleanUpInt: cleanUpInt,
		cleaner:    cleaner,
	}
}

func (c *CleanupJob) Run(ctx context.Context) {
	const op = "cleanUpJob.Run"

	log := c.log.With("op", op)

	ticker := time.NewTicker(c.cleanUpInt)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.cleaner.DeleteExpired(ctx); err != nil {
				log.Error("error", err)
			}
		}
	}
}

type Drainer interface {
	VisitsDrain(context.Context) (map[string]int64, error)
}

type VisitsProvider interface {
	AddVisitsBatch(ctx context.Context, visits map[string]int64) error
}

type VisitsDrainerJob struct {
	interval       time.Duration
	drainer        Drainer
	visitsProvider VisitsProvider
	log            logger.Logger
}

func NewVisitsDrainerJob(
	log logger.Logger, interval time.Duration,
	drainer Drainer, visitsProvider VisitsProvider) *VisitsDrainerJob {
	return &VisitsDrainerJob{
		log: log, interval: interval,
		drainer: drainer, visitsProvider: visitsProvider,
	}
}

func (v *VisitsDrainerJob) Run(ctx context.Context) {
	const op = "visitsDrainerJob.Run"

	log := v.log.With("op", op)

	ticker := time.NewTicker(v.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			visits, err := v.drainer.VisitsDrain(ctx)
			if err != nil {
				log.Error("drain", "error", err)
				continue
			}
			if err := v.visitsProvider.AddVisitsBatch(ctx, visits); err != nil {
				log.Error("add visits", "error", err)
			}
		}
	}
}

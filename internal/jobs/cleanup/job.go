package cleanup

import (
	"context"
	"time"

	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
)

type Cleaner interface {
	DeleteExpired(context.Context) error
}

type Job struct {
	cleanUpInt time.Duration
	cleaner    Cleaner
	log        logger.Logger
}

func NewJob(log logger.Logger, cleanUpInt time.Duration, cleaner Cleaner) *Job {
	return &Job{
		log:        log,
		cleanUpInt: cleanUpInt,
		cleaner:    cleaner,
	}
}

func (j *Job) Run(ctx context.Context) {
	const op = "cleanUpJob.Run"

	log := j.log.With("op", op)

	ticker := time.NewTicker(j.cleanUpInt)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := j.cleaner.DeleteExpired(ctx); err != nil {
				log.Error("error", err)
			}
		}
	}
}

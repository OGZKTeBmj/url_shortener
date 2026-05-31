package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlStorage struct {
	db *pgxpool.Pool
}

func NewUrlStorage(db *pgxpool.Pool) *UrlStorage {
	return &UrlStorage{
		db: db,
	}
}

func (u *UrlStorage) Save(ctx context.Context, short string, url string, tl time.Duration) error {
	const op = "postgres.urlStorage.Save"

	expiresAt := time.Now().Add(tl)

	if _, err := u.db.Exec(ctx, URLSaveQuery, short, url, expiresAt); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (u *UrlStorage) Get(ctx context.Context, short string) (string, error) {
	const op = "postgres.urlStorage.Get"

	var url string

	if err := u.db.QueryRow(ctx, URLGetQuery, short).Scan(&url); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", utils.ErrWrap(op, domain.ErrEntityNotFound)
		}
		return "", utils.ErrWrap(op, err)
	}
	return url, nil
}

func (u *UrlStorage) DeleteExpired(ctx context.Context) error {
	const op = "postgres.urlStorage.DeleteExpired"

	if _, err := u.db.Exec(ctx, DeleteExpiredQuery); err != nil {
		return utils.ErrWrap(op, err)
	}
	return nil
}

const (
	URLSaveQuery = `
	INSERT INTO short_url (short, url, expires_at)
	VALUES ($1, $2, $3);
	`
	URLGetQuery = `
	SELECT url FROM short_url
	WHERE short = $1
	AND (
		expires_at IS NULL
		OR expires_at > NOW()
	);
	`
	DeleteExpiredQuery = `
	DELETE FROM short_url
	WHERE expires_at IS NOT NULL
	AND expires_at < NOW(); 
	`
)

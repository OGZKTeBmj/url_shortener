package postgres

import (
	"context"
	"errors"
	"fmt"
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

func (u *UrlStorage) Save(
	ctx context.Context,
	short string,
	url string,
	userId domain.UUID,
	tl time.Duration) error {
	const op = "postgres.urlStorage.Save"

	var expiresAt *time.Time

	if tl > 0 {
		t := time.Now().Add(tl)
		expiresAt = &t
	}

	if _, err := u.db.Exec(ctx, URLSaveQuery, short, url, userId, expiresAt); err != nil {
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

func (u *UrlStorage) AddVisitsBatch(ctx context.Context, visits map[string]int64) error {
	const op = "postgres.urlStorage.AddVisitsBatch"

	if len(visits) == 0 {
		return nil
	}

	shorts := make([]string, 0, len(visits))
	vs := make([]int64, 0, len(visits))

	for short, visit := range visits {
		shorts = append(shorts, short)
		vs = append(vs, visit)
	}

	_, err := u.db.Exec(ctx, AddVisitsBatchQuery, shorts, vs)
	if err != nil {
		return utils.ErrWrap(op, err)
	}

	return nil
}

func (u *UrlStorage) GetURLsByUserID(ctx context.Context, userId domain.UUID) ([]domain.UserURL, error) {
	const op = "postgres.urlStorage.GetURLsByUserID"

	urls := make([]domain.UserURL, 0)

	rows, err := u.db.Query(ctx, GetURLsByUserIDQuery, userId)
	if err != nil {
		return nil, utils.ErrWrap(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var url domain.UserURL
		if err := rows.Scan(
			&url.Short,
			&url.URL,
			&url.Visits,
		); err != nil {
			return nil, utils.ErrWrap(op, err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, utils.ErrWrap(op, err)
	}
	fmt.Println(urls)
	return urls, nil
}

const (
	URLSaveQuery = `
	INSERT INTO short_url (short, url, user_id, expires_at)
	VALUES ($1, $2, $3, $4);
	`
	URLGetQuery = `
	SELECT url FROM short_url
	WHERE short = $1
	AND (
		expires_at IS NULL
		OR expires_at > NOW()
	);
	`
	GetURLsByUserIDQuery = `
	SELECT su.short, su.url, su.visits
	FROM short_url su
	JOIN users u ON u.id = su.user_id
	WHERE u.id = $1
  	AND (
    	su.expires_at IS NULL
    	OR su.expires_at > NOW()
  	);
	`
	DeleteExpiredQuery = `
	DELETE FROM short_url
	WHERE expires_at IS NOT NULL
	AND expires_at < NOW(); 
	`
	AddVisitsBatchQuery = `
		UPDATE short_url
		SET visits = short_url.visits + data.visits
		FROM (
			SELECT *
			FROM unnest($1::text[], $2::bigint[]) AS t(short, visits)
		) data
		WHERE short_url.short = data.short;
	`
)

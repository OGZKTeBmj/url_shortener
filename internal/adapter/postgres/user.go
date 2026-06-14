package postgres

import (
	"context"
	"errors"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) UserByName(ctx context.Context, name string) (*domain.User, error) {
	const op = "postgres.UserByName"

	var user domain.User

	if err := u.db.QueryRow(ctx, QueryUserByName, name).Scan(&user.UUID, &user.Name,
		&user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.User{}, domain.ErrEntityNotFound
		}
		return &domain.User{}, utils.ErrWrap(op, err)
	}
	return &user, nil
}

func (u *UserStorage) UserByID(ctx context.Context, id domain.UUID) (*domain.User, error) {
	const op = "postgres.UserByID"

	var user domain.User

	if err := u.db.QueryRow(ctx, QueryUserByID, id).Scan(&user.UUID, &user.Name,
		&user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.User{}, domain.ErrEntityNotFound
		}
		return &domain.User{}, utils.ErrWrap(op, err)
	}
	return &user, nil
}

func (u *UserStorage) SaveUser(ctx context.Context, name string, passwordhash []byte) (domain.UUID, error) {
	const op = "postgres.SaveUser"

	var id domain.UUID

	err := u.db.QueryRow(ctx, QuerySaveUser, name, passwordhash).
		Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.UUID{}, domain.ErrEntityAlreadyExists
			}
		}
		return domain.UUID{}, utils.ErrWrap(op, err)
	}
	return id, nil
}

const (
	QueryUserByName = `
	SELECT id, name, pass_hash FROM users
	WHERE name = $1
	`
	QueryUserByID = `
	SELECT id, name, pass_hash FROM users
	WHERE id = $1
	`
	QuerySaveUser = `
	INSERT INTO users (name, pass_hash)
	VALUES ($1, $2)
	RETURNING id
	`
)

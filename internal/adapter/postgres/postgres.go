package postgres

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string `envconfig:"USER"     required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Port     string `envconfig:"PORT"     required:"true"`
	Host     string `envconfig:"HOST"     required:"true"`
	DBName   string `envconfig:"DB_NAME"  required:"true"`
}

type Pool struct {
	pool *pgxpool.Pool
	cfg  Config
}

func New(cfg Config) *Pool {
	return &Pool{
		cfg: cfg,
	}
}

func (p *Pool) DB() *pgxpool.Pool {
	return p.pool
}

func (p *Pool) Connect(ctx context.Context) error {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(p.cfg.User),
		url.QueryEscape(p.cfg.Password),
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.DBName,
	)
	var err error

	for i := 1; i <= 5; i++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		pool, err := pgxpool.New(attemptCtx, connString)
		if err != nil {
			cancel()
			time.Sleep(2 * time.Second)
			continue
		}
		err = pool.Ping(attemptCtx)
		cancel()

		if err == nil {
			p.pool = pool
			return nil
		}
		pool.Close()
		time.Sleep(2 * time.Second)
	}
	return utils.ErrWrap(
		"failed to connect to postgres after retries",
		err,
	)
}

func (p *Pool) Close() {
	p.pool.Close()
}

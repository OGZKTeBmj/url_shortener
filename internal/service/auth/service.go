package auth

import (
	"context"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type UserRepository interface {
	UserByName(ctx context.Context, name string) (*domain.User, error)
	UserByID(ctx context.Context, id domain.UUID) (*domain.User, error)
	SaveUser(ctx context.Context, name string, passwordHash []byte) (domain.UUID, error)
}

type RTokenRepository interface {
	Save(ctx context.Context, hash string, session domain.RefreshSession, tl time.Duration) error
	Update(ctx context.Context, hash string, newHash string) error
	Delete(ctx context.Context, hash string) error
	GetSession(ctx context.Context, hash string) (*domain.RefreshSession, error)
}

type AuthConfig struct {
	AppName       string
	TTL           time.Duration
	JWTSecret     []byte
	RefreshSecret []byte
	RefreshTTL    time.Duration
}

type Auth struct {
	cfg              AuthConfig
	log              logger.Logger
	userRepository   UserRepository
	rTokenRepository RTokenRepository
}

func NewAuth(log logger.Logger, userRep UserRepository, rtRep RTokenRepository, cfg AuthConfig) *Auth {
	return &Auth{
		log:              log,
		userRepository:   userRep,
		rTokenRepository: rtRep,
		cfg:              cfg,
	}
}

type userClaims struct {
	Username string
	UserID   domain.UUID
	jwt.RegisteredClaims
}

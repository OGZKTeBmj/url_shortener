package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/internal/dto"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (a *Auth) Register(ctx context.Context, input dto.UserInput) (domain.UUID, error) {
	const op string = "authService.Register"

	log := a.log.With(
		"op", op,
		"user_name", input.Name,
	)

	log.Debug("Attempting to register user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "error", err)
		return domain.UUID{}, utils.ErrWrap(op, err)
	}

	id, err := a.userRepository.SaveUser(ctx, input.Name, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrEntityAlreadyExists) {
			log.Info("user already exists")
			return domain.UUID{}, utils.ErrWrap(op, domain.ErrEntityAlreadyExists)
		}

		log.Error("failed to save user", "error", err)
		return domain.UUID{}, utils.ErrWrap(op, err)
	}

	log.Info("user registered", "user_id", id)

	return id, nil
}

func (a *Auth) Login(ctx context.Context, input dto.UserInput) (string, string, error) {
	const op = "authService.Login"

	log := a.log.With(
		"op", op,
		"user_name", input.Name,
	)

	user, err := a.userRepository.UserByName(ctx, input.Name)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			log.Info("invalid credentials")
			return "", "", utils.ErrWrap(op, domain.ErrInvalidCredentails)
		}

		log.Error("failed to get user", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(input.Password)); err != nil {
		log.Info("invalid credentials")
		return "", "", utils.ErrWrap(op, domain.ErrInvalidCredentails)
	}

	accessToken, err := a.newToken(user)
	if err != nil {
		log.Error("failed to generate access token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Error("failed to generate refresh token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	session := domain.RefreshSession{
		UserID: user.UUID,
	}
	hashRefreshToken := utils.HashRefreshToken(refreshToken, a.cfg.RefreshSecret)
	err = a.rTokenRepository.Save(ctx, hashRefreshToken, session, a.cfg.RefreshTTL)
	if err != nil {
		log.Error("failed to save refresh token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	log.Info("user logged in", "user_id", user.UUID)

	return accessToken, refreshToken, nil
}

func (a *Auth) Logout(ctx context.Context, refreshToken string) error {
	const op = "authService.Logout"

	log := a.log.With("op", op)

	hash := utils.HashRefreshToken(refreshToken, a.cfg.RefreshSecret)
	err := a.rTokenRepository.Delete(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			log.Info("invalid token")
			return nil
		}
		log.Error("failed delete refresh token", "error", err)
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (a *Auth) Refresh(ctx context.Context, token string) (string, string, error) {
	const op = "authService.Refresh"

	log := a.log.With("op", op)

	hash := utils.HashRefreshToken(token, a.cfg.RefreshSecret)

	session, err := a.rTokenRepository.GetSession(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			log.Info("invalid refresh token")
			return "", "", utils.ErrWrap(op, domain.ErrInvalidCredentails)
		}

		log.Error("failed to get refresh token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	user, err := a.userRepository.UserByID(ctx, session.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	accessToken, err := a.newToken(user)
	if err != nil {
		log.Error("failed to generate access token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Error("failed to generate refresh token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	newHash := utils.HashRefreshToken(newRefreshToken, a.cfg.RefreshSecret)

	err = a.rTokenRepository.Update(ctx, hash, newHash)
	if err != nil {
		log.Error("failed to rotate refresh token", "error", err)
		return "", "", utils.ErrWrap(op, err)
	}

	log.Info("refresh success", "user_id", user.UUID)

	return accessToken, newRefreshToken, nil
}

func (a *Auth) parseToken(token string) (domain.UUID, error) {
	const op = "authService.parseToken"

	jwtToken, err := jwt.ParseWithClaims(token, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return a.cfg.JWTSecret, nil
	})

	if err != nil {
		return domain.UUID{}, utils.ErrWrap(op, err)
	}

	claims, ok := jwtToken.Claims.(*userClaims)
	if !ok || !jwtToken.Valid {
		return domain.UUID{}, utils.ErrWrap(op, domain.ErrInvalidToken)
	}

	return claims.UserID, nil
}

func (a *Auth) newToken(user *domain.User) (string, error) {
	const op = "authService.newToken"

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		Username: user.Name,
		UserID:   user.UUID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(user.UUID),
			Issuer:    a.cfg.AppName,
			ExpiresAt: jwt.NewNumericDate(now.Add(a.cfg.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	tokenString, err := token.SignedString(a.cfg.JWTSecret)
	if err != nil {
		return "", utils.ErrWrap(op, err)
	}

	return tokenString, nil
}

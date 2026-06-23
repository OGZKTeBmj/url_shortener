package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHead = "Authorization"
	ctxUserIdKey      = "user_id"
	ctxIsGuestKey     = "is_guest"
)

func MiddlewareLogger(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		log.Info("http request",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", ctx.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"ip", ctx.ClientIP(),
			"user_agent", ctx.Request.UserAgent(),
		)
	}
}

func (r *Router) MidddlewareAuthRequest(ctx *gin.Context) {
	token, err := extractToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	id, err := r.authService.ParseToken(token)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.Set(ctxUserIdKey, id)
	ctx.Set(ctxIsGuestKey, false)
}

func (r *Router) MiddlewareUserIdentity(ctx *gin.Context) {
	token, err := extractToken(ctx)
	if err == nil {
		id, err := r.authService.ParseToken(token)
		if err == nil {
			ctx.Set(ctxIsGuestKey, false)
			ctx.Set(ctxUserIdKey, id)
			return
		}
	}
	ctx.Set(ctxIsGuestKey, true)
}

func extractToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader(AuthorizationHead)
	if header == "" {
		return "", fmt.Errorf("empty auth header")
	}

	headersParts := strings.Split(header, " ")
	if len(headersParts) != 2 {
		return "", fmt.Errorf("invalid auth header")
	}
	return headersParts[1], nil
}

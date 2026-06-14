package http

import (
	"net/http"
	"time"

	"github.com/OGZKTeBmj/url_shortener/internal/service/auth"
	"github.com/OGZKTeBmj/url_shortener/internal/service/shortener"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	shortenerService *shortener.Shortener
	authService      *auth.Auth
}

func New(log logger.Logger, shortenerService *shortener.Shortener, authService *auth.Auth) *Router {
	engine := gin.New()

	engine.Use(MiddlewareLogger(log))

	router := &Router{
		engine:           engine,
		shortenerService: shortenerService,
		authService:      authService,
	}

	router.engine.GET("/:short", router.ShortRedirect)
	api := engine.Group("/api")
	{
		api.POST("/short", router.Short)
	}
	auth := engine.Group("/auth")
	{
		auth.POST("/sign-up", router.SignUp)
		auth.POST("/sign-in", router.SignIn)
		auth.POST("/refresh", router.Refresh)
		auth.POST("/logout", router.Logout)
	}
	return router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

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

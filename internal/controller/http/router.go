package http

import (
	"net/http"
	"time"

	service "github.com/OGZKTeBmj/url_shortener/internal/service/shortener"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	shortenerService *service.Shortener
}

func New(log logger.Logger, shortenerService *service.Shortener) *Router {
	engine := gin.New()

	engine.Use(MiddlewareLogger(log))

	router := &Router{
		engine:           engine,
		shortenerService: shortenerService,
	}

	router.engine.GET("/:short", router.ShortRedirect)
	api := engine.Group("/api")
	{
		api.POST("/short", router.Short)
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

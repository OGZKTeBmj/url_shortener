package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/internal/dto"
	"github.com/gin-gonic/gin"
)

func (r *Router) Short(ctx *gin.Context) {
	var input dto.ShortInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if err := input.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	input.IP = ctx.ClientIP()
	input.IsGuest = ctx.GetBool(ctxIsGuestKey)

	if !input.IsGuest {
		value, _ := ctx.Get(ctxUserIdKey)
		input.UserID, _ = value.(domain.UUID)
	}

	short, err := r.shortenerService.Short(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrRateLimitExceeded) {
			ctx.JSON(http.StatusTooManyRequests,
				gin.H{"message": "short URL creation limit exceeded. " +
					"Please try again later — the limit will reset soon"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"short": short})
}

func (r *Router) ShortRedirect(ctx *gin.Context) {
	short := ctx.Param("short")

	url, err := r.shortenerService.VisitUrl(ctx.Request.Context(), short)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "short code doesn't exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.Redirect(http.StatusFound, url)
}

func (r *Router) Urls(ctx *gin.Context) {
	value, exists := ctx.Get(ctxUserIdKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "user unauthorized"})
		return
	}

	user_id, ok := value.(domain.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	urls, err := r.shortenerService.Urls(ctx.Request.Context(), user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	fmt.Println(urls)

	response := make([]dto.UserURLOutput, len(urls))
	for i, url := range urls {
		response[i] = dto.UserURLOutput{
			Short:  url.Short,
			URL:    url.URL,
			Visits: url.Visits,
		}
	}
	ctx.JSON(http.StatusOK, response)
}

package http

import (
	"errors"
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
	input.IsGuest = true

	short, err := r.shortenerService.Short(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"short": short})
}

func (r *Router) ShortRedirect(ctx *gin.Context) {
	short := ctx.Param("short")

	url, err := r.shortenerService.GetUrl(ctx.Request.Context(), short)
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

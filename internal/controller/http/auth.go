package http

import (
	"errors"
	"net/http"

	"github.com/OGZKTeBmj/url_shortener/internal/domain"
	"github.com/OGZKTeBmj/url_shortener/internal/dto"
	"github.com/gin-gonic/gin"
)

func (r *Router) SignUp(ctx *gin.Context) {
	var userInput dto.UserInput
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := userInput.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	uuid, err := r.authService.Register(ctx.Request.Context(), userInput)
	if err != nil {
		if errors.Is(err, domain.ErrEntityAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "user is exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": uuid.String()})
}

func (r *Router) SignIn(ctx *gin.Context) {
	var request dto.UserInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if (len(request.Name) < 3) || (len(request.Password) < 8) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username lenght < 3 or password lenght < 8"})
		return
	}

	accessToken, refreshToken, err := r.authService.Login(ctx.Request.Context(), request)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentails) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credentials",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (r *Router) Logout(ctx *gin.Context) {
	var req dto.RefreshInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	err := r.authService.Logout(ctx.Request.Context(), req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}

func (r *Router) Refresh(ctx *gin.Context) {
	var request dto.RefreshInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	accessToken, refreshToken, err := r.authService.Refresh(ctx.Request.Context(), request.Token)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentails) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credentials",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

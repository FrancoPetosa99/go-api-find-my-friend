package controllers

import (
	"go-api-find-my-friend/internal/services"
	"go-api-find-my-friend/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidCredentials = errors.NewUnauthorizedError("invalid credentials")
	ErrInvalidBody        = errors.NewBadRequestError("invalid body")
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{authService: services.NewAuthService()}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var dto UserLoginDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidBody)
		return
	}

	token, err := c.authService.AuthenticateUser(dto.Email, dto.Password)
	if err == services.ErrInvalidCredentials {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

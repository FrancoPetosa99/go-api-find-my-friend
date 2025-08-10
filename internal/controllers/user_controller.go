package controllers

import (
	"net/http"
	"sync"

	"go-api-find-my-friend/internal/services"
	"go-api-find-my-friend/pkg/errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrCreateUserInvalidBody = errors.NewBadRequestError("invalid user body")
)

var (
	userControllerInstance *UserController
	userControllerOnce     sync.Once
)

type UserController struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserController() *UserController {
	userControllerOnce.Do(func() {
		userControllerInstance = &UserController{
			userService: services.NewUserService(),
			authService: services.NewAuthService(),
		}
	})
	return userControllerInstance
}

func (c *UserController) Register(ctx *gin.Context) {
	var dto services.UserCreateDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrCreateUserInvalidBody)
		return
	}

	user, err := c.userService.CreateUser(&dto)
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	token, _ := c.authService.GenerateToken(user.ID, user.Email)

	ctx.JSON(http.StatusCreated, gin.H{
		"user": UserCreateResponse{
			Name:     user.Name,
			LastName: user.LastName,
		},
		"auth_token": token,
	})
}

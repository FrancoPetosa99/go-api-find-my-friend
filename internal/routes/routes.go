package routes

import (
	"go-api-find-my-friend/internal/controllers"
	"go-api-find-my-friend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()
	petController := controllers.NewPetController()
	authController := controllers.NewAuthController()

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}

		users := v1.Group("/users")
		{
			users.POST("/", userController.Register)
		}

		pets := v1.Group("/pets")
		pets.Use(middleware.AuthMiddleware())
		{
			pets.POST("/", petController.CreatePet)
			pets.GET("/", petController.SearchPets)
			pets.GET("/:id", petController.GetPet)
			pets.PUT("/:id", petController.UpdatePet)
			pets.PATCH("/:id/mark-found", petController.UpdatePetMarkAsFound)
			pets.DELETE("/:id", petController.DeletePet)
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "API is running",
			"version": "1.0.0",
		})
	})
}

package main

import (
	"fmt"
	"log"

	"go-api-find-my-friend/internal/routes"
	"go-api-find-my-friend/pkg/config"
	"go-api-find-my-friend/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if config.IsProduction() {
		log.Printf("🚀 Running in PRODUCTION mode")
		database.CreateDB(config)
		database.Connect(config)
		database.AutoMigrate()
	} else {
		log.Printf("🔧 Running in DEVELOPMENT mode")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Servir archivos estáticos
	router.Static("/uploads", "./uploads")

	routes.SetupRoutes(router)

	port := fmt.Sprintf(":%s", config.Server.Port)
	log.Printf("Server starting on port %s", port)
	log.Printf("API available at http://localhost%s", port)
	log.Printf("Health check at http://localhost%s/api/v1/health", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

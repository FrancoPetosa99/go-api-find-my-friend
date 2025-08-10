package database

import (
	"log"

	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(config *config.Config) {
	dsn := config.GetDSN()

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func AutoMigrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Pet{})
	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}
	log.Println("Database migrated successfully")
}

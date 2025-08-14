package database

import (
	"log"
	"fmt"
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

	dbName := config.Database.Name
	createDB(dbName)

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func createDB(dbName string) {
	query := fmt.Sprintf("IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'%s') CREATE DATABASE [%s];", dbName, dbName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create database: ", err)
	}

	log.Println("Database exists or created successfully")
}

func AutoMigrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Pet{})
	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}
	log.Println("Database migrated successfully")
}

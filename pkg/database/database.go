package database

import (
	"fmt"
	"log"

	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(config *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func CreateDB(config *config.Config) {
	dbName := config.Database.Name

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?parseTime=true",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL for creating database: ", err)
	}

	createQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbName)
	if err := db.Exec(createQuery).Error; err != nil {
		log.Fatal("Failed to create database: ", err)
	}

	log.Println("Database exists or created successfully")
}

func AutoMigrate() {
	if err := DB.AutoMigrate(&models.User{}, &models.Pet{}); err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}
	log.Println("Database migrated successfully")
}
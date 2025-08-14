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
	createDB(config)
	
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

func createDB(cfg *config.Config) {
	connString := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Failed to connect to SQL Server: ", err)
	}
	defer db.Close()

	query := fmt.Sprintf("IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'%s') CREATE DATABASE [%s];", cfg.DBName, cfg.DBName)
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

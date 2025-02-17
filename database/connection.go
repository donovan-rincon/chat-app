package database

import (
	"fmt"
	"log"

	"chat-app/config"
	"chat-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	config.LoadConfig()
	dbConfig := config.AppConfig.DB

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	DB.AutoMigrate(&models.User{}, &models.UserMessage{}, &models.Chatroom{})
}

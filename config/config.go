package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		Port     string `json:"port"`
		SSLMode  string `json:"sslmode"`
	} `json:"db"`
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables")
	}

	AppConfig.DB.Host = os.Getenv("DB_HOST")
	AppConfig.DB.User = os.Getenv("DB_USER")
	AppConfig.DB.Password = os.Getenv("DB_PASSWORD")
	AppConfig.DB.DBName = os.Getenv("DB_NAME")
	AppConfig.DB.Port = os.Getenv("DB_PORT")
	AppConfig.DB.SSLMode = os.Getenv("DB_SSLMODE")
}

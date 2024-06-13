package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}

var AppConfig Config

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	AppConfig = Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
	}

	return nil
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

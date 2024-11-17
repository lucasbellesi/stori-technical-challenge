package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	ToEmail      string
}

var AppConfig Config

func LoadConfig() error {
	AppConfig = Config{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", "user@example.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "password"),
		FromEmail:    getEnv("FROM_EMAIL", "from@example.com"),
		ToEmail:      getEnv("TO_EMAIL", "to@example.com"),
	}

	// Validaciones adicionales
	if AppConfig.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if AppConfig.FromEmail == "" {
		return fmt.Errorf("FROM_EMAIL is required")
	}
	if AppConfig.ToEmail == "" {
		return fmt.Errorf("TO_EMAIL is required")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

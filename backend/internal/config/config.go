package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env            string
	HTTPPort       string
	DatabaseURL    string
	JWTSecret      string
	FrontendOrigin string

	SMTPHost      string
	SMTPPort      string
	SMTPFrom      string
	NotifyEmailTo string

	TelegramBotToken string
	TelegramChatID   string

	GarageEndpoint  string
	GarageRegion    string
	GarageAccessKey string
	GarageSecretKey string
	GarageBucket    string
}

func Load() (Config, error) {
	cfg := Config{
		Env:            getenv("APP_ENV", "local"),
		HTTPPort:       getenv("HTTP_PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		FrontendOrigin: getenv("FRONTEND_ORIGIN", "http://localhost:3000"),

		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPFrom:      os.Getenv("SMTP_FROM"),
		NotifyEmailTo: os.Getenv("NOTIFY_EMAIL_TO"),

		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),

		GarageEndpoint:  os.Getenv("GARAGE_ENDPOINT"),
		GarageRegion:    os.Getenv("GARAGE_REGION"),
		GarageAccessKey: os.Getenv("GARAGE_ACCESS_KEY"),
		GarageSecretKey: os.Getenv("GARAGE_SECRET_KEY"),
		GarageBucket:    os.Getenv("GARAGE_BUCKET"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.GarageEndpoint == "" {
		return Config{}, fmt.Errorf("GARAGE_ENDPOINT is required")
	}
	if cfg.GarageAccessKey == "" || cfg.GarageSecretKey == "" {
		return Config{}, fmt.Errorf("GARAGE_ACCESS_KEY and GARAGE_SECRET_KEY are required")
	}
	if cfg.GarageBucket == "" {
		return Config{}, fmt.Errorf("GARAGE_BUCKET is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

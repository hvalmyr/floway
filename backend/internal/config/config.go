package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env         string
	HTTPPort    string
	DatabaseURL string
	JWTSecret   string

	SMTPHost      string
	SMTPPort      string
	SMTPFrom      string
	NotifyEmailTo string

	TelegramBotToken string
	TelegramChatID   string
}

func Load() (Config, error) {
	cfg := Config{
		Env:         getenv("APP_ENV", "local"),
		HTTPPort:    getenv("HTTP_PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),

		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPFrom:      os.Getenv("SMTP_FROM"),
		NotifyEmailTo: os.Getenv("NOTIFY_EMAIL_TO"),

		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package config

import (
	"errors"
	"os"
)

type Config struct {
	BotToken   string
	DatabaseURL string
}

func Load() (*Config, error) {
	cfg := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.BotToken == "" {
		return nil, errors.New("BOT_TOKEN is required")
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return cfg, nil
}

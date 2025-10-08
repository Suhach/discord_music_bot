package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken    string
	SpotifyClientID string
	SpotifySecret   string
	YouTubeAPIKey   string
	SoundCloudID    string
}

func Load() (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env: %w", err)
	}

	cfg := &Config{
		DiscordToken:    os.Getenv("DISCORD_TOKEN"),
		SpotifyClientID: os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifySecret:   os.Getenv("SPOTIFY_SECRET"),
		YouTubeAPIKey:   os.Getenv("YOUTUBE_API_KEY"),
		SoundCloudID:    os.Getenv("SOUNDCLOUD_CLIENT_ID"),
	}

	// Проверяем обязательные поля
	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN не установлен")
	}

	return cfg, nil
}

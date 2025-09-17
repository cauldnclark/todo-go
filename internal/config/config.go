package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Redis    RedisConfig
	Google   GoogleConfig
}

type DatabaseConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	DatabaseURL string
}
type ServerConfig struct {
	Port      string
	Env       string
	IsProd    bool
	JWTSecret string
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func Load() (*Config, error) {
	// Load .env if available (non-fatal in production)
	_ = godotenv.Load()

	// Validate required env vars
	required := []string{
		"PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"REDIS_HOST",
		"REDIS_PORT",
		"JWT_SECRET",
		"ENV",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"GOOGLE_REDIRECT_URL",
	}
	for _, v := range required {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", v)
		}
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:        os.Getenv("DB_HOST"),
			Port:        os.Getenv("DB_PORT"),
			User:        os.Getenv("DB_USER"),
			Password:    os.Getenv("DB_PASSWORD"),
			Name:        os.Getenv("DB_NAME"),
			DatabaseURL: os.Getenv("DATABASE_URL"),
		},
		Server: ServerConfig{
			Port:      os.Getenv("PORT"),
			Env:       os.Getenv("ENV"),
			IsProd:    os.Getenv("ENV") == "production",
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Google: GoogleConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
	}

	return config, nil
}

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
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
type ServerConfig struct {
	Port   string
	Env    string
	IsProd bool
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
	if os.Getenv("ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error loading .env file")
		}
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Server: ServerConfig{
			Port:   os.Getenv("PORT"),
			Env:    os.Getenv("ENV"),
			IsProd: os.Getenv("ENV") == "production",
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

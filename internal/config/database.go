package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGXConnection(cfg *DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

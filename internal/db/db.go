package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolFromEnv() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	return pgxpool.New(context.Background(), dsn)
}

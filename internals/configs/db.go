package configs

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), os.Getenv("DB_URL_M"))
}

func PingDB(p *pgxpool.Pool) error {
	return p.Ping(context.Background())
}

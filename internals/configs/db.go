package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	db := os.Getenv("DB_NAME_M")

	connstring := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s", user, pwd, host, port, db,
	)
	return pgxpool.New(context.Background(), connstring)
}

func PingDB(p *pgxpool.Pool) error {
	return p.Ping(context.Background())
}

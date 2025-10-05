package configs

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	addr := os.Getenv("RDB_ADDR")

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

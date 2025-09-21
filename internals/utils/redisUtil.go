package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedisCache(rdb *redis.Client, ctx context.Context, key string, v any) error {
	rcmd := rdb.Get(ctx, key)
	if rcmd.Err() != nil {
		if rcmd.Err() != redis.Nil {
			PrintError("redis> SERVER ERROR", 16, rcmd.Err())
		}
	} else {
		// konversi menjadi tipe data []byte
		bites, err := rcmd.Bytes()
		if err != nil {
			PrintError("INTERNAL SERVER ERROR", 12, err)
		} else {
			if err := json.Unmarshal(bites, &v); err != nil {
				PrintError("UNABLE TO PARSE", 16, err)
			}
		}
	}

	return rcmd.Err()
}

func RedisCache(rdb *redis.Client, ctx context.Context, key string, v any, expiration time.Duration) {
	bites, err := json.Marshal(v)
	if err != nil {
		PrintError("redis> UNABLE TO PARSE KEY", 16, err)
	}

	if err := rdb.Set(ctx, key, bites, expiration).Err(); err != nil {
		PrintError("redis> UNABLE TO SET KEY", 20, err)
	}

	PrintError("redis> SET KEY SUCCESS", 20, nil)
}

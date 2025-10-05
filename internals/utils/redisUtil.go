package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func CacheGet(rdb *redis.Client, ctx context.Context, key string, result any) (bool, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("redis> KEY %s DOES NOT EXIST", key)
	} else if err != nil {
		return false, fmt.Errorf("redis> INTERNAL SERVER ERROR\n%s", err)
	}

	if err := json.Unmarshal([]byte(val), result); err != nil {
		return false, err
	}

	PrintError(
		fmt.Sprintf("redis> CACHE HIT %s", key), 20, nil,
	)
	return true, nil
}

func CacheSet(rdb *redis.Client, ctx context.Context, key string, value any, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	PrintError(
		fmt.Sprintf("redis> CACHE SET %s", key), 20, nil,
	)
	return rdb.Set(ctx, key, jsonData, expiration).Err()
}

func InvalidateCache(rdb *redis.Client, ctx context.Context, key string) error {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		PrintError(
			fmt.Sprintf("redis> ERROR INVALIDATE %s", key), 20, nil)
		return err
	}
	return nil
}

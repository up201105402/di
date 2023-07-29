package util

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectToRedis() (*redis.Client, error) {

	redisHost, exists := os.LookupEnv("REDIS_HOST")

	if !exists {
		panic("REDIS_HOST is not defined!")
	}

	redisPort, exists := os.LookupEnv("REDIS_PORT")

	if !exists {
		panic("REDIS_PORT is not defined!")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return rdb, nil
}

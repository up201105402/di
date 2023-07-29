package repository

import (
	"context"
	"di/model"
	"di/util/errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID uint, tokenID uint, expiresIn time.Duration) error {
	key := fmt.Sprintf("%d:%d", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID uint, tokenID uint) error {
	key := fmt.Sprintf("%d:%d", userID, tokenID)

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		return err
	}

	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID: %d/%d does not exist\n", userID, tokenID)
		return errors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID uint) error {
	pattern := fmt.Sprintf("%d*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}

	// check last value
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iter.Val())
	}

	if failCount > 0 {
		return errors.NewInternal("")
	}

	return nil
}

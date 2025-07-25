package redis

import (
	"commons/database"
	"commons/shared_types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisClient(addr, password string, redisIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       redisIndex,
	})
}

type RedisService struct {
	Conn *redis.Client
}

func (r *RedisService) Ping(ctx context.Context) error {
	r.Conn.Ping(ctx)
	return nil
}

func (r *RedisService) SetUserToken(ctx context.Context, email, token string, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("user-%s", email)
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return r.Conn.SetEX(ctx, cacheKey, data, expiration).Err()
}

func (r *RedisService) ValidateUserToken(ctx context.Context, email string) (*shared_types.RedisUserInfo, error) {
	cacheKey := fmt.Sprintf("user-%s", email)
	data, err := r.Conn.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, database.NotFound
	}
	if data == "" {
		return nil, database.NotFound
	}

	ttl := r.Conn.TTL(ctx, cacheKey).Val()
	return &shared_types.RedisUserInfo{
		Email: email,
		Token: data,
		TTL:   int64(ttl),
	}, nil
}

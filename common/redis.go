package commons

import "github.com/go-redis/redis/v8"

func NewRedisClient(addr, password string, redisIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       redisIndex,
	})
}

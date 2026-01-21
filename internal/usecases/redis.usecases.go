package usecases

import (
	"shortner-url/internal"
	"time"
)

type RedisUsecase struct{}

func NewRedisUsecase() *RedisUsecase {
	return &RedisUsecase{}
}

func (r *RedisUsecase) Set(key string, value string, ttl time.Duration) error {
	return internal.RedisClient.Set(internal.RedisCtx, key, value, ttl).Err()
}

func (r *RedisUsecase) Get(key string) (string, error) {
	return internal.RedisClient.Get(internal.RedisCtx, key).Result()
}

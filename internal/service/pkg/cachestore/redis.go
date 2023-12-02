package cachestore

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	redisC *redis.Client
}

func NewRedisStore(redisHost string, redisPort int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: "",
		DB:       0,
	})
	return &RedisStore{
		redisC: rdb,
	}
}

func (r *RedisStore) HSetDataPair(ctx context.Context, keyPrefix string, sno uint32, field string, value interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("%s_%d", keyPrefix, sno)
	if err := r.redisC.HSetNX(ctx, key, field, value).Err(); err != nil {
		return err
	}
	return r.redisC.Expire(ctx, key, ttl).Err()
}

func (r *RedisStore) HSetDataPairs(ctx context.Context, keyPrefix string, sno uint32, data map[string]interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("%s_%d", keyPrefix, sno)
	if err := r.redisC.HMSet(ctx, key, data).Err(); err != nil {
		return err
	}
	return r.redisC.Expire(ctx, key, ttl).Err()
}

func (r *RedisStore) HGetAll(ctx context.Context, keyPrefix string, sno uint32) (map[string]string, error) {
	key := fmt.Sprintf("%s_%d", keyPrefix, sno)
	return r.redisC.HGetAll(ctx, key).Result()
}

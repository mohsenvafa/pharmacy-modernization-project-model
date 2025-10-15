package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache struct {
	rdb    *redis.Client
	prefix string
	logger *zap.Logger
	hits   atomic.Int64
	misses atomic.Int64
	errors atomic.Int64
}

type RedisConfig struct {
	Addr     string // "host:6379"
	Password string
	DB       int
	Prefix   string // e.g. "rx:" to namespace keys per service/env
}

func NewRedisCache(config RedisConfig, logger *zap.Logger) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisCache{
		rdb:    rdb,
		prefix: config.Prefix,
		logger: logger,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	prefixedKey := r.prefix + key
	result, err := r.rdb.Get(ctx, prefixedKey).Result()
	if err != nil {
		if err == redis.Nil {
			r.misses.Add(1)
			return nil, ErrNotFound
		}
		r.errors.Add(1)
		return nil, err
	}

	r.hits.Add(1)
	return []byte(result), nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	prefixedKey := r.prefix + key
	err := r.rdb.Set(ctx, prefixedKey, value, ttl).Err()
	if err != nil {
		r.errors.Add(1)
	}
	return err
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	prefixedKey := r.prefix + key
	err := r.rdb.Del(ctx, prefixedKey).Err()
	if err != nil {
		r.errors.Add(1)
	}
	return err
}

func (r *RedisCache) Close() error {
	return r.rdb.Close()
}

func (r *RedisCache) Stats() CacheStats {
	hits := r.hits.Load()
	misses := r.misses.Load()
	total := hits + misses

	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return CacheStats{
		Hits:    hits,
		Misses:  misses,
		Errors:  r.errors.Load(),
		HitRate: hitRate,
	}
}

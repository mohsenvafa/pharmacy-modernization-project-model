package cache

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type CacheConfig struct {
	Strategy  string
	Redis     RedisConfig
	Memcached MemcachedConfig
	Memory    MemoryConfig
}

func NewCache(strategy string, config CacheConfig, logger *zap.Logger) (Cache, error) {
	switch strategy {
	case "hybrid":
		return NewHybridCacheFromConfig(config, logger)
	case "redis":
		return NewRedisCache(config.Redis, logger), nil
	case "memcached":
		return NewMemcachedCache(config.Memcached, logger), nil
	case "memory":
		return NewMemoryCache(config.Memory, logger)
	default:
		// Default to memory cache
		logger.Warn("Unknown cache strategy, falling back to memory cache", zap.String("strategy", strategy))
		return NewMemoryCache(MemoryConfig{
			MaxCost:     67108864, // 64MB
			BufferItems: 64,
			Metrics:     true,
			DefaultTTL:  30 * time.Minute,
		}, logger)
	}
}

func NewHybridCacheFromConfig(config CacheConfig, logger *zap.Logger) (Cache, error) {
	// Create local (memory) cache
	localCache, err := NewMemoryCache(config.Memory, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory cache: %w", err)
	}

	// Create shared cache (prefer Memcached, fallback to Redis)
	var sharedCache Cache
	if config.Memcached.Addr != "" {
		sharedCache = NewMemcachedCache(config.Memcached, logger)
		logger.Info("Hybrid cache using Memcached as shared tier")
	} else if config.Redis.Addr != "" {
		sharedCache = NewRedisCache(config.Redis, logger)
		logger.Info("Hybrid cache using Redis as shared tier")
	} else {
		return nil, fmt.Errorf("hybrid cache requires either Redis or Memcached configuration")
	}

	return NewHybridCache(localCache, sharedCache), nil
}

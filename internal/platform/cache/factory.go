package cache

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type CacheConfig struct {
	Strategy string
	Memory   MemoryConfig
	MongoDB  MongoDBConfig
}

func NewCache(strategy string, config CacheConfig, logger *zap.Logger) (Cache, error) {
	switch strategy {
	case "hybrid":
		return NewHybridCacheFromConfig(config, logger)
	case "mongodb":
		return NewMongoDBCache(config.MongoDB, logger)
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

	// Create shared cache (MongoDB)
	if config.MongoDB.Collection == nil {
		return nil, fmt.Errorf("hybrid cache requires MongoDB configuration")
	}

	sharedCache, err := NewMongoDBCache(config.MongoDB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB cache: %w", err)
	}
	logger.Info("Hybrid cache using MongoDB as shared tier")

	return NewHybridCache(localCache, sharedCache), nil
}

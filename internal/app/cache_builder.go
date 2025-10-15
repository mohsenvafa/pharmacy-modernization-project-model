package app

import (
	"time"

	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/platform/cache"
	"pharmacy-modernization-project-model/internal/platform/config"
	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
)

// CacheBuilder creates cache instances with flexible configuration
type CacheBuilder struct {
	config *config.Config
	logger *zap.Logger
}

// NewCacheBuilder creates a new cache builder
func NewCacheBuilder(config *config.Config, logger *zap.Logger) *CacheBuilder {
	return &CacheBuilder{
		config: config,
		logger: logger,
	}
}

// BuildCache creates a cache instance with the specified strategy and configuration
func (b *CacheBuilder) BuildCache(strategy string, instanceConfig CacheInstanceConfig) (cache.Cache, error) {
	// Validate strategy
	if strategy == "" {
		strategy = "memcached" // Default to memcached
		b.logger.Warn("Cache strategy not specified, using memcached as default")
	}

	// Build cache configuration
	config := b.buildCacheConfig(strategy, instanceConfig)

	// Validate configuration
	if err := b.validateCacheConfig(strategy, instanceConfig); err != nil {
		return nil, err
	}

	// Create cache service
	cacheService, err := cache.NewCache(strategy, config, b.logger)
	if err != nil {
		b.logger.Error("Failed to initialize cache service",
			zap.Error(err),
			zap.String("strategy", strategy))
		return nil, platformErrors.NewConfigurationError("cache", "strategy", "Failed to initialize cache service")
	}

	// Wrap with metrics middleware
	cacheService = cache.NewCacheMiddleware(cacheService, b.logger)

	b.logger.Info("Cache service created successfully",
		zap.String("strategy", strategy))

	return cacheService, nil
}

// buildCacheConfig builds cache configuration from instance config and global config
func (b *CacheBuilder) buildCacheConfig(strategy string, instanceConfig CacheInstanceConfig) cache.CacheConfig {
	// Use instance-specific config if available, otherwise fallback to global
	redisConfig := b.configToRedis(b.config.Cache.Redis)
	if instanceConfig.Redis.Addr != "" {
		redisConfig = instanceConfig.Redis
	}

	memcachedConfig := b.configToMemcached(b.config.Cache.Memcached)
	if instanceConfig.Memcached.Addr != "" {
		memcachedConfig = instanceConfig.Memcached
	}

	memoryConfig := b.configToMemory(b.config.Cache.Memory)
	if instanceConfig.Memory.MaxCost != 0 {
		memoryConfig = instanceConfig.Memory
	}

	return cache.CacheConfig{
		Strategy:  strategy,
		Redis:     redisConfig,
		Memcached: memcachedConfig,
		Memory:    memoryConfig,
	}
}

// Helper functions to convert config types
func (b *CacheBuilder) configToRedis(cfg config.RedisCacheConfig) cache.RedisConfig {
	return cache.RedisConfig{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		Prefix:   cfg.Prefix,
	}
}

func (b *CacheBuilder) configToMemcached(cfg config.MemcachedCacheConfig) cache.MemcachedConfig {
	return cache.MemcachedConfig{
		Addr:   cfg.Addr,
		Prefix: cfg.Prefix,
	}
}

func (b *CacheBuilder) configToMemory(cfg config.MemoryCacheConfig) cache.MemoryConfig {
	// Parse TTL string
	ttl, err := time.ParseDuration(cfg.DefaultTTL)
	if err != nil {
		b.logger.Warn("Invalid default TTL, using 30m", zap.String("ttl", cfg.DefaultTTL))
		ttl = 30 * time.Minute
	}

	return cache.MemoryConfig{
		MaxCost:     cfg.MaxCost,
		BufferItems: cfg.BufferItems,
		Metrics:     cfg.Metrics,
		DefaultTTL:  ttl,
	}
}

// validateCacheConfig validates cache configuration based on strategy
func (b *CacheBuilder) validateCacheConfig(strategy string, cfg CacheInstanceConfig) error {
	switch strategy {
	case "redis":
		if cfg.Redis.Addr == "" && b.config.Cache.Redis.Addr == "" {
			b.logger.Error("Redis address is required for Redis cache strategy")
			return platformErrors.NewConfigurationError("cache", "redis.addr", "Redis address is required for Redis cache strategy")
		}
	case "memcached":
		if cfg.Memcached.Addr == "" && b.config.Cache.Memcached.Addr == "" {
			b.logger.Error("Memcached address is required for Memcached cache strategy")
			return platformErrors.NewConfigurationError("cache", "memcached.addr", "Memcached address is required for Memcached cache strategy")
		}
	case "hybrid":
		if cfg.Redis.Addr == "" && cfg.Memcached.Addr == "" &&
			b.config.Cache.Redis.Addr == "" && b.config.Cache.Memcached.Addr == "" {
			b.logger.Error("Either Redis or Memcached address is required for hybrid cache strategy")
			return platformErrors.NewConfigurationError("cache", "hybrid", "Either Redis or Memcached address is required for hybrid cache strategy")
		}
	}

	return nil
}

// CacheInstanceConfig represents configuration for a specific cache instance
type CacheInstanceConfig struct {
	Redis     cache.RedisConfig
	Memcached cache.MemcachedConfig
	Memory    cache.MemoryConfig
}

// Helper methods for common cache configurations

// BuildMemcachedCache creates a Memcached cache with the given prefix
func (b *CacheBuilder) BuildMemcachedCache(prefix string) (cache.Cache, error) {
	return b.BuildCache("memcached", CacheInstanceConfig{
		Memcached: cache.MemcachedConfig{
			Addr:   b.config.Cache.Memcached.Addr,
			Prefix: prefix,
		},
	})
}

// BuildMemoryCache creates a memory cache with the given max cost
func (b *CacheBuilder) BuildMemoryCache(maxCost int64) (cache.Cache, error) {
	return b.BuildCache("memory", CacheInstanceConfig{
		Memory: cache.MemoryConfig{
			MaxCost:     maxCost,
			BufferItems: 64,
			Metrics:     true,
			DefaultTTL:  30 * time.Minute,
		},
	})
}

// BuildRedisCache creates a Redis cache with the given address and prefix
func (b *CacheBuilder) BuildRedisCache(addr string, prefix string) (cache.Cache, error) {
	return b.BuildCache("redis", CacheInstanceConfig{
		Redis: cache.RedisConfig{
			Addr:   addr,
			Prefix: prefix,
		},
	})
}

// BuildHybridCache creates a hybrid cache (memory + shared)
func (b *CacheBuilder) BuildHybridCache(memoryMaxCost int64) (cache.Cache, error) {
	return b.BuildCache("hybrid", CacheInstanceConfig{
		Memory: cache.MemoryConfig{
			MaxCost:     memoryMaxCost,
			BufferItems: 64,
			Metrics:     true,
			DefaultTTL:  30 * time.Minute,
		},
	})
}

package builder

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/platform/cache"
	"pharmacy-modernization-project-model/internal/platform/config"
	"pharmacy-modernization-project-model/internal/platform/database"
	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
)

// CacheBuilder creates cache instances with flexible configuration
type CacheBuilder struct {
	config            *config.Config
	logger            *zap.Logger
	cacheMongoConnMgr *database.ConnectionManager
}

// NewCacheBuilder creates a new cache builder
func NewCacheBuilder(config *config.Config, logger *zap.Logger) *CacheBuilder {
	return &CacheBuilder{
		config: config,
		logger: logger,
	}
}

// CreateCacheMongoDBConnection creates a separate MongoDB connection for caching
func CreateCacheMongoDBConnection(cfg *config.Config, logger *zap.Logger) (*database.ConnectionManager, error) {
	// Validate configuration
	if cfg.Cache.MongoDB.URI == "" {
		logger.Info("Cache MongoDB URI not configured, cache will use main MongoDB or fallback to memory")
		return nil, nil
	}

	if cfg.Cache.MongoDB.Database == "" {
		logger.Error("Cache MongoDB database name is required")
		return nil, platformErrors.NewConfigurationError("cache", "mongodb.database", "Cache MongoDB database name is required")
	}

	mongoConfig := database.MongoDBConfig{
		URI:      cfg.Cache.MongoDB.URI,
		Database: cfg.Cache.MongoDB.Database,
		Collections: map[string]string{
			"cache": cfg.Cache.MongoDB.Collection,
		},
		Connection: database.ConnectionConfig{
			MaxPoolSize:    cfg.Cache.MongoDB.Connection.MaxPoolSize,
			MinPoolSize:    cfg.Cache.MongoDB.Connection.MinPoolSize,
			MaxIdleTime:    cfg.Cache.MongoDB.Connection.MaxIdleTime,
			ConnectTimeout: cfg.Cache.MongoDB.Connection.ConnectTimeout,
			SocketTimeout:  cfg.Cache.MongoDB.Connection.SocketTimeout,
		},
		Options: database.OptionsConfig{
			RetryWrites: true,
			RetryReads:  true,
		},
	}

	connMgr, err := database.NewConnectionManager(mongoConfig, logger)
	if err != nil {
		logger.Error("Failed to initialize cache MongoDB connection",
			zap.Error(err),
			zap.String("uri", cfg.Cache.MongoDB.URI),
			zap.String("database", cfg.Cache.MongoDB.Database))
		return nil, platformErrors.NewConfigurationError("cache", "mongodb.connection", "Failed to initialize cache MongoDB connection")
	}

	logger.Info("Cache MongoDB connection established successfully",
		zap.String("database", cfg.Cache.MongoDB.Database))
	return connMgr, nil
}

// GetCacheMongoCollection returns the cache collection from cache MongoDB connection manager
func GetCacheMongoCollection(cacheMongoConnMgr *database.ConnectionManager, cfg *config.Config) *mongo.Collection {
	if cacheMongoConnMgr == nil {
		return nil
	}

	collectionName := cfg.Cache.MongoDB.Collection
	if collectionName == "" {
		collectionName = "cache"
	}

	return cacheMongoConnMgr.GetDatabase().Collection(collectionName)
}

// BuildCache creates a cache instance with the specified strategy and configuration
func (b *CacheBuilder) BuildCache(strategy string, instanceConfig CacheInstanceConfig) (cache.Cache, error) {
	// Validate strategy
	if strategy == "" {
		strategy = "memory" // Default to memory
		b.logger.Warn("Cache strategy not specified, using memory as default")
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
	memoryConfig := b.configToMemory(b.config.Cache.Memory)
	if instanceConfig.Memory.MaxCost != 0 {
		memoryConfig = instanceConfig.Memory
	}

	mongodbConfig := instanceConfig.MongoDB

	return cache.CacheConfig{
		Strategy: strategy,
		Memory:   memoryConfig,
		MongoDB:  mongodbConfig,
	}
}

// Helper functions to convert config types
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
	case "mongodb":
		if cfg.MongoDB.Collection == nil {
			b.logger.Error("MongoDB collection is required for MongoDB cache strategy")
			return platformErrors.NewConfigurationError("cache", "mongodb.collection", "MongoDB collection is required for MongoDB cache strategy")
		}
	case "hybrid":
		if cfg.MongoDB.Collection == nil {
			b.logger.Error("MongoDB is required for hybrid cache strategy")
			return platformErrors.NewConfigurationError("cache", "hybrid", "MongoDB is required for hybrid cache strategy")
		}
	}

	return nil
}

// CacheInstanceConfig represents configuration for a specific cache instance
type CacheInstanceConfig struct {
	Memory  cache.MemoryConfig
	MongoDB cache.MongoDBConfig
}

// Helper methods for common cache configurations

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

// BuildMongoDBCache creates a MongoDB cache with the given collection and prefix
func (b *CacheBuilder) BuildMongoDBCache(collection *mongo.Collection, prefix string) (cache.Cache, error) {
	if collection == nil {
		return nil, platformErrors.NewConfigurationError("cache", "mongodb.collection", "MongoDB collection cannot be nil")
	}

	return b.BuildCache("mongodb", CacheInstanceConfig{
		MongoDB: cache.MongoDBConfig{
			Collection: collection,
			Prefix:     prefix,
			TTLIndex:   true,
		},
	})
}

package cache

import (
	"context"
	"sync/atomic"
	"time"

	"pharmacy-modernization-project-model/internal/platform/sanitizer"

	"go.uber.org/zap"
)

type CacheMiddleware struct {
	cache      Cache
	logger     *zap.Logger
	totalOps   atomic.Int64
	totalNanos atomic.Int64
}

func NewCacheMiddleware(cache Cache, logger *zap.Logger) Cache {
	return &CacheMiddleware{
		cache:  cache,
		logger: logger,
	}
}

func (m *CacheMiddleware) Get(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()
	defer func() {
		m.totalOps.Add(1)
		m.totalNanos.Add(time.Since(start).Nanoseconds())
	}()

	value, err := m.cache.Get(ctx, key)
	if err != nil && err != ErrNotFound {
		m.logger.Debug("Cache get error",
			zap.String("key", sanitizer.ForLogging(key)),
			zap.Error(err),
			zap.Duration("latency", time.Since(start)))
	}

	return value, err
}

func (m *CacheMiddleware) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	start := time.Now()
	defer func() {
		m.totalOps.Add(1)
		m.totalNanos.Add(time.Since(start).Nanoseconds())
	}()

	err := m.cache.Set(ctx, key, value, ttl)
	if err != nil {
		m.logger.Warn("Cache set error",
			zap.String("key", sanitizer.ForLogging(key)),
			zap.Error(err),
			zap.Duration("latency", time.Since(start)))
	}

	return err
}

func (m *CacheMiddleware) Delete(ctx context.Context, key string) error {
	start := time.Now()
	defer func() {
		m.totalOps.Add(1)
		m.totalNanos.Add(time.Since(start).Nanoseconds())
	}()

	err := m.cache.Delete(ctx, key)
	if err != nil {
		m.logger.Warn("Cache delete error",
			zap.String("key", sanitizer.ForLogging(key)),
			zap.Error(err),
			zap.Duration("latency", time.Since(start)))
	}

	return err
}

func (m *CacheMiddleware) Close() error {
	return m.cache.Close()
}

func (m *CacheMiddleware) Stats() CacheStats {
	stats := m.cache.Stats()

	// Add average latency if we have operations
	ops := m.totalOps.Load()
	if ops > 0 {
		avgNanos := m.totalNanos.Load() / ops
		m.logger.Debug("Cache stats",
			zap.Int64("hits", stats.Hits),
			zap.Int64("misses", stats.Misses),
			zap.Float64("hit_rate", stats.HitRate),
			zap.Duration("avg_latency", time.Duration(avgNanos)))
	}

	return stats
}

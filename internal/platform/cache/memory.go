package cache

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto"
	"go.uber.org/zap"
)

type MemoryCache struct {
	rc     *ristretto.Cache
	config MemoryConfig
	logger *zap.Logger
}

type MemoryConfig struct {
	// Max keys * average value bytes ≈ MaxCost
	MaxCost     int64         // e.g. 64<<20 (~64MB)
	BufferItems int64         // e.g. 64
	Metrics     bool          // optional
	DefaultTTL  time.Duration // fallback when caller passes ttl<=0 (optional)
}

func NewMemoryCache(config MemoryConfig, logger *zap.Logger) (Cache, error) {
	rc, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: config.MaxCost * 10, // 10x recommended
		MaxCost:     config.MaxCost,
		BufferItems: config.BufferItems,
		Metrics:     config.Metrics,
	})
	if err != nil {
		return nil, err
	}

	return &MemoryCache{
		rc:     rc,
		config: config,
		logger: logger,
	}, nil
}

func (m *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	if value, found := m.rc.Get(key); found {
		return value.([]byte), nil
	}
	return nil, ErrNotFound
}

func (m *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// Use default TTL if none provided
	if ttl <= 0 {
		ttl = m.config.DefaultTTL
	}

	m.rc.Set(key, value, int64(ttl.Nanoseconds()))
	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.rc.Del(key)
	return nil
}

func (m *MemoryCache) Close() error {
	m.rc.Close()
	return nil
}

func (m *MemoryCache) Stats() CacheStats {
	metrics := m.rc.Metrics
	hits := metrics.Hits()
	misses := metrics.Misses()
	total := hits + misses

	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return CacheStats{
		Hits:      hits,
		Misses:    misses,
		Evictions: metrics.Evictions(),
		Errors:    metrics.Errors(),
		HitRate:   hitRate,
		Size:      metrics.KeysAdded() - metrics.KeysEvicted(),
		MaxSize:   int64(metrics.Capacity()),
	}
}

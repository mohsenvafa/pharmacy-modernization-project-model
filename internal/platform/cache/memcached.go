package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

type MemcachedCache struct {
	client *memcache.Client
	prefix string
	logger *zap.Logger
	hits   atomic.Int64
	misses atomic.Int64
	errors atomic.Int64
}

type MemcachedConfig struct {
	Addr   string // "host:11211"
	Prefix string // e.g. "rx:" to namespace keys per service/env
}

func NewMemcachedCache(config MemcachedConfig, logger *zap.Logger) Cache {
	return &MemcachedCache{
		client: memcache.New(config.Addr),
		prefix: config.Prefix,
		logger: logger,
	}
}

func (m *MemcachedCache) Get(ctx context.Context, key string) ([]byte, error) {
	prefixedKey := m.prefix + key
	item, err := m.client.Get(prefixedKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			m.misses.Add(1)
			return nil, ErrNotFound
		}
		m.errors.Add(1)
		return nil, err
	}

	m.hits.Add(1)
	return item.Value, nil
}

func (m *MemcachedCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	prefixedKey := m.prefix + key
	item := &memcache.Item{
		Key:        prefixedKey,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	}
	err := m.client.Set(item)
	if err != nil {
		m.errors.Add(1)
	}
	return err
}

func (m *MemcachedCache) Delete(ctx context.Context, key string) error {
	prefixedKey := m.prefix + key
	err := m.client.Delete(prefixedKey)
	if err != nil && err != memcache.ErrCacheMiss {
		m.errors.Add(1)
		return err
	}
	return nil
}

func (m *MemcachedCache) Close() error {
	// Memcached client doesn't have a Close method
	return nil
}

func (m *MemcachedCache) Stats() CacheStats {
	hits := m.hits.Load()
	misses := m.misses.Load()
	total := hits + misses

	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return CacheStats{
		Hits:    hits,
		Misses:  misses,
		Errors:  m.errors.Load(),
		HitRate: hitRate,
	}
}

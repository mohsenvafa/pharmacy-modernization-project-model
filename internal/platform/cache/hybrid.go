package cache

import (
	"context"
	"time"
)

type HybridCache struct {
	Local  Cache // fast tier (memory)
	Shared Cache // distributed tier (redis/memcached)
}

func NewHybridCache(local Cache, shared Cache) Cache {
	return &HybridCache{
		Local:  local,
		Shared: shared,
	}
}

func (h *HybridCache) Get(ctx context.Context, key string) ([]byte, error) {
	// 1) Try local cache first
	if b, err := h.Local.Get(ctx, key); err == nil {
		return b, nil
	}

	// 2) Try shared cache
	b, err := h.Shared.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// 3) Backfill local with shorter TTL to reduce memory churn
	_ = h.Local.Set(ctx, key, b, 30*time.Second)
	return b, nil
}

func (h *HybridCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// Set in local with shorter TTL
	_ = h.Local.Set(ctx, key, value, min(ttl, 30*time.Second))

	// Set in shared with full TTL
	return h.Shared.Set(ctx, key, value, ttl)
}

func (h *HybridCache) Delete(ctx context.Context, key string) error {
	// Delete from both
	_ = h.Local.Delete(ctx, key)
	return h.Shared.Delete(ctx, key)
}

func (h *HybridCache) Close() error {
	_ = h.Local.Close()
	return h.Shared.Close()
}

func (h *HybridCache) Stats() CacheStats {
	localStats := h.Local.Stats()
	sharedStats := h.Shared.Stats()

	totalHits := localStats.Hits + sharedStats.Hits
	totalMisses := localStats.Misses + sharedStats.Misses
	total := totalHits + totalMisses

	var hitRate float64
	if total > 0 {
		hitRate = float64(totalHits) / float64(total)
	}

	return CacheStats{
		Hits:      totalHits,
		Misses:    totalMisses,
		Evictions: localStats.Evictions + sharedStats.Evictions,
		Errors:    localStats.Errors + sharedStats.Errors,
		HitRate:   hitRate,
		Size:      localStats.Size + sharedStats.Size,
		MaxSize:   localStats.MaxSize + sharedStats.MaxSize,
	}
}

func min(a, b time.Duration) time.Duration {
	if a == 0 || b == 0 {
		return a
	}
	if a < b {
		return a
	}
	return b
}

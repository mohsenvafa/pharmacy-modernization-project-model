package cache

import (
	"context"
	"time"
)

type CacheHealthChecker struct {
	cache Cache
}

func NewCacheHealthChecker(cache Cache) *CacheHealthChecker {
	return &CacheHealthChecker{cache: cache}
}

func (h *CacheHealthChecker) Check(ctx context.Context) error {
	testKey := "health:check"
	testValue := []byte("ok")

	// Test set
	if err := h.cache.Set(ctx, testKey, testValue, time.Minute); err != nil {
		return err
	}

	// Test get
	if _, err := h.cache.Get(ctx, testKey); err != nil {
		return err
	}

	// Clean up
	_ = h.cache.Delete(ctx, testKey)

	return nil
}

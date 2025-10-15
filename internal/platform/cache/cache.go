package cache

import (
	"context"
	"errors"
	"time"
)

// Cache defines the unified cache interface
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
	Stats() CacheStats
}

// CacheStats provides cache performance metrics
type CacheStats struct {
	Hits      int64
	Misses    int64
	Evictions int64
	Errors    int64
	HitRate   float64
	Size      int64
	MaxSize   int64
}

var (
	ErrNotFound = errors.New("cache: key not found")
	ErrClosed   = errors.New("cache: cache is closed")
)

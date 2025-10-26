package cache

import (
	"context"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type prefixedKey string

type MongoDBCache struct {
	collection *mongo.Collection
	prefix     string
	logger     *zap.Logger
	hits       atomic.Int64
	misses     atomic.Int64
	errors     atomic.Int64
}

type MongoDBConfig struct {
	Collection *mongo.Collection // MongoDB collection for cache
	Prefix     string            // e.g. "rx:" to namespace keys per service/env
	TTLIndex   bool              // Create TTL index on expireAt field (default: true)
}

// cacheDocument represents a cached item in MongoDB
type cacheDocument struct {
	Key      prefixedKey `bson:"_id"`
	Value    []byte      `bson:"value"`
	ExpireAt time.Time   `bson:"expireAt"`
}

func NewMongoDBCache(config MongoDBConfig, logger *zap.Logger) (Cache, error) {
	cache := &MongoDBCache{
		collection: config.Collection,
		prefix:     config.Prefix,
		logger:     logger,
	}

	// Create TTL index if requested (default true)
	if config.TTLIndex {
		if err := cache.createTTLIndex(); err != nil {
			logger.Warn("Failed to create TTL index for MongoDB cache", zap.Error(err))
			// Don't fail - continue without TTL index
		}
	}

	return cache, nil
}

func (m *MongoDBCache) createTTLIndex() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expireAt", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := m.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	m.logger.Info("Created TTL index for MongoDB cache")
	return nil
}

func (m *MongoDBCache) Get(ctx context.Context, key string) ([]byte, error) {
	// Sanitize the key before using it in database queries
	sanitizedKey := SanitizeKey(key)
	prefixedKey := prefixedKey(m.prefix + sanitizedKey)

	var doc cacheDocument
	err := m.collection.FindOne(ctx, bson.M{"_id": prefixedKey}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.misses.Add(1)
			return nil, ErrNotFound
		}
		m.errors.Add(1)
		return nil, err
	}

	// Check if expired (in case TTL index hasn't processed yet)
	if time.Now().After(doc.ExpireAt) {
		m.misses.Add(1)
		// Delete expired document
		go func() {
			deleteCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			m.collection.DeleteOne(deleteCtx, bson.M{"_id": prefixedKey})
		}()
		return nil, ErrNotFound
	}

	m.hits.Add(1)
	return doc.Value, nil
}

func (m *MongoDBCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// Sanitize the key before using it in database queries
	sanitizedKey := SanitizeKey(key)
	prefixedKey := prefixedKey(m.prefix + sanitizedKey)
	expireAt := time.Now().Add(ttl)

	doc := cacheDocument{
		Key:      prefixedKey,
		Value:    value,
		ExpireAt: expireAt,
	}

	opts := options.Replace().SetUpsert(true)
	_, err := m.collection.ReplaceOne(ctx, bson.M{"_id": prefixedKey}, doc, opts)
	if err != nil {
		m.errors.Add(1)
		return err
	}

	return nil
}

func (m *MongoDBCache) Delete(ctx context.Context, key string) error {
	// Sanitize the key before using it in database queries
	sanitizedKey := SanitizeKey(key)
	prefixedKey := prefixedKey(m.prefix + sanitizedKey)

	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": prefixedKey})
	if err != nil {
		m.errors.Add(1)
		return err
	}

	return nil
}

func (m *MongoDBCache) Close() error {
	// MongoDB client is managed externally, so we don't close it here
	return nil
}

func (m *MongoDBCache) Stats() CacheStats {
	hits := m.hits.Load()
	misses := m.misses.Load()
	total := hits + misses

	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	// Get cache size (approximate)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := m.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		m.errors.Add(1)
	}

	return CacheStats{
		Hits:    hits,
		Misses:  misses,
		Errors:  m.errors.Load(),
		HitRate: hitRate,
		Size:    count,
	}
}

package redis

import (
	"context"
	common "url/db"

	redis "github.com/redis/go-redis/v9"
)

type RedisHelper struct {
	redisClient *redis.Client
}

var ctx = context.Background()

func NewRedisHelper() *RedisHelper {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       7,
	})
	return &RedisHelper{
		redisClient: redisClient,
	}
}

func (rh *RedisHelper) Get(key string) (string, error) {
	return rh.redisClient.Get(ctx, key).Result()
}

func (rh *RedisHelper) Set(key string, value string) error {
	return rh.redisClient.Set(ctx, key, value, 0).Err()
}

func (rh *RedisHelper) Delete(key string) error {
	return rh.redisClient.Del(ctx, key).Err()
}

// GetHash retrieves an Entry struct from Redis hash
func (rh *RedisHelper) GetHash(key string) (common.Entry, error) {
	// Get all fields from the hash
	result, err := rh.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return common.Entry{}, err
	}

	// Check if the hash exists
	if len(result) == 0 {
		return common.Entry{}, redis.Nil
	}

	// Convert hash fields back to Entry struct
	entry := common.Entry{
		URL:  result["url"],
		Hash: result["hash"],
		ID:   result["id"],
	}

	return entry, nil
}

// SetHash stores an Entry struct as Redis hash
func (rh *RedisHelper) SetHash(key string, entry common.Entry) error {
	// Convert Entry struct to hash fields
	fields := map[string]any{
		"url":  entry.URL,
		"hash": entry.Hash,
		"id":   entry.ID,
	}

	return rh.redisClient.HSet(ctx, key, fields).Err()
}

// Check if a hash exists
func (rh *RedisHelper) Exists(key string) (bool, error) {
	result, err := rh.redisClient.Exists(ctx, key).Result()
	return result > 0, err
}

func (rh *RedisHelper) Close() error {
	return rh.redisClient.Close()
}

func (rh *RedisHelper) FlushDB() error {
	return rh.redisClient.FlushDB(ctx).Err()
}

func (rh *RedisHelper) CountKeysOfPattern(pattern string) (int64, error) {
	keys, err := rh.redisClient.Keys(ctx, pattern).Result()
	return int64(len(keys)), err
}

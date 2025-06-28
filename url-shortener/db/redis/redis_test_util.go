package redis

import (
	"context"
	"testing"
	common "url/db"
)

type DatabaseTestSuite = common.DatabaseTestSuite

// NewRedisDatabaseTestSuite returns the test suite configuration for RedisDatabase
func NewRedisDatabaseTestSuite() DatabaseTestSuite {
	db := NewRedisDatabase()
	return DatabaseTestSuite{
		Name: "RedisDatabase",
		DB:   db,
		Validate: func(t *testing.T) {
			// Implementation-specific validation for in-memory database
			if db.redisHelper == nil {
				t.Error("Expected redisHelper to be initialized, got nil")
			}
			if db.redisHelper.redisClient == nil {
				t.Error("Expected redisClient to be initialized, got nil")
			}
			if err := db.redisHelper.redisClient.Ping(context.Background()); err.Err() != nil {
				t.Error("Expected redisClient to be able to ping, got error: ", err.Err().Error())
			}
		},
		Cleanup: func() error {
			return db.redisHelper.FlushDB()
		},
		Close: func() error {
			return db.Close()
		},
	}
}

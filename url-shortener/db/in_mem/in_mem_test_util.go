package in_mem

import (
	"testing"
	common "url/db"
)

type DatabaseTestSuite = common.DatabaseTestSuite

// NewInMemoryDatabaseTestSuite returns the test suite configuration for InMemoryDatabase
func NewInMemoryDatabaseTestSuite() DatabaseTestSuite {
	db := NewInMemoryDatabase()
	return DatabaseTestSuite{
		Name: "InMemoryDatabase",
		DB:   db,
		Validate: func(t *testing.T) {
			if db.entries == nil {
				t.Error("Expected entries map to be initialized, got nil")
			}
			if db.ids == nil {
				t.Error("Expected ids map to be initialized, got nil")
			}
		},
	}
}

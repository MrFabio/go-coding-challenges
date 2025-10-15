package tests

import (
	"testing"
	common "url-shortener/db"
	"url-shortener/db/in_mem"
	"url-shortener/db/redis"
)

type Database = common.Database
type DatabaseTestSuite = common.DatabaseTestSuite
type Entry = common.Entry

// TestAllImplementations runs all generic tests against all database implementations
func TestAllImplementations(t *testing.T) {
	suites := []common.DatabaseTestSuite{
		in_mem.NewInMemoryDatabaseTestSuite(),
		redis.NewRedisDatabaseTestSuite(),
	}

	for _, suite := range suites {
		runTestSuite(t, suite)
		if suite.Cleanup != nil {
			if err := suite.Cleanup(); err != nil {
				t.Errorf("Error cleaning up database: %v", err)
			}
		}
		if suite.Close != nil {
			if err := suite.Close(); err != nil {
				t.Errorf("Error closing database: %v", err)
			}
		}
	}
}

// runTestSuite runs all generic tests against a database implementation
func runTestSuite(t *testing.T, suite DatabaseTestSuite) {
	// Initialize database once for all tests in this suite
	db := suite.DB

	t.Run(suite.Name, func(t *testing.T) {
		// Run implementation-specific validation if provided
		if suite.Validate != nil {
			t.Run("ImplementationValidation", func(t *testing.T) {
				suite.Validate(t)
			})
		}

		t.Run("TestAddEntry", func(t *testing.T) {
			testAddEntry(t, db)
		})

		t.Run("TestAddEntryDuplicate", func(t *testing.T) {
			testAddEntryDuplicate(t, db)
		})

		t.Run("TestGetEntry", func(t *testing.T) {
			testGetEntry(t, db)
		})

		t.Run("TestGetEntryNonExistent", func(t *testing.T) {
			testGetEntryNonExistent(t, db)
		})

		t.Run("TestDeleteEntry", func(t *testing.T) {
			testDeleteEntry(t, db)
		})

		t.Run("TestDeleteEntryNonExistent", func(t *testing.T) {
			testDeleteEntryNonExistent(t, db)
		})

		t.Run("TestHasURLHash", func(t *testing.T) {
			testHasURLHash(t, db)
		})

		t.Run("TestMultipleEntries", func(t *testing.T) {
			testMultipleEntries(t, db)
		})
	})
}

// Generic test functions that work with any Database implementation

func testAddEntry(t *testing.T, db Database) {
	url := common.RandomURL()
	initialEntries := db.CountEntries()
	entry := db.AddEntry(url)

	if entry.URL != url {
		t.Errorf("Expected entry URL to be %s, got %s", url, entry.URL)
	}

	if entry.Hash != common.Hash(url) {
		t.Errorf("Expected entry hash to be %s, got %s", common.Hash(url), entry.Hash)
	}

	if entry.ID == "" {
		t.Error("Expected entry ID to be generated, got empty string")
	}

	// Verify entry is stored in database
	if db.CountEntries() != initialEntries+1 {
		t.Errorf("Expected %d entries in database, got %d", initialEntries+1, db.CountEntries())
	}

	// Verify the entry can be retrieved
	storedEntry, _ := db.GetEntry(entry.ID)
	if storedEntry.URL != url {
		t.Errorf("Expected stored entry URL to be %s, got %s", url, storedEntry.URL)
	}
}

func testAddEntryDuplicate(t *testing.T, db Database) {
	url := common.RandomURL()
	initialEntries := db.CountEntries()

	// Add the same URL twice
	entry1 := db.AddEntry(url)
	entry2 := db.AddEntry(url)

	// Both should return the same entry (same hash)
	if entry1.Hash != entry2.Hash {
		t.Errorf("Expected same hash for duplicate URLs, got %s and %s", entry1.Hash, entry2.Hash)
	}

	if entry1.ID != entry2.ID {
		t.Errorf("Expected same ID for duplicate URLs, got %s and %s", entry1.ID, entry2.ID)
	}

	// Should only have one entry in the database
	if db.CountEntries() != initialEntries+1 {
		t.Errorf("Expected %d entries in database for duplicate URLs, got %d", initialEntries+1, db.CountEntries())
	}
}

func testGetEntry(t *testing.T, db Database) {
	// Add an entry
	url := common.RandomURL()
	addedEntry := db.AddEntry(url)

	// Get the entry
	retrievedEntry, _ := db.GetEntry(addedEntry.ID)

	if retrievedEntry.URL != url {
		t.Errorf("Expected retrieved entry URL to be %s, got %s", url, retrievedEntry.URL)
	}

	if retrievedEntry.Hash != addedEntry.Hash {
		t.Errorf("Expected retrieved entry hash to be %s, got %s", addedEntry.Hash, retrievedEntry.Hash)
	}

	if retrievedEntry.ID != addedEntry.ID {
		t.Errorf("Expected retrieved entry ID to be %s, got %s", addedEntry.ID, retrievedEntry.ID)
	}
}

func testGetEntryNonExistent(t *testing.T, db Database) {
	// Try to get a non-existent entry
	entry, _ := db.GetEntry("non-existent-id")

	// Should return empty entry
	if entry.URL != "" {
		t.Errorf("Expected empty URL for non-existent entry, got %s", entry.URL)
	}

	if entry.Hash != "" {
		t.Errorf("Expected empty hash for non-existent entry, got %s", entry.Hash)
	}

	if entry.ID != "" {
		t.Errorf("Expected empty ID for non-existent entry, got %s", entry.ID)
	}
}

func testDeleteEntry(t *testing.T, db Database) {
	// Add an entry
	url := common.RandomURL()
	initialEntries := db.CountEntries()
	entry := db.AddEntry(url)

	// Verify entry exists
	if db.CountEntries() != initialEntries+1 {
		t.Errorf("Expected %d entries before deletion, got %d", initialEntries+1, db.CountEntries())
	}

	// Delete the entry
	db.DeleteEntry(entry.ID)

	// Verify entry is removed
	if db.CountEntries() != initialEntries {
		t.Errorf("Expected %d entries after deletion, got %d", initialEntries, db.CountEntries())
	}

	// Verify entry cannot be retrieved
	retrievedEntry, _ := db.GetEntry(entry.ID)
	if retrievedEntry.URL != "" {
		t.Errorf("Expected empty entry after deletion, got URL: %s", retrievedEntry.URL)
	}
}

func testDeleteEntryNonExistent(t *testing.T, db Database) {
	// Add an entry
	url := common.RandomURL()
	entry := db.AddEntry(url)

	initialEntries := db.CountEntries()

	// Try to delete a non-existent entry
	db.DeleteEntry("non-existent-id")

	// Should not affect the database
	if db.CountEntries() != initialEntries {
		t.Errorf("Expected %d entries after deleting non-existent entry, got %d", initialEntries, db.CountEntries())
	}

	// Original entry should still exist
	retrievedEntry, _ := db.GetEntry(entry.ID)
	if retrievedEntry.URL != url {
		t.Errorf("Expected original entry to still exist with URL %s, got %s", url, retrievedEntry.URL)
	}
}

func testHasURLHash(t *testing.T, db Database) {
	url := common.RandomURL()
	urlHash := common.Hash(url)

	// Initially should not have the hash
	if db.HasURLHash(urlHash) {
		t.Error("Expected hasURLHash to return false for non-existent hash")
	}

	// Add the entry
	entry := db.AddEntry(url)

	// Should now have the hash
	if !db.HasURLHash(urlHash) {
		t.Error("Expected hasURLHash to return true for existing hash")
	}

	// Should also work with the entry's hash
	if !db.HasURLHash(entry.Hash) {
		t.Error("Expected hasURLHash to return true for entry's hash")
	}

	// Test with non-existent hash
	if db.HasURLHash("non-existent-hash") {
		t.Error("Expected hasURLHash to return false for non-existent hash")
	}
}

func testMultipleEntries(t *testing.T, db Database) {
	urls := []string{
		common.RandomURL(),
		common.RandomURL(),
		common.RandomURL(),
		common.RandomURL(),
	}

	entries := make([]Entry, len(urls))
	currentEntries := db.CountEntries()

	// Add multiple entries
	for i, url := range urls {
		entries[i] = db.AddEntry(url)
	}

	// Verify all entries are stored
	if db.CountEntries() != currentEntries+len(urls) {
		t.Errorf("Expected %d entries in database, got %d", currentEntries+len(urls), db.CountEntries())
	}

	// Verify each entry can be retrieved
	for i, url := range urls {
		retrievedEntry, _ := db.GetEntry(entries[i].ID)
		if retrievedEntry.URL != url {
			t.Errorf("Expected entry %d URL to be %s, got %s", i, url, retrievedEntry.URL)
		}
	}

	// Delete one entry
	db.DeleteEntry(entries[0].ID)

	// Verify correct number of remaining entries
	if db.CountEntries() != currentEntries+len(urls)-1 {
		t.Errorf("Expected %d entries after deletion, got %d", currentEntries+len(urls)-1, db.CountEntries())
	}

	// Verify deleted entry cannot be retrieved
	deletedEntry, _ := db.GetEntry(entries[0].ID)
	if deletedEntry.URL != "" {
		t.Errorf("Expected deleted entry to be empty, got URL: %s", deletedEntry.URL)
	}

	// Verify other entries still exist
	for i := 1; i < len(urls); i++ {
		retrievedEntry, _ := db.GetEntry(entries[i].ID)
		if retrievedEntry.URL != urls[i] {
			t.Errorf("Expected remaining entry %d URL to be %s, got %s", i, urls[i], retrievedEntry.URL)
		}
	}
}

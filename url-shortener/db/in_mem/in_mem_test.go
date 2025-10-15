package in_mem

import (
	"testing"
	common "url-shortener/db"
)

func TestInMemoryDatabaseInitialization(t *testing.T) {
	db := NewInMemoryDatabase()

	// Verify the internal state maps are initialized
	if db.entries == nil {
		t.Error("Expected entries map to be initialized, got nil")
	}

	if db.ids == nil {
		t.Error("Expected ids map to be initialized, got nil")
	}

	if len(db.entries) != 0 {
		t.Errorf("Expected empty entries map")
	}

	if len(db.ids) != 0 {
		t.Errorf("Expected empty ids map")
	}
}

func TestInMemoryDatabaseInternalState(t *testing.T) {
	db := NewInMemoryDatabase()

	// Add an entry to the internal state map
	url := common.TestURL
	entry := db.AddEntry(url)

	// Verify the internal state map is consistent
	if db.entries[entry.ID] != entry {
		t.Error("Expected entry to be stored in entries map")
	}

	if db.ids[entry.Hash] != entry.ID {
		t.Errorf("Expected hash mapping to point to entry ID %s, got %s", entry.ID, db.ids[entry.Hash])
	}

	// Verify counts match internal state
	if len(db.entries) != db.CountEntries() {
		t.Errorf("Expected entries count to match internal map length")
	}

	if len(db.ids) != db.CountIds() {
		t.Errorf("Expected ids count to match internal map length")
	}
}

func TestInMemoryDatabaseDeleteInternalState(t *testing.T) {
	db := NewInMemoryDatabase()

	// Add an entry to the internal state map
	url := common.TestURL
	entry := db.AddEntry(url)

	// Verify entry exists in the internal state map
	if _, exists := db.entries[entry.ID]; !exists {
		t.Error("Expected entry to exist in entries map before deletion")
	}

	if _, exists := db.ids[entry.Hash]; !exists {
		t.Error("Expected hash mapping to exist in ids map before deletion")
	}

	// Delete the entry from the internal state map
	db.DeleteEntry(entry.ID)

	// Verify entry is removed from internal state map
	if _, exists := db.entries[entry.ID]; exists {
		t.Error("Expected entry to be removed from entries map after deletion")
	}

	if _, exists := db.ids[entry.Hash]; exists {
		t.Error("Expected hash mapping to be removed from ids map after deletion")
	}
}

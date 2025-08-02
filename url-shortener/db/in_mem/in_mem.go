package in_mem

import (
	"fmt"
	common "url/db"
)

type Entry = common.Entry
type Database = common.Database

// InMemoryDatabase implements Database interface using in-memory maps
type InMemoryDatabase struct {
	// key is id, value is entry
	entries map[string]Entry
	// key is hash, value is id - cache for faster lookup on get by id
	ids map[string]string
}

// NewInMemoryDatabase creates a new in-memory database instance
func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		entries: make(map[string]Entry),
		ids:     make(map[string]string),
	}
}

func (db *InMemoryDatabase) AddEntry(url string) Entry {
	hash := common.Hash(url)
	if db.HasURLHash(hash) {
		// entry already exists
		entry, _ := db.GetEntry(db.ids[hash])
		return entry
	}
	entry := common.GenerateEntry(url, hash)
	db.entries[entry.ID] = entry
	db.ids[hash] = entry.ID

	return entry
}

func (db *InMemoryDatabase) GetEntry(id string) (Entry, error) {
	return db.entries[id], nil
}

func (db *InMemoryDatabase) DeleteEntry(id string) {
	if entry, ok := db.entries[id]; ok {
		delete(db.entries, id)
		delete(db.ids, entry.Hash)
	}
}

func (db *InMemoryDatabase) HasURLHash(hash string) bool {
	_, ok := db.ids[hash]
	return ok
}

func (db *InMemoryDatabase) CountEntries() int {
	return len(db.entries)
}

func (db *InMemoryDatabase) CountIds() int {
	return len(db.ids)
}

func (db *InMemoryDatabase) String() {
	fmt.Println("db:")
	for _, e := range db.entries {
		e.String()
	}
}

func (db *InMemoryDatabase) Close() error {
	clear(db.entries)
	clear(db.ids)

	return nil
}

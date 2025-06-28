package common

// Database defines the interface for URL shortener storage
type Database interface {
	AddEntry(url string) Entry
	GetEntry(id string) (Entry, error)
	DeleteEntry(id string)
	HasURLHash(hash string) bool
	CountEntries() int
	String()
	Close() error
}

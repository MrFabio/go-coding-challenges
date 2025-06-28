package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// An Entry is the generic representation of a URL and its short id
type Entry struct {
	URL  string
	Hash string
	ID   string
}

func (e *Entry) String() {
	fmt.Println(e.ID, e.URL)
}

func Hash(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateId generates a random 6-character string, example: "f4b10g"
func GenerateId() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 6)
	for i := range id {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			// Fallback to a simple approach if crypto/rand fails
			id[i] = chars[i%len(chars)]
		} else {
			id[i] = chars[randomIndex.Int64()]
		}
	}
	return string(id)
}

func GenerateEntry(url string, hash string) Entry {
	return Entry{
		URL:  url,
		Hash: hash,
		ID:   GenerateId(),
	}
}

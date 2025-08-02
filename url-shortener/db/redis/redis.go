package redis

import (
	"fmt"
	"url/config"

	common "url/db"
)

type Entry = common.Entry
type Database = common.Database

// RedisDatabase implements Database interface using Redis
//
// It uses has:* for hash(url)->id and id:* for id->entry
type RedisDatabase struct {
	redisHelper *RedisHelper
}

func NewRedisDatabase(config *config.Config) *RedisDatabase {
	return &RedisDatabase{
		redisHelper: NewRedisHelper(config),
	}
}

func (db *RedisDatabase) AddEntry(url string) Entry {
	hash := common.Hash(url)
	if db.HasURLHash(hash) {
		// entry already exists
		id := db.GetIdByHash(hash)
		entry, _ := db.GetEntry(id)
		return entry
	}

	entry := common.GenerateEntry(url, hash)
	idKey, hashKey := generateKeys(&entry)
	if err := db.redisHelper.SetHash(idKey, entry); err != nil {
		return entry
	}
	if err := db.redisHelper.Set(hashKey, entry.ID); err != nil {
		return entry
	}

	return entry
}

func (db *RedisDatabase) HasURLHash(hash string) bool {
	exists, _ := db.redisHelper.Exists(generateKeyHash(hash))
	return exists
}

func (db *RedisDatabase) GetIdByHash(hash string) string {
	id, _ := db.redisHelper.Get(generateKeyHash(hash))
	return id
}

func (db *RedisDatabase) GetEntry(id string) (Entry, error) {
	return db.redisHelper.GetHash(generateKeyId(id))
}

func (db *RedisDatabase) DeleteEntry(id string) {
	entry, err := db.GetEntry(id)
	if err != nil {
		return
	}
	idKey, hashKey := generateKeys(&entry)
	if err := db.redisHelper.Delete(idKey); err != nil {
		return
	}
	if err := db.redisHelper.Delete(hashKey); err != nil {
		return
	}
}

func (db *RedisDatabase) CountEntries() int {
	count, _ := db.redisHelper.CountKeysOfPattern("id:*")
	return int(count)
}

func (db *RedisDatabase) String() {
	count, _ := db.redisHelper.CountKeysOfPattern("id:*")
	fmt.Printf("RedisDatabase: %d entries\n", count)
}

func (db *RedisDatabase) Close() error {
	return db.redisHelper.Close()
}

func generateKeys(entry *common.Entry) (string, string) {
	return generateKeyId(entry.ID), generateKeyHash(entry.Hash)
}

func generateKeyHash(hash string) string {
	return fmt.Sprintf("hash:%s", hash)
}

func generateKeyId(id string) string {
	return fmt.Sprintf("id:%s", id)
}

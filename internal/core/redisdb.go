package core

import (
	"log"
	"time"

	"github.com/nhtuan0700/godis/internal/config"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

type RedisDB struct {
	dict       map[string]*RedisObj
	expireDict map[string]uint64
	epool      data_structure.EvictionPool
}

func NewRedisDB() *RedisDB {
	return &RedisDB{
		dict:       make(map[string]*RedisObj),
		expireDict: make(map[string]uint64),
		epool:      *data_structure.NewEpool(config.EpoolMaxSize),
	}
}

// In redis, it will define a const for each type, and the RedisObj will have a field to indicate the type of value it holds.
// For simplicity, we just check the type when casting the value.
type RedisObj struct {
	value          any
	lastAccessTime uint32
}

func NewRedisObj(v any) *RedisObj {
	obj := &RedisObj{
		value:          v,
		lastAccessTime: uint32(time.Now().UnixMilli()),
	}

	return obj
}

func (db *RedisDB) Get(key string) *RedisObj {
	if obj, ok := db.dict[key]; ok {
		// delete epxired key in passive mode
		if db.HasExpired(key) {
			db.Delete(key)
			return nil
		}
		obj.lastAccessTime = uint32(time.Now().UnixMilli())
		return obj
	}

	return nil
}

func (db *RedisDB) Set(key string, obj *RedisObj, ttlMs uint64) {
	if len(db.dict) == config.MaxKeyNumber {
		db.evict()
	}

	db.dict[key] = obj

	if ttlMs > 0 {
		db.SetExpiry(key, ttlMs)
	}
}

func (db *RedisDB) Delete(key string) bool {
	delete(db.expireDict, key)

	_, ok := db.dict[key]
	if !ok {
		return false
	}

	delete(db.dict, key)
	return true
}

func (db *RedisDB) GetExpireDict() map[string]uint64 {
	return db.expireDict
}

func (db *RedisDB) SetExpiry(key string, ttl uint64) {
	db.expireDict[key] = uint64(time.Now().UnixMilli()) + ttl
}

func (db *RedisDB) GetExpiry(key string) (uint64, bool) {
	ttl, exist := db.expireDict[key]
	return ttl, exist
}

func (db *RedisDB) HasExpired(key string) bool {
	if ttl, exist := db.expireDict[key]; exist {
		return ttl <= uint64(time.Now().UnixMilli())
	}

	return false
}

func (db *RedisDB) evict() {
	switch config.EvictPolicy {
	case "allkeys-random":
		db.evictRandom()
	case "allkeys-lru":
		db.evictLru()
	}
}

// populateEpool push the new items with sampled size to the pool
func (db *RedisDB) populateEpool() {
	remain := config.LruSampledSize
	for k, v := range db.dict {
		db.epool.Push(k, v.lastAccessTime)
		remain--
		if remain == 0 {
			break
		}
	}

	log.Println("Epool: ")
	for _, item := range db.epool.Pool() {
		log.Println(item.Key(), item.LastAccessTime())
	}
}

func (db *RedisDB) evictLru() {
	db.populateEpool()

	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Println("Trigger LRU eviction")
	for i := 0; i < int(evictCount) && len(db.epool.Pool()) > 0; i++ {
		item := db.epool.Pop()
		log.Println("Delete key ", item.Key())
		db.Delete(item.Key())
	}
}

func (db *RedisDB) evictRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))

	for k := range db.dict {
		if evictCount == 0 {
			break
		}
		log.Println("Trigger evict random")
		log.Println("delete key: ", k)
		evictCount--
		db.Delete(k)
	}
}

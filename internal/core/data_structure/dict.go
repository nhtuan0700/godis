package data_structure

import (
	"log"
	"time"

	"github.com/nhtuan0700/godis/internal/config"
)

type Obj struct {
	Value          any
	LastAccessTime uint32
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func now() uint32 {
	return uint32(time.Now().Unix())
}

func NewDict() *Dict {
	return &Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
}

func (d *Dict) GetExpiredDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) NewObj(k string, v any, ttlMs uint64) *Obj {
	obj := &Obj{
		Value:          v,
		LastAccessTime: now(),
	}

	if ttlMs > 0 {
		d.SetExpiry(k, ttlMs)
	}

	return obj
}

func (d *Dict) SetExpiry(k string, ttlMs uint64) {
	_, ok := d.expiredDictStore[k]
	if !ok {
		HashKeySpaceStat.Expire++
	}
	d.expiredDictStore[k] = uint64(time.Now().UnixMilli()) + ttlMs
}

func (d *Dict) GetExpiry(k string) (uint64, bool) {
	exp, ok := d.expiredDictStore[k]
	return exp, ok
}

func (d *Dict) HasExpired(k string) bool {
	if expired, ok := d.expiredDictStore[k]; ok {
		return expired <= uint64(time.Now().UnixMilli())
	}

	return false
}

func (d *Dict) Get(k string) *Obj {
	if obj, ok := d.dictStore[k]; ok {
		// delete epxired key in passive mode
		if d.HasExpired(k) {
			d.Del(k)
			return nil
		}
		obj.LastAccessTime = now()
		return obj
	}

	return nil
}

func (d *Dict) Set(k string, obj *Obj) {
	if len(d.dictStore) == config.MaxKeyNumber {
		d.evict()
	}

	_, ok := d.dictStore[k]
	if !ok {
		HashKeySpaceStat.Key++
	}
	d.dictStore[k] = obj
}

func (d *Dict) Del(k string) bool {
	if _, ok := d.expiredDictStore[k]; ok {
		delete(d.expiredDictStore, k)
		HashKeySpaceStat.Expire--
	}
	if _, ok := d.dictStore[k]; ok {
		delete(d.dictStore, k)
		HashKeySpaceStat.Key--
		return true
	}

	return false
}

func (d *Dict) evict() {
	switch config.EvictPolicy {
	case "allkeys-random":
		d.evictRandom()
	case "allkeys-lru":
		d.evictLru()
	}
}

// populateEpool push the new items with sampled size to the pool
func (d *Dict) populateEpool() {
	remain := config.LruSampledSize
	for k, v := range d.dictStore {
		epool.Push(k, v.LastAccessTime)
		remain--
		if remain == 0 {
			break
		}
	}

	log.Println("Epool: ")
	for _, item := range epool.pool {
		log.Println(item.key, item.lastAccessTime)
	}
}

func (d *Dict) evictLru() {
	d.populateEpool()

	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Println("Trigger LRU eviction")
	for i := 0; i < int(evictCount) && len(epool.pool) > 0; i++ {
		item := epool.Pop()
		log.Println("Delete key ", item.key)
		d.Del(item.key)
	}
}

func (d *Dict) evictRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))

	for k := range d.dictStore {
		if evictCount == 0 {
			break
		}
		log.Println("Trigger evict random")
		log.Println("delete key: ", k)
		evictCount--
		d.Del(k)
	}
}

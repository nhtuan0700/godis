package data_structure

import (
	"sort"

	"github.com/nhtuan0700/godis/internal/config"
)

type EvictionCandidate struct {
	key            string
	lastAccessTime uint32
}

type EvictionPool struct {
	pool []*EvictionCandidate
}

// ByLastAccessTime used for sort the pool.
type ByLastAccessTime []*EvictionCandidate

func (b ByLastAccessTime) Len() int {
	return len(b)
}

func (b ByLastAccessTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByLastAccessTime) Less(i, j int) bool {
	return b[i].lastAccessTime < b[j].lastAccessTime
}

// Push add a new item to the pool, maintains the lassAccessTime accenting order (old items are on the left).
// If pool size > EpoolMaxSize, remove the newest item.
func (p *EvictionPool) Push(key string, lastAccessTime uint32) {
	newItem := &EvictionCandidate{
		key:            key,
		lastAccessTime: lastAccessTime,
	}

	// Note: In redis implementation, it does not explicity check if a key is already in the eviction pool
	// before attempting to insert it. This could lead to a key being in the pool twice
	// if it's sampled and inserted a second time. However, since the pool is very small (EpoolMaxSize is 16)
	// and the random sampling is just small fraction of the total keys, the probabiltity of this happening is extremely low.
	// Ref: https://github.com/redis/redis/blob/unstable/src/evict.c#L126
	exist := false
	for i := 0; i < len(p.pool); i++ {
		if p.pool[i].key == key {
			exist = true
			p.pool[i].lastAccessTime = lastAccessTime
		}
	}

	if !exist {
		p.pool = append(p.pool, newItem)
	}

	sort.Sort(ByLastAccessTime(p.pool))
	if len(p.pool) > config.EpoolMaxSize {
		lastIndex := len(p.pool) - 1
		p.pool = p.pool[:lastIndex]
	}
}

// Remove the oldest item in the pool.
func (p *EvictionPool) Pop() *EvictionCandidate {
	if len(p.pool) == 0 {
		return nil
	}

	oldestItem := p.pool[0]
	p.pool = p.pool[1:]

	return oldestItem
}

func newEpool(n int) *EvictionPool {
	return &EvictionPool{
		pool: make([]*EvictionCandidate, n),
	}
}

var epool *EvictionPool = newEpool(0)

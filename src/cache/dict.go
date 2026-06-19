package cache

import (
	"hash/fnv"
	"im_cacher/src/aof"
	"im_cacher/src/protocol/resp"
	"sync"
	"time"
)

type dictEntry struct {
	Key            string
	Value          resp.Value
	ExpiresAt      int64
	LastAccessedAt int64
	Next           *dictEntry
}

type Dict struct {
	buckets    []*dictEntry
	size       uint32
	sizeMask   uint32
	maxEntries int
	currentLen int
	mu         sync.RWMutex
	aof        *aof.AOF
}

func NewDict(initialSize uint32, aof *aof.AOF, maxEntries int) *Dict {
	return &Dict{
		buckets:    make([]*dictEntry, initialSize),
		size:       initialSize,
		sizeMask:   initialSize - 1,
		mu:         sync.RWMutex{},
		maxEntries: maxEntries,
		aof:        aof,
	}
}

func (d *Dict) calculateBucketIndex(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	hashValue := h.Sum32()

	return hashValue & d.sizeMask
}

func (d *Dict) Set(key string, value resp.Value, ttl int) {
	index := d.calculateBucketIndex(key)

	var expiresAt int64 = 0
	if ttl > 0 {
		expiresAt = time.Now().Add(time.Second * time.Duration(ttl)).UnixNano()
	}

	current := d.buckets[index]
	for current != nil {
		if current.Key == key {
			// if u want to set value without restarting ttl -> TTL = -1
			// for example : for rate limiting you need to increase counter of requests
			// and the expire date must remain
			// in such cases ttl must be -1
			if ttl == -1 {
				d.mu.Lock()
				current.Value = value
				current.LastAccessedAt = time.Now().UnixNano()
				d.mu.Unlock()
				return
			}
			d.mu.Lock()
			current.Value = value
			current.ExpiresAt = expiresAt
			current.LastAccessedAt = time.Now().UnixNano()
			d.mu.Unlock()
			return
		}
		current = current.Next
	}

	for d.currentLen >= d.maxEntries {
		d.evictApproximatedLRU()
	}

	newEntry := &dictEntry{
		Key:            key,
		Value:          value,
		ExpiresAt:      expiresAt,
		LastAccessedAt: time.Now().UnixNano(),
	}

	d.mu.Lock()
	newEntry.Next = d.buckets[index]
	d.buckets[index] = newEntry
	d.currentLen++
	d.mu.Unlock()
}

func (d *Dict) Get(key string) (*resp.Value, bool) {
	index := d.calculateBucketIndex(key)

	dictElem := d.buckets[index]
	for dictElem != nil {
		if dictElem.Key == key {
			if dictElem.ExpiresAt < time.Now().UnixNano() {
				return nil, false
			}
			dictElem.LastAccessedAt = time.Now().UnixNano()
			return &dictElem.Value, true
		}

		dictElem = dictElem.Next
	}

	return nil, false
}

func (d *Dict) Del(key string) {
	index := d.calculateBucketIndex(key)

	dictElem := d.buckets[index]
	var preElem *dictEntry
	for dictElem != nil {
		if dictElem.Key == key {
			if preElem == nil {
				d.mu.Lock()
				d.currentLen--
				d.buckets[index] = dictElem.Next
				d.mu.Unlock()
				return
			}
			d.mu.Lock()
			d.currentLen--
			preElem.Next = dictElem.Next
			d.mu.Unlock()
			return
		}

		preElem = dictElem
		dictElem = dictElem.Next
	}
}

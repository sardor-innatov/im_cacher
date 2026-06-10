package cache

import (
	"im_cacher/src/protocol/resp"
	"math/rand/v2"
	"time"
)

func (d *Dict) StartEvictionLoop(interval time.Duration) (Close func()) {

	var stopCh chan struct{}

	go func(stopCh chan struct{}) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				d.evictExpired()
			case <-stopCh:
				return
			}
		}
	}(stopCh)

	return func() {
		stopCh <- struct{}{}
	}
}

func (d *Dict) deleteKeyWithoutLock(key string) {
	index := d.calculateBucketIndex(key)

	dictElem := d.buckets[index]
	var preElem *dictEntry
	for dictElem != nil {

		if dictElem.ExpiresAt == 0 && time.Now().UnixNano() < dictElem.ExpiresAt {
			break
		}

		if dictElem.Key == key {

			if preElem == nil {
				d.buckets[index] = dictElem.Next
				break
			}
			preElem.Next = dictElem.Next
			break
		}

		preElem = dictElem
		dictElem = dictElem.Next
	}
}

func (d *Dict) evictExpired() {
	var keysToDelete []string

	var dictElem *dictEntry
	d.mu.RLock()
	for _, bucket := range d.buckets {
		dictElem = bucket
		for dictElem != nil {
			if dictElem.ExpiresAt > 0 && dictElem.ExpiresAt < time.Now().UnixNano() {
				keysToDelete = append(keysToDelete, dictElem.Key)
			}
			dictElem = dictElem.Next
		}
	}
	d.mu.RUnlock()

	if len(keysToDelete) > 0 {
		d.mu.Lock()
		for _, key := range keysToDelete {
			d.deleteKeyWithoutLock(key)
			d.currentLen--

			d.aof.AppendDelCommand(key)
		}
		d.mu.Unlock()
	}

}

func (d *Dict) evictApproximatedLRU() {
	if d.currentLen == 0 {
		return
	}

	const sampleSize = 5
	var bestCandidate *dictEntry
	var oldestTime int64 = 0

	for i := 0; i < sampleSize; i++ {
		randomIndex := rand.Uint32() & d.sizeMask
		bucket := d.buckets[randomIndex]

		if bucket == nil {
			continue
		}

		candidate := bucket

		if bestCandidate == nil || candidate.LastAccessedAt < oldestTime {
			bestCandidate = candidate
			oldestTime = candidate.LastAccessedAt
		}
	}

	if bestCandidate != nil {
		d.deleteKeyWithoutLock(bestCandidate.Key)
		d.currentLen--

		d.aof.AppendDelCommand(bestCandidate.Key)
	}
}

func (d *Dict) Iterator(rawValue *resp.Value) {

	arrValue := rawValue.Array()

	switch arrValue[0].String() {
	case "set":

		ttl := time.Second * time.Duration(arrValue[1].Integer())

		key := arrValue[2].String()

		d.Set(key, arrValue[3], ttl)

	case "del":

		key := arrValue[2].String()

		d.Del(key)
	}

}

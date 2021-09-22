package xcache

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex

	// The item's key.
	key interface{}
	// The item's data.
	data interface{}
	// How long will the item live in the cache when not being accessed/kept alive.
	lifeSpan time.Duration

	// Creation timestamp.
	createdOn time.Time
	// Last access timestamp.
	accessedOn time.Time
	// How often the item was accessed.
	accessCount int64

	// Callback method triggered right before removing the item from the cache
	aboutToExpire []func(key interface{})
}

//func NewCacheItem(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
//	t := time.Now()
//	return &CacheItem{
//		key:           key,
//		lifeSpan:      lifeSpan,
//		createdOn:     t,
//		accessedOn:    t,
//		accessCount:   0,
//		aboutToExpire: nil,
//		data:          data,
//	}
//}

func NewCacheItem(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	t := time.Now()
	return &CacheItem{
		key:           key,
		lifeSpan:      lifeSpan,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
		aboutToExpire: nil,
		data:          data,
	}
}

func (i *CacheItem) KeepLive() {
	i.Lock()
	defer i.Unlock()
	i.accessedOn = time.Now()
	i.accessCount++
}

func (i *CacheItem) LifeSpan() time.Duration {
	//i.RUnlock()
	//defer i.Unlock()
	return i.lifeSpan
}

func (i *CacheItem) AccessedOn() time.Time {
	i.RLock()
	defer i.RUnlock()
	return i.accessedOn
}

func (i *CacheItem) CreateOn() time.Time {
	return i.createdOn
}

func (i *CacheItem) AccessedCount() int64 {
	i.RLock()
	defer i.RUnlock()
	return i.accessCount
}

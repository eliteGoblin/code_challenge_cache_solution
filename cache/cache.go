package cache

import (
	"sync"
	"time"
)

type DistributedCache struct {
	rwmutex sync.RWMutex
	data    map[string]interface{}
}

func New() *DistributedCache {
	return &DistributedCache{
		data: make(map[string]interface{}),
	}
}

func (dc *DistributedCache) Value(key string) (interface{}, error) {
	// simulate 100ms roundtrip to the distributed cache
	time.Sleep(100 * time.Millisecond)
	dc.rwmutex.RLock()
	defer dc.rwmutex.RUnlock()
	return dc.data[key], nil
}

func (dc *DistributedCache) Store(key string, value interface{}) error {
	// simulate 100ms roundtrip to the distributed cache
	time.Sleep(100 * time.Millisecond)
	dc.rwmutex.Lock()
	defer dc.rwmutex.Unlock()
	dc.data[key] = value
	return nil
}

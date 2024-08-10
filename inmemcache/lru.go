package inmemcache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

const (
	// DefaultCapacity defines the default capacity of the in-memory cache
	DefaultCapacity = 1000
)

// LRUCache represents a thread-safe least-recently-used (LRU) cache
type LRUCache struct {
	mutex    sync.Mutex
	mapKV    map[string]*list.Element
	list     *list.List
	capacity int
}

// element is a key-value pair stored in the LRUCache
type element struct {
	key   string
	value interface{}
}

// New initializes and returns a new LRUCache with the specified capacity.
// Returns an error if the capacity is less than or equal to zero.
func New(capacity int) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, errors.New("invalid capacity: must be greater than zero")
	}

	return &LRUCache{
		mapKV:    make(map[string]*list.Element),
		list:     list.New(),
		capacity: capacity,
	}, nil
}

// Value retrieves the value associated with the given key from the cache.
// If the key exists, the associated element is moved to the front (most recently used).
// Returns nil if the key is not found.
func (cache *LRUCache) Value(key string) (interface{}, error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	node, ok := cache.mapKV[key]
	if !ok {
		return nil, nil // Cache miss, return nil without error
	}

	e, ok := node.Value.(*element)
	if !ok {
		return nil, fmt.Errorf("invalid value type in Value: %+v", node.Value)
	}

	// Move the accessed element to the front of the list
	cache.list.MoveToFront(node)
	return e.value, nil
}

// Store adds or updates the value associated with the given key in the cache.
// If the key already exists, its value is updated and the element is moved to the front.
// If the cache exceeds its capacity, the least recently used element is removed.
func (cache *LRUCache) Store(key string, value interface{}) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// If the key already exists, update the value and move it to the front
	if node, ok := cache.mapKV[key]; ok {
		node.Value = cache.newElement(key, value)
		cache.list.MoveToFront(node)
		return nil
	}

	// If the cache is full, remove the least recently used element (from the back)
	if cache.list.Len() >= cache.capacity {
		if err := cache.removeOldest(); err != nil {
			return err
		}
	}

	// Add the new element to the front of the list
	node := cache.list.PushFront(cache.newElement(key, value))
	cache.mapKV[key] = node
	return nil
}

// removeOldest removes the oldest (least recently used) element from the cache.
// This method is only called when the cache exceeds its capacity.
func (cache *LRUCache) removeOldest() error {
	tail := cache.list.Back()
	if tail == nil {
		return errors.New("cache is empty, cannot remove element")
	}

	e, ok := tail.Value.(*element)
	if !ok {
		return fmt.Errorf("invalid value type in removeOldest: %+v", tail.Value)
	}

	cache.list.Remove(tail)
	delete(cache.mapKV, e.key)
	return nil
}

// newElement creates a new key-value pair element.
func (cache *LRUCache) newElement(key string, value interface{}) *element {
	return &element{
		key:   key,
		value: value,
	}
}

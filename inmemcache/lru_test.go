package inmemcache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New(0)
	assert.NotNil(t, err)
	_, err = New(-1)
	assert.NotNil(t, err)
	_, err = New(1)
	assert.Nil(t, err)
}

func TestStore(t *testing.T) {
	for capacity := 1; capacity <= 5; capacity += 1 {
		cache, err := New(capacity)
		assert.Nil(t, err)
		// cache will only store last one, because capacity is only one
		for i := 0; i < 100; i++ {
			err = cache.Store(fmt.Sprintf("%d", i), i)
			assert.Nil(t, err)
		}
	}
}

func TestGetEmpty(t *testing.T) {
	cache, err := New(10)
	assert.Nil(t, err)
	for i := 0; i <= 10; i++ {
		value, err := cache.Value(fmt.Sprintf("%d", i))
		assert.Nil(t, err)
		assert.Nil(t, value)
	}
}

var setGetTests = []struct {
	capacity int         // capacity of cache
	key      string      // key try to get
	found    bool        // if found
	value    interface{} // value of key
}{
	{
		capacity: 1,
		key:      "99",
		found:    true,
		value:    "99",
	},
	{
		capacity: 1,
		key:      "100",
		found:    false,
	},
	{
		capacity: 1,
		key:      "98",
		found:    false,
	},
	{
		capacity: 2,
		key:      "99",
		found:    true,
		value:    "99",
	},
	{
		capacity: 2,
		key:      "98",
		found:    true,
		value:    "98",
	},
	{
		capacity: 2,
		key:      "97",
		found:    false,
	},
	{
		capacity: 2,
		key:      "100",
		found:    false,
	},
}

func TestStoreAndGet(t *testing.T) {
	for i, test := range setGetTests {
		t.Run(fmt.Sprintf("case %d", i),
			func(t *testing.T) {
				cache, err := New(test.capacity)
				assert.Nil(t, err)
				// store key [0, 99] to cache
				for i := 0; i < 100; i++ {
					key := fmt.Sprintf("%d", i)
					value := key
					err = cache.Store(key, value)
					assert.Nil(t, err)
				}
				v, err := cache.Value(test.key)
				assert.Nil(t, err)
				if test.found {
					assert.NotNil(t, v)
				} else {
					assert.Nil(t, v)
				}
			})
	}
}

func TestStoreAndGetNil(t *testing.T) {
	cache, err := New(5)
	assert.Nil(t, err)
	err = cache.Store("nil", nil)
	assert.Nil(t, err)
	err = cache.Store("non-nil", 123)
	assert.Nil(t, err)
	value, err := cache.Value("nil")
	assert.Nil(t, err)
	assert.True(t, value == nil)
	value, err = cache.Value("non-nil")
	assert.Nil(t, err)
	assert.True(t, value != nil)
	value, err = cache.Value("not-exist")
	assert.Nil(t, err)
	assert.Nil(t, value)
}

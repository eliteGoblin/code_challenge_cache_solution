package datasource

import (
	"fmt"
)

type DataRetrieve struct {
	database      Cache
	externalCache Cache
	localCache    Cache
}

func NewDataRetrieve(
	localCache Cache,
	externalCache Cache,
	database Cache) *DataRetrieve {
	return &DataRetrieve{
		database:      database,
		externalCache: externalCache,
		localCache:    localCache,
	}
}

func (selfPtr *DataRetrieve) Value(key string) (value interface{}, err error) {
	value, err = readThroughCaches(
		key, selfPtr.localCache, selfPtr.externalCache, selfPtr.database)
	if err != nil {
		return nil, fmt.Errorf("error getting value of %s: %w", key, err)
	}
	if notExist(value) {
		return nil, nil
	}
	return value, nil
}

func readThroughCaches(key string, caches ...Cache) (value interface{}, err error) {
	if len(caches) == 0 || caches[0] == nil {
		return nonExist(0), nil
	}
	value, err = caches[0].Value(key)
	if err != nil {
		// fail immediately to prevent overload low level caches
		return nil, fmt.Errorf("error getting value of %s: %w", key, err)
	}
	if value != nil {
		return value, nil
	}
	// key not found and does not know if exist
	// try to find in low-level caches
	value, err = readThroughCaches(key, caches[1:]...)
	if err != nil {
		return nil, fmt.Errorf("error getting value of %s: %w", key, err)
	}

	if len(caches) > 1 {
		// last cache is source of data, does not need to store it
		err = caches[0].Store(key, value)
		if err != nil {
			return nil, fmt.Errorf("error getting value of %s: %w", key, err)
		}
	}
	return value, nil
}

func notExist(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(nonExist)
	return ok
}

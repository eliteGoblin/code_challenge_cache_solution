package datasource

import (
	"cache_solution/datasource/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewDataRetrieve(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	databaseMock := mock.NewMockCache(mockCtrl)
	cacheMock := mock.NewMockCache(mockCtrl)
	localCache := mock.NewMockCache(mockCtrl)
	retrieve := NewDataRetrieve(databaseMock, cacheMock, localCache)
	assert.NotNil(t, retrieve)
}

func TestNotExist(t *testing.T) {
	var v interface{}
	v = nil
	assert.False(t, notExist(v))
	v = "random"
	assert.False(t, notExist(v))
	v = nonExist(123)
	assert.True(t, notExist(v))
}

// query a key only exist in database
func TestReadThroughCacheCacheNonEmpty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	database := mock.NewMockCache(mockCtrl)
	external := mock.NewMockCache(mockCtrl)
	local := mock.NewMockCache(mockCtrl)

	// call sequence
	// 1. local: value
	// 2. external: value
	// 3. database: value(found the key)
	// 4. external: store
	// 5. local: store
	// return value to caller
	localValue := local.EXPECT().Value("key").Return(nil, nil)
	externalValue := external.EXPECT().Value("key").Return(nil, nil).After(localValue)
	databaseValue := database.EXPECT().Value("key").Return("value", nil).After(externalValue)
	externalStore := external.EXPECT().Store("key", "value").Return(nil).After(databaseValue)
	local.EXPECT().Store("key", "value").Return(nil).After(externalStore)
	value, err := readThroughCaches("key", local, external, database)
	assert.Nil(t, err)
	assert.Equal(t, "value", value)
}

// query a key exist in database and external cache
func TestReadThroughCacheLocal(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	external := mock.NewMockCache(mockCtrl)
	local := mock.NewMockCache(mockCtrl)

	// call sequence
	// 1. local: value
	// 2. external: value
	// 3. database: value(found the key)
	// 4. external: store
	// 5. local: store
	// return value to caller
	localValue := local.EXPECT().Value("key").Return(nil, nil)
	externalValue := external.EXPECT().Value("key").Return("value", nil).After(localValue)
	local.EXPECT().Store("key", "value").Return(nil).After(externalValue)
	value, err := readThroughCaches("key", local, external, nil)
	assert.Nil(t, err)
	assert.Equal(t, "value", value)
}

// query a non-exist key
func TestReadThroughCacheCacheEmpty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	database := mock.NewMockCache(mockCtrl)
	external := mock.NewMockCache(mockCtrl)
	local := mock.NewMockCache(mockCtrl)

	// call sequence
	// 1. local: value
	// 2. external: value
	// 3. database: value, key still not found
	// 4. external: store
	// 5. local: store
	// return value to caller
	nonExistMark := nonExist(0)
	localValue := local.EXPECT().Value("key").Return(nil, nil)
	externalValue := external.EXPECT().Value("key").Return(nil, nil).After(localValue)
	databaseValue := database.EXPECT().Value("key").Return(nil, nil).After(externalValue)
	externalStore := external.EXPECT().Store("key", nonExistMark).Return(nil).After(databaseValue)
	local.EXPECT().Store("key", nonExistMark).Return(nil).After(externalStore)
	value, err := readThroughCaches("key", local, external, database)
	assert.Nil(t, err)
	assert.Equal(t, nonExistMark, value)
}

func TestValueNonExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	local := mock.NewMockCache(mockCtrl)
	local.EXPECT().Value("key").Return("value", nil)
	retriever := NewDataRetrieve(local, nil, nil)
	value, err := retriever.Value("key")
	assert.Nil(t, err)
	assert.Equal(t, "value", value)
}

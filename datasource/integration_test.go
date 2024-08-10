//go:build integration
// +build integration

package datasource

import (
	"cache_solution/cache"
	"cache_solution/database"
	"cache_solution/inmemcache"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	database := database.New()
	external := cache.New()
	local, err := inmemcache.New(5)
	assert.Nil(t, err)
	retriever := NewDataRetrieve(local, external, database)
	// setup database
	database.Store("k1", "v1")
	value, err := retriever.Value("k1")
	assert.Nil(t, err)
	assert.Equal(t, "v1", value)
	value, err = retriever.Value("non-exist-key")
	assert.Nil(t, err)
	assert.Nil(t, value)
}

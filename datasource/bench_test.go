package datasource

import (
	"cache_solution/cache"
	"cache_solution/database"
	"cache_solution/inmemcache"
	"cache_solution/pkg/keygen"
	"testing"
)

func BenchmarkDataRetrieveValue(b *testing.B) {
	db := database.New()
	database.FillDatabase(db)
	external := cache.New()
	local, _ := inmemcache.New(inmemcache.DefaultCapacity)
	retriever := NewDataRetrieve(local, external, db)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			retriever.Value(keygen.RandomKey(0, 9))
		}
	})
}

func BenchmarkDataRetrieveValueLRU(b *testing.B) {
	db := database.New()
	database.FillDatabase(db)
	external := cache.New()
	local, _ := inmemcache.New(3)
	retriever := NewDataRetrieve(local, external, db)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			retriever.Value(keygen.RandomKey(0, 9))
		}
	})
}

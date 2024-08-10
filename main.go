package main

import (
	"cache_solution/cache"
	"cache_solution/database"
	"cache_solution/datasource"
	"cache_solution/inmemcache"
	"cache_solution/pkg/keygen"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize the database and fill it with initial data
	db := database.New()
	database.FillDatabase(db)
	logrus.Info("Database initialization completed")

	// Initialize the local in-memory cache with default capacity
	localCache, err := inmemcache.New(inmemcache.DefaultCapacity)
	if err != nil {
		logrus.Fatalf("Failed to initialize local cache: %+v", err)
	}

	// Initialize the external distributed cache
	externalCache := cache.New()

	// Create the data retriever which will handle cache lookups
	retriever := datasource.NewDataRetrieve(localCache, externalCache, db)

	// Define the number of concurrent workers and iterations per worker
	const (
		numWorkers    = 10
		numIterations = 50
	)

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		// Launch workers as goroutines
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				startTime := time.Now()

				// Generate a random key for retrieval
				key := keygen.RandomKey(0, 9)
				value, err := retriever.Value(key)
				if err != nil {
					logrus.Errorf("Worker %d: Error retrieving key '%s': %+v", workerID, key, err)
					continue
				}

				// Log the response time and result
				logrus.Infof("Worker %d: Request '%s', response '%s', time: %.2f ms",
					workerID, key, value.(string), float64(time.Since(startTime))/float64(time.Millisecond))
			}
		}(i)
	}

	// Wait for all workers to complete
	wg.Wait()
	logrus.Info("All workers have completed their tasks")
}

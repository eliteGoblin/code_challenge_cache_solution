# cache_solution

## Summary

cache_solution implement cache service which have 3 layers:

*  local in-memory cache
*  external, distributed cache(simulated with R/W latency of 100ms)
*  database (simulated with R/W latency of 100ms)

Use read/write through policy to load cache, steps:

1. See if local cache hit, if miss
2. See if distributed cache hit: If hit, populate local cache and return
3. If miss, query database
4. If database hit, populate distributed and local cache with key/value pair
5. If data not exist, mark keys in distributed and local cache as "NON-exist"

Others:
*  Dependency using go module
*  Error propagation using golang.org/x/xerrors

## Build & Run

```shell
go build .
./cache_solution
```

## Test

### Unit testing

```shell
go test -race -cover ./...
```

```
?   	cache_solution	[no test files]
ok  	cache_solution/datasource	1.016s	coverage: 80.8% of statements
?   	cache_solution/datasource/mock	[no test files]
?   	cache_solution/infra/cache	[no test files]
?   	cache_solution/infra/database	[no test files]
ok  	cache_solution/infra/inmemcache	1.028s	coverage: 83.8% of statements
```

### Integration test

```shell
go test -tags=integration -race ./...
```

```
?   	cache_solution	[no test files]
ok  	cache_solution/datasource	2.916s
?   	cache_solution/datasource/mock	[no test files]
?   	cache_solution/infra/cache	[no test files]
?   	cache_solution/infra/database	[no test files]
ok  	cache_solution/infra/inmemcache	(cached)
```

## Benchmark

```shell
cd datasource
go test -cpu=4 -bench . -benchtime=1m
```

### Unlimited LRU Capacity

```
goos: linux
goarch: amd64
pkg: cache_solution/datasource
BenchmarkDataRetrieveValue-4   	200000000	       525 ns/op
PASS
ok  	cache_solution/datasource	149.238s
```

### Limited LRU Capacity

* test with: 
  +  LRU fit for 3 keys
  +  Unlimited distributed cache
  +  Total data 10 keys
  +  Total random data(no pattern of LRU)
  
```golang
go test -cpu=4 -bench=DataRetrieveValueLRU -benchtime=1m
```

```
goos: linux
goarch: amd64
pkg: cache_solution/datasource
BenchmarkDataRetrieveValueLRU-4   	    5000	  17118441 ns/op (17ms)
PASS
ok  	cache_solution/datasource	145.243s
```
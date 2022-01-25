# go-wlru
Thread-safe LRU cache with permanency and context-based expiration

[![Go](https://github.com/heucuva/go-wlru/actions/workflows/go.yml/badge.svg)](https://github.com/heucuva/go-wlru/actions/workflows/go.yml)


## Operational Complexity (Time)

| Operation | Best        | Average     | Worst       |
|-----------|-------------|-------------|-------------|
| Access    | Θ(1)        | Θ(1)        | O(1)        |
| Search    | Θ(1)        | Θ(1)        | O(n)        |
| Insertion | Θ(1)        | Θ(1)        | O(n)        |
| Deletion  | Θ(1)        | Θ(1)        | O(n)        |
| Snapshot  | Θ(n)        | Θ(n)        | Θ(n)        |

## Operation Complexity (Space)

| Complexity | Value           |
|------------|-----------------|
| Best       | Ω(2n)           |
| Average    | Ω(2n)           |
| Worst      | Ω(n + n log(n)) |

## Usage
This is a simple example LRU cache structure made with API request lookup caching in mind. If you decide to use this, do so at your own peril.

## Thread Safety
It should be thread-safe on all operations.

## Benchmarks

### BenchmarkSet100k

```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkSet100k$ github.com/heucuva/go-wlru

goos: linux
goarch: amd64
pkg: github.com/heucuva/go-wlru
cpu: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz
BenchmarkSet100k-4   	1000000000	         0.05225 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/heucuva/go-wlru	0.436s
```

As long as the data set is relatively small in size (<131072 items), the performance stays fairly good. Once this size is passed, the time to perform access/search/insert/delete operations increase significantly - this is due to the design of the underlying `sync.Map` implementation, as shown below:

### BenchmarkSet1M

```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkSet1M$ github.com/heucuva/go-wlru

goos: linux
goarch: amd64
pkg: github.com/heucuva/go-wlru
cpu: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz
BenchmarkSet1M-4   	1000000000	         0.7867 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/heucuva/go-wlru	42.616s
```
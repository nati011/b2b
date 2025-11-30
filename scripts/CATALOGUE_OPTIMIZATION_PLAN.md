# Catalogue API Optimization Plan

## Benchmark Results Summary

### Current Performance Metrics

| Benchmark | Time/Op | Memory/Op | Allocs/Op | Throughput |
|-----------|---------|-----------|-----------|------------|
| HTTP Handler (Sequential) | ~2.57ms | 23.7 KB | 222 | ~389 req/s |
| HTTP Handler (Parallel) | ~2.85ms | 47.5 KB | 400 | ~351 req/s |
| Service Layer (Sequential) | ~910μs | 4.1 KB | 73 | ~1,099 ops/s |
| Service Layer (Parallel) | ~475μs | 6.1 KB | 90 | ~2,105 ops/s |

### Key Findings

1. **Service layer is 2.8x faster** than HTTP handler (910μs vs 2.57ms)
2. **Parallel service is 2x faster** than sequential (475μs vs 910μs)
3. **High allocation overhead** in HTTP handler (222 allocations vs 73)
4. **Parallel HTTP performs worse** than sequential (likely due to contention)
5. **Memory usage is high** - 23.7KB per request in HTTP handler

## Optimization Opportunities

### 1. Database Query Optimization (High Priority)

#### Current Issues:
- Multiple sequential database queries
- N+1 query problem in configurable products loop
- No query result caching
- Potential missing database indexes

#### Recommendations:
```sql
-- Add indexes for common queries
CREATE INDEX IF NOT EXISTS idx_product_distributor_active 
ON product(distributor_id, is_active) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_configurable_product_active 
ON configurable_product(is_available) WHERE is_available = true;

CREATE INDEX IF NOT EXISTS idx_product_category 
ON product_category(product_id, category_id);
```

#### Code Changes:
- Use JOIN queries instead of multiple sequential queries
- Implement batch fetching for configurable products
- Add query result caching with TTL
- Use database connection pooling optimization

### 2. Reduce Memory Allocations (High Priority)

#### Current Issues:
- 222 allocations per HTTP request
- Multiple slice allocations in loops
- Unnecessary data copying

#### Recommendations:
- Pre-allocate slices with known capacity
- Use object pooling for frequently allocated objects
- Reduce intermediate data structures
- Reuse buffers where possible

#### Code Changes:
```go
// Before
var resp []GetCatalogueResponse
for _, j := range cp.List {
    var configurables []CatalogueResponse
    // ...
}

// After
resp := make([]GetCatalogueResponse, 0, len(cp.List))
for _, j := range cp.List {
    configurables := make([]CatalogueResponse, 0, len(j.Products))
    // ...
}
```

### 3. Implement Response Caching (Medium Priority)

#### Strategy:
- Cache full catalogue response for 5-10 minutes
- Invalidate on product updates
- Use Redis or in-memory cache
- Cache key: "catalogue:all"

#### Implementation:
```go
type CachedCatalogueService struct {
    service Provider
    cache   *cache.Cache
}

func (c *CachedCatalogueService) GetAll(ctx context.Context) (GetAllCatalogueResponse, error) {
    cacheKey := "catalogue:all"
    if cached, found := c.cache.Get(cacheKey); found {
        return cached.(GetAllCatalogueResponse), nil
    }
    
    resp, err := c.service.GetAll(ctx)
    if err == nil {
        c.cache.Set(cacheKey, resp, 5*time.Minute)
    }
    return resp, err
}
```

### 4. Optimize Data Transformation (Medium Priority)

#### Current Issues:
- Multiple loops transforming the same data
- Unnecessary image array copying
- Redundant map operations

#### Recommendations:
- Combine transformation loops where possible
- Use pointers for large structs
- Minimize data copying
- Use more efficient data structures

### 5. Parallel Processing (Low Priority)

#### Current Issues:
- Sequential processing of configurable products and standalone products
- No parallel fetching of related data

#### Recommendations:
- Use goroutines for independent operations
- Parallel fetch configurable products and standalone products
- Use sync.WaitGroup for coordination

#### Code Changes:
```go
var wg sync.WaitGroup
var cpResp []GetCatalogueResponse
var prResp GetAllResponse
var cpErr, prErr error

wg.Add(2)

// Fetch configurable products in parallel
go func() {
    defer wg.Done()
    cp, err := c.ConfigurableProductservice.Catalogue(ctx)
    // ... process cp
}()

// Fetch standalone products in parallel
go func() {
    defer wg.Done()
    pr, err := c.ProductService.Catalogue(ctx)
    // ... process pr
}()

wg.Wait()
```

### 6. HTTP Handler Optimization (Low Priority)

#### Current Issues:
- JSON encoding overhead
- Response serialization could be optimized
- No response compression

#### Recommendations:
- Enable gzip compression for responses
- Use more efficient JSON encoding
- Consider protocol buffers for internal APIs
- Add response size limits

## Implementation Priority

### Phase 1: Quick Wins (1-2 days)
1. ✅ Fix migration path issue (DONE)
2. Add database indexes
3. Pre-allocate slices
4. Reduce unnecessary allocations

### Phase 2: Performance Improvements (3-5 days)
1. Implement response caching
2. Optimize database queries (JOINs, batch fetching)
3. Parallel processing for independent operations
4. Optimize data transformations

### Phase 3: Advanced Optimizations (1 week)
1. Implement query result caching
2. Add response compression
3. Database connection pooling optimization
4. Consider read replicas for catalogue queries

## Expected Performance Improvements

| Optimization | Expected Improvement | Target Time/Op |
|-------------|---------------------|----------------|
| Database Indexes | 20-30% | ~1.8ms |
| Memory Optimization | 15-20% | ~2.0ms |
| Response Caching | 80-90% (cache hits) | ~0.1ms |
| Query Optimization | 30-40% | ~1.5ms |
| Parallel Processing | 20-30% | ~1.8ms |
| **Combined** | **60-70%** | **~0.8ms** |

## Monitoring & Metrics

### Key Metrics to Track:
- Request latency (p50, p95, p99)
- Memory usage per request
- Cache hit rate
- Database query time
- Throughput (requests per second)

### Tools:
- Go pprof for profiling
- Database query logs
- Application metrics (Prometheus)
- APM tools (if available)

## Testing Strategy

1. Run benchmarks before and after each optimization
2. Load testing with realistic data volumes
3. Memory profiling to identify leaks
4. Database query analysis
5. Cache effectiveness testing

## Next Steps

1. Review and approve optimization plan
2. Start with Phase 1 quick wins
3. Measure improvements after each phase
4. Iterate based on results


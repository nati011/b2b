# Cache

Caching infrastructure using Redis or Redis-like tools.

## Overview

The Cache module provides caching capabilities for the system using Redis or Redis-compatible storage. It offers high-performance caching for frequently accessed data, reducing database load and improving system responsiveness. The cache supports various caching patterns and strategies suitable for banking operations.

## Domain Boundaries

**Owns**:
- Cache infrastructure and implementation
- Cache key management and naming conventions
- Cache expiration and TTL management
- Cache invalidation strategies
- Cache serialization and deserialization
- Cache connection pooling and management
- Cache monitoring and metrics

**References** (via IDs/identifiers, not domain objects):
- Data from core domains (cached data, not owned by cache)
- Cache keys (references to cached data)

**Does NOT Own**:
- Cached data (owned by respective domains - cache stores copies)
- Business logic (cache is infrastructure, not business logic)
- Data persistence (cache is ephemeral, not primary storage)

**Interaction Patterns**:
- Used by: All core domains (for caching frequently accessed data)
- Provides: High-performance caching infrastructure
- Supports: Various caching patterns and strategies

## Responsibilities

- Provide high-performance caching using Redis/Redis-like tools
- Manage cache keys and naming conventions
- Handle cache expiration and TTL
- Support cache invalidation strategies
- Serialize and deserialize cached data
- Manage cache connections and pooling
- Monitor cache performance and health
- Handle cache failures gracefully

## Technology

The cache module uses **Redis or Redis-compatible** storage as the underlying technology:

- **Redis**: Primary choice for in-memory data structure store
- **Redis-Compatible**: Support for Redis-compatible alternatives (e.g., KeyDB, Dragonfly)
- **Features**: Leverages Redis features like TTL, pub/sub, transactions, Lua scripting

## Caching Patterns

### Cache-Aside (Lazy Loading)
- Application checks cache first
- On cache miss, loads from database
- Stores result in cache for future requests
- Most common pattern for read-heavy operations

### Write-Through
- Writes to both cache and database simultaneously
- Ensures cache and database are always in sync
- Suitable for write-heavy operations with consistency requirements

### Write-Back (Write-Behind)
- Writes to cache immediately
- Writes to database asynchronously
- Higher performance but requires careful handling of failures

### Refresh-Ahead
- Proactively refreshes cache before expiration
- Reduces cache miss latency
- Suitable for predictable access patterns

## Use Cases

### Domain Data Caching

#### User Domain
- Cache user profiles and authentication data
- Cache user role assignments
- Cache user status information
- Reduce database queries for user lookups

#### Client Domain
- Cache client/member records
- Cache KYC data
- Cache client hierarchy relationships
- Improve client lookup performance

#### Portfolio Domain
- Cache account balances (with appropriate TTL)
- Cache product configurations
- Cache account metadata
- Reduce database load for account queries

#### Organization Domain
- Cache branch information
- Cache staff assignments
- Cache currency and code lookups
- Cache calendar/holiday data

#### Accounting Domain
- Cache GL account information
- Cache account-to-GL mappings
- Cache exchange rates
- Improve accounting query performance

### Session and State Caching
- User session data
- Authentication tokens
- Temporary transaction state
- Request context data

### Reference Data Caching
- System codes and lookups
- Currency information
- Branch configurations
- Product definitions

### Query Result Caching
- Complex query results
- Report data
- Aggregated statistics
- Frequently accessed calculations

## Cache Key Management

### Key Naming Conventions
- **Format**: `{domain}:{entity}:{identifier}:{version}`
- **Examples**:
  - `user:profile:user-123:v1`
  - `portfolio:account:acc-456:v1`
  - `organization:branch:br-789:v1`
  - `accounting:gl-mapping:acc-456:v1`

### Key Versioning
- Version keys to support cache invalidation
- Increment version when data structure changes
- Support multiple versions during migration

### Key Expiration
- Set appropriate TTL based on data volatility
- Short TTL for frequently changing data
- Longer TTL for relatively static data
- No expiration for truly static reference data

## Cache Expiration Strategies

### Time-Based Expiration (TTL)
- **Short TTL** (seconds/minutes): Frequently changing data
- **Medium TTL** (hours): Moderately changing data
- **Long TTL** (days): Rarely changing data
- **No TTL**: Static reference data

### Event-Based Invalidation
- Invalidate cache on data updates
- Listen to domain events for cache invalidation
- Support pattern-based invalidation (e.g., all user caches)

### Manual Invalidation
- Explicit cache invalidation via API
- Support for administrative cache clearing
- Selective invalidation by pattern

## Cache Serialization

- **JSON**: Human-readable, widely supported
- **MessagePack**: More compact, faster serialization
- **Protocol Buffers**: Type-safe, efficient binary format
- **Custom**: Domain-specific serialization if needed

## Cache Operations

### Basic Operations
- **Get**: Retrieve cached value
- **Set**: Store value in cache
- **Delete**: Remove cached value
- **Exists**: Check if key exists
- **Expire**: Set TTL on key

### Advanced Operations
- **Increment/Decrement**: Atomic counter operations
- **Batch Operations**: Get/Set multiple keys
- **Pattern Matching**: Find keys by pattern
- **Transactions**: Atomic multi-key operations

## Cache Invalidation

### Invalidation Strategies

#### Immediate Invalidation
- Invalidate cache immediately on data update
- Ensures consistency but may cause cache thrashing

#### Lazy Invalidation
- Mark cache as stale on update
- Invalidate on next access
- Balances consistency and performance

#### Time-Based Invalidation
- Rely on TTL for expiration
- Simple but may serve stale data temporarily

#### Event-Driven Invalidation
- Listen to domain events
- Invalidate related caches on events
- Example: `AccountUpdated` event invalidates account cache

### Invalidation Patterns
- **Single Key**: Invalidate specific cache key
- **Pattern-Based**: Invalidate all keys matching pattern
- **Tag-Based**: Invalidate all keys with specific tag
- **Cascade**: Invalidate related caches

## Cache Failure Handling

### Graceful Degradation
- Continue operation if cache is unavailable
- Fall back to database queries
- Log cache failures for monitoring

### Circuit Breaker
- Detect cache failures
- Open circuit after threshold failures
- Bypass cache during circuit open
- Attempt recovery after timeout

### Retry Logic
- Retry cache operations on transient failures
- Exponential backoff for retries
- Limit retry attempts

## Security

### Data Protection
- **Encryption at Rest**: Encrypt sensitive cached data
- **Encryption in Transit**: Use TLS for Redis connections
- **Access Control**: Restrict Redis access to authorized services
- **Data Masking**: Mask sensitive data in cache (if required)

### Key Security
- **Key Isolation**: Separate cache keys by tenant/organization
- **Key Validation**: Validate cache keys before operations
- **Key Expiration**: Ensure sensitive data expires promptly

## Performance Considerations

### Connection Pooling
- Maintain connection pool to Redis
- Reuse connections for efficiency
- Configure pool size based on load

### Pipeline Operations
- Batch multiple operations
- Reduce network round trips
- Improve throughput

### Compression
- Compress large cached values
- Reduce memory usage
- Trade-off: CPU for memory

## Monitoring and Observability

### Metrics
- **Hit Rate**: Percentage of cache hits
- **Miss Rate**: Percentage of cache misses
- **Latency**: Cache operation latency
- **Error Rate**: Cache operation errors
- **Memory Usage**: Redis memory consumption
- **Connection Pool**: Connection pool statistics

### Health Checks
- Redis connection health
- Cache availability
- Performance degradation detection

### Alerts
- Cache unavailability
- High miss rates
- Performance degradation
- Memory threshold breaches

## Best Practices

1. **Key Design**: Use clear, consistent key naming conventions
2. **TTL Management**: Set appropriate TTLs based on data volatility
3. **Invalidation**: Implement proper cache invalidation strategies
4. **Failure Handling**: Always handle cache failures gracefully
5. **Monitoring**: Monitor cache performance and health
6. **Security**: Secure cache connections and sensitive data
7. **Testing**: Test cache behavior under various scenarios
8. **Documentation**: Document cache keys and TTLs

## Relationship with Other Modules

### Core Domains
- **Used By**: All core domains use cache for performance optimization
- **Purpose**: Reduce database load and improve response times
- **Flow**: Domain → Cache lookup → Database (on miss) → Cache store

### Database Module
- **Fallback**: Database is fallback when cache misses
- **Source of Truth**: Database is always source of truth
- **Flow**: Cache miss → Database query → Cache store

### Event Bus
- **Integration**: Cache listens to domain events for invalidation
- **Purpose**: Keep cache in sync with data changes
- **Flow**: Domain event → Cache invalidation → Cache cleared

## Example Scenarios

### Scenario 1: User Profile Caching
1. Request for user profile
2. Check cache for `user:profile:user-123:v1`
3. Cache hit → Return cached data
4. Cache miss → Query database → Store in cache → Return data

### Scenario 2: Cache Invalidation
1. User profile updated
2. Domain publishes `UserUpdated` event
3. Cache module listens to event
4. Invalidates `user:profile:user-123:v1`
5. Next request will fetch fresh data from database

### Scenario 3: Account Balance Caching
1. Request for account balance
2. Check cache for `portfolio:balance:acc-456:v1`
3. Cache hit with short TTL (30 seconds) → Return cached balance
4. Balance updated → Invalidate cache immediately
5. Next request fetches fresh balance

### Scenario 4: Reference Data Caching
1. Request for currency information
2. Check cache for `organization:currency:USD:v1`
3. Cache hit (long TTL, 24 hours) → Return cached data
4. Currency updated → Invalidate cache
5. Cache repopulated on next access

## Configuration

### Redis Configuration
- **Host**: Redis server host
- **Port**: Redis server port
- **Password**: Redis authentication password
- **Database**: Redis database number
- **TLS**: Enable TLS for secure connections
- **Pool Size**: Connection pool size
- **Timeout**: Connection and operation timeouts

### Cache Configuration
- **Default TTL**: Default expiration time
- **Max Memory**: Maximum cache memory usage
- **Eviction Policy**: Cache eviction policy (LRU, LFU, etc.)
- **Compression**: Enable compression for large values


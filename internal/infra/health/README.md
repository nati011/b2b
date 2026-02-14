# Health Checks

Health check infrastructure for system components.

## Overview

The Health module provides infrastructure for health checks across system components. It enables monitoring systems to verify system health, detect failures, and trigger alerts. Health checks verify the availability and status of critical system components including databases, caches, external services, and application services.

## Domain Boundaries

**Owns**:
- Health check infrastructure and implementation
- Health check definitions and configurations
- Health status aggregation
- Health check execution and scheduling
- Health status reporting
- Dependency health tracking

**References** (via IDs/identifiers, not domain objects):
- System components (for health verification)
- External services (for health verification)

**Does NOT Own**:
- System components (health checks verify, don't own components)
- Business logic (health checks are infrastructure, not business logic)
- Component status (health checks report status, don't manage it)

**Interaction Patterns**:
- Used by: Monitoring systems (for health verification)
- Provides: Health check infrastructure for system monitoring
- Verifies: Health of databases, caches, services, external dependencies

## Responsibilities

- Execute health checks on system components
- Aggregate health status across components
- Report overall system health
- Provide health check API endpoint
- Support liveness and readiness probes
- Track dependency health
- Support health check scheduling

## Health Check Types

### Liveness Probe
- Indicates if application is running
- Used by orchestrators (Kubernetes) to restart containers
- Should return healthy if application is alive
- Example: Application process is running

### Readiness Probe
- Indicates if application is ready to serve traffic
- Used by load balancers to route traffic
- Should return healthy if application can handle requests
- Example: Database connection available, cache available

### Startup Probe
- Indicates if application has finished starting
- Used during application startup
- Different from readiness (startup is one-time)
- Example: Application initialization complete

## Health Check Components

### Database Health
- Verify database connection
- Execute simple query (e.g., SELECT 1)
- Check connection pool status
- Verify transaction capability

### Cache Health
- Verify Redis/cache connection
- Execute simple operation (e.g., PING)
- Check cache availability
- Verify cache performance

### External Service Health
- Verify external API availability
- Check external service connectivity
- Monitor external service response time
- Track external service failures

### Application Health
- Verify application is running
- Check application resources (memory, CPU)
- Verify application can process requests
- Check application configuration

## Health Status

### Status Levels
- **Healthy**: Component is healthy and operational
- **Degraded**: Component is functional but with issues
- **Unhealthy**: Component is not operational
- **Unknown**: Health status cannot be determined

### Overall Health
- **Healthy**: All critical components healthy
- **Degraded**: Some non-critical components unhealthy
- **Unhealthy**: Critical components unhealthy

## Use Cases

### System Monitoring
- Monitor overall system health
- Detect component failures
- Trigger alerts on health degradation
- Support automated recovery

### Load Balancer Integration
- Provide health status to load balancers
- Enable/disable traffic routing based on health
- Support graceful shutdown
- Enable zero-downtime deployments

### Container Orchestration
- Provide liveness/readiness probes for Kubernetes
- Support container restart on failure
- Enable rolling deployments
- Support health-based scaling

### Dependency Monitoring
- Monitor external service health
- Track dependency availability
- Support circuit breaker patterns
- Enable graceful degradation

## Health Check Endpoints

### `/health`
- Overall system health
- Aggregated health status
- Quick health check

### `/health/live`
- Liveness probe
- Application is running
- Used by orchestrators

### `/health/ready`
- Readiness probe
- Application is ready
- Used by load balancers

### `/health/detailed`
- Detailed health information
- Component-level health
- Dependency health status

## Health Check Configuration

### Check Intervals
- **Liveness**: Check every 10-30 seconds
- **Readiness**: Check every 5-10 seconds
- **Startup**: Check every 1-5 seconds during startup

### Timeouts
- **Connection Timeout**: Timeout for connection checks
- **Query Timeout**: Timeout for query execution
- **Overall Timeout**: Maximum time for health check

### Failure Thresholds
- **Consecutive Failures**: Number of failures before unhealthy
- **Failure Rate**: Percentage of failures before unhealthy
- **Recovery Threshold**: Successes needed to recover

## Best Practices

1. **Fast Checks**: Health checks should be fast (< 1 second)
2. **Lightweight**: Minimal resource usage for health checks
3. **Accurate**: Health checks should accurately reflect component status
4. **Non-Blocking**: Health checks should not block application
5. **Comprehensive**: Check all critical dependencies
6. **Documentation**: Document health check endpoints and meanings
7. **Monitoring**: Monitor health check results

## Relationship with Other Modules

### Database Module
- **Verifies**: Database health via connection check
- **Purpose**: Ensure database is available
- **Flow**: Health Check → Database → Connection Test

### Cache Module
- **Verifies**: Cache health via Redis connection
- **Purpose**: Ensure cache is available
- **Flow**: Health Check → Cache → PING Test

### HTTP Server
- **Exposes**: Health check endpoints
- **Purpose**: Allow external systems to check health
- **Flow**: HTTP Request → Health Check → Status Response

### Monitoring Module
- **Integration**: Health status feeds into monitoring
- **Purpose**: Enable health-based alerting
- **Flow**: Health Check → Status → Monitoring → Alerts

## Example Health Check Responses

### Healthy
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "components": {
    "database": "healthy",
    "cache": "healthy",
    "application": "healthy"
  }
}
```

### Unhealthy
```json
{
  "status": "unhealthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "components": {
    "database": "healthy",
    "cache": "unhealthy",
    "application": "healthy"
  },
  "errors": {
    "cache": "connection timeout"
  }
}
```


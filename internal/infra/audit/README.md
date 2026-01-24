# Audit Module

Comprehensive audit trail and event logging for all system activities.

## Overview

The Audit module provides a centralized audit trail system that captures and stores all significant events and activities across the entire system. It ensures complete traceability, compliance, and accountability by maintaining an immutable record of all operations, changes, and events.

## Domain Boundaries

**Owns**:
- Audit records and audit trails
- Event logging and storage
- Audit query and retrieval
- Audit retention policies
- Audit report generation

**References** (via IDs/identifiers, not domain objects):
- User IDs (who performed the action)
- Resource IDs (what was affected)
- Event IDs (from event bus)
- Transaction IDs (for financial operations)
- Account IDs (for account-related activities)

**Does NOT Own**:
- Business logic (audit is infrastructure, not business logic)
- Domain events (audit consumes events, doesn't create them)
- User entities (references user IDs only)
- Resources being audited (references resource IDs only)

**Interaction Patterns**:
- **Consumes**: All domain events from event bus
- **Consumes**: HTTP request/response events from middleware
- **Consumes**: Database transaction events
- **Consumes**: Authentication/authorization events
- **Used By**: All domains (for audit trail queries)
- **Used By**: Compliance and reporting modules

## Features

### Event Auditing
- **Domain Event Auditing**: Captures all domain events (UserCreated, AccountCreated, TransactionSettled, etc.)
- **HTTP Request Auditing**: Logs all HTTP requests and responses
- **Database Change Auditing**: Tracks all database modifications
- **Authentication Auditing**: Records all authentication attempts and authorization decisions
- **Business Operation Auditing**: Captures all business-critical operations

### Audit Record Structure
- **Timestamp**: Precise timestamp of the event
- **Actor**: User or system that performed the action
- **Action**: Type of action performed (create, update, delete, view, etc.)
- **Resource**: Resource type and ID affected
- **Event Type**: Category of event (domain event, HTTP request, database change, etc.)
- **Event Data**: Complete event payload and context
- **IP Address**: Source IP address (for HTTP requests)
- **User Agent**: Client user agent (for HTTP requests)
- **Request ID**: Correlation ID for request tracing
- **Outcome**: Success or failure status
- **Error Details**: Error information if operation failed

### Audit Query and Retrieval
- Query audit records by date range
- Filter by user, resource, action type, or event type
- Search audit records by content
- Generate audit reports
- Export audit trails for compliance

### Audit Retention
- Configurable retention policies
- Automatic archival of old records
- Compliance with regulatory requirements
- Secure storage of audit data

## Architecture

### Event Sources

The audit module receives events from multiple sources:

1. **Event Bus**: All domain events published through the internal event bus
2. **HTTP Middleware**: All HTTP requests and responses
3. **Database Triggers**: Database-level change tracking (optional)
4. **Service Layer**: Direct audit calls from service methods
5. **Authentication Middleware**: Login attempts, authorization decisions

### Event Flow

```
Event Source (Domain/HTTP/Database)
    ↓
Audit Service
    ↓
Event Validation & Enrichment
    ↓
Audit Repository
    ↓
Persistent Storage (Database)
```

### Storage Strategy

- **Primary Storage**: PostgreSQL for structured audit records
- **Archival Storage**: Long-term storage for compliance (configurable)
- **Indexing**: Optimized indexes for common query patterns
- **Partitioning**: Time-based partitioning for performance

## Audit Record Types

### Domain Events
- User lifecycle events (created, updated, activated, deactivated)
- Account events (created, updated, closed)
- Transaction events (initiated, settled, reversed)
- Product events (created, updated, activated)
- Role and permission changes
- Configuration changes

### HTTP Requests
- All API requests and responses
- Request parameters and payloads
- Response status codes and bodies
- Authentication and authorization decisions
- Error responses and exceptions

### Database Changes
- All INSERT, UPDATE, DELETE operations
- Before and after values (for updates)
- Transaction context
- Database user and connection info

### Authentication Events
- Login attempts (successful and failed)
- Logout events
- Password changes
- Token generation and validation
- Authorization decisions

### Business Operations
- Loan approvals and rejections
- Account openings and closures
- Transaction processing
- Report generation
- Bulk operations

## Compliance and Security

### Regulatory Compliance
- **SOX Compliance**: Financial transaction auditing
- **GDPR Compliance**: Data access and modification tracking
- **PCI DSS Compliance**: Payment card data access tracking
- **Banking Regulations**: Financial services audit requirements

### Security Features
- **Immutable Records**: Audit records cannot be modified or deleted
- **Tamper Detection**: Cryptographic hashing for integrity verification
- **Access Control**: Restricted access to audit data
- **Encryption**: Encrypted storage of sensitive audit data
- **Retention Policies**: Automatic enforcement of retention rules

## Query and Reporting

### Query Capabilities
- Filter by date range
- Filter by user/actor
- Filter by resource type and ID
- Filter by action type
- Filter by event type
- Filter by outcome (success/failure)
- Full-text search on event data

### Report Generation
- User activity reports
- Resource change history
- Failed operation reports
- Compliance reports
- Custom audit reports

## Performance Considerations

### Optimization Strategies
- **Asynchronous Processing**: Non-blocking audit record creation
- **Batch Insertion**: Batch multiple audit records for efficiency
- **Partitioning**: Time-based table partitioning
- **Indexing**: Strategic indexes on common query fields
- **Archival**: Move old records to archival storage

### Scalability
- **Horizontal Scaling**: Support for distributed audit storage
- **Sharding**: Shard audit data by date or resource type
- **Caching**: Cache frequently accessed audit records
- **Compression**: Compress archived audit data

## Integration Points

### Event Bus Integration
- Subscribes to all domain events
- Automatically captures event payloads
- Enriches events with audit metadata
- Stores events in audit repository

### HTTP Middleware Integration
- Intercepts all HTTP requests
- Captures request/response data
- Extracts user context
- Records authentication/authorization decisions

### Database Integration
- Optional database trigger integration
- Change data capture (CDC) support
- Transaction-level auditing

## Usage Examples

### Querying Audit Records

```go
// Query all user-related events in date range
auditRecords, err := auditService.Query(ctx, audit.QueryParams{
    ResourceType: "users",
    StartDate: time.Now().AddDate(0, -1, 0),
    EndDate: time.Now(),
})

// Query specific user's activity
userAudit, err := auditService.QueryByActor(ctx, userID, audit.QueryParams{
    StartDate: time.Now().AddDate(0, -1, 0),
})

// Query resource change history
resourceHistory, err := auditService.QueryByResource(ctx, "loan-products", productID)
```

### Generating Audit Reports

```go
// Generate compliance report
report, err := auditService.GenerateComplianceReport(ctx, audit.ComplianceReportParams{
    ReportType: "sox",
    StartDate: startDate,
    EndDate: endDate,
})

// Generate user activity report
activityReport, err := auditService.GenerateUserActivityReport(ctx, userID, dateRange)
```

## Future Enhancements

- Real-time audit monitoring and alerting
- Machine learning for anomaly detection
- Advanced analytics and visualization
- Integration with external audit systems
- Automated compliance checking
- Audit record correlation and analysis

## Relationship with Other Modules

### Event Bus
- **Consumes**: All events from event bus
- **Purpose**: Capture all domain events for audit trail

### All Core Domains
- **Consumes**: Events from all domains
- **Purpose**: Maintain complete audit trail of all operations

### HTTP Middleware
- **Integrates**: With HTTP middleware for request/response auditing
- **Purpose**: Capture all API interactions

### Authentication/Authorization
- **Integrates**: With auth middleware
- **Purpose**: Track all authentication and authorization events

### Report Module
- **Used By**: Report module for audit report generation
- **Purpose**: Provide audit data for reporting

## Best Practices

1. **Comprehensive Coverage**: Audit all significant operations
2. **Performance**: Use asynchronous processing to avoid blocking business operations
3. **Retention**: Implement appropriate retention policies
4. **Security**: Protect audit data with proper access controls
5. **Compliance**: Ensure audit records meet regulatory requirements
6. **Queryability**: Design audit schema for efficient querying
7. **Immutability**: Never allow modification or deletion of audit records
8. **Correlation**: Include correlation IDs for request tracing


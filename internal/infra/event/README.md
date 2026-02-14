# Event Bus

Internal and external event bus infrastructure.

## Overview

The Event Bus module provides infrastructure for event-driven communication within the system and with external systems. It supports both internal domain events (for inter-domain communication) and external integration events (for webhooks, API integrations, and external system communication).

## Domain Boundaries

**Owns**:
- Event bus infrastructure and implementation
- Event publishing and subscription mechanisms
- Event routing and delivery
- Event persistence and replay
- External event integration (webhooks, API events)
- Event transformation and mapping
- Event delivery guarantees and retry logic

**References** (via IDs/identifiers, not domain objects):
- Domain event types (from core domains - references only)
- External system identifiers (for external integrations)

**Does NOT Own**:
- Domain events (owned by core domains - event bus transports them)
- Business logic (event bus is infrastructure, not business logic)
- Event handlers (owned by consuming domains)

**Interaction Patterns**:
- Used by: All core domains (for publishing and consuming events)
- Provides: Event bus infrastructure for asynchronous communication
- Supports: Internal domain events and external integration events

## Responsibilities

### Internal Event Bus
- Publish domain events from core domains
- Subscribe to domain events for cross-domain communication
- Route events to appropriate subscribers
- Ensure event delivery guarantees
- Support event replay and recovery
- Maintain event ordering where required

### External Event Bus
- Publish events to external systems (webhooks, APIs)
- Receive events from external systems
- Transform events between internal and external formats
- Handle external system authentication and security
- Support retry logic for external event delivery
- Maintain external event delivery status

## Internal Event Bus

The internal event bus facilitates asynchronous communication between core domains using domain events.

### Features
- **Event Publishing**: Domains publish events through the event bus
- **Event Subscription**: Domains subscribe to events they need to consume
- **Event Routing**: Events are routed to all registered subscribers
- **Delivery Guarantees**: At-least-once or exactly-once delivery (configurable)
- **Event Persistence**: Events can be persisted for replay and audit
- **Event Ordering**: Maintains event order within event streams (if required)

### Use Cases

#### Domain Event Publishing
- Portfolio publishes `AccountCreated` event
- Settlement publishes `TransactionSettled` event
- Accounting publishes `JournalEntryPosted` event
- Organization publishes `BranchCreated` event

#### Domain Event Consumption
- Accounting consumes `AccountCreated` to create GL mappings
- Accounting consumes `TransactionSettled` to create journal entries
- Settlement consumes `AccountCreated` to initialize transaction tracking
- Organization consumes `UserCreated` to create staff records

### Event Flow
```
Domain A (Publisher)
    ↓
Publishes Event
    ↓
Event Bus (Internal)
    ↓
Routes to Subscribers
    ↓
Domain B (Subscriber)
Domain C (Subscriber)
```

## External Event Bus

The external event bus handles integration with external systems through webhooks, APIs, and event streams.

### Features
- **Webhook Publishing**: Publish events to external systems via webhooks
- **Webhook Receiving**: Receive webhooks from external systems
- **API Integration**: Integrate with external APIs for event exchange
- **Event Transformation**: Transform between internal and external event formats
- **Authentication**: Handle authentication for external systems
- **Retry Logic**: Retry failed external event deliveries
- **Delivery Status**: Track delivery status of external events

### Use Cases

#### Outbound Events (System → External)
- Publish transaction events to payment gateways
- Send account events to external reporting systems
- Notify external systems of settlement completions
- Integrate with mobile money providers

#### Inbound Events (External → System)
- Receive webhooks from payment gateways
- Process events from mobile money providers
- Handle events from external banking systems
- Integrate with third-party services

### Event Flow
```
Internal Domain Event
    ↓
Event Bus (External)
    ↓
Transform to External Format
    ↓
Publish to External System (Webhook/API)
    ↓
External System
```

## Event Types

### Internal Domain Events
- **AccountCreated**: Account created in Portfolio
- **TransactionSettled**: Transaction settled in Settlement
- **JournalEntryPosted**: Journal entry posted in Accounting
- **UserCreated**: User created in User domain
- **BranchCreated**: Branch created in Organization

### External Integration Events
- **PaymentGatewayWebhook**: Webhook from payment gateway
- **MobileMoneyEvent**: Event from mobile money provider
- **BankTransferEvent**: Event from external bank
- **ThirdPartyNotification**: Notification from third-party service

## Event Delivery Guarantees

### Internal Events
- **At-Least-Once**: Events are delivered at least once (may have duplicates)
- **Exactly-Once**: Events are delivered exactly once (idempotent processing)
- **Ordered Delivery**: Events are delivered in order (within event stream)

### External Events
- **Best Effort**: Attempts delivery, may fail (retry logic applies)
- **Guaranteed Delivery**: Ensures delivery with retries
- **Idempotent Delivery**: External systems handle duplicate events

## Event Persistence

Events can be persisted for:
- **Audit Trail**: Complete event history for auditing
- **Event Replay**: Replay events for recovery or reprocessing
- **Debugging**: Troubleshoot event delivery issues
- **Analytics**: Analyze event patterns and flows

## Event Transformation

The event bus supports event transformation:
- **Format Conversion**: Convert between internal and external formats
- **Schema Mapping**: Map event schemas between systems
- **Data Enrichment**: Add additional data to events
- **Data Filtering**: Filter sensitive data from external events

## Security

### Internal Events
- **Authentication**: Verify publisher identity
- **Authorization**: Check publisher permissions
- **Encryption**: Encrypt events in transit (if required)

### External Events
- **Webhook Authentication**: Verify webhook signatures
- **API Authentication**: Handle API keys, OAuth, etc.
- **Data Encryption**: Encrypt sensitive data in external events
- **Rate Limiting**: Limit external event rates

## Error Handling

### Internal Events
- **Retry Logic**: Retry failed event deliveries
- **Dead Letter Queue**: Handle permanently failed events
- **Error Notifications**: Notify administrators of failures

### External Events
- **Retry Logic**: Retry failed external deliveries
- **Exponential Backoff**: Use exponential backoff for retries
- **Failure Notifications**: Notify of external delivery failures
- **Fallback Mechanisms**: Alternative delivery methods

## Monitoring and Observability

- **Event Metrics**: Track event publishing and consumption rates
- **Delivery Metrics**: Monitor event delivery success/failure rates
- **Latency Metrics**: Track event processing latency
- **Error Tracking**: Monitor and alert on event errors
- **Event Tracing**: Trace events through the system

## Best Practices

1. **Event Naming**: Use clear, descriptive event names
2. **Event Versioning**: Version events for backward compatibility
3. **Idempotency**: Design event handlers to be idempotent
4. **Error Handling**: Handle errors gracefully in event handlers
5. **Monitoring**: Monitor event bus health and performance
6. **Security**: Secure both internal and external event channels
7. **Documentation**: Document event schemas and contracts

## Relationship with Other Modules

### Core Domains
- **Used By**: All core domains use event bus for publishing/consuming events
- **Purpose**: Enables asynchronous cross-domain communication
- **Flow**: Domain → Event Bus → Domain

### Settlement Domain
- **Publishes**: Transaction events for external systems
- **Consumes**: External payment events
- **Purpose**: Integrate with payment gateways and external systems

### Integration Sub-Modules
- **Settlement/integration**: Uses external event bus for payment gateway integration
- **Purpose**: Handle external system events and webhooks

## Example Scenarios

### Scenario 1: Internal Domain Event
1. Portfolio creates account
2. Portfolio publishes `AccountCreated` event via internal event bus
3. Event bus routes to subscribers
4. Accounting consumes event and creates GL mapping
5. Settlement consumes event and initializes transaction tracking

### Scenario 2: External Webhook (Outbound)
1. Settlement completes transaction
2. Settlement publishes `TransactionSettled` event
3. External event bus transforms event to webhook format
4. External event bus publishes webhook to payment gateway
5. Payment gateway receives webhook notification

### Scenario 3: External Webhook (Inbound)
1. Payment gateway sends webhook to system
2. External event bus receives and validates webhook
3. External event bus transforms to internal event format
4. Internal event bus publishes internal event
5. Settlement consumes event and processes payment


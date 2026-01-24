# Order

## Overview
The order domain handles purchasing workflows, capturing order items, pricing
snapshots, and status transitions from creation through fulfillment.

## Domain Boundaries
**Owns**:
- Order entity and line items
- Pricing, tax, and totals (snapshot at order time)
- Status transitions (pending, confirmed, shipped, delivered, cancelled)
- Fulfillment milestones and audit trail

**References** (via IDs/identifiers, not domain objects):
- Customer IDs
- Supplier IDs
- Product or listing IDs
- Payment or transaction IDs (external system)

**Does NOT Own**:
- Customer profiles and addresses of record
- Product catalog, listings, or inventory
- Supplier compliance data

**Interaction Patterns**:
- Requests availability and pricing checks from Product
- Sends fulfillment tasks to Supplier
- Stores customer and address snapshots on creation

## Responsibilities
- Create and validate orders
- Capture immutable price and item snapshots
- Manage order state changes and cancellations
- Provide order history to Customer views

## Use Cases
- Create a new order from a cart
- Cancel or modify an order before fulfillment
- Update fulfillment status and tracking info
- Generate order summaries and invoices

## Relationship with Other Modules
- Customer: consumes customer data and provides history
- Product: validates availability and pricing snapshots
- Supplier: routes fulfillment and SLA tracking

## Example
Customer places an order:
1) Order module receives `customer_id` and selected listing IDs.
2) Product module validates availability and returns price snapshot.
3) Order is created with immutable totals and status `pending`.

## Best Practices
- Keep order snapshots immutable for auditability
- Ensure idempotency for order creation
- Use IDs for cross-module references

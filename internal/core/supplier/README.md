# Supplier

## Overview
The supplier domain models sellers and their business profiles, onboarding
status, and fulfillment capabilities. It is the source of truth for supplier
identity and compliance in the marketplace.

## Domain Boundaries
**Owns**:
- Supplier profile (legal name, status)
- Business registration and compliance state
- Contacts and support channels
- Fulfillment capabilities (service areas, lead times, warehouses)

**References** (via IDs/identifiers, not domain objects):
- Product IDs and listing IDs
- Order IDs for performance metrics
- Payout account IDs (external system)

**Does NOT Own**:
- Product catalog definitions or pricing rules
- Inventory reservations
- Order lifecycle or customer data

**Interaction Patterns**:
- Product references `supplier_id` for listings
- Order routes fulfillment tasks to suppliers
- Customer views supplier profiles and ratings

## Responsibilities
- Onboard and approve suppliers
- Maintain compliance and verification status
- Store supplier contacts and operational details
- Provide supplier data to Product and Order flows

## Use Cases
- Register a new supplier and complete onboarding
- Update compliance documents or business details
- Suspend or reactivate a supplier account
- Expose supplier info for catalog and order views

## Relationship with Other Modules
- Product: supplier listings and ownership metadata
- Order: fulfillment, SLA tracking, and performance analytics
- Customer: supplier discovery and trust signals

## Example
Supplier publishes a listing:
1) Supplier is approved and active.
2) Product module creates a listing referencing `supplier_id`.
3) Order module can now route orders to the supplier.

## Best Practices
- Keep compliance and onboarding data here
- Use IDs across boundaries; avoid sharing domain objects
- Store immutable supplier snapshots for historical orders

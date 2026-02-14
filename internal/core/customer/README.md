# Customer

## Overview
The customer domain models buyer organizations and the people who act on their
behalf. It is the source of truth for customer identity, contact details, and
eligibility within the marketplace.

## Domain Boundaries
**Owns**:
- Customer profile (legal name, status, tax info)
- Contacts and roles (buyer, approver)
- Addresses (shipping, billing)
- Customer preferences and segments (tier, payment terms)

**References** (via IDs/identifiers, not domain objects):
- Auth/user IDs
- Order IDs
- Supplier IDs (preferred suppliers)
- Product IDs (saved items, frequently ordered)

**Does NOT Own**:
- Product catalog, pricing, or inventory
- Supplier onboarding or compliance
- Order lifecycle, payment processing

**Interaction Patterns**:
- Reads product and supplier data for discovery
- Sends customer profile and address snapshot to Order at checkout
- Receives order history summaries from Order for account views

## Responsibilities
- Register and maintain customer organizations
- Manage contacts, roles, and permissions for buyer teams
- Store addresses and default preferences
- Provide eligibility signals for ordering (status, tier)

## Use Cases
- Create or update a customer profile
- Add or change shipping and billing addresses
- Validate customer status before placing orders
- Deactivate or suspend a customer account

## Relationship with Other Modules
- Product: customer-specific visibility or tier-based pricing
- Supplier: supplier discovery and preferred supplier lists
- Order: customer data snapshot and order history
- Auth: map user accounts to customer contacts

## Example
Customer places an order:
1) Customer module returns the active customer profile and default address.
2) Order module creates an order referencing `customer_id`.
3) Product module validates listing availability and pricing.

## Best Practices
- Keep PII and customer records inside this module
- Use IDs and snapshots across module boundaries
- Validate customer status on each order attempt

# Product

## Overview
The product domain manages the marketplace catalog and supplier listings. It
defines product attributes, categorization, and the offer details used for
discovery and ordering.

## Domain Boundaries
**Owns**:
- Product definitions (name, description, attributes, units)
- Catalog taxonomy (categories)
- Supplier listings (price, MOQ, availability status)

**References** (via IDs/identifiers, not domain objects):
- Supplier IDs (listing owner)
- Inventory or stock IDs (if managed elsewhere)
- Order IDs for demand analytics

**Does NOT Own**:
- Supplier onboarding and compliance
- Order lifecycle, payment, or customer profiles

**Interaction Patterns**:
- Exposes searchable catalog to Customer
- Validates listing availability for Order
- Accepts supplier updates to listings

## Responsibilities
- Create and maintain product definitions and listings
- Manage pricing, availability, and listing visibility
- Provide search and filtering capabilities
- Supply product snapshots for order creation

## Use Cases
- Add a new product and publish it to the catalog
- Create a supplier listing with price and MOQ
- Update pricing or availability for a listing
- Retire a product or listing from sale

## Relationship with Other Modules
- Supplier: owns seller identity for listings
- Order: uses product snapshots and validates availability
- Customer: consumes catalog for discovery and ordering

## Example
Supplier adds a new catalog listing:
1) Product module creates a product definition.
2) A listing is created referencing `supplier_id`.
3) Customer module can surface the listing in search.

## Best Practices
- Keep product data normalized and versioned
- Store order-time snapshots for price and attributes
- Avoid embedding supplier or customer objects

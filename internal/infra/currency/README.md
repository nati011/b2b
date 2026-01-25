# Currency

Currency management and currency configuration.

## Overview

The Currency sub-module manages currencies used in the system, including currency definitions, exchange rates, and currency-related configurations. It supports multi-currency operations and currency conversions.

## Architecture

This module uses a simplified architecture without separate layers:

- **currency.go**: Core currency model, domain logic, enums, errors, and validation
- **handler.go**: HTTP handlers with embedded business logic, database access, and DTOs
- **routes.go**: HTTP route definitions and registration

The module consolidates handler, service, and repository logic into a single handler file for simplicity, making it suitable for a straightforward service.

## Domain Boundaries

**Owns**:
- Currency entities and domain logic
- Currency definitions and codes
- Exchange rates and rate history
- Currency configurations
- Base currency settings
- Currency conversion rules

**References** (via IDs/identifiers, not domain objects):
- Organization IDs (for organization-level currency settings)

**Does NOT Own**:
- Account balances (owned by Portfolio - balances may be in different currencies)
- Transactions (owned by Settlement - transactions may involve currency conversion)
- GL accounts (owned by Accounting - GL accounts may be currency-specific)

**Interaction Patterns**:
- Publishes events: `CurrencyCreated`, `ExchangeRateUpdated`, `BaseCurrencyChanged`
- Consumed by: All domains (for currency validation and conversion)
- Provides read-only queries for currency lookup and exchange rates
- Used as reference data throughout the system

## Responsibilities

- Define and manage currencies
- Manage exchange rates and rate history
- Configure base currency
- Support currency conversion
- Provide currency validation
- Support multi-currency operations
- Maintain currency audit trail

## Currency Configuration

Currencies are configured with:

- **Currency Code**: ISO 4217 currency code (e.g., USD, EUR, KES)
- **Currency Name**: Display name of the currency
- **Symbol**: Currency symbol (e.g., $, €, KSh)
- **Decimal Places**: Number of decimal places for currency
- **Base Currency**: Organization's base currency
- **Active Status**: Whether currency is active

## Exchange Rates

Exchange rates can be managed:

- **Spot Rates**: Current exchange rates
- **Historical Rates**: Historical exchange rate data
- **Rate Sources**: Source of exchange rates (manual, external API)
- **Rate Updates**: Frequency of rate updates
- **Rate Validation**: Validation rules for exchange rates

## Use Cases

### Currency Management
- Define currencies used in the system
- Configure base currency for organization
- Activate/deactivate currencies
- Update currency information

### Exchange Rate Management
- Configure exchange rates
- Update exchange rates (manual or automated)
- Maintain exchange rate history
- Support multiple rate sources

### Currency Conversion
- Convert amounts between currencies
- Apply exchange rates for conversions
- Support transaction currency conversion
- Handle multi-currency account balances

### Multi-Currency Operations
- Support accounts in different currencies
- Process transactions in multiple currencies
- Generate reports in base currency
- Handle currency conversion in accounting

## Relationship with Other Domains

### Portfolio Domain
- **References**: Portfolio may reference Currency IDs for account currencies
- **Purpose**: Accounts may be denominated in different currencies
- **Flow**: Account → Currency → Account balance in currency

### Settlement Domain
- **Queries**: Settlement queries currency for conversion rates
- **Purpose**: Transactions may involve currency conversion
- **Flow**: Transaction → Currency conversion → Settlement

### Accounting Domain
- **References**: Accounting may reference Currency IDs for multi-currency GL accounts
- **Purpose**: GL accounts may be currency-specific
- **Flow**: GL Account → Currency → Currency-specific accounting

### All Domains
- **Queries**: All domains query currency for validation and conversion
- **Purpose**: Ensure currency operations are valid
- **Flow**: Operation → Currency query → Currency validation/conversion

## Example Currency Scenarios

### Scenario 1: Multi-Currency Account
1. Member opens account in USD
2. Account balance maintained in USD
3. Transactions processed in USD
4. Reports converted to base currency (KES)
5. Accounting entries in base currency

### Scenario 2: Currency Conversion
1. Member deposits USD 100
2. Exchange rate: 1 USD = 150 KES
3. Convert to base currency: 15,000 KES
4. Account balance updated
5. Accounting entry in base currency

### Scenario 3: Exchange Rate Update
1. Exchange rate changes: 1 USD = 155 KES
2. Update exchange rate in currency module
3. New conversions use updated rate
4. Historical transactions maintain old rate
5. Rate history maintained for audit

## Currency Types

The Currency module supports:

- **Base Currency**: Organization's primary currency
- **Foreign Currency**: Other currencies used in operations
- **Local Currency**: Local currency for branches
- **Reporting Currency**: Currency used for reporting

## Best Practices

1. **Base Currency**: Define clear base currency for organization
2. **Rate Management**: Regularly update exchange rates
3. **Rate History**: Maintain complete exchange rate history
4. **Validation**: Validate currency codes and rates
5. **Conversion Rules**: Define clear currency conversion rules
6. **Multi-Currency**: Support multi-currency operations consistently
7. **Audit Trail**: Maintain complete audit trail for currency changes


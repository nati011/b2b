# Referral

## Overview
The referral domain manages affiliate marketing and influencer referral programmes, similar to social media affiliate marketing. It enables affiliates, influencers, and social media marketers (who may or may not be customers) to share referral links and codes, track conversions from their social media campaigns, and earn commissions when people they refer become customers and place orders. The program is designed for external marketing and customer acquisition through social media channels, not for existing customers to refer each other.

## Domain Boundaries
**Owns**:
- Referral codes and links (unique identifiers for affiliates/influencers - can be customers or non-customers)
- Referral relationships (affiliate/influencer → new customer mapping)
- Referral rewards and earnings (commissions, points, credits, discounts)
- Affiliate program configuration (rules, tiers, commission rates, expiration)
- Referral tracking and attribution (order-to-referral mapping, UTM parameters)
- New customer acquisition tracking (social media referral → customer conversion)
- Affiliate/influencer profiles and performance metrics

**References** (via IDs/identifiers, not domain objects):
- Affiliate/Influencer IDs (referrer - can be customer or non-customer affiliate)
- Order IDs (attributed orders for commission calculation - only from new customers)
- Product IDs (for product-specific affiliate programs)
- Campaign IDs (for tracking social media campaigns)

**Does NOT Own**:
- Customer profiles or authentication
- Order lifecycle or payment processing
- Product catalog or pricing
- Reward redemption or payment processing

**Interaction Patterns**:
- Receives order completion events from Order module (only for new customers)
- Validates referral codes/links during new customer registration or first order placement
- Tracks referral source (social media platform, campaign, UTM parameters)
- Calculates and records affiliate commissions based on order totals from newly acquired customers
- Provides affiliate dashboard and statistics (for affiliates/influencers, not just customers)
- Tracks new customer acquisition through social media and affiliate channels
- Handles affiliate onboarding and profile management

## Responsibilities
- Generate and manage unique referral codes and links for affiliates/influencers (can be customers or non-customers)
- Track referral relationships between affiliates/influencers and new customers
- Validate referral codes/links during new customer registration (not for existing customers)
- Track referral source attribution (social media platform, campaign, UTM parameters)
- Attribute orders from newly acquired customers to affiliate referral codes
- Calculate affiliate commissions based on program rules (only for new customer orders)
- Manage affiliate program tiers, commission rates, and expiration
- Track new customer acquisition and conversion metrics from social media channels
- Prevent existing customers from using referral codes (acquisition-only program)
- Manage affiliate/influencer profiles, onboarding, and performance tracking

## Standard E-Commerce Referral Features

### 1. Referral Code & Link Generation
- **Unique Code Generation**: Generate unique, shareable referral codes for affiliates/influencers
- **Custom Codes**: Allow affiliates to create custom referral codes (branded links)
- **Code Formats**: Support various formats (alphanumeric, UUID, short codes)
- **Referral Links**: Generate trackable referral links with UTM parameters for social media
- **Social Media Links**: Create platform-specific links (Instagram, TikTok, Facebook, etc.)
- **Code Validation**: Validate referral codes for format and expiration
- **Link Shortening**: Provide shortened URLs for social media sharing

### 2. Referral Tracking & Attribution
- **Affiliate-Customer Relationships**: Track which affiliate/influencer referred which new customer
- **Order Attribution**: Link orders to referral codes/links at checkout
- **Social Media Attribution**: Track which social media platform generated the referral
- **UTM Parameter Tracking**: Track campaign, source, medium, content via UTM parameters
- **Multi-Touch Attribution**: Support first-touch, last-touch, or multi-touch attribution models
- **Cookie/Token Tracking**: Track referrals via cookies, tokens, or URL parameters
- **Referral Link Sharing**: Generate shareable referral links optimized for social media
- **Click Tracking**: Track link clicks and conversion rates per affiliate/campaign
- **Platform-Specific Tracking**: Track referrals from Instagram, TikTok, Facebook, Twitter, etc.

### 3. Commission Calculation
- **Percentage-Based Commissions**: Calculate commissions as percentage of order total (e.g., 10% commission)
- **Fixed Amount Commissions**: Support fixed commission amounts per referral or order
- **Tiered Commission Structures**: Implement tiered commission rates (e.g., 10% for first 10 sales, 15% for next 20)
- **Performance-Based Rates**: Higher commission rates for top-performing affiliates
- **Minimum Order Threshold**: Require minimum order value before commissions are earned
- **Maximum Commission Caps**: Set maximum commission limits per referral or time period
- **Commission Currency**: Support commissions in points, credits, cash, or discounts
- **Recurring Commissions**: Support recurring commissions for subscription orders

### 4. Affiliate Program Rules
- **Program Status**: Enable/disable affiliate programs globally or per affiliate
- **Eligibility Rules**: 
  - Affiliates/Influencers: Can be customers or non-customers (anyone can become an affiliate)
  - Referees: Only non-customers (new prospects) can use referral codes/links
  - Existing customers cannot use referral codes (acquisition-only program)
- **Affiliate Onboarding**: Manage affiliate registration, approval, and onboarding process
- **Commission Conditions**: Set conditions for commission eligibility (order status, payment confirmation)
- **Expiration Policies**: Set expiration dates for referral codes or commissions
- **Geographic Restrictions**: Limit referrals to specific regions or countries
- **Product Exclusions**: Exclude certain products or categories from affiliate commissions
- **New Customer Validation**: Ensure referee is a new customer, not an existing one
- **Social Media Platform Rules**: Different rules/rates for different social media platforms

### 5. Commission Distribution
- **Immediate Commissions**: Award commissions immediately upon order completion
- **Delayed Commissions**: Award commissions after order delivery or payment confirmation
- **Pending Commissions**: Track commissions in pending state until conditions are met
- **Commission Reversal**: Handle commission reversals for cancelled or refunded orders
- **Commission History**: Maintain audit trail of all commission transactions
- **Payout Management**: Support various payout methods (bank transfer, PayPal, credits, etc.)
- **Payout Schedules**: Configure payout schedules (weekly, monthly, on-demand)
- **Minimum Payout Thresholds**: Set minimum commission amounts before payout

### 6. Affiliate Analytics & Reporting
- **Affiliate Dashboard**: Provide affiliate/influencer dashboard with real-time statistics
- **Conversion Tracking**: Track referral-to-order conversion rates per affiliate/campaign
- **Earnings Reports**: Show total commissions, pending commissions, and payout history
- **Affiliate Performance**: Track top affiliates, most successful referral codes/links
- **Customer Activity**: Track referred customer behavior (orders placed, lifetime value)
- **Social Media Analytics**: Track performance by social media platform (Instagram, TikTok, etc.)
- **Campaign Performance**: Track performance by campaign, UTM parameters
- **Click-Through Rates**: Track link clicks and conversion rates
- **Revenue Attribution**: Track total revenue generated per affiliate

### 7. Multi-Level Referrals
- **Single-Tier Programs**: Primary focus - referrer (existing customer) earns rewards when new customer places orders
- **Referee Incentives**: New customers may receive sign-up bonuses or first-order discounts (not referral rewards)
- **Multi-Level Marketing**: Support multi-level referral structures where referrers can earn from their downline's referrals
- **Downline Tracking**: Track referral networks and hierarchies (referrer → new customer → their referrals)
- **Commission Splits**: Distribute rewards across multiple levels (referrer and their upline, if applicable)
- **Note**: Referees are new customers and do not earn referral rewards (acquisition-only program)

### 8. Referral Program Management
- **Program Campaigns**: Create time-limited referral campaigns
- **Promotional Periods**: Run special promotions with enhanced rewards
- **A/B Testing**: Test different reward structures and messaging
- **Program Analytics**: Track program performance, ROI, and effectiveness

### 9. Fraud Prevention
- **Self-Referral Prevention**: Prevent affiliates from referring themselves
- **Existing Customer Prevention**: Prevent existing customers from using referral codes
- **New Customer Validation**: Verify that referee is a new customer before awarding commissions
- **Duplicate Detection**: Detect and prevent duplicate referral claims
- **Fraud Monitoring**: Monitor for suspicious referral patterns (e.g., existing customers trying to use codes, fake accounts)
- **Validation Rules**: Validate referral codes, relationships, and customer status
- **IP Address Tracking**: Track and flag suspicious IP patterns
- **Device Fingerprinting**: Detect and prevent fraud through device tracking
- **Bot Detection**: Identify and block bot-generated referrals

### 10. Integration Features
- **Social Media Integration**: Deep integration with Instagram, TikTok, Facebook, Twitter APIs
- **Social Sharing Tools**: Enable easy sharing of referral links via social media platforms
- **Email Integration**: Send referral invitation emails with tracking
- **API Access**: Provide APIs for referral code validation, commission calculation, and affiliate management
- **Webhook Support**: Send webhooks for referral events (new referral, commission earned, payout processed)
- **Affiliate Portal**: Self-service portal for affiliates to manage their accounts, view stats, request payouts
- **Marketing Tools**: Provide marketing materials, banners, and content for affiliates
- **Social Media Widgets**: Embeddable widgets for affiliates to display on their social media profiles

## Use Cases
- Onboard a new affiliate/influencer (can be customer or non-customer)
- Generate a unique referral code/link for an affiliate/influencer
- Share referral links on social media platforms (Instagram, TikTok, Facebook, etc.)
- Track clicks and conversions from social media campaigns
- Validate a referral code/link during new customer registration (not for existing customers)
- Validate referral code during first order placement by a new customer
- Attribute an order from a newly acquired customer to an affiliate referral code
- Calculate and record affiliate commissions when a new customer's order is completed
- Track affiliate performance and earnings (conversion rates, revenue generated)
- Track new customer acquisition metrics through social media and affiliate channels
- Manage affiliate program rules, commission rates, and configurations
- Handle commission reversals for cancelled orders
- Generate affiliate analytics and reports (acquisition-focused, social media performance)
- Process affiliate payouts (commissions)
- Prevent existing customers from using referral codes

## Relationship with Other Modules
- **Customer**: 
  - Validates that referee is a new customer (not existing)
  - May store affiliate information if affiliate is also a customer
  - Displays affiliate dashboard for affiliates (if they are customers)
- **Order**: Receives order completion events from newly acquired customers, attributes orders to affiliate referral codes
- **Product**: May exclude certain products from affiliate commissions
- **Auth/User**: Manages affiliate accounts and authentication (affiliates may have separate accounts)

## Example
Social media affiliate marketing flow:
1) Influencer/Affiliate (can be customer or non-customer) signs up for affiliate program
2) Affiliate receives unique referral code and shareable link (e.g., `yourstore.com/ref/ABC123`)
3) Affiliate shares referral link on Instagram/TikTok/Facebook with their followers
4) Follower (Person B - non-customer) clicks the link and visits the store
5) Person B registers as a new customer using the referral code (or link automatically tracks)
6) Referral module validates that Person B is a new customer (not existing) and records the referral relationship with UTM tracking (source: Instagram, campaign: summer2024)
7) New Customer B places their first order
8) When New Customer B's order is completed, Referral module calculates affiliate commission (e.g., 10% of order total)
9) Commission is credited to Affiliate's account (pending payout)
10) Affiliate can view performance statistics in their dashboard (clicks, conversions, earnings)
11) If an existing customer tries to use a referral code, the system rejects it (acquisition-only program)
12) Affiliate requests payout when they reach minimum threshold

## Best Practices
- **New Customer Validation**: Always validate that referee is a new customer before accepting referral code
- **Existing Customer Rejection**: Reject referral code usage by existing customers (acquisition-only)
- **Affiliate Onboarding**: Streamline affiliate signup and approval process for social media influencers
- Validate referral codes before order attribution
- Use immutable snapshots for commission calculations (based on order total at time of order)
- Implement idempotency for commission calculations to prevent double-crediting
- Maintain audit trail of all referral and commission transactions
- Handle edge cases (cancelled orders, refunds, expired codes, existing customer attempts)
- Track new customer acquisition metrics separately from general customer metrics
- Support affiliate commissions only (referees are new customers, not eligible for commissions)
- Track referral source (social media platform, campaign, UTM parameters) for acquisition analytics
- Monitor for fraud attempts (existing customers trying to use referral codes, fake accounts)
- Provide easy-to-use affiliate tools and marketing materials for social media sharing
- Support multiple payout methods and schedules for affiliates
- Track and optimize performance by social media platform


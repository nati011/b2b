package domain

import "time"

// ReferralRelationship represents the relationship between an affiliate and a new customer.
// This tracks which affiliate referred which new customer.
type ReferralRelationship struct {
	ID             int64
	AffiliateID    int64
	CustomerID     int64  // The new customer who was referred
	ReferralCodeID int64  // The referral code used
	Source         string // Social media platform (Instagram, TikTok, Facebook, etc.)
	Campaign       string // Campaign name/ID
	UTMSource      string // UTM source parameter
	UTMMedium      string // UTM medium parameter
	UTMCampaign    string // UTM campaign parameter
	UTMContent     string // UTM content parameter
	ReferralLink   string // The full referral link used
	IPAddress      string // IP address of the referral
	UserAgent      string // User agent of the referral
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewReferralRelationship creates a new referral relationship.
func NewReferralRelationship(affiliateID, customerID, referralCodeID int64, source, campaign, utmSource, utmMedium, utmCampaign, utmContent, referralLink, ipAddress, userAgent string) (*ReferralRelationship, error) {
	if err := ValidateReferralRelationshipInput(affiliateID, customerID, referralCodeID); err != nil {
		return nil, err
	}

	now := time.Now()
	return &ReferralRelationship{
		AffiliateID:    affiliateID,
		CustomerID:     customerID,
		ReferralCodeID: referralCodeID,
		Source:         source,
		Campaign:       campaign,
		UTMSource:      utmSource,
		UTMMedium:      utmMedium,
		UTMCampaign:    utmCampaign,
		UTMContent:     utmContent,
		ReferralLink:   referralLink,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Update modifies referral relationship tracking data.
func (rr *ReferralRelationship) Update(source, campaign, utmSource, utmMedium, utmCampaign, utmContent string) {
	rr.Source = source
	rr.Campaign = campaign
	rr.UTMSource = utmSource
	rr.UTMMedium = utmMedium
	rr.UTMCampaign = utmCampaign
	rr.UTMContent = utmContent
	rr.UpdatedAt = time.Now()
}

package domain

import "time"

// AffiliateStatus represents the lifecycle status of an affiliate.
type AffiliateStatus string

const (
	AffiliateStatusActive    AffiliateStatus = "active"
	AffiliateStatusInactive  AffiliateStatus = "inactive"
	AffiliateStatusSuspended AffiliateStatus = "suspended"
)

// Affiliate represents an affiliate/influencer profile.
// Affiliates can be customers or non-customers.
type Affiliate struct {
	ID                int64
	AffiliateID       *int64 // Optional: ID if affiliate is also a customer
	Email             string
	PhoneNumber       string
	FullName          string
	Status            AffiliateStatus
	SocialMediaHandle string // Instagram, TikTok, etc.
	Platform          string // Primary social media platform
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewAffiliate creates a new affiliate with validated inputs.
func NewAffiliate(affiliateID *int64, email, phoneNumber, fullName, status, socialMediaHandle, platform string) (*Affiliate, error) {
	parsedStatus, err := ParseAffiliateStatus(status)
	if err != nil {
		return nil, err
	}

	if err := ValidateAffiliateInput(email, phoneNumber, fullName); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Affiliate{
		AffiliateID:       affiliateID,
		Email:             email,
		PhoneNumber:       phoneNumber,
		FullName:          fullName,
		Status:            parsedStatus,
		SocialMediaHandle: socialMediaHandle,
		Platform:          platform,
		IsActive:          parsedStatus == AffiliateStatusActive,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// Update refreshes the affiliate fields in-place.
func (a *Affiliate) Update(email, phoneNumber, fullName, status, socialMediaHandle, platform string) error {
	parsedStatus, err := ParseAffiliateStatus(status)
	if err != nil {
		return err
	}

	if err := ValidateAffiliateInput(email, phoneNumber, fullName); err != nil {
		return err
	}

	a.Email = email
	a.PhoneNumber = phoneNumber
	a.FullName = fullName
	a.Status = parsedStatus
	a.SocialMediaHandle = socialMediaHandle
	a.Platform = platform
	a.IsActive = parsedStatus == AffiliateStatusActive
	a.UpdatedAt = time.Now()
	return nil
}

package domain

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

// ReferralCodeStatus represents the status of a referral code.
type ReferralCodeStatus string

const (
	ReferralCodeStatusActive   ReferralCodeStatus = "active"
	ReferralCodeStatusInactive ReferralCodeStatus = "inactive"
	ReferralCodeStatusExpired  ReferralCodeStatus = "expired"
)

// ReferralCode represents a unique referral code/link for an affiliate.
type ReferralCode struct {
	ID          int64
	AffiliateID int64
	Code        string // Unique referral code (e.g., "ABC123")
	IsCustom    bool   // Whether the code was custom-created by affiliate
	Status      ReferralCodeStatus
	ExpiresAt   *time.Time // Optional expiration date
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewReferralCode creates a new referral code.
// If customCode is provided, it will be used; otherwise, a random code is generated.
func NewReferralCode(affiliateID int64, customCode, status string, expiresAt *time.Time) (*ReferralCode, error) {
	parsedStatus, err := ParseReferralCodeStatus(status)
	if err != nil {
		return nil, err
	}

	var code string
	isCustom := false
	if customCode != "" {
		code = customCode
		isCustom = true
		if err := ValidateReferralCodeFormat(code); err != nil {
			return nil, err
		}
	} else {
		code, err = GenerateReferralCode()
		if err != nil {
			return nil, err
		}
	}

	now := time.Now()
	referralCode := &ReferralCode{
		AffiliateID: affiliateID,
		Code:        code,
		IsCustom:    isCustom,
		Status:      parsedStatus,
		ExpiresAt:   expiresAt,
		IsActive:    parsedStatus == ReferralCodeStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Check if already expired
	if expiresAt != nil && expiresAt.Before(now) {
		referralCode.Status = ReferralCodeStatusExpired
		referralCode.IsActive = false
	}

	return referralCode, nil
}

// Update modifies referral code status.
func (rc *ReferralCode) Update(status string, expiresAt *time.Time) error {
	parsedStatus, err := ParseReferralCodeStatus(status)
	if err != nil {
		return err
	}

	rc.Status = parsedStatus
	rc.ExpiresAt = expiresAt
	now := time.Now()

	// Check if expired
	if expiresAt != nil && expiresAt.Before(now) {
		rc.Status = ReferralCodeStatusExpired
		rc.IsActive = false
	} else {
		rc.IsActive = parsedStatus == ReferralCodeStatusActive
	}

	rc.UpdatedAt = now
	return nil
}

// IsExpired checks if the referral code has expired.
func (rc *ReferralCode) IsExpired() bool {
	if rc.ExpiresAt == nil {
		return false
	}
	return rc.ExpiresAt.Before(time.Now())
}

// GenerateReferralCode generates a random alphanumeric referral code.
func GenerateReferralCode() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := base64.URLEncoding.EncodeToString(bytes)
	// Remove padding and make it uppercase, alphanumeric only
	if len(code) > 8 {
		code = code[:8]
	}
	// Convert to alphanumeric only (remove special chars)
	var result strings.Builder
	for _, c := range code {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			result.WriteRune(c)
		}
	}
	// Ensure we have at least 6 characters
	if result.Len() < 6 {
		return GenerateReferralCode()
	}
	final := result.String()
	if len(final) > 8 {
		return final[:8], nil
	}
	return final, nil
}

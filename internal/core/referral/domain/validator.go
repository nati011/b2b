package domain

import (
	"regexp"
	"strings"

	"marketplace/pkg/validate"
)

// ParseAffiliateStatus normalizes and validates an affiliate status value.
func ParseAffiliateStatus(value string) (AffiliateStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return AffiliateStatusInactive, nil
	}

	validStatuses := []string{
		string(AffiliateStatusActive),
		string(AffiliateStatusInactive),
		string(AffiliateStatusSuspended),
	}

	for _, status := range validStatuses {
		if normalized == status {
			return AffiliateStatus(normalized), nil
		}
	}

	return "", validate.NewJSONError([]string{"invalid affiliate status"})
}

// ValidateAffiliateInput validates core affiliate fields.
func ValidateAffiliateInput(email, phoneNumber, fullName string) error {
	if strings.TrimSpace(email) == "" && strings.TrimSpace(phoneNumber) == "" {
		return validate.NewJSONError([]string{"either email or phone number must be provided"})
	}

	result := validate.New().
		And(validate.NonEmpty(fullName)).
		And(validate.MaxLen(fullName, 255))

	if email != "" {
		result.And(validate.EmailValid(email)).
			And(validate.MaxLen(email, 255))
	}

	if phoneNumber != "" {
		result.And(validate.MinLen(phoneNumber, 10)).
			And(validate.MaxLen(phoneNumber, 50))
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validate.NewJSONError(validation.Message)
	}

	return nil
}

// ParseReferralCodeStatus normalizes and validates a referral code status value.
func ParseReferralCodeStatus(value string) (ReferralCodeStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return ReferralCodeStatusActive, nil
	}

	validStatuses := []string{
		string(ReferralCodeStatusActive),
		string(ReferralCodeStatusInactive),
		string(ReferralCodeStatusExpired),
	}

	for _, status := range validStatuses {
		if normalized == status {
			return ReferralCodeStatus(normalized), nil
		}
	}

	return "", validate.NewJSONError([]string{"invalid referral code status"})
}

// ValidateReferralCodeFormat validates the format of a referral code.
func ValidateReferralCodeFormat(code string) error {
	// Allow alphanumeric codes, 4-20 characters
	matched, err := regexp.MatchString(`^[A-Za-z0-9]{4,20}$`, code)
	if err != nil {
		return validate.NewJSONError([]string{"invalid referral code format"})
	}
	if !matched {
		return validate.NewJSONError([]string{"referral code must be 4-20 alphanumeric characters"})
	}
	return nil
}

// ParseCommissionStatus normalizes and validates a commission status value.
func ParseCommissionStatus(value string) (CommissionStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return CommissionStatusPending, nil
	}

	validStatuses := []string{
		string(CommissionStatusPending),
		string(CommissionStatusApproved),
		string(CommissionStatusPaid),
		string(CommissionStatusReversed),
		string(CommissionStatusCancelled),
	}

	for _, status := range validStatuses {
		if normalized == status {
			return CommissionStatus(normalized), nil
		}
	}

	return "", validate.NewJSONError([]string{"invalid commission status"})
}

// ValidateReferralRelationshipInput validates referral relationship inputs.
func ValidateReferralRelationshipInput(affiliateID, customerID, referralCodeID int64) error {
	if affiliateID <= 0 {
		return validate.NewJSONError([]string{"affiliate ID is required"})
	}
	if customerID <= 0 {
		return validate.NewJSONError([]string{"customer ID is required"})
	}
	if referralCodeID <= 0 {
		return validate.NewJSONError([]string{"referral code ID is required"})
	}
	return nil
}

// ValidateCommissionInput validates commission inputs.
func ValidateCommissionInput(affiliateID, referralRelationshipID, orderID, customerID int64, commissionRate float64) error {
	if affiliateID <= 0 {
		return validate.NewJSONError([]string{"affiliate ID is required"})
	}
	if referralRelationshipID <= 0 {
		return validate.NewJSONError([]string{"referral relationship ID is required"})
	}
	if orderID <= 0 {
		return validate.NewJSONError([]string{"order ID is required"})
	}
	if customerID <= 0 {
		return validate.NewJSONError([]string{"customer ID is required"})
	}
	if commissionRate < 0 || commissionRate > 100 {
		return validate.NewJSONError([]string{"commission rate must be between 0 and 100"})
	}
	return nil
}

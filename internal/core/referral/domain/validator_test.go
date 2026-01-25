package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAffiliateStatus(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      AffiliateStatus
		wantError bool
	}{
		{"valid active", "active", AffiliateStatusActive, false},
		{"valid inactive", "inactive", AffiliateStatusInactive, false},
		{"valid suspended", "suspended", AffiliateStatusSuspended, false},
		{"empty defaults to inactive", "", AffiliateStatusInactive, false},
		{"case insensitive", "ACTIVE", AffiliateStatusActive, false},
		{"invalid status", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAffiliateStatus(tt.input)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestValidateAffiliateInput(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		phone     string
		fullName  string
		wantError bool
	}{
		{"valid with email", "test@example.com", "", "John Doe", false},
		{"valid with phone", "", "0912345678", "John Doe", false},
		{"valid with both", "test@example.com", "0912345678", "John Doe", false},
		{"missing both email and phone", "", "", "John Doe", true},
		{"missing full name", "test@example.com", "", "", true},
		{"invalid email", "invalid-email", "", "John Doe", true},
		{"short phone", "test@example.com", "123", "John Doe", true},
		{"long full name", "test@example.com", "", string(make([]byte, 300)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAffiliateInput(tt.email, tt.phone, tt.fullName)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseReferralCodeStatus(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      ReferralCodeStatus
		wantError bool
	}{
		{"valid active", "active", ReferralCodeStatusActive, false},
		{"valid inactive", "inactive", ReferralCodeStatusInactive, false},
		{"valid expired", "expired", ReferralCodeStatusExpired, false},
		{"empty defaults to active", "", ReferralCodeStatusActive, false},
		{"case insensitive", "ACTIVE", ReferralCodeStatusActive, false},
		{"invalid status", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseReferralCodeStatus(tt.input)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestValidateReferralCodeFormat(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantError bool
	}{
		{"valid alphanumeric", "ABC123", false},
		{"valid lowercase", "abc123", false},
		{"valid mixed case", "AbC123", false},
		{"valid numbers only", "123456", false},
		{"valid letters only", "ABCDEF", false},
		{"valid max length", "ABCDEFGHIJKLMNOPQRST", false}, // 20 chars is max
		{"too short", "ABC", true},
		{"too long", "ABCDEFGHIJKLMNOPQRSTU", true}, // 21 chars exceeds max
		{"contains special chars", "ABC-123", true},
		{"contains spaces", "ABC 123", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReferralCodeFormat(tt.code)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseCommissionStatus(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      CommissionStatus
		wantError bool
	}{
		{"valid pending", "pending", CommissionStatusPending, false},
		{"valid approved", "approved", CommissionStatusApproved, false},
		{"valid paid", "paid", CommissionStatusPaid, false},
		{"valid reversed", "reversed", CommissionStatusReversed, false},
		{"valid cancelled", "cancelled", CommissionStatusCancelled, false},
		{"empty defaults to pending", "", CommissionStatusPending, false},
		{"case insensitive", "PENDING", CommissionStatusPending, false},
		{"invalid status", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCommissionStatus(tt.input)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestValidateReferralRelationshipInput(t *testing.T) {
	tests := []struct {
		name          string
		affiliateID   int64
		customerID    int64
		referralCodeID int64
		wantError     bool
	}{
		{"valid", 1, 2, 3, false},
		{"zero affiliate ID", 0, 2, 3, true},
		{"zero customer ID", 1, 0, 3, true},
		{"zero referral code ID", 1, 2, 0, true},
		{"all zero", 0, 0, 0, true},
		{"negative affiliate ID", -1, 2, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReferralRelationshipInput(tt.affiliateID, tt.customerID, tt.referralCodeID)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCommissionInput(t *testing.T) {
	tests := []struct {
		name                  string
		affiliateID           int64
		referralRelationshipID int64
		orderID               int64
		customerID            int64
		commissionRate        float64
		wantError             bool
	}{
		{"valid", 1, 2, 3, 4, 10.0, false},
		{"zero affiliate ID", 0, 2, 3, 4, 10.0, true},
		{"zero referral relationship ID", 1, 0, 3, 4, 10.0, true},
		{"zero order ID", 1, 2, 0, 4, 10.0, true},
		{"zero customer ID", 1, 2, 3, 0, 10.0, true},
		{"negative commission rate", 1, 2, 3, 4, -1.0, true},
		{"commission rate over 100", 1, 2, 3, 4, 101.0, true},
		{"zero commission rate", 1, 2, 3, 4, 0.0, false},
		{"100 commission rate", 1, 2, 3, 4, 100.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommissionInput(tt.affiliateID, tt.referralRelationshipID, tt.orderID, tt.customerID, tt.commissionRate)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}


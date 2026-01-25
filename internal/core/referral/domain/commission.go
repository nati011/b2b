package domain

import (
	"time"

	goodmoney "github.com/the-nucleus-project/good_money"
)

// CommissionStatus represents the status of a commission.
type CommissionStatus string

const (
	CommissionStatusPending   CommissionStatus = "pending"
	CommissionStatusApproved  CommissionStatus = "approved"
	CommissionStatusPaid      CommissionStatus = "paid"
	CommissionStatusReversed  CommissionStatus = "reversed"
	CommissionStatusCancelled CommissionStatus = "cancelled"
)

// Commission represents a commission earned by an affiliate from a referred customer's order.
type Commission struct {
	ID                     int64
	AffiliateID            int64
	ReferralRelationshipID int64
	OrderID                int64
	CustomerID             int64 // The new customer who placed the order
	Amount                 *goodmoney.Money
	CommissionRate         float64 // Percentage (e.g., 10.0 for 10%)
	Status                 CommissionStatus
	OrderTotal             *goodmoney.Money // Snapshot of order total at time of commission calculation
	Notes                  string
	PaidAt                 *time.Time
	ReversedAt             *time.Time
	IsActive               bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// NewCommission creates a new commission record.
func NewCommission(affiliateID, referralRelationshipID, orderID, customerID int64, amount *goodmoney.Money, commissionRate float64, orderTotal *goodmoney.Money, notes string) (*Commission, error) {
	if err := ValidateCommissionInput(affiliateID, referralRelationshipID, orderID, customerID, commissionRate); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Commission{
		AffiliateID:            affiliateID,
		ReferralRelationshipID: referralRelationshipID,
		OrderID:                orderID,
		CustomerID:             customerID,
		Amount:                 amount,
		CommissionRate:         commissionRate,
		Status:                 CommissionStatusPending,
		OrderTotal:             orderTotal,
		Notes:                  notes,
		IsActive:               true,
		CreatedAt:              now,
		UpdatedAt:              now,
	}, nil
}

// UpdateStatus updates the commission status.
func (c *Commission) UpdateStatus(status string) error {
	parsedStatus, err := ParseCommissionStatus(status)
	if err != nil {
		return err
	}

	c.Status = parsedStatus
	c.UpdatedAt = time.Now()

	// Set timestamps based on status
	now := time.Now()
	if parsedStatus == CommissionStatusPaid && c.PaidAt == nil {
		c.PaidAt = &now
	}
	if parsedStatus == CommissionStatusReversed && c.ReversedAt == nil {
		c.ReversedAt = &now
	}

	return nil
}

// Reverse reverses a commission (e.g., for cancelled orders).
func (c *Commission) Reverse(reason string) error {
	if c.Status == CommissionStatusReversed || c.Status == CommissionStatusCancelled {
		return nil // Already reversed
	}

	c.Status = CommissionStatusReversed
	now := time.Now()
	c.ReversedAt = &now
	if reason != "" {
		c.Notes = reason
	}
	c.UpdatedAt = now
	return nil
}

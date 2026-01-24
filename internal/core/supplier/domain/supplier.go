package domain

import "time"

// SupplierStatus represents the lifecycle status of a supplier.
type SupplierStatus string

const (
	SupplierStatusActive    SupplierStatus = "active"
	SupplierStatusInactive  SupplierStatus = "inactive"
	SupplierStatusSuspended SupplierStatus = "suspended"
)

// Supplier represents a supplier profile record.
type Supplier struct {
	ID           int64
	BusinessName string
	Status       SupplierStatus
	SupportEmail string
	SupportPhone string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewSupplier creates a new supplier with validated inputs.
func NewSupplier(businessName, status, supportEmail, supportPhone string) (*Supplier, error) {
	parsedStatus, err := ParseSupplierStatus(status)
	if err != nil {
		return nil, err
	}

	if err := ValidateSupplierInput(businessName, supportEmail, supportPhone); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Supplier{
		BusinessName: businessName,
		Status:       parsedStatus,
		SupportEmail: supportEmail,
		SupportPhone: supportPhone,
		IsActive:     parsedStatus == SupplierStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update refreshes the supplier fields in-place.
func (s *Supplier) Update(businessName, status, supportEmail, supportPhone string) error {
	parsedStatus, err := ParseSupplierStatus(status)
	if err != nil {
		return err
	}

	if err := ValidateSupplierInput(businessName, supportEmail, supportPhone); err != nil {
		return err
	}

	s.BusinessName = businessName
	s.Status = parsedStatus
	s.SupportEmail = supportEmail
	s.SupportPhone = supportPhone
	s.IsActive = parsedStatus == SupplierStatusActive
	s.UpdatedAt = time.Now()
	return nil
}

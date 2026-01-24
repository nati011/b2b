package domain

import "time"

// CustomerStatus represents the lifecycle status of a customer.
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusSuspended CustomerStatus = "suspended"
)

// Customer represents a customer profile record.
type Customer struct {
	ID          int64
	FullName    string
	Status      CustomerStatus
	City        string
	Region      string
	Woreda      string
	PhoneNumber string
	Email       string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCustomer creates a new customer with validated inputs.
func NewCustomer(fullName, status, city, region, woreda, phoneNumber, email string) (*Customer, error) {
	parsedStatus, err := ParseCustomerStatus(status)
	if err != nil {
		return nil, err
	}

	if err := ValidateCustomerInput(fullName, city, region, woreda, phoneNumber, email); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Customer{
		FullName:    fullName,
		Status:      parsedStatus,
		City:        city,
		Region:      region,
		Woreda:      woreda,
		PhoneNumber: phoneNumber,
		Email:       email,
		IsActive:    parsedStatus == CustomerStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update refreshes the customer fields in-place.
func (c *Customer) Update(fullName, status, city, region, woreda, phoneNumber, email string) error {
	parsedStatus, err := ParseCustomerStatus(status)
	if err != nil {
		return err
	}

	if err := ValidateCustomerInput(fullName, city, region, woreda, phoneNumber, email); err != nil {
		return err
	}

	c.FullName = fullName
	c.Status = parsedStatus
	c.City = city
	c.Region = region
	c.Woreda = woreda
	c.PhoneNumber = phoneNumber
	c.Email = email
	c.IsActive = parsedStatus == CustomerStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

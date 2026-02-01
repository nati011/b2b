package domain

import (
	"encoding/json"
	"time"
)

// Product represents a catalog product and its snapshot metadata.
type Product struct {
	ID                int64
	Name              string
	Description       string
	ExternalID        string
	Attributes        json.RawMessage
	Unit              string
	IsActive          bool
	SupplierID        int64
	Price             *float64
	TotalQuantity     int
	ReservedQuantity  int
	AvailableQuantity int
	CategoryIDs       []int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ProductMetadata holds optional product fields.
type ProductMetadata struct {
	Description      string
	ExternalID       string
	Attributes       json.RawMessage
	Unit             string
	IsActive         bool
	Price            *float64
	TotalQuantity    int
	ReservedQuantity int
	CategoryIDs      []int64
}

// NewProduct creates a new product entity after validation.
func NewProduct(supplierID int64, name string, metadata ProductMetadata) (*Product, error) {
	if err := ValidateProductInput(supplierID, name, metadata); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Product{
		Name:             name,
		Description:      metadata.Description,
		ExternalID:       metadata.ExternalID,
		Attributes:       metadata.Attributes,
		Unit:             metadata.Unit,
		IsActive:         metadata.IsActive,
		SupplierID:       supplierID,
		Price:            metadata.Price,
		TotalQuantity:    metadata.TotalQuantity,
		ReservedQuantity: metadata.ReservedQuantity,
		CategoryIDs:      metadata.CategoryIDs,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Update modifies product fields in-place.
func (p *Product) Update(name string, metadata ProductMetadata) error {
	if err := ValidateProductInput(p.SupplierID, name, metadata); err != nil {
		return err
	}

	p.Name = name
	p.Description = metadata.Description
	p.ExternalID = metadata.ExternalID
	p.Attributes = metadata.Attributes
	p.Unit = metadata.Unit
	p.IsActive = metadata.IsActive
	p.Price = metadata.Price
	p.TotalQuantity = metadata.TotalQuantity
	p.ReservedQuantity = metadata.ReservedQuantity
	p.CategoryIDs = metadata.CategoryIDs
	p.UpdatedAt = time.Now()
	return nil
}

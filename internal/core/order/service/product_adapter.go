package order

import (
	"context"
	productservice "marketplace/internal/core/product/service"
)

// productServiceAdapter adapts productservice.Service to ProductService interface.
type productServiceAdapter struct {
	service *productservice.Service
}

// NewProductServiceAdapter creates a new adapter.
func NewProductServiceAdapter(service *productservice.Service) ProductService {
	return &productServiceAdapter{service: service}
}

// Get retrieves a product by ID and returns it as a Product.
func (a *productServiceAdapter) Get(ctx context.Context, id int64) (*Product, error) {
	p, err := a.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:         p.ID,
		SupplierID: p.SupplierID,
	}, nil
}




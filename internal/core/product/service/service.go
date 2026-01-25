package service

import (
	"context"
	"errors"
	"marketplace/internal/core/product/domain"
	"marketplace/pkg/logger"

	goodmoney "github.com/the-nucleus-project/good_money"
)

var (
	// ErrProductNotFound indicates a product lookup failure.
	ErrProductNotFound = errors.New("product not found")
)

// ProductInput captures product payload data.
type ProductInput struct {
	SupplierID       int64
	Name             string
	Description      string
	ExternalID       string
	Attributes       []byte
	Unit             string
	IsActive         bool
	Price            *float64
	Currency         string
	TotalQuantity    int
	ReservedQuantity int
	CategoryIDs      []int64
}

// ListQuery defines filters for product listing.
type ListQuery struct {
	SupplierID int64
	CategoryID int64
	IsActive   *bool
	Limit      int
	Offset     int
}

// Repository defines product persistence operations.
type Repository interface {
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context, query ListQuery) ([]*domain.Product, int, error)
}

// Service coordinates product business logic.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// Create registers a new product.
func (s *Service) Create(ctx context.Context, input ProductInput) (*domain.Product, error) {
	// Default currency to ETB if not provided
	currency := input.Currency
	if currency == "" {
		currency = "ETB"
	}

	var price *goodmoney.Money
	if input.Price != nil {
		m, err := goodmoney.New(*input.Price, currency)
		if err != nil {
			return nil, err
		}
		price = m
	}

	metadata := domain.ProductMetadata{
		Description:      input.Description,
		ExternalID:       input.ExternalID,
		Attributes:       input.Attributes,
		Unit:             input.Unit,
		IsActive:         input.IsActive,
		Price:            price,
		TotalQuantity:    input.TotalQuantity,
		ReservedQuantity: input.ReservedQuantity,
		CategoryIDs:      input.CategoryIDs,
	}

	product, err := domain.NewProduct(input.SupplierID, input.Name, metadata)
	if err != nil {
		logger.Warn("Product creation failed: validation error", "supplier_id", input.SupplierID, "error", err)
		return nil, err
	}

	if err := s.repository.Create(ctx, product); err != nil {
		logger.Error("Product creation failed: repository error", "supplier_id", input.SupplierID, "error", err)
		return nil, err
	}

	logger.Info("Product created successfully", "product_id", product.ID, "supplier_id", product.SupplierID)
	return product, nil
}

// Update modifies an existing product.
func (s *Service) Update(ctx context.Context, id int64, input ProductInput) (*domain.Product, error) {
	product, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Product update failed: not found", "product_id", id)
		return nil, err
	}

	// Use existing product's currency if available, otherwise default to ETB
	currency := input.Currency
	if currency == "" && product.Price != nil {
		currency = product.Price.Currency()
	}
	if currency == "" {
		currency = "ETB"
	}

	var price *goodmoney.Money
	if input.Price != nil {
		m, err := goodmoney.New(*input.Price, currency)
		if err != nil {
			return nil, err
		}
		price = m
	}

	metadata := domain.ProductMetadata{
		Description:      input.Description,
		ExternalID:       input.ExternalID,
		Attributes:       input.Attributes,
		Unit:             input.Unit,
		IsActive:         input.IsActive,
		Price:            price,
		TotalQuantity:    input.TotalQuantity,
		ReservedQuantity: input.ReservedQuantity,
		CategoryIDs:      input.CategoryIDs,
	}

	if err := product.Update(input.Name, metadata); err != nil {
		logger.Warn("Product update failed: validation error", "product_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, product); err != nil {
		logger.Error("Product update failed: repository error", "product_id", id, "error", err)
		return nil, err
	}

	logger.Info("Product updated successfully", "product_id", id)
	return product, nil
}

// Delete removes a product (soft delete).
func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrProductNotFound) {
			logger.Debug("Product deletion failed: not found", "product_id", id)
		} else {
			logger.Error("Product deletion failed: repository error", "product_id", id, "error", err)
		}
		return err
	}

	logger.Info("Product deleted successfully", "product_id", id)
	return nil
}

// Get retrieves a product by ID.
func (s *Service) Get(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			logger.Debug("Product retrieval failed: not found", "product_id", id)
		} else {
			logger.Error("Product retrieval failed: repository error", "product_id", id, "error", err)
		}
		return nil, err
	}

	logger.Debug("Product retrieved successfully", "product_id", id)
	return product, nil
}

// List retrieves products with optional filters.
func (s *Service) List(ctx context.Context, query ListQuery) ([]*domain.Product, int, error) {
	products, total, err := s.repository.List(ctx, query)
	if err != nil {
		logger.Error("Product list failed: repository error", "error", err)
		return nil, 0, err
	}
	logger.Debug("Product list completed", "count", len(products), "total", total)
	return products, total, nil
}

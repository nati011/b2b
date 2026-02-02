package service

import (
	"context"
	"errors"
	"time"

	"marketplace/internal/core/product/domain"
	"marketplace/pkg/logger"
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
	RecordPriceChange(ctx context.Context, productID int64, oldPrice, newPrice *float64, userID *int64, userEmail, reason string) error
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
	metadata := domain.ProductMetadata{
		Description:      input.Description,
		ExternalID:       input.ExternalID,
		Attributes:       input.Attributes,
		Unit:             input.Unit,
		IsActive:         input.IsActive,
		Price:            input.Price,
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

	metadata := domain.ProductMetadata{
		Description:      input.Description,
		ExternalID:       input.ExternalID,
		Attributes:       input.Attributes,
		Unit:             input.Unit,
		IsActive:         input.IsActive,
		Price:            input.Price,
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

// CreateGRN creates a goods receiving note and increases product quantity.
func (s *Service) CreateGRN(ctx context.Context, productID int64, quantity int, notes string) (*domain.Product, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	product, err := s.repository.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			logger.Debug("GRN creation failed: product not found", "product_id", productID)
		} else {
			logger.Error("GRN creation failed: repository error", "product_id", productID, "error", err)
		}
		return nil, err
	}

	// Increase total quantity
	product.TotalQuantity += quantity
	product.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, product); err != nil {
		logger.Error("GRN creation failed: failed to update product", "product_id", productID, "error", err)
		return nil, err
	}

	logger.Info("GRN created successfully", "product_id", productID, "quantity", quantity, "notes", notes)
	return product, nil
}

// PriceUpdateInput captures price update payload data.
type PriceUpdateInput struct {
	ProductID int64
	NewPrice  float64
	Reason    string
	UserID    *int64
	UserEmail string
}

// UpdatePrice updates product price and records the change in price history.
func (s *Service) UpdatePrice(ctx context.Context, input PriceUpdateInput) (*domain.Product, error) {
	if input.NewPrice < 0 {
		return nil, errors.New("price must be non-negative")
	}

	product, err := s.repository.FindByID(ctx, input.ProductID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			logger.Debug("Price update failed: product not found", "product_id", input.ProductID)
		} else {
			logger.Error("Price update failed: repository error", "product_id", input.ProductID, "error", err)
		}
		return nil, err
	}

	oldPrice := product.Price
	newPrice := &input.NewPrice

	// Update product price
	product.Price = newPrice
	product.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, product); err != nil {
		logger.Error("Price update failed: failed to update product", "product_id", input.ProductID, "error", err)
		return nil, err
	}

	// Record price change in history
	if err := s.repository.RecordPriceChange(ctx, input.ProductID, oldPrice, newPrice, input.UserID, input.UserEmail, input.Reason); err != nil {
		logger.Warn("Price update succeeded but history recording failed", "product_id", input.ProductID, "error", err)
		// Don't fail the update if history recording fails, but log it
	}

	logger.Info("Price updated successfully", "product_id", input.ProductID, "old_price", oldPrice, "new_price", newPrice, "reason", input.Reason)
	return product, nil
}

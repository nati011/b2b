package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace/internal/core/product/domain"

	"github.com/stretchr/testify/require"
)

type mockProductRepository struct {
	products      map[int64]*domain.Product
	nextID        int64
	createFunc    func(context.Context, *domain.Product) error
	updateFunc    func(context.Context, *domain.Product) error
	deleteFunc    func(context.Context, int64) error
	findByIDFunc  func(context.Context, int64) (*domain.Product, error)
	listFunc      func(context.Context, ListQuery) ([]*domain.Product, int, error)
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[int64]*domain.Product),
		nextID:   1,
	}
}

func (m *mockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, product)
	}
	if product.ID == 0 {
		product.ID = m.nextID
		m.nextID++
	}
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, product)
	}
	if _, ok := m.products[product.ID]; !ok {
		return ErrProductNotFound
	}
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	if _, ok := m.products[id]; !ok {
		return ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *mockProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	product, ok := m.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (m *mockProductRepository) RecordPriceChange(ctx context.Context, productID int64, oldPrice, newPrice *float64, userID *int64, userEmail, reason string) error {
	// Mock implementation - just return nil
	return nil
}

func (m *mockProductRepository) List(ctx context.Context, query ListQuery) ([]*domain.Product, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, query)
	}
	var results []*domain.Product
	for _, product := range m.products {
		matches := true
		if query.SupplierID != 0 && product.SupplierID != query.SupplierID {
			matches = false
		}
		if query.IsActive != nil && product.IsActive != *query.IsActive {
			matches = false
		}
		if matches {
			results = append(results, product)
		}
	}
	return results, len(results), nil
}

func TestProductServiceCreateSuccess(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	attributes, _ := json.Marshal(map[string]interface{}{
		"color": "red",
		"size":  "large",
	})
	price := 99.99

	result, err := service.Create(context.Background(), ProductInput{
		SupplierID:       1,
		Name:             "Test Product",
		Description:      "A test product",
		ExternalID:       "EXT-001",
		Attributes:       attributes,
		Unit:             "piece",
		IsActive:         true,
		Price:            &price,
		TotalQuantity:    100,
		ReservedQuantity: 10,
		CategoryIDs:      []int64{1, 2},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotZero(t, result.ID)
	require.Equal(t, "Test Product", result.Name)
	require.Equal(t, int64(1), result.SupplierID)
	require.True(t, result.IsActive)
	require.Equal(t, 100, result.TotalQuantity)
}

func TestProductServiceCreateValidationError(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	_, err := service.Create(context.Background(), ProductInput{
		SupplierID: 0, // Invalid supplier ID
		Name:       "", // Empty name
	})
	require.Error(t, err)
}

func TestProductServiceGetSuccess(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	product := &domain.Product{
		ID:           1,
		Name:         "Existing Product",
		SupplierID:   1,
		IsActive:     true,
		TotalQuantity: 50,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.products[1] = product

	result, err := service.Get(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, "Existing Product", result.Name)
}

func TestProductServiceGetNotFound(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	_, err := service.Get(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrProductNotFound))
}

func TestProductServiceUpdateSuccess(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	product := &domain.Product{
		ID:           1,
		Name:         "Original Name",
		SupplierID:   1,
		IsActive:     true,
		TotalQuantity: 50,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.products[1] = product

	price := 149.99
	result, err := service.Update(context.Background(), 1, ProductInput{
		SupplierID:       1,
		Name:              "Updated Name",
		Description:       "Updated description",
		IsActive:          true,
		Price:             &price,
		TotalQuantity:     75,
		ReservedQuantity:  5,
		CategoryIDs:       []int64{1},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "Updated Name", result.Name)
	require.Equal(t, 75, result.TotalQuantity)
}

func TestProductServiceUpdateNotFound(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	_, err := service.Update(context.Background(), 999, ProductInput{
		SupplierID: 1,
		Name:       "Test",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrProductNotFound))
}

func TestProductServiceDeleteSuccess(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	product := &domain.Product{
		ID:         1,
		Name:       "To Delete",
		SupplierID: 1,
	}
	repo.products[1] = product

	err := service.Delete(context.Background(), 1)
	require.NoError(t, err)
	_, exists := repo.products[1]
	require.False(t, exists)
}

func TestProductServiceDeleteNotFound(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	err := service.Delete(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrProductNotFound))
}

func TestProductServiceListSuccess(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	product1 := &domain.Product{
		ID:         1,
		Name:       "Product 1",
		SupplierID: 1,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	product2 := &domain.Product{
		ID:         2,
		Name:       "Product 2",
		SupplierID: 1,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	product3 := &domain.Product{
		ID:         3,
		Name:       "Product 3",
		SupplierID: 2,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	repo.products[1] = product1
	repo.products[2] = product2
	repo.products[3] = product3

	products, total, err := service.List(context.Background(), ListQuery{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Len(t, products, 3)
}

func TestProductServiceListWithSupplierFilter(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	product1 := &domain.Product{
		ID:         1,
		Name:       "Product 1",
		SupplierID: 1,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	product2 := &domain.Product{
		ID:         2,
		Name:       "Product 2",
		SupplierID: 2,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	repo.products[1] = product1
	repo.products[2] = product2

	products, total, err := service.List(context.Background(), ListQuery{
		SupplierID: 1,
		Limit:      10,
		Offset:     0,
	})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, products, 1)
	require.Equal(t, int64(1), products[0].SupplierID)
}

func TestProductServiceListWithActiveFilter(t *testing.T) {
	repo := newMockProductRepository()
	service := NewService(repo)

	active := true
	product1 := &domain.Product{
		ID:         1,
		Name:       "Active Product",
		SupplierID: 1,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	product2 := &domain.Product{
		ID:         2,
		Name:       "Inactive Product",
		SupplierID: 1,
		IsActive:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	repo.products[1] = product1
	repo.products[2] = product2

	products, total, err := service.List(context.Background(), ListQuery{
		IsActive: &active,
		Limit:    10,
		Offset:   0,
	})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, products, 1)
	require.True(t, products[0].IsActive)
}


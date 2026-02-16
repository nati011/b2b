package order

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace/internal/core/order/domain"

	"github.com/stretchr/testify/require"
)

type mockOrderRepository struct {
	orders             map[int64]*domain.Order
	nextID             int64
	createFunc         func(context.Context, *domain.Order) error
	updateFunc         func(context.Context, *domain.Order) error
	findByIDFunc       func(context.Context, int64) (*domain.Order, error)
	listByCustomerFunc func(context.Context, CustomerOrderQuery) ([]*domain.Order, int, error)
	listBySupplierFunc func(context.Context, SupplierOrderQuery) ([]*domain.Order, int, error)
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order)
	}
	if order.ID == 0 {
		order.ID = m.nextID
		m.nextID++
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, order)
	}
	if _, ok := m.orders[order.ID]; !ok {
		return ErrOrderNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	order, ok := m.orders[orderID]
	if !ok {
		return ErrOrderNotFound
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
}

func (m *mockOrderRepository) UpdatePaymentStatus(ctx context.Context, orderID int64, paymentStatus string) error {
	order, ok := m.orders[orderID]
	if !ok {
		return ErrOrderNotFound
	}
	order.PaymentStatus = paymentStatus
	order.UpdatedAt = time.Now()
	return nil
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	order, ok := m.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (m *mockOrderRepository) ListByCustomer(ctx context.Context, query CustomerOrderQuery) ([]*domain.Order, int, error) {
	if m.listByCustomerFunc != nil {
		return m.listByCustomerFunc(ctx, query)
	}
	var results []*domain.Order
	for _, order := range m.orders {
		if order.CustomerID == query.CustomerID {
			if query.Status == "" || string(order.Status) == query.Status {
				results = append(results, order)
			}
		}
	}
	return results, len(results), nil
}

func (m *mockOrderRepository) ListBySupplier(ctx context.Context, query SupplierOrderQuery) ([]*domain.Order, int, error) {
	if m.listBySupplierFunc != nil {
		return m.listBySupplierFunc(ctx, query)
	}
	// Mock implementation - return empty list for now
	return []*domain.Order{}, 0, nil
}

func (m *mockOrderRepository) BelongsToSupplier(ctx context.Context, orderID int64, supplierID int64) (bool, error) {
	// Mock implementation - return true for testing
	return true, nil
}

type mockProductService struct {
	products map[int64]*Product
}

func newMockProductService() *mockProductService {
	return &mockProductService{
		products: make(map[int64]*Product),
	}
}

func (m *mockProductService) Get(ctx context.Context, id int64) (*Product, error) {
	product, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (m *mockProductService) setProduct(id int64, supplierID int64) {
	m.products[id] = &Product{
		ID:         id,
		SupplierID: supplierID,
	}
}

func TestOrderServiceCreateSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	productService.setProduct(1, 10) // Product ID 1 has supplier ID 10
	service := NewService(repo, productService)

	customerSnapshot, _ := json.Marshal(map[string]interface{}{
		"id":   1,
		"name": "Test Customer",
	})

	total := 100.50
	result, err := service.Create(context.Background(), OrderInput{
		CustomerID:         1,
		Status:             "pending",
		PaymentStatus:      "unpaid",
		DeliveryStatus:     "pending",
		ConfirmationStatus: "pending",
		Total:              &total,
		ReferralCode:       "",
		CustomerSnapshot:   customerSnapshot,
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 2, Price: &total},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotZero(t, result.ID)
	require.Equal(t, domain.OrderStatusPending, result.Status)
	require.Equal(t, int64(1), result.CustomerID)
	require.Equal(t, int64(10), result.SupplierID) // Should be set from product
	require.Empty(t, result.ReferralCode)
	require.Len(t, result.Items, 1)
}

func TestOrderServiceCreateWithReferralCode(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	productService.setProduct(1, 10)
	service := NewService(repo, productService)

	customerSnapshot, _ := json.Marshal(map[string]interface{}{
		"id":   1,
		"name": "Test Customer",
	})
	total := 50.0

	result, err := service.Create(context.Background(), OrderInput{
		CustomerID:         1,
		Status:             "pending",
		PaymentStatus:      "unpaid",
		DeliveryStatus:     "pending",
		ConfirmationStatus: "pending",
		Total:              &total,
		ReferralCode:       "REF123",
		CustomerSnapshot:   customerSnapshot,
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1, Price: &total},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "REF123", result.ReferralCode)
	stored, ok := repo.orders[result.ID]
	require.True(t, ok)
	require.Equal(t, "REF123", stored.ReferralCode)
}

func TestOrderServiceCreateValidationError(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	_, err := service.Create(context.Background(), OrderInput{
		CustomerID:   0, // Invalid customer ID
		Status:       "pending",
		ReferralCode: "",
	})
	require.Error(t, err)
}

func TestOrderServiceGetSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	order := &domain.Order{
		ID:           1,
		CustomerID:   1,
		Status:       domain.OrderStatusPending,
		ReferralCode: "REF456",
	}
	repo.orders[1] = order

	result, err := service.Get(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, int64(1), result.CustomerID)
	require.Equal(t, "REF456", result.ReferralCode)
}

func TestOrderServiceGetNotFound(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	_, err := service.Get(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrOrderNotFound))
}

func TestOrderServiceUpdateStatusSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	order := &domain.Order{
		ID:         1,
		CustomerID: 1,
		Status:     domain.OrderStatusPending,
	}
	repo.orders[1] = order

	result, err := service.UpdateStatus(context.Background(), 1, "confirmed")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, domain.OrderStatusConfirmed, result.Status)
}

func TestOrderServiceUpdateStatusNotFound(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	_, err := service.UpdateStatus(context.Background(), 999, "confirmed")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrOrderNotFound))
}

func TestOrderServiceUpdateStatusInvalidStatus(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	order := &domain.Order{
		ID:         1,
		CustomerID: 1,
		Status:     domain.OrderStatusPending,
	}
	repo.orders[1] = order

	_, err := service.UpdateStatus(context.Background(), 1, "invalid_status")
	require.Error(t, err)
}

func TestOrderServiceListByCustomerSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	order1 := &domain.Order{
		ID:           1,
		CustomerID:   1,
		Status:       domain.OrderStatusPending,
		ReferralCode: "REF1",
	}
	order2 := &domain.Order{
		ID:           2,
		CustomerID:   1,
		Status:       domain.OrderStatusConfirmed,
		ReferralCode: "",
	}
	order3 := &domain.Order{
		ID:           3,
		CustomerID:   2,
		Status:       domain.OrderStatusPending,
		ReferralCode: "REF3",
	}
	repo.orders[1] = order1
	repo.orders[2] = order2
	repo.orders[3] = order3

	orders, total, err := service.ListByCustomer(context.Background(), CustomerOrderQuery{
		CustomerID: 1,
		Limit:      10,
		Offset:     0,
	})
	require.NoError(t, err)
	require.Equal(t, 2, total)
	require.Len(t, orders, 2)
	require.Equal(t, "REF1", orders[0].ReferralCode)
	require.Equal(t, "", orders[1].ReferralCode)
}

func TestOrderServiceListByCustomerWithStatusFilter(t *testing.T) {
	repo := newMockOrderRepository()
	productService := newMockProductService()
	service := NewService(repo, productService)

	order1 := &domain.Order{
		ID:           1,
		CustomerID:   1,
		Status:       domain.OrderStatusPending,
		ReferralCode: "PENDING_REF",
	}
	order2 := &domain.Order{
		ID:           2,
		CustomerID:   1,
		Status:       domain.OrderStatusConfirmed,
		ReferralCode: "CONFIRMED_REF",
	}
	repo.orders[1] = order1
	repo.orders[2] = order2

	orders, total, err := service.ListByCustomer(context.Background(), CustomerOrderQuery{
		CustomerID: 1,
		Status:     "pending",
		Limit:      10,
		Offset:     0,
	})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, orders, 1)
	require.Equal(t, domain.OrderStatusPending, orders[0].Status)
	require.Equal(t, "PENDING_REF", orders[0].ReferralCode)
}

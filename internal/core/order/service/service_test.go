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
	orders        map[int64]*domain.Order
	nextID        int64
	createFunc    func(context.Context, *domain.Order) error
	updateFunc    func(context.Context, *domain.Order) error
	findByIDFunc  func(context.Context, int64) (*domain.Order, error)
	listByCustomerFunc func(context.Context, CustomerOrderQuery) ([]*domain.Order, int, error)
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

func TestOrderServiceCreateSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

	customerSnapshot, _ := json.Marshal(map[string]interface{}{
		"id":   1,
		"name": "Test Customer",
	})
	shippingSnapshot, _ := json.Marshal(map[string]interface{}{
		"address": "123 Main St",
		"city":    "Addis Ababa",
	})

	total := 100.50
	result, err := service.Create(context.Background(), OrderInput{
		CustomerID:              1,
		Status:                  "pending",
		PaymentStatus:           "unpaid",
		DeliveryStatus:          "pending",
		ConfirmationStatus:      "pending",
		Total:                   &total,
		Currency:                "ETB",
		CustomerSnapshot:        customerSnapshot,
		ShippingAddressSnapshot: shippingSnapshot,
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 2, Price: &total},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotZero(t, result.ID)
	require.Equal(t, domain.OrderStatusPending, result.Status)
	require.Equal(t, int64(1), result.CustomerID)
	require.Len(t, result.Items, 1)
}

func TestOrderServiceCreateValidationError(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

	_, err := service.Create(context.Background(), OrderInput{
		CustomerID: 0, // Invalid customer ID
		Status:     "pending",
	})
	require.Error(t, err)
}

func TestOrderServiceGetSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

	order := &domain.Order{
		ID:         1,
		CustomerID: 1,
		Status:     domain.OrderStatusPending,
	}
	repo.orders[1] = order

	result, err := service.Get(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, int64(1), result.CustomerID)
}

func TestOrderServiceGetNotFound(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

	_, err := service.Get(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrOrderNotFound))
}

func TestOrderServiceUpdateStatusSuccess(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

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
	service := NewService(repo)

	_, err := service.UpdateStatus(context.Background(), 999, "confirmed")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrOrderNotFound))
}

func TestOrderServiceUpdateStatusInvalidStatus(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

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
	service := NewService(repo)

	order1 := &domain.Order{
		ID:         1,
		CustomerID: 1,
		Status:     domain.OrderStatusPending,
	}
	order2 := &domain.Order{
		ID:         2,
		CustomerID: 1,
		Status:     domain.OrderStatusConfirmed,
	}
	order3 := &domain.Order{
		ID:         3,
		CustomerID: 2,
		Status:     domain.OrderStatusPending,
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
}

func TestOrderServiceListByCustomerWithStatusFilter(t *testing.T) {
	repo := newMockOrderRepository()
	service := NewService(repo)

	order1 := &domain.Order{
		ID:         1,
		CustomerID: 1,
		Status:     domain.OrderStatusPending,
	}
	order2 := &domain.Order{
		ID:         2,
		CustomerID: 1,
		Status:     domain.OrderStatusConfirmed,
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
}


package order

import (
	"context"
	"encoding/json"
	"errors"

	"marketplace/internal/core/order/domain"
	"marketplace/pkg/logger"
)

var (
	// ErrOrderNotFound indicates an order lookup failure.
	ErrOrderNotFound = errors.New("order not found")
)

// OrderItemInput captures line item inputs for order creation.
type OrderItemInput struct {
	ProductID int64
	Quantity  int
	Price     *float64
}

// OrderInput captures order fields for creation or updates.
type OrderInput struct {
	CustomerID              int64
	Status                  string
	PaymentStatus           string
	DeliveryStatus          string
	ConfirmationStatus      string
	Total                   *float64
	Currency                string
	CustomerSnapshot        json.RawMessage
	ShippingAddressSnapshot json.RawMessage
	BillingAddressSnapshot  json.RawMessage
	Items                   []OrderItemInput
}

// CustomerOrderQuery captures filters for customer order history.
type CustomerOrderQuery struct {
	CustomerID int64
	Status     string
	Limit      int
	Offset     int
}

// Repository defines order persistence operations.
type Repository interface {
	Create(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error
	FindByID(ctx context.Context, id int64) (*domain.Order, error)
	ListByCustomer(ctx context.Context, query CustomerOrderQuery) ([]*domain.Order, int, error)
}

// Service coordinates order business logic.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// Create registers a new order with line items.
func (s *Service) Create(ctx context.Context, input OrderInput) (*domain.Order, error) {
	items := make([]domain.OrderItem, len(input.Items))
	for i, item := range input.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	metadata := domain.OrderMetadata{
		PaymentStatus:           input.PaymentStatus,
		DeliveryStatus:          input.DeliveryStatus,
		ConfirmationStatus:      input.ConfirmationStatus,
		Total:                   input.Total,
		Currency:                input.Currency,
		CustomerSnapshot:        input.CustomerSnapshot,
		ShippingAddressSnapshot: input.ShippingAddressSnapshot,
		BillingAddressSnapshot:  input.BillingAddressSnapshot,
	}

	order, err := domain.NewOrder(
		input.CustomerID,
		input.Status,
		metadata,
		items,
	)
	if err != nil {
		logger.Warn("Order creation failed: validation error", "customer_id", input.CustomerID, "error", err)
		return nil, err
	}

	if err := s.repository.Create(ctx, order); err != nil {
		logger.Error("Order creation failed: repository error", "customer_id", input.CustomerID, "error", err)
		return nil, err
	}

	logger.Info("Order created successfully", "order_id", order.ID, "customer_id", order.CustomerID)
	return order, nil
}

// Get retrieves an order by ID.
func (s *Service) Get(ctx context.Context, id int64) (*domain.Order, error) {
	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			logger.Debug("Order retrieval failed: not found", "order_id", id)
		} else {
			logger.Error("Order retrieval failed: repository error", "order_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Order retrieved successfully", "order_id", id)
	return order, nil
}

// UpdateStatus updates the order status.
func (s *Service) UpdateStatus(ctx context.Context, id int64, status string) (*domain.Order, error) {
	parsedStatus, err := domain.ParseOrderStatus(status)
	if err != nil {
		return nil, err
	}

	if err := s.repository.UpdateStatus(ctx, id, parsedStatus); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			logger.Debug("Order status update failed: not found", "order_id", id)
		} else {
			logger.Error("Order status update failed: repository error", "order_id", id, "error", err)
		}
		return nil, err
	}

	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info("Order status updated successfully", "order_id", id, "status", parsedStatus)
	return order, nil
}

// ListByCustomer fetches order history for a customer.
func (s *Service) ListByCustomer(ctx context.Context, query CustomerOrderQuery) ([]*domain.Order, int, error) {
	orders, total, err := s.repository.ListByCustomer(ctx, query)
	if err != nil {
		logger.Error("Order history lookup failed", "customer_id", query.CustomerID, "error", err)
		return nil, 0, err
	}
	logger.Debug("Order history retrieved", "customer_id", query.CustomerID, "count", len(orders))
	return orders, total, nil
}

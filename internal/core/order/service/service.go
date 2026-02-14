package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"marketplace/internal/core/order/domain"
	"marketplace/pkg/logger"
)

// ProductService defines the interface for product operations needed by order service.
type ProductService interface {
	Get(ctx context.Context, id int64) (*Product, error)
}

// Product represents a product from the product domain (to avoid circular dependency).
type Product struct {
	ID         int64
	SupplierID int64
}

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
	CustomerSnapshot        json.RawMessage
	Items                   []OrderItemInput
}

// CustomerOrderQuery captures filters for customer order history.
type CustomerOrderQuery struct {
	CustomerID int64  // Required: customer ID to filter orders
	Status     string // Optional: filter by order status
	Limit      int
	Offset     int
}

// SupplierOrderQuery captures filters for supplier order history.
type SupplierOrderQuery struct {
	SupplierID int64
	Status     string
	Limit      int
	Offset     int
}

// Repository defines order persistence operations.
type Repository interface {
	Create(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error
	UpdatePaymentStatus(ctx context.Context, orderID int64, paymentStatus string) error
	FindByID(ctx context.Context, id int64) (*domain.Order, error)
	ListByCustomer(ctx context.Context, query CustomerOrderQuery) ([]*domain.Order, int, error)
	ListBySupplier(ctx context.Context, query SupplierOrderQuery) ([]*domain.Order, int, error)
	BelongsToSupplier(ctx context.Context, orderID int64, supplierID int64) (bool, error)
}

// Service coordinates order business logic.
type Service struct {
	repository    Repository
	productService ProductService
}

func NewService(repository Repository, productService ProductService) *Service {
	return &Service{
		repository:     repository,
		productService: productService,
	}
}

// Create registers a new order with line items.
func (s *Service) Create(ctx context.Context, input OrderInput) (*domain.Order, error) {
	// Validate that items exist
	if len(input.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Validate first item has a product ID
	if input.Items[0].ProductID <= 0 {
		return nil, errors.New("first order item must have a valid product_id")
	}

	items := make([]domain.OrderItem, len(input.Items))
	for i, item := range input.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// Get supplier_id from the first product in the order (required)
	product, err := s.productService.Get(ctx, input.Items[0].ProductID)
	if err != nil {
		logger.Error("Failed to get product for supplier_id lookup", "product_id", input.Items[0].ProductID, "error", err)
		return nil, fmt.Errorf("failed to get product information: %w", err)
	}

	supplierID := product.SupplierID
	if supplierID == 0 {
		logger.Error("Product has no supplier_id", "product_id", input.Items[0].ProductID)
		return nil, errors.New("product does not have a supplier_id")
	}

	logger.Debug("Retrieved supplier_id from product", "product_id", input.Items[0].ProductID, "supplier_id", supplierID)

	metadata := domain.OrderMetadata{
		PaymentStatus:           input.PaymentStatus,
		DeliveryStatus:          input.DeliveryStatus,
		ConfirmationStatus:      input.ConfirmationStatus,
		Total:                   input.Total,
		CustomerSnapshot:        input.CustomerSnapshot,
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

	// Set supplier_id from product
	order.SupplierID = supplierID

	if err := s.repository.Create(ctx, order); err != nil {
		logger.Error("Order creation failed: repository error", "customer_id", input.CustomerID, "error", err)
		return nil, err
	}

	logger.Info("Order created successfully", "order_id", order.ID, "customer_id", order.CustomerID, "supplier_id", order.SupplierID)
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

// UpdatePaymentStatus updates the order payment status.
func (s *Service) UpdatePaymentStatus(ctx context.Context, id int64, paymentStatus string) (*domain.Order, error) {
	if paymentStatus == "" {
		return nil, errors.New("payment_status is required")
	}

	if err := s.repository.UpdatePaymentStatus(ctx, id, paymentStatus); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			logger.Debug("Order payment status update failed: not found", "order_id", id)
		} else {
			logger.Error("Order payment status update failed: repository error", "order_id", id, "error", err)
		}
		return nil, err
	}

	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info("Order payment status updated successfully", "order_id", id, "payment_status", paymentStatus)
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

// ListBySupplier fetches orders filtered by supplier_id.
func (s *Service) ListBySupplier(ctx context.Context, query SupplierOrderQuery) ([]*domain.Order, int, error) {
	orders, total, err := s.repository.ListBySupplier(ctx, query)
	if err != nil {
		logger.Error("Supplier order history lookup failed", "supplier_id", query.SupplierID, "error", err)
		return nil, 0, err
	}
	logger.Debug("Supplier order history retrieved", "supplier_id", query.SupplierID, "count", len(orders))
	return orders, total, nil
}

// BelongsToSupplier checks if an order belongs to a specific supplier by supplier_id.
func (s *Service) BelongsToSupplier(ctx context.Context, orderID int64, supplierID int64) (bool, error) {
	belongs, err := s.repository.BelongsToSupplier(ctx, orderID, supplierID)
	if err != nil {
		logger.Error("Failed to check if order belongs to supplier", "order_id", orderID, "supplier_id", supplierID, "error", err)
		return false, err
	}
	return belongs, nil
}

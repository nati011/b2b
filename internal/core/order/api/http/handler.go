package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/order/domain"
	orderservice "marketplace/internal/core/order/service"
	supplierservice "marketplace/internal/core/supplier/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/logger"
)

// OrderItemRequest represents an order item payload.
type OrderItemRequest struct {
	ID        int64    `json:"id,omitempty"`
	ProductID int64    `json:"product_id,omitempty"`
	Quantity  int      `json:"quantity"`
	Price     *float64 `json:"price,omitempty"`
}

// CreateOrderRequest represents the payload to create an order.
type CreateOrderRequest struct {
	CustomerID         int64              `json:"customer_id"`
	Status            string             `json:"status,omitempty"`
	PaymentStatus     string             `json:"payment_status,omitempty"`
	DeliveryStatus    string             `json:"delivery_status,omitempty"`
	ConfirmationStatus string             `json:"confirmation_status,omitempty"`
	Total             *float64           `json:"total,omitempty"`
	ReferralCode      string             `json:"referral_code,omitempty"`
	CustomerSnapshot  json.RawMessage    `json:"customer_snapshot,omitempty"`
	Items             []OrderItemRequest `json:"items"`
}

// OrderItemResponse represents an order item returned to clients.
type OrderItemResponse struct {
	ProductID int64    `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Price     *float64 `json:"price,omitempty"`
}

// OrderResponse represents an order returned to clients.
type OrderResponse struct {
	ID                 int64               `json:"id"`
	CustomerID         int64               `json:"customer_id"`
	SupplierID         int64               `json:"supplier_id"`
	Status             string              `json:"status"`
	PaymentStatus      string              `json:"payment_status,omitempty"`
	DeliveryStatus     string              `json:"delivery_status,omitempty"`
	ConfirmationStatus string              `json:"confirmation_status,omitempty"`
	Total              *float64            `json:"total,omitempty"`
	ReferralCode       string              `json:"referral_code,omitempty"`
	CustomerSnapshot   json.RawMessage     `json:"customer_snapshot,omitempty"`
	CartSnapshot       json.RawMessage     `json:"cart_snapshot"`
	Items              []OrderItemResponse `json:"items,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

// OrderListResponse wraps order history results.
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// OrderHandler exposes HTTP endpoints for managing orders.
type OrderHandler struct {
	service         *orderservice.Service
	supplierService *supplierservice.SupplierService
}

func NewOrderHandler(service *orderservice.Service, supplierService *supplierservice.SupplierService) *OrderHandler {
	return &OrderHandler{
		service:         service,
		supplierService: supplierService,
	}
}

// CreateOrder handles POST /order requests.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	items, err := toOrderItemInputs(req.Items)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	input := orderservice.OrderInput{
		CustomerID:         req.CustomerID,
		Status:             req.Status,
		PaymentStatus:      req.PaymentStatus,
		DeliveryStatus:     req.DeliveryStatus,
		ConfirmationStatus: req.ConfirmationStatus,
		Total:              req.Total,
		ReferralCode:       strings.TrimSpace(req.ReferralCode),
		CustomerSnapshot:   req.CustomerSnapshot,
		Items:              items,
	}

	order, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToOrderResponse(order, true))
}

// GetOrder handles GET /order?id= requests.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id, err := parseIDQuery(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	order, err := h.service.Get(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// If user is a supplier, verify the order belongs to their supplier_id
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		belongs, err := h.service.BelongsToSupplier(ctx, id, supplierID)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if !belongs {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: order does not belong to your supplier account"))
			return
		}
		logger.DebugContext(ctx, "Supplier accessing order", "supplier_id", supplierID, "order_id", order.ID)
	}

	httputil.JSON(w, http.StatusOK, ToOrderResponse(order, true))
}

// UpdateOrderStatus handles PATCH /order?id=&command= requests.
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPatch) {
		return
	}

	id, err := parseIDQuery(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	command := strings.TrimSpace(r.URL.Query().Get("command"))
	if command == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("command is required"))
		return
	}

	ctx := r.Context()
	// If user is a supplier, verify the order belongs to their supplier_id before updating
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		belongs, err := h.service.BelongsToSupplier(ctx, id, supplierID)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if !belongs {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: order does not belong to your supplier account"))
			return
		}
	}

	order, err := h.service.UpdateStatus(ctx, id, command)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToOrderResponse(order, true))
}

// UpdatePaymentStatus handles PATCH /order?id=&payment_status= requests.
func (h *OrderHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPatch) {
		return
	}

	id, err := parseIDQuery(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	paymentStatus := strings.TrimSpace(r.URL.Query().Get("payment_status"))
	if paymentStatus == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("payment_status is required"))
		return
	}

	ctx := r.Context()
	// If user is a supplier, verify the order belongs to their supplier_id before updating
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		belongs, err := h.service.BelongsToSupplier(ctx, id, supplierID)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if !belongs {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: order does not belong to your supplier account"))
			return
		}
	}

	order, err := h.service.UpdatePaymentStatus(ctx, id, paymentStatus)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToOrderResponse(order, true))
}

// ListCustomerOrders handles GET /orders/customer requests.
// Requires customer_id parameter. Supports optional status filter.
func (h *OrderHandler) ListCustomerOrders(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx := r.Context()
	customerID := parseOptionalInt64(r.URL.Query().Get("customer_id"))
	if customerID == 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("customer_id is required"))
		return
	}

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseOptionalInt(r.URL.Query().Get("limit"), 10)
	offset := parseOptionalInt(r.URL.Query().Get("offset"), 0)

	query := orderservice.CustomerOrderQuery{
		CustomerID: customerID,
		Status:     status,
		Limit:      limit,
		Offset:     offset,
	}

	orders, total, err := h.service.ListByCustomer(ctx, query)
	if err != nil {
		h.writeError(w, err)
		return
	}

	items := make([]OrderResponse, len(orders))
	for i, order := range orders {
		items[i] = ToOrderResponse(order, true)
	}

	httputil.JSON(w, http.StatusOK, OrderListResponse{
		Orders: items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// ListSupplierOrders handles GET /order/supplier requests.
// Automatically uses the authenticated user's supplier_id to filter orders
func (h *OrderHandler) ListSupplierOrders(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx := r.Context()

	// Get supplier_id from authenticated user
	supplierID := h.getSupplierIDFromUser(ctx)
	if supplierID == 0 {
		httputil.Error(w, http.StatusForbidden, errors.New("access denied: user is not associated with a supplier account"))
		return
	}

	logger.DebugContext(ctx, "ListSupplierOrders: filtering by supplier_id", "supplier_id", supplierID)

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit := parseOptionalInt(r.URL.Query().Get("limit"), 10)
	offset := parseOptionalInt(r.URL.Query().Get("offset"), 0)

	query := orderservice.SupplierOrderQuery{
		SupplierID: supplierID,
		Status:     status,
		Limit:      limit,
		Offset:     offset,
	}

	orders, total, err := h.service.ListBySupplier(ctx, query)
	if err != nil {
		logger.ErrorContext(ctx, "ListSupplierOrders: error fetching orders", "supplier_id", supplierID, "error", err)
		h.writeError(w, err)
		return
	}

	logger.DebugContext(ctx, "ListSupplierOrders: orders retrieved", "supplier_id", supplierID, "count", len(orders), "total", total)
	if err != nil {
		h.writeError(w, err)
		return
	}

	items := make([]OrderResponse, len(orders))
	for i, order := range orders {
		items[i] = ToOrderResponse(order, true)
	}

	httputil.JSON(w, http.StatusOK, OrderListResponse{
		Orders: items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// ToOrderResponse converts a domain order into a DTO.
func ToOrderResponse(order *domain.Order, includeItems bool) OrderResponse {
	// Ensure cart_snapshot is always included, even if empty
	cartSnapshot := order.CartSnapshot
	if len(cartSnapshot) == 0 {
		// Set to empty JSON array if cart_snapshot is empty
		cartSnapshot = []byte("[]")
	}

	resp := OrderResponse{
		ID:                 order.ID,
		CustomerID:         order.CustomerID,
		SupplierID:         order.SupplierID,
		Status:             string(order.Status),
		PaymentStatus:      order.PaymentStatus,
		DeliveryStatus:     order.DeliveryStatus,
		ConfirmationStatus: order.ConfirmationStatus,
		Total:              order.Total,
		ReferralCode:       order.ReferralCode,
		CustomerSnapshot:   order.CustomerSnapshot,
		CartSnapshot:       cartSnapshot,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}

	if includeItems && len(order.Items) > 0 {
		resp.Items = make([]OrderItemResponse, len(order.Items))
		for i, item := range order.Items {
			resp.Items[i] = OrderItemResponse{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}
	}

	return resp
}

func (h *OrderHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, orderservice.ErrOrderNotFound):
		httputil.Error(w, http.StatusNotFound, err)
	default:
		httputil.Error(w, http.StatusBadRequest, err)
	}
}

func toOrderItemInputs(items []OrderItemRequest) ([]orderservice.OrderItemInput, error) {
	result := make([]orderservice.OrderItemInput, len(items))
	for i, item := range items {
		productID := item.ProductID
		if productID == 0 {
			productID = item.ID
		}
		if productID == 0 {
			return nil, errors.New("order item product_id is required")
		}
		result[i] = orderservice.OrderItemInput{
			ProductID: productID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	return result, nil
}

func parseIDQuery(r *http.Request, key string) (int64, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, errors.New("id is required")
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func parseOptionalInt64(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func parseOptionalInt(value string, fallback int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// getSupplierIDFromUser attempts to get supplier_id from the authenticated user's email.
// Returns 0 if user is not found, not authenticated, or is not linked to a supplier.
func (h *OrderHandler) getSupplierIDFromUser(ctx context.Context) int64 {
	user := httputil.UserFromContext(ctx)
	if user == nil {
		return 0
	}

	// Get user email
	email := user.Email.String()
	if email == "" {
		return 0
	}

	// Look up supplier by email
	supplier, err := h.supplierService.GetByEmail(ctx, email)
	if err != nil {
		return 0
	}

	return supplier.ID
}

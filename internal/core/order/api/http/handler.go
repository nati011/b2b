package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/order/domain"
	orderservice "marketplace/internal/core/order/service"
	httputil "marketplace/pkg/http"
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
	CustomerID              int64              `json:"customer_id"`
	Status                  string             `json:"status,omitempty"`
	PaymentStatus           string             `json:"payment_status,omitempty"`
	DeliveryStatus          string             `json:"delivery_status,omitempty"`
	ConfirmationStatus      string             `json:"confirmation_status,omitempty"`
	Total                   *float64           `json:"total,omitempty"`
	Currency                string             `json:"currency,omitempty"`
	CustomerSnapshot        json.RawMessage    `json:"customer_snapshot,omitempty"`
	ShippingAddressSnapshot json.RawMessage    `json:"shipping_address_snapshot,omitempty"`
	BillingAddressSnapshot  json.RawMessage    `json:"billing_address_snapshot,omitempty"`
	Items                   []OrderItemRequest `json:"items"`
}

// OrderItemResponse represents an order item returned to clients.
type OrderItemResponse struct {
	ProductID int64    `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Price     *float64 `json:"price,omitempty"`
}

// OrderResponse represents an order returned to clients.
type OrderResponse struct {
	ID                      int64               `json:"id"`
	CustomerID              int64               `json:"customer_id"`
	Status                  string              `json:"status"`
	PaymentStatus           string              `json:"payment_status,omitempty"`
	DeliveryStatus          string              `json:"delivery_status,omitempty"`
	ConfirmationStatus      string              `json:"confirmation_status,omitempty"`
	Total                   *float64            `json:"total,omitempty"`
	Currency                string              `json:"currency,omitempty"`
	CustomerSnapshot        json.RawMessage     `json:"customer_snapshot,omitempty"`
	ShippingAddressSnapshot json.RawMessage     `json:"shipping_address_snapshot,omitempty"`
	BillingAddressSnapshot  json.RawMessage     `json:"billing_address_snapshot,omitempty"`
	Items                   []OrderItemResponse `json:"items,omitempty"`
	CreatedAt               time.Time           `json:"created_at"`
	UpdatedAt               time.Time           `json:"updated_at"`
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
	service *orderservice.Service
}

func NewOrderHandler(service *orderservice.Service) *OrderHandler {
	return &OrderHandler{service: service}
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
		CustomerID:              req.CustomerID,
		Status:                  req.Status,
		PaymentStatus:           req.PaymentStatus,
		DeliveryStatus:          req.DeliveryStatus,
		ConfirmationStatus:      req.ConfirmationStatus,
		Total:                   req.Total,
		Currency:                req.Currency,
		CustomerSnapshot:        req.CustomerSnapshot,
		ShippingAddressSnapshot: req.ShippingAddressSnapshot,
		BillingAddressSnapshot:  req.BillingAddressSnapshot,
		Items:                   items,
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

	order, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
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

	order, err := h.service.UpdateStatus(r.Context(), id, command)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToOrderResponse(order, true))
}

// ListCustomerOrders handles GET /orders/customer requests.
func (h *OrderHandler) ListCustomerOrders(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	customerID := parseOptionalInt64(r.URL.Query().Get("customer_id"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))

	limit := parseOptionalInt(r.URL.Query().Get("limit"), 10)
	offset := parseOptionalInt(r.URL.Query().Get("offset"), 0)

	query := orderservice.CustomerOrderQuery{
		CustomerID: customerID,
		Status:     status,
		Limit:      limit,
		Offset:     offset,
	}

	orders, total, err := h.service.ListByCustomer(r.Context(), query)
	if err != nil {
		h.writeError(w, err)
		return
	}

	items := make([]OrderResponse, len(orders))
	for i, order := range orders {
		items[i] = ToOrderResponse(order, false)
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
	resp := OrderResponse{
		ID:                      order.ID,
		CustomerID:              order.CustomerID,
		Status:                  string(order.Status),
		PaymentStatus:           order.PaymentStatus,
		DeliveryStatus:          order.DeliveryStatus,
		ConfirmationStatus:      order.ConfirmationStatus,
		Total:                   order.Total,
		Currency:                order.Currency,
		CustomerSnapshot:        order.CustomerSnapshot,
		ShippingAddressSnapshot: order.ShippingAddressSnapshot,
		BillingAddressSnapshot:  order.BillingAddressSnapshot,
		CreatedAt:               order.CreatedAt,
		UpdatedAt:               order.UpdatedAt,
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

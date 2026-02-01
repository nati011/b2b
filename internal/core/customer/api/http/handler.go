package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/customer/domain"
	customerservice "marketplace/internal/core/customer/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/pagination"
)

// CreateCustomerRequest represents the payload to create a customer.
type CreateCustomerRequest struct {
	FullName    string `json:"full_name"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	City        string `json:"city,omitempty"`
	Region      string `json:"region,omitempty"`
	Woreda      string `json:"woreda,omitempty"`
	Status      string `json:"status,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
	Password    string `json:"password,omitempty"`
}

// UpdateCustomerRequest represents the payload to update a customer.
type UpdateCustomerRequest struct {
	FullName    string `json:"full_name"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	City        string `json:"city,omitempty"`
	Region      string `json:"region,omitempty"`
	Woreda      string `json:"woreda,omitempty"`
	Status      string `json:"status,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// CustomerResponse represents a customer returned to clients.
type CustomerResponse struct {
	ID          int64     `json:"id"`
	FullName    string    `json:"full_name"`
	Status      string    `json:"status"`
	City        string    `json:"city,omitempty"`
	Region      string    `json:"region,omitempty"`
	Woreda      string    `json:"woreda,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Email       string    `json:"email,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CustomerListResponse is a paginated collection of customers.
type CustomerListResponse = pagination.PageResult[CustomerResponse]

// CreateCustomerResponse wraps customer creation responses.
type CreateCustomerResponse struct {
	Message  string           `json:"message"`
	Customer CustomerResponse `json:"customer"`
}

// CustomerHandler exposes HTTP endpoints for managing customers.
type CustomerHandler struct {
	service *customerservice.CustomerService
}

func NewCustomerHandler(service *customerservice.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// CreateCustomer handles POST /customer requests.
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	fullName := resolveFullName(req.FullName, req.FirstName, req.LastName)
	if fullName == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("full_name is required"))
		return
	}

	phoneNumber := resolvePhone(req.PhoneNumber, req.Phone)
	status := resolveStatus(req.Status, req.IsActive, string(domain.CustomerStatusInactive))

	ctx := r.Context()
	customer, err := h.service.Create(ctx, customerservice.CustomerInput{
		FullName:    fullName,
		Status:      status,
		City:        req.City,
		Region:      req.Region,
		Woreda:      req.Woreda,
		PhoneNumber: phoneNumber,
		Email:       req.Email,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, CreateCustomerResponse{
		Message:  "customer created successfully",
		Customer: ToCustomerResponse(customer),
	})
}

// GetCustomer handles GET /customer/:id.
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id, err := h.extractID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	customer, err := h.service.Get(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToCustomerResponse(customer))
}

// UpdateCustomer handles PUT /customer/:id.
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id, err := h.extractID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	var req UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	fullName := resolveFullName(req.FullName, req.FirstName, req.LastName)
	if fullName == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("full_name is required"))
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" && req.IsActive != nil {
		if *req.IsActive {
			status = string(domain.CustomerStatusActive)
		} else {
			status = string(domain.CustomerStatusInactive)
		}
	}
	if status == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("status is required"))
		return
	}

	phoneNumber := resolvePhone(req.PhoneNumber, req.Phone)

	ctx := r.Context()
	customer, err := h.service.Update(ctx, id, customerservice.CustomerInput{
		FullName:    fullName,
		Status:      status,
		City:        req.City,
		Region:      req.Region,
		Woreda:      req.Woreda,
		PhoneNumber: phoneNumber,
		Email:       req.Email,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToCustomerResponse(customer))
}

// DeleteCustomer handles DELETE /customer/:id.
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := h.extractID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCustomers handles GET /customer.
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	// Check if email query parameter is provided
	email := strings.TrimSpace(r.URL.Query().Get("email"))
	if email != "" {
		// If email is provided, return customer by email
		ctx := r.Context()
		customer, err := h.service.GetByEmail(ctx, email)
		if err != nil {
			h.writeError(w, err)
			return
		}
		httputil.JSON(w, http.StatusOK, ToCustomerResponse(customer))
		return
	}

	pageReq := pagination.FromRequest(r)
	ctx := r.Context()
	result, err := h.service.List(ctx, pageReq)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToCustomerListResponse(result))
}

// ToCustomerResponse converts a domain customer into a DTO.
func ToCustomerResponse(customer *domain.Customer) CustomerResponse {
	return CustomerResponse{
		ID:          customer.ID,
		FullName:    customer.FullName,
		Status:      string(customer.Status),
		City:        customer.City,
		Region:      customer.Region,
		Woreda:      customer.Woreda,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
		IsActive:    customer.IsActive,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}
}

// ToCustomerListResponse converts a paginated domain result into a DTO.
func ToCustomerListResponse(result pagination.PageResult[*domain.Customer]) CustomerListResponse {
	items := make([]CustomerResponse, len(result.Items))
	for i, customer := range result.Items {
		items[i] = ToCustomerResponse(customer)
	}
	return pagination.PageResult[CustomerResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}
}

func (h *CustomerHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, customerservice.ErrCustomerEmailExists),
		errors.Is(err, customerservice.ErrCustomerPhoneExists):
		httputil.Error(w, http.StatusConflict, err)
	case errors.Is(err, customerservice.ErrCustomerNotFound):
		httputil.Error(w, http.StatusNotFound, err)
	default:
		httputil.Error(w, http.StatusBadRequest, err)
	}
}

func (h *CustomerHandler) extractID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == ResourceCustomers {
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, errors.New("invalid customer id")
		}
		return id, nil
	}
	return 0, errors.New("customer identifier is required")
}

func resolveFullName(fullName, firstName, lastName string) string {
	if strings.TrimSpace(fullName) != "" {
		return strings.TrimSpace(fullName)
	}
	return strings.TrimSpace(strings.TrimSpace(firstName + " " + lastName))
}

func resolvePhone(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return strings.TrimSpace(primary)
	}
	return strings.TrimSpace(fallback)
}

func resolveStatus(status string, isActive *bool, defaultStatus string) string {
	status = strings.TrimSpace(status)
	if status != "" {
		return status
	}
	if isActive != nil {
		if *isActive {
			return string(domain.CustomerStatusActive)
		}
		return string(domain.CustomerStatusInactive)
	}
	return defaultStatus
}

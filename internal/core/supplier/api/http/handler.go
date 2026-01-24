package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/supplier/domain"
	supplierservice "marketplace/internal/core/supplier/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/pagination"
)

// CreateSupplierRequest represents the payload to create a supplier.
type CreateSupplierRequest struct {
	BusinessName string `json:"business_name"`
	Status      string `json:"status,omitempty"`
	SupportEmail string `json:"support_email,omitempty"`
	SupportPhone string `json:"support_phone,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// UpdateSupplierRequest represents the payload to update a supplier.
type UpdateSupplierRequest struct {
	BusinessName string `json:"business_name"`
	Status       string `json:"status,omitempty"`
	SupportEmail string `json:"support_email,omitempty"`
	SupportPhone string `json:"support_phone,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

// SupplierResponse represents a supplier returned to clients.
type SupplierResponse struct {
	ID           int64     `json:"id"`
	BusinessName string    `json:"business_name"`
	Status       string    `json:"status"`
	SupportEmail string    `json:"support_email,omitempty"`
	SupportPhone string    `json:"support_phone,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SupplierListResponse is a paginated collection of suppliers.
type SupplierListResponse = pagination.PageResult[SupplierResponse]

// CreateSupplierResponse wraps supplier creation responses.
type CreateSupplierResponse struct {
	Message   string           `json:"message"`
	Supplier  SupplierResponse `json:"supplier"`
}

// SupplierHandler exposes HTTP endpoints for managing suppliers.
type SupplierHandler struct {
	service *supplierservice.SupplierService
}

func NewSupplierHandler(service *supplierservice.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

// CreateSupplier handles POST /supplier requests.
func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(req.BusinessName) == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("business_name is required"))
		return
	}

	status := resolveStatus(req.Status, req.IsActive, string(domain.SupplierStatusInactive))

	ctx := r.Context()
	supplier, err := h.service.Create(ctx, supplierservice.SupplierInput{
		BusinessName: req.BusinessName,
		Status:       status,
		SupportEmail: req.SupportEmail,
		SupportPhone: req.SupportPhone,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, CreateSupplierResponse{
		Message:  "supplier created successfully",
		Supplier: ToSupplierResponse(supplier),
	})
}

// GetSupplier handles GET /supplier/:id.
func (h *SupplierHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id, err := h.extractID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	supplier, err := h.service.Get(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToSupplierResponse(supplier))
}

// UpdateSupplier handles PUT /supplier/:id.
func (h *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id, err := h.extractID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	var req UpdateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(req.BusinessName) == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("business_name is required"))
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" && req.IsActive != nil {
		if *req.IsActive {
			status = string(domain.SupplierStatusActive)
		} else {
			status = string(domain.SupplierStatusInactive)
		}
	}
	if status == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("status is required"))
		return
	}

	ctx := r.Context()
	supplier, err := h.service.Update(ctx, id, supplierservice.SupplierInput{
		BusinessName: req.BusinessName,
		Status:       status,
		SupportEmail: req.SupportEmail,
		SupportPhone: req.SupportPhone,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToSupplierResponse(supplier))
}

// DeleteSupplier handles DELETE /supplier/:id.
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
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

// ListSuppliers handles GET /supplier.
func (h *SupplierHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	pageReq := pagination.FromRequest(r)
	ctx := r.Context()
	result, err := h.service.List(ctx, pageReq)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToSupplierListResponse(result))
}

// ToSupplierResponse converts a domain supplier into a DTO.
func ToSupplierResponse(supplier *domain.Supplier) SupplierResponse {
	return SupplierResponse{
		ID:           supplier.ID,
		BusinessName: supplier.BusinessName,
		Status:       string(supplier.Status),
		SupportEmail: supplier.SupportEmail,
		SupportPhone: supplier.SupportPhone,
		IsActive:     supplier.IsActive,
		CreatedAt:    supplier.CreatedAt,
		UpdatedAt:    supplier.UpdatedAt,
	}
}

// ToSupplierListResponse converts a paginated domain result into a DTO.
func ToSupplierListResponse(result pagination.PageResult[*domain.Supplier]) SupplierListResponse {
	items := make([]SupplierResponse, len(result.Items))
	for i, supplier := range result.Items {
		items[i] = ToSupplierResponse(supplier)
	}
	return pagination.PageResult[SupplierResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}
}

func (h *SupplierHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, supplierservice.ErrSupplierEmailExists),
		errors.Is(err, supplierservice.ErrSupplierPhoneExists):
		httputil.Error(w, http.StatusConflict, err)
	case errors.Is(err, supplierservice.ErrSupplierNotFound):
		httputil.Error(w, http.StatusNotFound, err)
	default:
		httputil.Error(w, http.StatusBadRequest, err)
	}
}

func (h *SupplierHandler) extractID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == ResourceSuppliers {
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, errors.New("invalid supplier id")
		}
		return id, nil
	}
	return 0, errors.New("supplier identifier is required")
}

func resolveStatus(status string, isActive *bool, defaultStatus string) string {
	status = strings.TrimSpace(status)
	if status != "" {
		return status
	}
	if isActive != nil {
		if *isActive {
			return string(domain.SupplierStatusActive)
		}
		return string(domain.SupplierStatusInactive)
	}
	return defaultStatus
}


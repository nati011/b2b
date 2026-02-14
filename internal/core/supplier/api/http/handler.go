package http

import (
	"context"
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
	
	// If user is a supplier, verify they can only access their own supplier record
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		if id != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: supplier record does not belong to your account"))
			return
		}
	}

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

	ctx := r.Context()
	
	// If user is a supplier, verify they can only update their own supplier record
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		if id != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: supplier record does not belong to your account"))
			return
		}
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

	ctx := r.Context()
	
	// If user is a supplier, only return their own supplier record
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		supplier, err := h.service.Get(ctx, supplierID)
		if err != nil {
			h.writeError(w, err)
			return
		}
		
		// Return as a single-item paginated result
		result := pagination.PageResult[*domain.Supplier]{
			Items:      []*domain.Supplier{supplier},
			Total:      1,
			Page:       1,
			Limit:      1,
			TotalPages: 1,
			HasNext:    false,
			HasPrev:    false,
		}
		httputil.JSON(w, http.StatusOK, ToSupplierListResponse(result))
		return
	}

	// Otherwise, return all suppliers (admin view)
	pageReq := pagination.FromRequest(r)
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	result, err := h.service.List(ctx, pageReq, search)
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

// getSupplierIDFromUser attempts to get supplier_id from the authenticated user's email.
// Returns 0 if user is not found, not authenticated, or is not linked to a supplier.
func (h *SupplierHandler) getSupplierIDFromUser(ctx context.Context) int64 {
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
	supplier, err := h.service.GetByEmail(ctx, email)
	if err != nil {
		return 0
	}

	return supplier.ID
}

// BankAccountRequest represents the payload for bank account operations.
type BankAccountRequest struct {
	SupplierID        int64  `json:"supplier_id"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
	BranchName        string `json:"branch_name,omitempty"`
	AccountType       string `json:"account_type,omitempty"`
	IsPrimary         bool   `json:"is_primary"`
}

// BankAccountResponse represents a bank account returned to clients.
type BankAccountResponse struct {
	ID                int64     `json:"id"`
	SupplierID        int64     `json:"supplier_id"`
	BankName          string    `json:"bank_name"`
	AccountNumber     string    `json:"account_number"`
	AccountHolderName string    `json:"account_holder_name"`
	BranchName        string    `json:"branch_name,omitempty"`
	AccountType       string    `json:"account_type"`
	IsPrimary         bool      `json:"is_primary"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateBankAccount handles POST /supplier/bank-account requests.
func (h *SupplierHandler) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req BankAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	
	// If user is a supplier, ensure they can only create accounts for themselves
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		if req.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: can only create bank accounts for your own supplier account"))
			return
		}
	}

	account, err := h.service.CreateBankAccount(ctx, supplierservice.BankAccountInput{
		SupplierID:        req.SupplierID,
		BankName:          req.BankName,
		AccountNumber:     req.AccountNumber,
		AccountHolderName: req.AccountHolderName,
		BranchName:        req.BranchName,
		AccountType:       req.AccountType,
		IsPrimary:         req.IsPrimary,
	})
	if err != nil {
		h.writeBankAccountError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToBankAccountResponse(account))
}

// ListBankAccounts handles GET /supplier/bank-account?supplier_id= requests.
func (h *SupplierHandler) ListBankAccounts(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	supplierIDStr := r.URL.Query().Get("supplier_id")
	if supplierIDStr == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("supplier_id is required"))
		return
	}

	supplierID, err := strconv.ParseInt(supplierIDStr, 10, 64)
	if err != nil || supplierID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("invalid supplier_id"))
		return
	}

	ctx := r.Context()
	
	// Verify supplier exists
	_, err = h.service.Get(ctx, supplierID)
	if err != nil {
		if errors.Is(err, supplierservice.ErrSupplierNotFound) {
			httputil.Error(w, http.StatusNotFound, errors.New("supplier not found"))
			return
		}
		h.writeError(w, err)
		return
	}
	
	// If user is a supplier, ensure they can only view their own accounts
	if userSupplierID := h.getSupplierIDFromUser(ctx); userSupplierID > 0 {
		if supplierID != userSupplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: can only view bank accounts for your own supplier account"))
			return
		}
	}

	accounts, err := h.service.ListBankAccounts(ctx, supplierID)
	if err != nil {
		h.writeBankAccountError(w, err)
		return
	}

	responses := make([]BankAccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = ToBankAccountResponse(account)
	}

	httputil.JSON(w, http.StatusOK, responses)
}

// UpdateBankAccount handles PUT /supplier/bank-account/:id requests.
func (h *SupplierHandler) UpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id, err := h.extractBankAccountID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	
	// Get existing account to check supplier_id
	existingAccount, err := h.service.GetBankAccount(ctx, id)
	if err != nil {
		h.writeBankAccountError(w, err)
		return
	}

	// If user is a supplier, ensure they can only update their own accounts
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		if existingAccount.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: can only update bank accounts for your own supplier account"))
			return
		}
	}

	var req BankAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	// Ensure supplier_id matches
	if req.SupplierID != existingAccount.SupplierID {
		httputil.Error(w, http.StatusBadRequest, errors.New("supplier_id cannot be changed"))
		return
	}

	account, err := h.service.UpdateBankAccount(ctx, id, supplierservice.BankAccountInput{
		SupplierID:        req.SupplierID,
		BankName:          req.BankName,
		AccountNumber:     req.AccountNumber,
		AccountHolderName: req.AccountHolderName,
		BranchName:        req.BranchName,
		AccountType:       req.AccountType,
		IsPrimary:         req.IsPrimary,
	})
	if err != nil {
		h.writeBankAccountError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToBankAccountResponse(account))
}

// DeleteBankAccount handles DELETE /supplier/bank-account/:id requests.
func (h *SupplierHandler) DeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := h.extractBankAccountID(r.URL.Path)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	
	// Get existing account to check supplier_id
	existingAccount, err := h.service.GetBankAccount(ctx, id)
	if err != nil {
		h.writeBankAccountError(w, err)
		return
	}

	// If user is a supplier, ensure they can only delete their own accounts
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		if existingAccount.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: can only delete bank accounts for your own supplier account"))
			return
		}
	}

	if err := h.service.DeleteBankAccount(ctx, id); err != nil {
		h.writeBankAccountError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToBankAccountResponse converts a domain bank account into a DTO.
func ToBankAccountResponse(account *domain.BankAccount) BankAccountResponse {
	return BankAccountResponse{
		ID:                account.ID,
		SupplierID:        account.SupplierID,
		BankName:          account.BankName,
		AccountNumber:     account.AccountNumber,
		AccountHolderName: account.AccountHolderName,
		BranchName:        account.BranchName,
		AccountType:       string(account.AccountType),
		IsPrimary:         account.IsPrimary,
		IsActive:          account.IsActive,
		CreatedAt:         account.CreatedAt,
		UpdatedAt:         account.UpdatedAt,
	}
}

func (h *SupplierHandler) writeBankAccountError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, supplierservice.ErrBankAccountNotFound):
		httputil.Error(w, http.StatusNotFound, err)
	default:
		httputil.Error(w, http.StatusBadRequest, err)
	}
}

func (h *SupplierHandler) extractBankAccountID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Expected path: /supplier/bank-account/:id
	if len(parts) >= 3 && parts[0] == ResourceSuppliers && parts[1] == "bank-account" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return 0, errors.New("invalid bank account id")
		}
		return id, nil
	}
	return 0, errors.New("bank account identifier is required")
}


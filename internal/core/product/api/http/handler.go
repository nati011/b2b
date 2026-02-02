package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/product/domain"
	productservice "marketplace/internal/core/product/service"
	supplierservice "marketplace/internal/core/supplier/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/logger"
)

// CreateProductRequest represents the payload to create a product.
type CreateProductRequest struct {
	Name          string          `json:"name"`
	Description   string          `json:"description,omitempty"`
	ExternalID    string          `json:"external_id,omitempty"`
	Attributes    json.RawMessage `json:"attributes,omitempty"`
	Unit          string          `json:"unit,omitempty"`
	IsActive      bool            `json:"is_active"`
	SupplierID    *int64          `json:"supplier_id,omitempty"` // Optional: auto-set from context for suppliers
	Price         *float64        `json:"price,omitempty"`
	TotalQuantity int             `json:"total_quantity,omitempty"`
	CategoryIDs   []int64         `json:"category_ids,omitempty"`
}

// UpdateProductRequest represents the payload to update a product.
type UpdateProductRequest struct {
	Name          string          `json:"name"`
	Description   string          `json:"description,omitempty"`
	ExternalID    string          `json:"external_id,omitempty"`
	Attributes    json.RawMessage `json:"attributes,omitempty"`
	Unit          string          `json:"unit,omitempty"`
	IsActive      bool            `json:"is_active"`
	Price         *float64        `json:"price,omitempty"`
	TotalQuantity int             `json:"total_quantity,omitempty"`
	CategoryIDs   []int64         `json:"category_ids,omitempty"`
}

// ProductResponse represents a product returned to clients.
type ProductResponse struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description,omitempty"`
	ExternalID        string          `json:"external_id,omitempty"`
	Attributes        json.RawMessage `json:"attributes,omitempty"`
	Unit              string          `json:"unit,omitempty"`
	IsActive          bool            `json:"is_active"`
	SupplierID        int64           `json:"supplier_id"`
	Price             *float64        `json:"price,omitempty"`
	TotalQuantity     int             `json:"total_quantity"`
	ReservedQuantity  int             `json:"reserved_quantity"`
	AvailableQuantity int             `json:"available_quantity"`
	CategoryIDs       []int64         `json:"category_ids,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// ProductListResponse wraps product list results.
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// ProductHandler exposes HTTP endpoints for managing products.
type ProductHandler struct {
	service        *productservice.Service
	supplierService *supplierservice.SupplierService
}

func NewProductHandler(service *productservice.Service, supplierService *supplierservice.SupplierService) *ProductHandler {
	return &ProductHandler{
		service:         service,
		supplierService: supplierService,
	}
}

// CreateProduct handles POST /product requests.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	var supplierID int64

	// Always try to get supplier_id from the logged-in user's context first
	if id := h.getSupplierIDFromUser(ctx); id > 0 {
		supplierID = id
		logger.DebugContext(ctx, "Using supplier_id from logged-in user context", "supplier_id", supplierID)
		// Ignore any supplier_id provided in request body - use the one from context
	} else {
		// If not found in context, check if provided in request body (for admins)
		if req.SupplierID != nil && *req.SupplierID > 0 {
			supplierID = *req.SupplierID
			logger.DebugContext(ctx, "Using supplier_id from request body", "supplier_id", supplierID)
		} else {
			// supplier_id must be provided - either from context or request body
			httputil.Error(w, http.StatusBadRequest, errors.New("supplier_id must be provided"))
			return
		}
	}

	input := productservice.ProductInput{
		SupplierID:    supplierID,
		Name:          req.Name,
		Description:   req.Description,
		ExternalID:    req.ExternalID,
		Attributes:    req.Attributes,
		Unit:          req.Unit,
		IsActive:      req.IsActive,
		Price:         req.Price,
		TotalQuantity: req.TotalQuantity,
		// ReservedQuantity defaults to 0 for new products
		ReservedQuantity: 0,
		CategoryIDs:       req.CategoryIDs,
	}

	product, err := h.service.Create(ctx, input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToProductResponse(product))
}

// GetProduct handles GET /product?id= requests.
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id, err := parseIDQuery(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	product, err := h.service.Get(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// If user is a supplier, verify they can only access their own products
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		if product.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: product does not belong to your supplier account"))
			return
		}
	}

	httputil.JSON(w, http.StatusOK, ToProductResponse(product))
}

// UpdateProduct handles PUT /product?id= requests.
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id, err := parseIDQuery(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	// If user is a supplier, verify they can only update their own products
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		product, err := h.service.Get(ctx, id)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if product.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: product does not belong to your supplier account"))
			return
		}
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	// Get existing product to preserve reserved_quantity
	existingProduct, err := h.service.Get(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	input := productservice.ProductInput{
		Name:          req.Name,
		Description:   req.Description,
		ExternalID:    req.ExternalID,
		Attributes:    req.Attributes,
		Unit:          req.Unit,
		IsActive:      req.IsActive,
		Price:         req.Price,
		TotalQuantity: req.TotalQuantity,
		// Preserve existing reserved_quantity to avoid changing available_quantity
		ReservedQuantity: existingProduct.ReservedQuantity,
		CategoryIDs:       req.CategoryIDs,
	}

	product, err := h.service.Update(ctx, id, input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToProductResponse(product))
}

// DeleteProduct handles DELETE /product?id= requests.
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := parseIDQuery(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	// If user is a supplier, verify they can only delete their own products
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		product, err := h.service.Get(ctx, id)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if product.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: product does not belong to your supplier account"))
			return
		}
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateGRN handles POST /product/grn requests.
func (h *ProductHandler) CreateGRN(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req struct {
		ProductID int64  `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Notes     string `json:"notes,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.ProductID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("product_id is required"))
		return
	}

	if req.Quantity <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("quantity must be greater than 0"))
		return
	}

	ctx := r.Context()
	// If user is a supplier, verify they can only create GRN for their own products
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		product, err := h.service.Get(ctx, req.ProductID)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if product.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: product does not belong to your supplier account"))
			return
		}
	}

	product, err := h.service.CreateGRN(ctx, req.ProductID, req.Quantity, req.Notes)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToProductResponse(product))
}

// UpdatePriceRequest represents the payload to update product price.
type UpdatePriceRequest struct {
	ProductID int64   `json:"product_id"`
	NewPrice  float64 `json:"new_price"`
	Reason    string  `json:"reason,omitempty"`
}

// UpdatePrice handles PATCH /product/price requests.
func (h *ProductHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPatch) {
		return
	}

	var req UpdatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.NewPrice < 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("price must be non-negative"))
		return
	}

	ctx := r.Context()
	
	// Get user info for audit trail
	var userID *int64
	userEmail := ""
	user := httputil.UserFromContext(ctx)
	if user != nil {
		// User ID is a string, try to parse it as int64
		if user.ID != "" {
			if id, err := strconv.ParseInt(user.ID, 10, 64); err == nil {
				userID = &id
			}
		}
		// Email is a domain.Email type with String() method
		emailStr := user.Email.String()
		if emailStr != "" {
			userEmail = emailStr
		}
	}

	// Verify supplier ownership if user is a supplier
	if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
		product, err := h.service.Get(ctx, req.ProductID)
		if err != nil {
			h.writeError(w, err)
			return
		}
		if product.SupplierID != supplierID {
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: you can only update prices for your own products"))
			return
		}
	}

	product, err := h.service.UpdatePrice(ctx, productservice.PriceUpdateInput{
		ProductID: req.ProductID,
		NewPrice:  req.NewPrice,
		Reason:    req.Reason,
		UserID:    userID,
		UserEmail: userEmail,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToProductResponse(product))
}

// ListProducts handles GET /products requests.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	h.listProducts(w, r)
}

// ListCatalogue handles GET /catalogue requests.
func (h *ProductHandler) ListCatalogue(w http.ResponseWriter, r *http.Request) {
	h.listProducts(w, r)
}

// ListSupplierProducts handles GET /products/supplier requests.
// This is a dedicated endpoint for suppliers to fetch their own products.
// Authentication is required and products are automatically filtered by supplier.
func (h *ProductHandler) ListSupplierProducts(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx := r.Context()
	
	// This endpoint requires authentication - verify user is authenticated
	user := httputil.UserFromContext(ctx)
	if user == nil {
		httputil.Error(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	// Get supplier ID from authenticated user
	supplierID := h.getSupplierIDFromUser(ctx)
	if supplierID == 0 {
		httputil.Error(w, http.StatusForbidden, errors.New("access denied: user is not associated with a supplier"))
		return
	}

	query := productservice.ListQuery{
		SupplierID: supplierID, // Always filter by supplier for this endpoint
		CategoryID: parseOptionalInt64(r.URL.Query().Get("category_id")),
		Limit:      parseOptionalInt(r.URL.Query().Get("limit"), 10),
		Offset:     parseOptionalInt(r.URL.Query().Get("offset"), 0),
	}
	if isActive, ok := parseOptionalBool(r.URL.Query().Get("is_active")); ok {
		query.IsActive = &isActive
	}

	// Note: supplier_id query parameter is ignored for this endpoint - always uses authenticated supplier
	logger.DebugContext(ctx, "Supplier products endpoint - filtering by supplier", "supplier_id", supplierID)

	products, total, err := h.service.List(ctx, query)
	if err != nil {
		h.writeError(w, err)
		return
	}

	items := make([]ProductResponse, len(products))
	for i, product := range products {
		items[i] = ToProductResponse(product)
	}

	httputil.JSON(w, http.StatusOK, ProductListResponse{
		Products: items,
		Total:    total,
		Limit:    query.Limit,
		Offset:   query.Offset,
	})
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx := r.Context()
	user := httputil.UserFromContext(ctx)

	query := productservice.ListQuery{
		SupplierID: parseOptionalInt64(r.URL.Query().Get("supplier_id")),
		CategoryID: parseOptionalInt64(r.URL.Query().Get("category_id")),
		Limit:      parseOptionalInt(r.URL.Query().Get("limit"), 10),
		Offset:     parseOptionalInt(r.URL.Query().Get("offset"), 0),
	}
	if isActive, ok := parseOptionalBool(r.URL.Query().Get("is_active")); ok {
		query.IsActive = &isActive
	}

	// If supplier_id is not provided in query, try to get it from authenticated user
	// This ensures suppliers can only view their own products
	// The auth middleware now authenticates even public routes when credentials are provided
	if query.SupplierID == 0 {
		if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 {
			query.SupplierID = supplierID
			logger.DebugContext(ctx, "Auto-filtering products by supplier", "supplier_id", supplierID)
		} else if user != nil {
			// User is authenticated but not a supplier - return all products (admin behavior)
			logger.DebugContext(ctx, "Authenticated non-supplier user - returning all products")
		} else {
			// No user context - public access without credentials, return all products
			logger.DebugContext(ctx, "Public access - returning all products")
		}
	} else {
		// If supplier_id is explicitly provided, verify the user has permission to view that supplier's products
		// For supplier users, they can only view their own products
		if supplierID := h.getSupplierIDFromUser(ctx); supplierID > 0 && query.SupplierID != supplierID {
			// Supplier user trying to access another supplier's products - deny access
			httputil.Error(w, http.StatusForbidden, errors.New("access denied: you can only view your own products"))
			return
		}
	}

	products, total, err := h.service.List(ctx, query)
	if err != nil {
		h.writeError(w, err)
		return
	}

	items := make([]ProductResponse, len(products))
	for i, product := range products {
		items[i] = ToProductResponse(product)
	}

	httputil.JSON(w, http.StatusOK, ProductListResponse{
		Products: items,
		Total:    total,
		Limit:    query.Limit,
		Offset:   query.Offset,
	})
}

// ToProductResponse converts a domain product into a DTO.
func ToProductResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:                product.ID,
		Name:              product.Name,
		Description:       product.Description,
		ExternalID:        product.ExternalID,
		Attributes:        product.Attributes,
		Unit:              product.Unit,
		IsActive:          product.IsActive,
		SupplierID:        product.SupplierID,
		Price:             product.Price,
		TotalQuantity:     product.TotalQuantity,
		ReservedQuantity:  product.ReservedQuantity,
		AvailableQuantity: product.AvailableQuantity,
		CategoryIDs:       product.CategoryIDs,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

func (h *ProductHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, productservice.ErrProductNotFound):
		httputil.Error(w, http.StatusNotFound, err)
	default:
		httputil.Error(w, http.StatusBadRequest, err)
	}
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

func parseOptionalBool(value string) (bool, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return false, false
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, false
	}
	return parsed, true
}

// getSupplierIDFromUser attempts to get supplier_id from the authenticated user's email.
// Returns 0 if user is not found, not authenticated, or is not linked to a supplier.
func (h *ProductHandler) getSupplierIDFromUser(ctx context.Context) int64 {
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

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/product/domain"
	productservice "marketplace/internal/core/product/service"
	httputil "marketplace/pkg/http"
)

// CreateProductRequest represents the payload to create a product.
type CreateProductRequest struct {
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	ExternalID       string          `json:"external_id,omitempty"`
	Attributes       json.RawMessage `json:"attributes,omitempty"`
	Unit             string          `json:"unit,omitempty"`
	IsActive         bool            `json:"is_active"`
	SupplierID       int64           `json:"supplier_id"`
	Price            *float64        `json:"price,omitempty"`
	TotalQuantity    int             `json:"total_quantity,omitempty"`
	ReservedQuantity int             `json:"reserved_quantity,omitempty"`
	CategoryIDs      []int64         `json:"category_ids,omitempty"`
}

// UpdateProductRequest represents the payload to update a product.
type UpdateProductRequest struct {
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	ExternalID       string          `json:"external_id,omitempty"`
	Attributes       json.RawMessage `json:"attributes,omitempty"`
	Unit             string          `json:"unit,omitempty"`
	IsActive         bool            `json:"is_active"`
	Price            *float64        `json:"price,omitempty"`
	TotalQuantity    int             `json:"total_quantity,omitempty"`
	ReservedQuantity int             `json:"reserved_quantity,omitempty"`
	CategoryIDs      []int64         `json:"category_ids,omitempty"`
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
	service *productservice.Service
}

func NewProductHandler(service *productservice.Service) *ProductHandler {
	return &ProductHandler{service: service}
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

	input := productservice.ProductInput{
		SupplierID:       req.SupplierID,
		Name:             req.Name,
		Description:      req.Description,
		ExternalID:       req.ExternalID,
		Attributes:       req.Attributes,
		Unit:             req.Unit,
		IsActive:         req.IsActive,
		Price:            req.Price,
		TotalQuantity:    req.TotalQuantity,
		ReservedQuantity: req.ReservedQuantity,
		CategoryIDs:      req.CategoryIDs,
	}

	product, err := h.service.Create(r.Context(), input)
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

	product, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
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

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	input := productservice.ProductInput{
		Name:             req.Name,
		Description:      req.Description,
		ExternalID:       req.ExternalID,
		Attributes:       req.Attributes,
		Unit:             req.Unit,
		IsActive:         req.IsActive,
		Price:            req.Price,
		TotalQuantity:    req.TotalQuantity,
		ReservedQuantity: req.ReservedQuantity,
		CategoryIDs:      req.CategoryIDs,
	}

	product, err := h.service.Update(r.Context(), id, input)
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

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListProducts handles GET /products requests.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	h.listProducts(w, r)
}

// ListCatalogue handles GET /catalogue requests.
func (h *ProductHandler) ListCatalogue(w http.ResponseWriter, r *http.Request) {
	h.listProducts(w, r)
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	query := productservice.ListQuery{
		SupplierID: parseOptionalInt64(r.URL.Query().Get("supplier_id")),
		CategoryID: parseOptionalInt64(r.URL.Query().Get("category_id")),
		Limit:      parseOptionalInt(r.URL.Query().Get("limit"), 10),
		Offset:     parseOptionalInt(r.URL.Query().Get("offset"), 0),
	}
	if isActive, ok := parseOptionalBool(r.URL.Query().Get("is_active")); ok {
		query.IsActive = &isActive
	}

	products, total, err := h.service.List(r.Context(), query)
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

package currency

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	httputil "core/pkg/http"
	"core/pkg/logger"

	goodmoney "github.com/nucleus-proj/goodmoney"
)

const (
	errValidationFailed        = "validation failed"
	errValidationFailedMessage = "One or more validation errors occurred"
)

// DTOs

// CreateCurrencyRequest represents a request to create a currency
type CreateCurrencyRequest struct {
	Code string `json:"code"`
}

// CurrencyResponse represents a currency in API responses
type CurrencyResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

// CurrencyListResponse represents a list response (for backward compatibility, but only returns one currency)
type CurrencyListResponse struct {
	Items []CurrencyResponse `json:"items"`
}

// CurrencyDropdownItem represents a currency option for dropdowns
type CurrencyDropdownItem struct {
	Code   string `json:"code"`
	Digits int    `json:"digits"` // MinorUnit from goodmoney
}

// CurrencyDropdownResponse represents a list of currencies for dropdowns
type CurrencyDropdownResponse struct {
	Items []CurrencyDropdownItem `json:"items"`
}

// ToCurrencyResponse converts a Currency to a CurrencyResponse
func ToCurrencyResponse(currency *Currency) CurrencyResponse {
	return CurrencyResponse{
		ID:        currency.ID,
		Code:      currency.Code,
		CreatedAt: currency.CreatedAt,
		UpdatedAt: currency.UpdatedAt,
		CreatedBy: currency.CreatedBy,
		UpdatedBy: currency.UpdatedBy,
	}
}

// CurrencyHandler handles HTTP requests for currency operations
type CurrencyHandler struct {
	db *sql.DB
	// cachedDropdownItems caches the currency dropdown items since they're static
	cachedDropdownItems []CurrencyDropdownItem
	cacheOnce           sync.Once
}

// NewCurrencyHandler creates a new currency handler
func NewCurrencyHandler(db *sql.DB) *CurrencyHandler {
	return &CurrencyHandler{
		db: db,
	}
}

// CreateCurrency handles POST /organization/currencies
// @action name=create desc="Configure currency during registration (can only be done once)"
func (h *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.InfoContext(ctx, "Received create currency request")

	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.WarnContext(ctx, "Failed to decode create currency request", "error", err)
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	// Normalize currency code (trim whitespace and convert to uppercase)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))

	// Validate request
	if req.Code == "" {
		logger.WarnContext(ctx, "Create currency validation failed", "error", "code is required")
		response := ValidationErrorResponse{
			Error:   errValidationFailed,
			Message: errValidationFailedMessage,
			Errors:  []string{"code is required"},
		}
		httputil.JSON(w, http.StatusBadRequest, response)
		return
	}

	// Get created_by from context (set by auth middleware)
	createdBy := h.getUserIDFromContext(ctx)

	currency, err := h.createCurrency(ctx, req.Code, createdBy)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create currency", "error", err)
		h.writeError(w, err)
		return
	}

	logger.InfoContext(ctx, "Currency created successfully", "currency_id", currency.ID, "code", currency.Code)
	httputil.JSON(w, http.StatusCreated, ToCurrencyResponse(currency))
}

// GetCurrency handles GET /organization/currencies
// @action name=view desc="Get the configured currency"
func (h *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	logger.DebugContext(ctx, "Received get currency request")

	currency, err := h.getCurrency(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get currency", "error", err)
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToCurrencyResponse(currency))
}

// ListCurrencies handles GET /organization/currencies (returns single currency for backward compatibility)
// @action name=view desc="Get the configured currency"
func (h *CurrencyHandler) ListCurrencies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	currency, err := h.getCurrency(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get currency", "error", err)
		h.writeError(w, err)
		return
	}

	// Return as list for backward compatibility
	response := CurrencyListResponse{
		Items: []CurrencyResponse{ToCurrencyResponse(currency)},
	}

	httputil.JSON(w, http.StatusOK, response)
}

// GetCurrencyDropdown handles GET /organization/currencies/dropdown
// @action name=view desc="Get currency dropdown options"
func (h *CurrencyHandler) GetCurrencyDropdown(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	logger.DebugContext(ctx, "Received get currency dropdown request")

	currencies, err := h.getSupportedCurrencies(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get supported currencies", "error", err)
		h.writeError(w, err)
		return
	}

	response := CurrencyDropdownResponse{
		Items: currencies,
	}

	httputil.JSON(w, http.StatusOK, response)
}

// ValidationErrorResponse represents a validation error response with multiple errors
type ValidationErrorResponse struct {
	Error   string   `json:"error"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// writeError writes an error response with appropriate HTTP status code
func (h *CurrencyHandler) writeError(w http.ResponseWriter, err error) {
	status := h.mapError(err)
	httputil.Error(w, status, err)
}

// mapError maps service errors to HTTP status codes
func (h *CurrencyHandler) mapError(err error) int {
	// Check service errors
	if errors.Is(err, ErrCurrencyNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrCurrencyCodeConflict) {
		return http.StatusConflict
	}

	// Check domain errors
	var domainErr *CurrencyError
	if errors.As(err, &domainErr) {
		if domainErr.Code == ErrorCodeCurrencyAlreadyConfigured {
			return http.StatusForbidden // 403 - operation not allowed
		}
		if domainErr.Code == ErrorCodeCurrencyNotSupported {
			return http.StatusBadRequest
		}
		if domainErr.Code == ErrorCodeValidationFailed ||
			domainErr.Code == ErrorCodeInvalidEnumValue {
			return http.StatusBadRequest
		}
		// Default for any other domain error
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}

// getUserIDFromContext extracts the user ID from the request context
func (h *CurrencyHandler) getUserIDFromContext(ctx context.Context) string {
	user := httputil.UserFromContext(ctx)
	if user != nil {
		return user.ID
	}
	return "system" // Default fallback
}

// Business logic and database access methods

// createCurrency creates a new currency (only allowed if no currency exists)
func (h *CurrencyHandler) createCurrency(ctx context.Context, code string, createdBy string) (*Currency, error) {
	logger.InfoContext(ctx, "Creating currency", "code", code)

	// Check if currency has already been configured
	existing, err := h.findCurrency(ctx)
	if err != nil && !errors.Is(err, ErrCurrencyNotFound) {
		logger.ErrorContext(ctx, "Failed to check if currency exists", "error", err)
		return nil, fmt.Errorf("failed to check if currency exists: %w", err)
	}
	if existing != nil {
		logger.WarnContext(ctx, "Currency already configured", "existing_code", existing.Code)
		return nil, ErrCurrencyAlreadyConfigured
	}

	// Create the currency using domain factory
	currency, err := NewCurrency(NewCurrencyParams{Code: code})
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create currency domain object", "error", err)
		return nil, fmt.Errorf("failed to create currency: %w", err)
	}

	// Set created by
	currency.CreatedBy = createdBy
	currency.UpdatedBy = createdBy

	// Save to database
	if err := h.insertCurrency(ctx, currency); err != nil {
		logger.ErrorContext(ctx, "Failed to save currency to database", "error", err, "currency_id", currency.ID)
		return nil, fmt.Errorf("failed to save currency: %w", err)
	}

	logger.InfoContext(ctx, "Currency created successfully", "currency_id", currency.ID, "code", currency.Code)
	return currency, nil
}

// getCurrency retrieves the configured currency
func (h *CurrencyHandler) getCurrency(ctx context.Context) (*Currency, error) {
	return h.findCurrency(ctx)
}

// Database access methods

// insertCurrency inserts a new currency into the database
func (h *CurrencyHandler) insertCurrency(ctx context.Context, currency *Currency) error {
	logger.DebugContext(ctx, "Inserting currency in database", "currency_id", currency.ID)

	query := `
		INSERT INTO currencies (
			id, code, singleton, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, 1, $3, $4, $5, $6
		)
	`

	_, err := h.db.ExecContext(ctx, query,
		currency.ID,
		currency.Code,
		currency.CreatedAt,
		currency.UpdatedAt,
		currency.CreatedBy,
		currency.UpdatedBy,
	)

	if err != nil {
		logger.ErrorContext(ctx, "Failed to insert currency", "error", err, "currency_id", currency.ID)
		// Check if error is due to unique constraint violation (currency already exists)
		if strings.Contains(err.Error(), "idx_currency_singleton") || strings.Contains(err.Error(), "duplicate key") {
			return ErrCurrencyAlreadyConfigured
		}
		return err
	}

	logger.DebugContext(ctx, "Currency inserted in database successfully", "currency_id", currency.ID)
	return nil
}

// findCurrency retrieves the single configured currency
func (h *CurrencyHandler) findCurrency(ctx context.Context) (*Currency, error) {
	query := `
		SELECT id, code, singleton, created_at, updated_at, created_by, updated_by
		FROM currencies
		WHERE singleton = 1
		LIMIT 1
	`

	return h.scanCurrencyFromRow(ctx, query)
}

// scanCurrencyFromRow scans a single currency from a database row
func (h *CurrencyHandler) scanCurrencyFromRow(ctx context.Context, query string, args ...interface{}) (*Currency, error) {
	var currency Currency
	var singleton int // Read singleton but don't use it

	err := h.db.QueryRowContext(ctx, query, args...).Scan(
		&currency.ID,
		&currency.Code,
		&singleton,
		&currency.CreatedAt,
		&currency.UpdatedAt,
		&currency.CreatedBy,
		&currency.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCurrencyNotFound
		}
		return nil, err
	}

	return &currency, nil
}

// getSupportedCurrencies returns all currencies supported by goodmoney
// Results are cached since the currency map is static
func (h *CurrencyHandler) getSupportedCurrencies(ctx context.Context) ([]CurrencyDropdownItem, error) {
	h.cacheOnce.Do(func() {
		// Pre-allocate with known capacity for better performance
		items := make([]CurrencyDropdownItem, 0, len(goodmoney.CurrencyMap))

		// Iterate over CurrencyMap from goodmoney to get all supported currencies
		for code, currency := range goodmoney.CurrencyMap {
			item := CurrencyDropdownItem{
				Code:   code,
				Digits: currency.MinorUnit,
			}
			items = append(items, item)
		}

		// Sort by code for consistent ordering
		sort.Slice(items, func(i, j int) bool {
			return items[i].Code < items[j].Code
		})

		h.cachedDropdownItems = items
		logger.DebugContext(ctx, "Cached currency dropdown items", "count", len(h.cachedDropdownItems))
	})

	return h.cachedDropdownItems, nil
}

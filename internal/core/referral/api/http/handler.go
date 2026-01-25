package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/internal/core/referral/domain"
	referralservice "marketplace/internal/core/referral/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/pagination"
)

// CreateAffiliateRequest represents the payload to create an affiliate.
type CreateAffiliateRequest struct {
	AffiliateID      *int64 `json:"affiliate_id,omitempty"` // Optional: if affiliate is also a customer
	Email            string `json:"email,omitempty"`
	PhoneNumber      string `json:"phone_number,omitempty"`
	FullName         string `json:"full_name"`
	Status           string `json:"status,omitempty"`
	SocialMediaHandle string `json:"social_media_handle,omitempty"`
	Platform         string `json:"platform,omitempty"`
}

// UpdateAffiliateRequest represents the payload to update an affiliate.
type UpdateAffiliateRequest struct {
	Email            string `json:"email,omitempty"`
	PhoneNumber      string `json:"phone_number,omitempty"`
	FullName         string `json:"full_name"`
	Status           string `json:"status,omitempty"`
	SocialMediaHandle string `json:"social_media_handle,omitempty"`
	Platform         string `json:"platform,omitempty"`
}

// AffiliateResponse represents an affiliate returned to clients.
type AffiliateResponse struct {
	ID               int64     `json:"id"`
	AffiliateID      *int64    `json:"affiliate_id,omitempty"`
	Email            string    `json:"email,omitempty"`
	PhoneNumber      string    `json:"phone_number,omitempty"`
	FullName         string    `json:"full_name"`
	Status           string    `json:"status"`
	SocialMediaHandle string    `json:"social_media_handle,omitempty"`
	Platform         string    `json:"platform,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AffiliateListResponse is a paginated collection of affiliates.
type AffiliateListResponse = pagination.PageResult[AffiliateResponse]

// CreateAffiliateResponse wraps affiliate creation responses.
type CreateAffiliateResponse struct {
	Message   string            `json:"message"`
	Affiliate AffiliateResponse `json:"affiliate"`
}

// CreateReferralCodeRequest represents the payload to create a referral code.
type CreateReferralCodeRequest struct {
	AffiliateID int64   `json:"affiliate_id"`
	CustomCode  string  `json:"custom_code,omitempty"`
	Status      string  `json:"status,omitempty"`
	ExpiresAt   *string `json:"expires_at,omitempty"` // ISO 8601 date string
}

// ReferralCodeResponse represents a referral code returned to clients.
type ReferralCodeResponse struct {
	ID          int64      `json:"id"`
	AffiliateID int64      `json:"affiliate_id"`
	Code        string     `json:"code"`
	IsCustom    bool       `json:"is_custom"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateReferralRelationshipRequest represents the payload to create a referral relationship.
type CreateReferralRelationshipRequest struct {
	AffiliateID    int64  `json:"affiliate_id"`
	CustomerID     int64  `json:"customer_id"`
	ReferralCodeID int64  `json:"referral_code_id"`
	Source         string `json:"source,omitempty"`
	Campaign       string `json:"campaign,omitempty"`
	UTMSource      string `json:"utm_source,omitempty"`
	UTMMedium      string `json:"utm_medium,omitempty"`
	UTMCampaign    string `json:"utm_campaign,omitempty"`
	UTMContent     string `json:"utm_content,omitempty"`
	ReferralLink   string `json:"referral_link,omitempty"`
	IPAddress      string `json:"ip_address,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
}

// ReferralRelationshipResponse represents a referral relationship returned to clients.
type ReferralRelationshipResponse struct {
	ID              int64     `json:"id"`
	AffiliateID     int64     `json:"affiliate_id"`
	CustomerID      int64     `json:"customer_id"`
	ReferralCodeID  int64     `json:"referral_code_id"`
	Source          string    `json:"source,omitempty"`
	Campaign        string    `json:"campaign,omitempty"`
	UTMSource       string    `json:"utm_source,omitempty"`
	UTMMedium       string    `json:"utm_medium,omitempty"`
	UTMCampaign     string    `json:"utm_campaign,omitempty"`
	UTMContent      string    `json:"utm_content,omitempty"`
	ReferralLink    string    `json:"referral_link,omitempty"`
	IPAddress       string    `json:"ip_address,omitempty"`
	UserAgent       string    `json:"user_agent,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateCommissionRequest represents the payload to create a commission.
type CreateCommissionRequest struct {
	AffiliateID            int64   `json:"affiliate_id"`
	ReferralRelationshipID int64   `json:"referral_relationship_id"`
	OrderID                int64   `json:"order_id"`
	CustomerID             int64   `json:"customer_id"`
	CommissionRate         float64 `json:"commission_rate"`
	OrderTotal             string  `json:"order_total"` // JSON string of money
	Notes                  string  `json:"notes,omitempty"`
}

// CommissionResponse represents a commission returned to clients.
type CommissionResponse struct {
	ID                    int64      `json:"id"`
	AffiliateID           int64      `json:"affiliate_id"`
	ReferralRelationshipID int64      `json:"referral_relationship_id"`
	OrderID               int64      `json:"order_id"`
	CustomerID            int64      `json:"customer_id"`
	Amount                *float64   `json:"amount,omitempty"` // Simplified for JSON
	CommissionRate        float64    `json:"commission_rate"`
	Status                string     `json:"status"`
	OrderTotal            *float64   `json:"order_total,omitempty"` // Simplified for JSON
	Notes                 string     `json:"notes,omitempty"`
	PaidAt                *time.Time `json:"paid_at,omitempty"`
	ReversedAt            *time.Time `json:"reversed_at,omitempty"`
	IsActive              bool       `json:"is_active"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// ValidateReferralCodeRequest represents the payload to validate a referral code.
type ValidateReferralCodeRequest struct {
	Code       string `json:"code"`
	CustomerID int64  `json:"customer_id"`
}

// ReferralHandler exposes HTTP endpoints for managing referrals.
type ReferralHandler struct {
	service *referralservice.ReferralService
}

func NewReferralHandler(service *referralservice.ReferralService) *ReferralHandler {
	return &ReferralHandler{service: service}
}

// CreateAffiliate handles POST /referral/affiliate requests.
func (h *ReferralHandler) CreateAffiliate(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateAffiliateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(req.FullName) == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("full_name is required"))
		return
	}

	status := resolveStatus(req.Status, string(domain.AffiliateStatusInactive))

	ctx := r.Context()
	affiliate, err := h.service.CreateAffiliate(ctx, referralservice.AffiliateInput{
		AffiliateID:      req.AffiliateID,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		FullName:         req.FullName,
		Status:           status,
		SocialMediaHandle: req.SocialMediaHandle,
		Platform:         req.Platform,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, CreateAffiliateResponse{
		Message:   "affiliate created successfully",
		Affiliate: ToAffiliateResponse(affiliate),
	})
}

// GetAffiliate handles GET /referral/affiliate/:id.
func (h *ReferralHandler) GetAffiliate(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id, err := h.extractID(r.URL.Path, "affiliate")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	affiliate, err := h.service.GetAffiliate(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToAffiliateResponse(affiliate))
}

// ListAffiliates handles GET /referral/affiliate.
func (h *ReferralHandler) ListAffiliates(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	pageReq := pagination.FromRequest(r)
	ctx := r.Context()
	result, err := h.service.ListAffiliates(ctx, pageReq)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToAffiliateListResponse(result))
}

// CreateReferralCode handles POST /referral/code requests.
func (h *ReferralHandler) CreateReferralCode(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateReferralCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.AffiliateID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("affiliate_id is required"))
		return
	}

	status := resolveStatus(req.Status, string(domain.ReferralCodeStatusActive))

	ctx := r.Context()
	code, err := h.service.CreateReferralCode(ctx, referralservice.ReferralCodeInput{
		AffiliateID: req.AffiliateID,
		CustomCode:  req.CustomCode,
		Status:      status,
		ExpiresAt:   req.ExpiresAt,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToReferralCodeResponse(code))
}

// GetReferralCode handles GET /referral/code/:id.
func (h *ReferralHandler) GetReferralCode(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id, err := h.extractID(r.URL.Path, "code")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	code, err := h.service.GetReferralCode(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToReferralCodeResponse(code))
}

// ValidateReferralCode handles POST /referral/validate requests.
func (h *ReferralHandler) ValidateReferralCode(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req ValidateReferralCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("code is required"))
		return
	}

	if req.CustomerID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("customer_id is required"))
		return
	}

	ctx := r.Context()
	code, err := h.service.ValidateReferralCode(ctx, req.Code, req.CustomerID, nil) // CustomerChecker would be injected
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToReferralCodeResponse(code))
}

// CreateReferralRelationship handles POST /referral/relationship requests.
func (h *ReferralHandler) CreateReferralRelationship(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateReferralRelationshipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.AffiliateID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("affiliate_id is required"))
		return
	}
	if req.CustomerID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("customer_id is required"))
		return
	}
	if req.ReferralCodeID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("referral_code_id is required"))
		return
	}

	// Extract IP and User Agent from request
	ipAddress := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ipAddress = strings.Split(forwarded, ",")[0]
	}
	userAgent := r.Header.Get("User-Agent")

	if req.IPAddress == "" {
		req.IPAddress = ipAddress
	}
	if req.UserAgent == "" {
		req.UserAgent = userAgent
	}

	ctx := r.Context()
	relationship, err := h.service.CreateReferralRelationship(ctx, referralservice.ReferralRelationshipInput{
		AffiliateID:    req.AffiliateID,
		CustomerID:     req.CustomerID,
		ReferralCodeID: req.ReferralCodeID,
		Source:         req.Source,
		Campaign:       req.Campaign,
		UTMSource:      req.UTMSource,
		UTMMedium:      req.UTMMedium,
		UTMCampaign:    req.UTMCampaign,
		UTMContent:     req.UTMContent,
		ReferralLink:   req.ReferralLink,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToReferralRelationshipResponse(relationship))
}

// CreateCommission handles POST /referral/commission requests.
func (h *ReferralHandler) CreateCommission(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateCommissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.AffiliateID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("affiliate_id is required"))
		return
	}
	if req.ReferralRelationshipID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("referral_relationship_id is required"))
		return
	}
	if req.OrderID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("order_id is required"))
		return
	}
	if req.CustomerID <= 0 {
		httputil.Error(w, http.StatusBadRequest, errors.New("customer_id is required"))
		return
	}
	if req.CommissionRate < 0 || req.CommissionRate > 100 {
		httputil.Error(w, http.StatusBadRequest, errors.New("commission_rate must be between 0 and 100"))
		return
	}

	ctx := r.Context()
	commission, err := h.service.CreateCommission(ctx, referralservice.CommissionInput{
		AffiliateID:            req.AffiliateID,
		ReferralRelationshipID: req.ReferralRelationshipID,
		OrderID:                req.OrderID,
		CustomerID:             req.CustomerID,
		CommissionRate:         req.CommissionRate,
		OrderTotal:             req.OrderTotal,
		Notes:                  req.Notes,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToCommissionResponse(commission))
}

// Response conversion functions

func ToAffiliateResponse(affiliate *domain.Affiliate) AffiliateResponse {
	return AffiliateResponse{
		ID:               affiliate.ID,
		AffiliateID:      affiliate.AffiliateID,
		Email:            affiliate.Email,
		PhoneNumber:      affiliate.PhoneNumber,
		FullName:         affiliate.FullName,
		Status:           string(affiliate.Status),
		SocialMediaHandle: affiliate.SocialMediaHandle,
		Platform:         affiliate.Platform,
		IsActive:         affiliate.IsActive,
		CreatedAt:        affiliate.CreatedAt,
		UpdatedAt:        affiliate.UpdatedAt,
	}
}

func ToAffiliateListResponse(result pagination.PageResult[*domain.Affiliate]) AffiliateListResponse {
	items := make([]AffiliateResponse, len(result.Items))
	for i, affiliate := range result.Items {
		items[i] = ToAffiliateResponse(affiliate)
	}
	return pagination.PageResult[AffiliateResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}
}

func ToReferralCodeResponse(code *domain.ReferralCode) ReferralCodeResponse {
	return ReferralCodeResponse{
		ID:          code.ID,
		AffiliateID: code.AffiliateID,
		Code:        code.Code,
		IsCustom:    code.IsCustom,
		Status:      string(code.Status),
		ExpiresAt:   code.ExpiresAt,
		IsActive:    code.IsActive,
		CreatedAt:   code.CreatedAt,
		UpdatedAt:   code.UpdatedAt,
	}
}

func ToReferralRelationshipResponse(relationship *domain.ReferralRelationship) ReferralRelationshipResponse {
	return ReferralRelationshipResponse{
		ID:              relationship.ID,
		AffiliateID:     relationship.AffiliateID,
		CustomerID:      relationship.CustomerID,
		ReferralCodeID:  relationship.ReferralCodeID,
		Source:          relationship.Source,
		Campaign:        relationship.Campaign,
		UTMSource:       relationship.UTMSource,
		UTMMedium:       relationship.UTMMedium,
		UTMCampaign:     relationship.UTMCampaign,
		UTMContent:      relationship.UTMContent,
		ReferralLink:    relationship.ReferralLink,
		IPAddress:       relationship.IPAddress,
		UserAgent:       relationship.UserAgent,
		IsActive:        relationship.IsActive,
		CreatedAt:       relationship.CreatedAt,
		UpdatedAt:       relationship.UpdatedAt,
	}
}

func ToCommissionResponse(commission *domain.Commission) CommissionResponse {
	var amount *float64
	if commission.Amount != nil {
		val := float64(commission.Amount.Amount()) / 100.0 // Convert from cents
		amount = &val
	}

	var orderTotal *float64
	if commission.OrderTotal != nil {
		val := float64(commission.OrderTotal.Amount()) / 100.0 // Convert from cents
		orderTotal = &val
	}

	return CommissionResponse{
		ID:                    commission.ID,
		AffiliateID:           commission.AffiliateID,
		ReferralRelationshipID: commission.ReferralRelationshipID,
		OrderID:                commission.OrderID,
		CustomerID:            commission.CustomerID,
		Amount:                amount,
		CommissionRate:        commission.CommissionRate,
		Status:                string(commission.Status),
		OrderTotal:            orderTotal,
		Notes:                 commission.Notes,
		PaidAt:                commission.PaidAt,
		ReversedAt:           commission.ReversedAt,
		IsActive:              commission.IsActive,
		CreatedAt:             commission.CreatedAt,
		UpdatedAt:             commission.UpdatedAt,
	}
}

func (h *ReferralHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, referralservice.ErrAffiliateEmailExists),
		errors.Is(err, referralservice.ErrAffiliatePhoneExists),
		errors.Is(err, referralservice.ErrReferralCodeExists):
		httputil.Error(w, http.StatusConflict, err)
	case errors.Is(err, referralservice.ErrAffiliateNotFound),
		errors.Is(err, referralservice.ErrReferralCodeNotFound),
		errors.Is(err, referralservice.ErrReferralRelationshipNotFound),
		errors.Is(err, referralservice.ErrCommissionNotFound):
		httputil.Error(w, http.StatusNotFound, err)
	case errors.Is(err, referralservice.ErrExistingCustomerReferral),
		errors.Is(err, referralservice.ErrSelfReferral):
		httputil.Error(w, http.StatusForbidden, err)
	default:
		httputil.Error(w, http.StatusBadRequest, err)
	}
}

func (h *ReferralHandler) extractID(path string, resource string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Expected format: /referral/{resource}/{id}
	if len(parts) >= 3 && parts[0] == "referral" && parts[1] == resource {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return 0, errors.New("invalid " + resource + " id")
		}
		return id, nil
	}
	return 0, errors.New(resource + " identifier is required")
}

func resolveStatus(status string, defaultStatus string) string {
	status = strings.TrimSpace(status)
	if status != "" {
		return status
	}
	return defaultStatus
}


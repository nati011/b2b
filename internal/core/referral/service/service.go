package referral

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"marketplace/internal/core/referral/domain"
	"marketplace/pkg/logger"
	"marketplace/pkg/pagination"

	goodmoney "github.com/the-nucleus-project/good_money"
)

var (
	// ErrAffiliateNotFound indicates an affiliate lookup failure.
	ErrAffiliateNotFound = errors.New("affiliate not found")
	// ErrAffiliateEmailExists indicates an affiliate with the email already exists.
	ErrAffiliateEmailExists = errors.New("affiliate email already exists")
	// ErrAffiliatePhoneExists indicates an affiliate with the phone number already exists.
	ErrAffiliatePhoneExists = errors.New("affiliate phone number already exists")
	// ErrReferralCodeNotFound indicates a referral code lookup failure.
	ErrReferralCodeNotFound = errors.New("referral code not found")
	// ErrReferralCodeExists indicates a referral code already exists.
	ErrReferralCodeExists = errors.New("referral code already exists")
	// ErrReferralRelationshipNotFound indicates a referral relationship lookup failure.
	ErrReferralRelationshipNotFound = errors.New("referral relationship not found")
	// ErrCommissionNotFound indicates a commission lookup failure.
	ErrCommissionNotFound = errors.New("commission not found")
	// ErrExistingCustomerReferral indicates an existing customer tried to use a referral code.
	ErrExistingCustomerReferral = errors.New("existing customers cannot use referral codes")
	// ErrSelfReferral indicates an affiliate tried to refer themselves.
	ErrSelfReferral = errors.New("affiliates cannot refer themselves")
)

// Repository defines the referral persistence contract.
type Repository interface {
	// Affiliate operations
	CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error
	FindAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error)
	FindAffiliateByEmail(ctx context.Context, email string) (*domain.Affiliate, error)
	FindAffiliateByPhone(ctx context.Context, phone string) (*domain.Affiliate, error)
	FindAffiliateByCustomerID(ctx context.Context, customerID int64) (*domain.Affiliate, error)
	DeleteAffiliate(ctx context.Context, id int64) error
	FindAllAffiliates(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Affiliate], error)

	// Referral code operations
	CreateReferralCode(ctx context.Context, code *domain.ReferralCode) error
	UpdateReferralCode(ctx context.Context, code *domain.ReferralCode) error
	FindReferralCodeByID(ctx context.Context, id int64) (*domain.ReferralCode, error)
	FindReferralCodeByCode(ctx context.Context, code string) (*domain.ReferralCode, error)
	FindReferralCodesByAffiliateID(ctx context.Context, affiliateID int64) ([]*domain.ReferralCode, error)
	DeleteReferralCode(ctx context.Context, id int64) error

	// Referral relationship operations
	CreateReferralRelationship(ctx context.Context, relationship *domain.ReferralRelationship) error
	UpdateReferralRelationship(ctx context.Context, relationship *domain.ReferralRelationship) error
	FindReferralRelationshipByID(ctx context.Context, id int64) (*domain.ReferralRelationship, error)
	FindReferralRelationshipByCustomerID(ctx context.Context, customerID int64) (*domain.ReferralRelationship, error)
	FindReferralRelationshipsByAffiliateID(ctx context.Context, affiliateID int64, pageReq pagination.PageRequest) (pagination.PageResult[*domain.ReferralRelationship], error)
	DeleteReferralRelationship(ctx context.Context, id int64) error

	// Commission operations
	CreateCommission(ctx context.Context, commission *domain.Commission) error
	UpdateCommission(ctx context.Context, commission *domain.Commission) error
	FindCommissionByID(ctx context.Context, id int64) (*domain.Commission, error)
	FindCommissionsByAffiliateID(ctx context.Context, affiliateID int64, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Commission], error)
	FindCommissionsByOrderID(ctx context.Context, orderID int64) ([]*domain.Commission, error)
	FindCommissionsByReferralRelationshipID(ctx context.Context, referralRelationshipID int64) ([]*domain.Commission, error)
	DeleteCommission(ctx context.Context, id int64) error
}

// AffiliateInput captures incoming affiliate profile fields.
type AffiliateInput struct {
	AffiliateID       *int64 // Optional: if affiliate is also a customer
	Email             string
	PhoneNumber       string
	FullName          string
	Status            string
	SocialMediaHandle string
	Platform          string
}

// ReferralCodeInput captures incoming referral code fields.
type ReferralCodeInput struct {
	AffiliateID int64
	CustomCode  string
	Status      string
	ExpiresAt   *string // ISO 8601 date string
}

// ReferralRelationshipInput captures incoming referral relationship fields.
type ReferralRelationshipInput struct {
	AffiliateID    int64
	CustomerID     int64
	ReferralCodeID int64
	Source         string
	Campaign       string
	UTMSource      string
	UTMMedium      string
	UTMCampaign    string
	UTMContent     string
	ReferralLink   string
	IPAddress      string
	UserAgent      string
}

// CommissionInput captures incoming commission fields.
type CommissionInput struct {
	AffiliateID            int64
	ReferralRelationshipID int64
	OrderID                int64
	CustomerID             int64
	CommissionRate         float64
	OrderTotal             string // JSON string of money
	Notes                  string
}

// ReferralService coordinates referral business logic.
type ReferralService struct {
	repository Repository
}

func NewReferralService(repository Repository) *ReferralService {
	return &ReferralService{repository: repository}
}

// CreateAffiliate registers a new affiliate.
func (s *ReferralService) CreateAffiliate(ctx context.Context, input AffiliateInput) (*domain.Affiliate, error) {
	if err := s.ensureEmailAvailable(ctx, input.Email, nil); err != nil {
		return nil, err
	}
	if err := s.ensurePhoneAvailable(ctx, input.PhoneNumber, nil); err != nil {
		return nil, err
	}

	affiliate, err := domain.NewAffiliate(
		input.AffiliateID,
		input.Email,
		input.PhoneNumber,
		input.FullName,
		input.Status,
		input.SocialMediaHandle,
		input.Platform,
	)
	if err != nil {
		logger.Warn("Affiliate creation failed: validation error", "full_name", input.FullName, "error", err)
		return nil, err
	}

	if err := s.repository.CreateAffiliate(ctx, affiliate); err != nil {
		logger.Error("Affiliate creation failed: repository error", "full_name", input.FullName, "error", err)
		return nil, err
	}

	logger.Info("Affiliate created successfully", "affiliate_id", affiliate.ID)
	return affiliate, nil
}

// GetAffiliate fetches an affiliate by identifier.
func (s *ReferralService) GetAffiliate(ctx context.Context, id int64) (*domain.Affiliate, error) {
	affiliate, err := s.repository.FindAffiliateByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAffiliateNotFound) {
			logger.Debug("Affiliate retrieval failed: affiliate not found", "affiliate_id", id)
		} else {
			logger.Error("Affiliate retrieval failed: repository error", "affiliate_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Affiliate retrieved successfully", "affiliate_id", id)
	return affiliate, nil
}

// UpdateAffiliate modifies an existing affiliate.
func (s *ReferralService) UpdateAffiliate(ctx context.Context, id int64, input AffiliateInput) (*domain.Affiliate, error) {
	affiliate, err := s.repository.FindAffiliateByID(ctx, id)
	if err != nil {
		logger.Debug("Affiliate update failed: affiliate not found", "affiliate_id", id)
		return nil, err
	}

	if input.Email != "" && input.Email != affiliate.Email {
		if err := s.ensureEmailAvailable(ctx, input.Email, &id); err != nil {
			return nil, err
		}
	}
	if input.PhoneNumber != "" && input.PhoneNumber != affiliate.PhoneNumber {
		if err := s.ensurePhoneAvailable(ctx, input.PhoneNumber, &id); err != nil {
			return nil, err
		}
	}

	if err := affiliate.Update(
		input.Email,
		input.PhoneNumber,
		input.FullName,
		input.Status,
		input.SocialMediaHandle,
		input.Platform,
	); err != nil {
		logger.Warn("Affiliate update failed: validation error", "affiliate_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.UpdateAffiliate(ctx, affiliate); err != nil {
		logger.Error("Affiliate update failed: repository error", "affiliate_id", id, "error", err)
		return nil, err
	}

	logger.Info("Affiliate updated successfully", "affiliate_id", id)
	return affiliate, nil
}

// DeleteAffiliate removes an affiliate (soft delete).
func (s *ReferralService) DeleteAffiliate(ctx context.Context, id int64) error {
	if err := s.repository.DeleteAffiliate(ctx, id); err != nil {
		if errors.Is(err, ErrAffiliateNotFound) {
			logger.Debug("Affiliate deletion failed: affiliate not found", "affiliate_id", id)
		} else {
			logger.Error("Affiliate deletion failed: repository error", "affiliate_id", id, "error", err)
		}
		return err
	}

	logger.Info("Affiliate deleted successfully", "affiliate_id", id)
	return nil
}

// ListAffiliates retrieves a paginated list of affiliates.
func (s *ReferralService) ListAffiliates(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Affiliate], error) {
	result, err := s.repository.FindAllAffiliates(ctx, pageReq)
	if err != nil {
		logger.Error("Affiliate pagination failed: repository error", "page", pageReq.Page, "limit", pageReq.Limit, "error", err)
		return pagination.PageResult[*domain.Affiliate]{}, err
	}
	logger.Debug("Affiliate pagination completed", "page", result.Page, "total", result.Total, "items", len(result.Items))
	return result, nil
}

func (s *ReferralService) ensureEmailAvailable(ctx context.Context, email string, excludeID *int64) error {
	if email == "" {
		return nil
	}
	existing, err := s.repository.FindAffiliateByEmail(ctx, email)
	if err == nil && existing != nil {
		if excludeID == nil || existing.ID != *excludeID {
			logger.Warn("Affiliate email conflict detected", "email", email)
			return ErrAffiliateEmailExists
		}
	}
	if err != nil && !errors.Is(err, ErrAffiliateNotFound) {
		logger.Error("Affiliate email lookup failed", "email", email, "error", err)
		return err
	}
	return nil
}

func (s *ReferralService) ensurePhoneAvailable(ctx context.Context, phoneNumber string, excludeID *int64) error {
	if phoneNumber == "" {
		return nil
	}
	existing, err := s.repository.FindAffiliateByPhone(ctx, phoneNumber)
	if err == nil && existing != nil {
		if excludeID == nil || existing.ID != *excludeID {
			logger.Warn("Affiliate phone conflict detected", "phone_number", phoneNumber)
			return ErrAffiliatePhoneExists
		}
	}
	if err != nil && !errors.Is(err, ErrAffiliateNotFound) {
		logger.Error("Affiliate phone lookup failed", "phone_number", phoneNumber, "error", err)
		return err
	}
	return nil
}

// CreateReferralCode creates a new referral code for an affiliate.
func (s *ReferralService) CreateReferralCode(ctx context.Context, input ReferralCodeInput) (*domain.ReferralCode, error) {
	// Verify affiliate exists
	_, err := s.repository.FindAffiliateByID(ctx, input.AffiliateID)
	if err != nil {
		if errors.Is(err, ErrAffiliateNotFound) {
			logger.Debug("Referral code creation failed: affiliate not found", "affiliate_id", input.AffiliateID)
			return nil, ErrAffiliateNotFound
		}
		return nil, err
	}

	// Check if custom code already exists
	if input.CustomCode != "" {
		existing, err := s.repository.FindReferralCodeByCode(ctx, input.CustomCode)
		if err == nil && existing != nil {
			logger.Warn("Referral code already exists", "code", input.CustomCode)
			return nil, ErrReferralCodeExists
		}
		if err != nil && !errors.Is(err, ErrReferralCodeNotFound) {
			return nil, err
		}
	}

	var expiresAt *time.Time
	if input.ExpiresAt != nil && *input.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, *input.ExpiresAt)
		if err != nil {
			return nil, errors.New("invalid expiration date format, use RFC3339")
		}
		expiresAt = &parsed
	}

	code, err := domain.NewReferralCode(
		input.AffiliateID,
		input.CustomCode,
		input.Status,
		expiresAt,
	)
	if err != nil {
		logger.Warn("Referral code creation failed: validation error", "affiliate_id", input.AffiliateID, "error", err)
		return nil, err
	}

	if err := s.repository.CreateReferralCode(ctx, code); err != nil {
		logger.Error("Referral code creation failed: repository error", "affiliate_id", input.AffiliateID, "error", err)
		return nil, err
	}

	logger.Info("Referral code created successfully", "referral_code_id", code.ID, "code", code.Code)
	return code, nil
}

// GetReferralCode fetches a referral code by identifier.
func (s *ReferralService) GetReferralCode(ctx context.Context, id int64) (*domain.ReferralCode, error) {
	code, err := s.repository.FindReferralCodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrReferralCodeNotFound) {
			logger.Debug("Referral code retrieval failed: code not found", "referral_code_id", id)
		} else {
			logger.Error("Referral code retrieval failed: repository error", "referral_code_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Referral code retrieved successfully", "referral_code_id", id)
	return code, nil
}

// GetReferralCodeByCode fetches a referral code by code string.
func (s *ReferralService) GetReferralCodeByCode(ctx context.Context, code string) (*domain.ReferralCode, error) {
	referralCode, err := s.repository.FindReferralCodeByCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrReferralCodeNotFound) {
			logger.Debug("Referral code retrieval failed: code not found", "code", code)
		} else {
			logger.Error("Referral code retrieval failed: repository error", "code", code, "error", err)
		}
		return nil, err
	}
	logger.Debug("Referral code retrieved successfully", "code", code)
	return referralCode, nil
}

// ValidateReferralCode validates a referral code for use by a new customer.
// This checks if the code exists, is active, not expired, and the customer is new.
func (s *ReferralService) ValidateReferralCode(ctx context.Context, code string, customerID int64, customerService CustomerChecker) (*domain.ReferralCode, error) {
	referralCode, err := s.repository.FindReferralCodeByCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrReferralCodeNotFound) {
			logger.Debug("Referral code validation failed: code not found", "code", code)
			return nil, ErrReferralCodeNotFound
		}
		return nil, err
	}

	// Check if code is active
	if !referralCode.IsActive || referralCode.Status != domain.ReferralCodeStatusActive {
		logger.Warn("Referral code is not active", "code", code, "status", referralCode.Status)
		return nil, errors.New("referral code is not active")
	}

	// Check if expired
	if referralCode.IsExpired() {
		logger.Warn("Referral code has expired", "code", code)
		return nil, errors.New("referral code has expired")
	}

	// Check if customer is new (not existing)
	if customerService != nil {
		isNew, err := customerService.IsNewCustomer(ctx, customerID)
		if err != nil {
			logger.Error("Failed to check customer status", "customer_id", customerID, "error", err)
			return nil, err
		}
		if !isNew {
			logger.Warn("Existing customer attempted to use referral code", "customer_id", customerID, "code", code)
			return nil, ErrExistingCustomerReferral
		}
	}

	// Check if customer is trying to refer themselves
	if customerService != nil {
		affiliate, err := s.repository.FindAffiliateByID(ctx, referralCode.AffiliateID)
		if err == nil && affiliate != nil && affiliate.AffiliateID != nil && *affiliate.AffiliateID == customerID {
			logger.Warn("Customer attempted to refer themselves", "customer_id", customerID, "affiliate_id", affiliate.ID)
			return nil, ErrSelfReferral
		}
	}

	return referralCode, nil
}

// CreateReferralRelationship creates a new referral relationship.
func (s *ReferralService) CreateReferralRelationship(ctx context.Context, input ReferralRelationshipInput) (*domain.ReferralRelationship, error) {
	// Verify affiliate exists
	_, err := s.repository.FindAffiliateByID(ctx, input.AffiliateID)
	if err != nil {
		if errors.Is(err, ErrAffiliateNotFound) {
			logger.Debug("Referral relationship creation failed: affiliate not found", "affiliate_id", input.AffiliateID)
			return nil, ErrAffiliateNotFound
		}
		return nil, err
	}

	// Verify referral code exists
	_, err = s.repository.FindReferralCodeByID(ctx, input.ReferralCodeID)
	if err != nil {
		if errors.Is(err, ErrReferralCodeNotFound) {
			logger.Debug("Referral relationship creation failed: referral code not found", "referral_code_id", input.ReferralCodeID)
			return nil, ErrReferralCodeNotFound
		}
		return nil, err
	}

	// Check if relationship already exists for this customer
	existing, err := s.repository.FindReferralRelationshipByCustomerID(ctx, input.CustomerID)
	if err == nil && existing != nil {
		logger.Warn("Referral relationship already exists for customer", "customer_id", input.CustomerID)
		return existing, nil // Return existing relationship
	}
	if err != nil && !errors.Is(err, ErrReferralRelationshipNotFound) {
		return nil, err
	}

	relationship, err := domain.NewReferralRelationship(
		input.AffiliateID,
		input.CustomerID,
		input.ReferralCodeID,
		input.Source,
		input.Campaign,
		input.UTMSource,
		input.UTMMedium,
		input.UTMCampaign,
		input.UTMContent,
		input.ReferralLink,
		input.IPAddress,
		input.UserAgent,
	)
	if err != nil {
		logger.Warn("Referral relationship creation failed: validation error", "customer_id", input.CustomerID, "error", err)
		return nil, err
	}

	if err := s.repository.CreateReferralRelationship(ctx, relationship); err != nil {
		logger.Error("Referral relationship creation failed: repository error", "customer_id", input.CustomerID, "error", err)
		return nil, err
	}

	logger.Info("Referral relationship created successfully", "referral_relationship_id", relationship.ID)
	return relationship, nil
}

// GetReferralRelationship fetches a referral relationship by identifier.
func (s *ReferralService) GetReferralRelationship(ctx context.Context, id int64) (*domain.ReferralRelationship, error) {
	relationship, err := s.repository.FindReferralRelationshipByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrReferralRelationshipNotFound) {
			logger.Debug("Referral relationship retrieval failed: relationship not found", "referral_relationship_id", id)
		} else {
			logger.Error("Referral relationship retrieval failed: repository error", "referral_relationship_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Referral relationship retrieved successfully", "referral_relationship_id", id)
	return relationship, nil
}

// GetReferralRelationshipByCustomerID fetches a referral relationship by customer ID.
func (s *ReferralService) GetReferralRelationshipByCustomerID(ctx context.Context, customerID int64) (*domain.ReferralRelationship, error) {
	relationship, err := s.repository.FindReferralRelationshipByCustomerID(ctx, customerID)
	if err != nil {
		if errors.Is(err, ErrReferralRelationshipNotFound) {
			logger.Debug("Referral relationship retrieval failed: relationship not found", "customer_id", customerID)
		} else {
			logger.Error("Referral relationship retrieval failed: repository error", "customer_id", customerID, "error", err)
		}
		return nil, err
	}
	logger.Debug("Referral relationship retrieved successfully", "customer_id", customerID)
	return relationship, nil
}

// CreateCommission creates a new commission record.
// This is typically called when an order from a referred customer is completed.
func (s *ReferralService) CreateCommission(ctx context.Context, input CommissionInput) (*domain.Commission, error) {
	// Verify referral relationship exists
	relationship, err := s.repository.FindReferralRelationshipByID(ctx, input.ReferralRelationshipID)
	if err != nil {
		if errors.Is(err, ErrReferralRelationshipNotFound) {
			logger.Debug("Commission creation failed: referral relationship not found", "referral_relationship_id", input.ReferralRelationshipID)
			return nil, ErrReferralRelationshipNotFound
		}
		return nil, err
	}

	// Verify customer matches relationship
	if relationship.CustomerID != input.CustomerID {
		return nil, errors.New("customer ID does not match referral relationship")
	}

	// Parse order total from JSON
	var orderTotal *goodmoney.Money
	if input.OrderTotal != "" {
		var money goodmoney.Money
		if err := json.Unmarshal([]byte(input.OrderTotal), &money); err != nil {
			return nil, errors.New("invalid order total format")
		}
		orderTotal = &money
	}

	// Calculate commission amount
	var amount *goodmoney.Money
	if orderTotal != nil && input.CommissionRate > 0 {
		// Calculate percentage commission
		// Convert commission rate percentage to decimal and multiply by order total
		// goodmoney stores amounts in cents, so we work with the raw amount
		orderAmountCents := float64(orderTotal.Amount())
		commissionValueCents := orderAmountCents * input.CommissionRate / 100.0
		commissionMoney, err := goodmoney.New(commissionValueCents, orderTotal.Currency())
		if err != nil {
			return nil, errors.New("failed to create commission amount: " + err.Error())
		}
		amount = commissionMoney
	}

	commission, err := domain.NewCommission(
		input.AffiliateID,
		input.ReferralRelationshipID,
		input.OrderID,
		input.CustomerID,
		amount,
		input.CommissionRate,
		orderTotal,
		input.Notes,
	)
	if err != nil {
		logger.Warn("Commission creation failed: validation error", "order_id", input.OrderID, "error", err)
		return nil, err
	}

	if err := s.repository.CreateCommission(ctx, commission); err != nil {
		logger.Error("Commission creation failed: repository error", "order_id", input.OrderID, "error", err)
		return nil, err
	}

	logger.Info("Commission created successfully", "commission_id", commission.ID, "order_id", input.OrderID)
	return commission, nil
}

// GetCommission fetches a commission by identifier.
func (s *ReferralService) GetCommission(ctx context.Context, id int64) (*domain.Commission, error) {
	commission, err := s.repository.FindCommissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCommissionNotFound) {
			logger.Debug("Commission retrieval failed: commission not found", "commission_id", id)
		} else {
			logger.Error("Commission retrieval failed: repository error", "commission_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Commission retrieved successfully", "commission_id", id)
	return commission, nil
}

// ReverseCommission reverses a commission (e.g., for cancelled orders).
func (s *ReferralService) ReverseCommission(ctx context.Context, id int64, reason string) error {
	commission, err := s.repository.FindCommissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCommissionNotFound) {
			logger.Debug("Commission reversal failed: commission not found", "commission_id", id)
			return ErrCommissionNotFound
		}
		return err
	}

	if err := commission.Reverse(reason); err != nil {
		return err
	}

	if err := s.repository.UpdateCommission(ctx, commission); err != nil {
		logger.Error("Commission reversal failed: repository error", "commission_id", id, "error", err)
		return err
	}

	logger.Info("Commission reversed successfully", "commission_id", id)
	return nil
}

// CustomerChecker is an interface for checking if a customer is new.
// This should be implemented by the customer service to avoid circular dependencies.
type CustomerChecker interface {
	IsNewCustomer(ctx context.Context, customerID int64) (bool, error)
}

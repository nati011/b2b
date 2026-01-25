package referral

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	goodmoney "github.com/the-nucleus-project/good_money"
	"marketplace/internal/core/referral/domain"
	referralservice "marketplace/internal/core/referral/service"
	"marketplace/pkg/pagination"
)

var (
	// Errors are imported from the service package.
	ErrAffiliateNotFound            = referralservice.ErrAffiliateNotFound
	ErrReferralCodeNotFound         = referralservice.ErrReferralCodeNotFound
	ErrReferralRelationshipNotFound = referralservice.ErrReferralRelationshipNotFound
	ErrCommissionNotFound           = referralservice.ErrCommissionNotFound
)

// ReferralRepository persists referral data.
type ReferralRepository struct {
	db *sql.DB
}

func NewReferralRepository(db *sql.DB) *ReferralRepository {
	return &ReferralRepository{db: db}
}

// ========== Affiliate Operations ==========

// CreateAffiliate inserts a new affiliate record.
func (r *ReferralRepository) CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `
		INSERT INTO affiliates (
			affiliate_id,
			email,
			phone_number,
			full_name,
			status,
			social_media_handle,
			platform,
			is_active,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, FALSE)
		RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		nullableInt64(affiliate.AffiliateID),
		nullableString(affiliate.Email),
		nullableString(affiliate.PhoneNumber),
		affiliate.FullName,
		string(affiliate.Status),
		nullableString(affiliate.SocialMediaHandle),
		nullableString(affiliate.Platform),
		affiliate.IsActive,
		affiliate.CreatedAt,
		affiliate.UpdatedAt,
	).Scan(&affiliate.ID)
}

// UpdateAffiliate modifies an existing affiliate record.
func (r *ReferralRepository) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	query := `
		UPDATE affiliates
		SET email = $2,
		    phone_number = $3,
		    full_name = $4,
		    status = $5,
		    social_media_handle = $6,
		    platform = $7,
		    is_active = $8,
		    last_modified = $9
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		affiliate.ID,
		nullableString(affiliate.Email),
		nullableString(affiliate.PhoneNumber),
		affiliate.FullName,
		string(affiliate.Status),
		nullableString(affiliate.SocialMediaHandle),
		nullableString(affiliate.Platform),
		affiliate.IsActive,
		affiliate.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAffiliateNotFound
	}

	return nil
}

// FindAffiliateByID fetches an affiliate by identifier.
func (r *ReferralRepository) FindAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error) {
	query := `
		SELECT id, affiliate_id, email, phone_number, full_name, status, social_media_handle, platform, is_active, created_date, last_modified
		FROM affiliates
		WHERE id = $1 AND is_deleted = FALSE
	`
	affiliate, err := scanAffiliate(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAffiliateNotFound
		}
		return nil, err
	}
	return affiliate, nil
}

// FindAffiliateByEmail fetches an affiliate by email address.
func (r *ReferralRepository) FindAffiliateByEmail(ctx context.Context, email string) (*domain.Affiliate, error) {
	query := `
		SELECT id, affiliate_id, email, phone_number, full_name, status, social_media_handle, platform, is_active, created_date, last_modified
		FROM affiliates
		WHERE email = $1 AND is_deleted = FALSE
	`
	affiliate, err := scanAffiliate(r.db.QueryRowContext(ctx, query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAffiliateNotFound
		}
		return nil, err
	}
	return affiliate, nil
}

// FindAffiliateByPhone fetches an affiliate by phone number.
func (r *ReferralRepository) FindAffiliateByPhone(ctx context.Context, phone string) (*domain.Affiliate, error) {
	query := `
		SELECT id, affiliate_id, email, phone_number, full_name, status, social_media_handle, platform, is_active, created_date, last_modified
		FROM affiliates
		WHERE phone_number = $1 AND is_deleted = FALSE
	`
	affiliate, err := scanAffiliate(r.db.QueryRowContext(ctx, query, phone))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAffiliateNotFound
		}
		return nil, err
	}
	return affiliate, nil
}

// FindAffiliateByCustomerID fetches an affiliate by customer ID (if affiliate is also a customer).
func (r *ReferralRepository) FindAffiliateByCustomerID(ctx context.Context, customerID int64) (*domain.Affiliate, error) {
	query := `
		SELECT id, affiliate_id, email, phone_number, full_name, status, social_media_handle, platform, is_active, created_date, last_modified
		FROM affiliates
		WHERE affiliate_id = $1 AND is_deleted = FALSE
	`
	affiliate, err := scanAffiliate(r.db.QueryRowContext(ctx, query, customerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAffiliateNotFound
		}
		return nil, err
	}
	return affiliate, nil
}

// DeleteAffiliate performs a soft delete on the affiliate.
func (r *ReferralRepository) DeleteAffiliate(ctx context.Context, id int64) error {
	query := `
		UPDATE affiliates
		SET is_deleted = TRUE, last_modified = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAffiliateNotFound
	}

	return nil
}

// FindAllAffiliates returns a paginated list of affiliates.
func (r *ReferralRepository) FindAllAffiliates(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Affiliate], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM affiliates WHERE is_deleted = FALSE`).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Affiliate]{}, err
	}

	query := `
		SELECT id, affiliate_id, email, phone_number, full_name, status, social_media_handle, platform, is_active, created_date, last_modified
		FROM affiliates
		WHERE is_deleted = FALSE
		ORDER BY created_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.Affiliate]{}, err
	}
	defer rows.Close()

	var affiliates []*domain.Affiliate
	for rows.Next() {
		affiliate, err := scanAffiliate(rows)
		if err != nil {
			return pagination.PageResult[*domain.Affiliate]{}, err
		}
		affiliates = append(affiliates, affiliate)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.Affiliate]{}, err
	}

	return pagination.NewPageResult(affiliates, total, pageReq), nil
}

// ========== Referral Code Operations ==========

// CreateReferralCode inserts a new referral code record.
func (r *ReferralRepository) CreateReferralCode(ctx context.Context, code *domain.ReferralCode) error {
	query := `
		INSERT INTO referral_codes (
			affiliate_id,
			code,
			is_custom,
			status,
			expires_at,
			is_active,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, FALSE)
		RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		code.AffiliateID,
		code.Code,
		code.IsCustom,
		string(code.Status),
		nullableTime(code.ExpiresAt),
		code.IsActive,
		code.CreatedAt,
		code.UpdatedAt,
	).Scan(&code.ID)
}

// UpdateReferralCode modifies an existing referral code record.
func (r *ReferralRepository) UpdateReferralCode(ctx context.Context, code *domain.ReferralCode) error {
	query := `
		UPDATE referral_codes
		SET status = $2,
		    expires_at = $3,
		    is_active = $4,
		    last_modified = $5
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		code.ID,
		string(code.Status),
		nullableTime(code.ExpiresAt),
		code.IsActive,
		code.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrReferralCodeNotFound
	}

	return nil
}

// FindReferralCodeByID fetches a referral code by identifier.
func (r *ReferralRepository) FindReferralCodeByID(ctx context.Context, id int64) (*domain.ReferralCode, error) {
	query := `
		SELECT id, affiliate_id, code, is_custom, status, expires_at, is_active, created_date, last_modified
		FROM referral_codes
		WHERE id = $1 AND is_deleted = FALSE
	`
	code, err := scanReferralCode(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReferralCodeNotFound
		}
		return nil, err
	}
	return code, nil
}

// FindReferralCodeByCode fetches a referral code by code string.
func (r *ReferralRepository) FindReferralCodeByCode(ctx context.Context, code string) (*domain.ReferralCode, error) {
	query := `
		SELECT id, affiliate_id, code, is_custom, status, expires_at, is_active, created_date, last_modified
		FROM referral_codes
		WHERE code = $1 AND is_deleted = FALSE
	`
	referralCode, err := scanReferralCode(r.db.QueryRowContext(ctx, query, code))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReferralCodeNotFound
		}
		return nil, err
	}
	return referralCode, nil
}

// FindReferralCodesByAffiliateID fetches all referral codes for an affiliate.
func (r *ReferralRepository) FindReferralCodesByAffiliateID(ctx context.Context, affiliateID int64) ([]*domain.ReferralCode, error) {
	query := `
		SELECT id, affiliate_id, code, is_custom, status, expires_at, is_active, created_date, last_modified
		FROM referral_codes
		WHERE affiliate_id = $1 AND is_deleted = FALSE
		ORDER BY created_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []*domain.ReferralCode
	for rows.Next() {
		code, err := scanReferralCode(rows)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return codes, nil
}

// DeleteReferralCode performs a soft delete on the referral code.
func (r *ReferralRepository) DeleteReferralCode(ctx context.Context, id int64) error {
	query := `
		UPDATE referral_codes
		SET is_deleted = TRUE, last_modified = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrReferralCodeNotFound
	}

	return nil
}

// ========== Referral Relationship Operations ==========

// CreateReferralRelationship inserts a new referral relationship record.
func (r *ReferralRepository) CreateReferralRelationship(ctx context.Context, relationship *domain.ReferralRelationship) error {
	query := `
		INSERT INTO referral_relationships (
			affiliate_id,
			customer_id,
			referral_code_id,
			source,
			campaign,
			utm_source,
			utm_medium,
			utm_campaign,
			utm_content,
			referral_link,
			ip_address,
			user_agent,
			is_active,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, FALSE)
		RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		relationship.AffiliateID,
		relationship.CustomerID,
		relationship.ReferralCodeID,
		nullableString(relationship.Source),
		nullableString(relationship.Campaign),
		nullableString(relationship.UTMSource),
		nullableString(relationship.UTMMedium),
		nullableString(relationship.UTMCampaign),
		nullableString(relationship.UTMContent),
		nullableString(relationship.ReferralLink),
		nullableString(relationship.IPAddress),
		nullableString(relationship.UserAgent),
		relationship.IsActive,
		relationship.CreatedAt,
		relationship.UpdatedAt,
	).Scan(&relationship.ID)
}

// UpdateReferralRelationship modifies an existing referral relationship record.
func (r *ReferralRepository) UpdateReferralRelationship(ctx context.Context, relationship *domain.ReferralRelationship) error {
	query := `
		UPDATE referral_relationships
		SET source = $2,
		    campaign = $3,
		    utm_source = $4,
		    utm_medium = $5,
		    utm_campaign = $6,
		    utm_content = $7,
		    last_modified = $8
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		relationship.ID,
		nullableString(relationship.Source),
		nullableString(relationship.Campaign),
		nullableString(relationship.UTMSource),
		nullableString(relationship.UTMMedium),
		nullableString(relationship.UTMCampaign),
		nullableString(relationship.UTMContent),
		relationship.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrReferralRelationshipNotFound
	}

	return nil
}

// FindReferralRelationshipByID fetches a referral relationship by identifier.
func (r *ReferralRepository) FindReferralRelationshipByID(ctx context.Context, id int64) (*domain.ReferralRelationship, error) {
	query := `
		SELECT id, affiliate_id, customer_id, referral_code_id, source, campaign, utm_source, utm_medium, utm_campaign, utm_content, referral_link, ip_address, user_agent, is_active, created_date, last_modified
		FROM referral_relationships
		WHERE id = $1 AND is_deleted = FALSE
	`
	relationship, err := scanReferralRelationship(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReferralRelationshipNotFound
		}
		return nil, err
	}
	return relationship, nil
}

// FindReferralRelationshipByCustomerID fetches a referral relationship by customer ID.
func (r *ReferralRepository) FindReferralRelationshipByCustomerID(ctx context.Context, customerID int64) (*domain.ReferralRelationship, error) {
	query := `
		SELECT id, affiliate_id, customer_id, referral_code_id, source, campaign, utm_source, utm_medium, utm_campaign, utm_content, referral_link, ip_address, user_agent, is_active, created_date, last_modified
		FROM referral_relationships
		WHERE customer_id = $1 AND is_deleted = FALSE
		ORDER BY created_date DESC
		LIMIT 1
	`
	relationship, err := scanReferralRelationship(r.db.QueryRowContext(ctx, query, customerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReferralRelationshipNotFound
		}
		return nil, err
	}
	return relationship, nil
}

// FindReferralRelationshipsByAffiliateID returns a paginated list of referral relationships for an affiliate.
func (r *ReferralRepository) FindReferralRelationshipsByAffiliateID(ctx context.Context, affiliateID int64, pageReq pagination.PageRequest) (pagination.PageResult[*domain.ReferralRelationship], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM referral_relationships WHERE affiliate_id = $1 AND is_deleted = FALSE`, affiliateID).Scan(&total); err != nil {
		return pagination.PageResult[*domain.ReferralRelationship]{}, err
	}

	query := `
		SELECT id, affiliate_id, customer_id, referral_code_id, source, campaign, utm_source, utm_medium, utm_campaign, utm_content, referral_link, ip_address, user_agent, is_active, created_date, last_modified
		FROM referral_relationships
		WHERE affiliate_id = $1 AND is_deleted = FALSE
		ORDER BY created_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.ReferralRelationship]{}, err
	}
	defer rows.Close()

	var relationships []*domain.ReferralRelationship
	for rows.Next() {
		relationship, err := scanReferralRelationship(rows)
		if err != nil {
			return pagination.PageResult[*domain.ReferralRelationship]{}, err
		}
		relationships = append(relationships, relationship)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.ReferralRelationship]{}, err
	}

	return pagination.NewPageResult(relationships, total, pageReq), nil
}

// DeleteReferralRelationship performs a soft delete on the referral relationship.
func (r *ReferralRepository) DeleteReferralRelationship(ctx context.Context, id int64) error {
	query := `
		UPDATE referral_relationships
		SET is_deleted = TRUE, last_modified = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrReferralRelationshipNotFound
	}

	return nil
}

// ========== Commission Operations ==========

// CreateCommission inserts a new commission record.
func (r *ReferralRepository) CreateCommission(ctx context.Context, commission *domain.Commission) error {
	query := `
		INSERT INTO commissions (
			affiliate_id,
			referral_relationship_id,
			order_id,
			customer_id,
			amount,
			commission_rate,
			status,
			order_total,
			notes,
			paid_at,
			reversed_at,
			is_active,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, FALSE)
		RETURNING id
	`

	var amountJSON []byte
	if commission.Amount != nil {
		var err error
		amountJSON, err = json.Marshal(commission.Amount)
		if err != nil {
			return err
		}
	}

	var orderTotalJSON []byte
	if commission.OrderTotal != nil {
		var err error
		orderTotalJSON, err = json.Marshal(commission.OrderTotal)
		if err != nil {
			return err
		}
	}

	return r.db.QueryRowContext(
		ctx,
		query,
		commission.AffiliateID,
		commission.ReferralRelationshipID,
		commission.OrderID,
		commission.CustomerID,
		nullableBytes(amountJSON),
		commission.CommissionRate,
		string(commission.Status),
		nullableBytes(orderTotalJSON),
		nullableString(commission.Notes),
		nullableTime(commission.PaidAt),
		nullableTime(commission.ReversedAt),
		commission.IsActive,
		commission.CreatedAt,
		commission.UpdatedAt,
	).Scan(&commission.ID)
}

// UpdateCommission modifies an existing commission record.
func (r *ReferralRepository) UpdateCommission(ctx context.Context, commission *domain.Commission) error {
	query := `
		UPDATE commissions
		SET status = $2,
		    notes = $3,
		    paid_at = $4,
		    reversed_at = $5,
		    last_modified = $6
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		commission.ID,
		string(commission.Status),
		nullableString(commission.Notes),
		nullableTime(commission.PaidAt),
		nullableTime(commission.ReversedAt),
		commission.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCommissionNotFound
	}

	return nil
}

// FindCommissionByID fetches a commission by identifier.
func (r *ReferralRepository) FindCommissionByID(ctx context.Context, id int64) (*domain.Commission, error) {
	query := `
		SELECT id, affiliate_id, referral_relationship_id, order_id, customer_id, amount, commission_rate, status, order_total, notes, paid_at, reversed_at, is_active, created_date, last_modified
		FROM commissions
		WHERE id = $1 AND is_deleted = FALSE
	`
	commission, err := scanCommission(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommissionNotFound
		}
		return nil, err
	}
	return commission, nil
}

// FindCommissionsByAffiliateID returns a paginated list of commissions for an affiliate.
func (r *ReferralRepository) FindCommissionsByAffiliateID(ctx context.Context, affiliateID int64, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Commission], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM commissions WHERE affiliate_id = $1 AND is_deleted = FALSE`, affiliateID).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Commission]{}, err
	}

	query := `
		SELECT id, affiliate_id, referral_relationship_id, order_id, customer_id, amount, commission_rate, status, order_total, notes, paid_at, reversed_at, is_active, created_date, last_modified
		FROM commissions
		WHERE affiliate_id = $1 AND is_deleted = FALSE
		ORDER BY created_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, affiliateID, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.Commission]{}, err
	}
	defer rows.Close()

	var commissions []*domain.Commission
	for rows.Next() {
		commission, err := scanCommission(rows)
		if err != nil {
			return pagination.PageResult[*domain.Commission]{}, err
		}
		commissions = append(commissions, commission)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.Commission]{}, err
	}

	return pagination.NewPageResult(commissions, total, pageReq), nil
}

// FindCommissionsByOrderID fetches all commissions for an order.
func (r *ReferralRepository) FindCommissionsByOrderID(ctx context.Context, orderID int64) ([]*domain.Commission, error) {
	query := `
		SELECT id, affiliate_id, referral_relationship_id, order_id, customer_id, amount, commission_rate, status, order_total, notes, paid_at, reversed_at, is_active, created_date, last_modified
		FROM commissions
		WHERE order_id = $1 AND is_deleted = FALSE
		ORDER BY created_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []*domain.Commission
	for rows.Next() {
		commission, err := scanCommission(rows)
		if err != nil {
			return nil, err
		}
		commissions = append(commissions, commission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commissions, nil
}

// FindCommissionsByReferralRelationshipID fetches all commissions for a referral relationship.
func (r *ReferralRepository) FindCommissionsByReferralRelationshipID(ctx context.Context, referralRelationshipID int64) ([]*domain.Commission, error) {
	query := `
		SELECT id, affiliate_id, referral_relationship_id, order_id, customer_id, amount, commission_rate, status, order_total, notes, paid_at, reversed_at, is_active, created_date, last_modified
		FROM commissions
		WHERE referral_relationship_id = $1 AND is_deleted = FALSE
		ORDER BY created_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, referralRelationshipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []*domain.Commission
	for rows.Next() {
		commission, err := scanCommission(rows)
		if err != nil {
			return nil, err
		}
		commissions = append(commissions, commission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commissions, nil
}

// DeleteCommission performs a soft delete on the commission.
func (r *ReferralRepository) DeleteCommission(ctx context.Context, id int64) error {
	query := `
		UPDATE commissions
		SET is_deleted = TRUE, last_modified = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCommissionNotFound
	}

	return nil
}

// ========== Helper Functions ==========

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

func nullableInt64(value *int64) *int64 {
	return value
}

func nullableTime(value *time.Time) *time.Time {
	return value
}

func nullableBytes(value []byte) []byte {
	if len(value) == 0 {
		return nil
	}
	return value
}

type affiliateScanner interface {
	Scan(dest ...interface{}) error
}

func scanAffiliate(scanner affiliateScanner) (*domain.Affiliate, error) {
	var affiliate domain.Affiliate
	var affiliateID sql.NullInt64
	var email, phoneNumber, status, socialMediaHandle, platform sql.NullString
	if err := scanner.Scan(
		&affiliate.ID,
		&affiliateID,
		&email,
		&phoneNumber,
		&affiliate.FullName,
		&status,
		&socialMediaHandle,
		&platform,
		&affiliate.IsActive,
		&affiliate.CreatedAt,
		&affiliate.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if affiliateID.Valid {
		affiliate.AffiliateID = &affiliateID.Int64
	}
	if email.Valid {
		affiliate.Email = email.String
	}
	if phoneNumber.Valid {
		affiliate.PhoneNumber = phoneNumber.String
	}
	if status.Valid {
		affiliate.Status = domain.AffiliateStatus(status.String)
	}
	if socialMediaHandle.Valid {
		affiliate.SocialMediaHandle = socialMediaHandle.String
	}
	if platform.Valid {
		affiliate.Platform = platform.String
	}

	return &affiliate, nil
}

type referralCodeScanner interface {
	Scan(dest ...interface{}) error
}

func scanReferralCode(scanner referralCodeScanner) (*domain.ReferralCode, error) {
	var code domain.ReferralCode
	var status sql.NullString
	var expiresAt sql.NullTime
	if err := scanner.Scan(
		&code.ID,
		&code.AffiliateID,
		&code.Code,
		&code.IsCustom,
		&status,
		&expiresAt,
		&code.IsActive,
		&code.CreatedAt,
		&code.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if status.Valid {
		code.Status = domain.ReferralCodeStatus(status.String)
	}
	if expiresAt.Valid {
		code.ExpiresAt = &expiresAt.Time
	}

	return &code, nil
}

type referralRelationshipScanner interface {
	Scan(dest ...interface{}) error
}

func scanReferralRelationship(scanner referralRelationshipScanner) (*domain.ReferralRelationship, error) {
	var relationship domain.ReferralRelationship
	var source, campaign, utmSource, utmMedium, utmCampaign, utmContent, referralLink, ipAddress, userAgent sql.NullString
	if err := scanner.Scan(
		&relationship.ID,
		&relationship.AffiliateID,
		&relationship.CustomerID,
		&relationship.ReferralCodeID,
		&source,
		&campaign,
		&utmSource,
		&utmMedium,
		&utmCampaign,
		&utmContent,
		&referralLink,
		&ipAddress,
		&userAgent,
		&relationship.IsActive,
		&relationship.CreatedAt,
		&relationship.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if source.Valid {
		relationship.Source = source.String
	}
	if campaign.Valid {
		relationship.Campaign = campaign.String
	}
	if utmSource.Valid {
		relationship.UTMSource = utmSource.String
	}
	if utmMedium.Valid {
		relationship.UTMMedium = utmMedium.String
	}
	if utmCampaign.Valid {
		relationship.UTMCampaign = utmCampaign.String
	}
	if utmContent.Valid {
		relationship.UTMContent = utmContent.String
	}
	if referralLink.Valid {
		relationship.ReferralLink = referralLink.String
	}
	if ipAddress.Valid {
		relationship.IPAddress = ipAddress.String
	}
	if userAgent.Valid {
		relationship.UserAgent = userAgent.String
	}

	return &relationship, nil
}

type commissionScanner interface {
	Scan(dest ...interface{}) error
}

func scanCommission(scanner commissionScanner) (*domain.Commission, error) {
	var commission domain.Commission
	var status, notes sql.NullString
	var amountJSON, orderTotalJSON []byte
	var paidAt, reversedAt sql.NullTime
	if err := scanner.Scan(
		&commission.ID,
		&commission.AffiliateID,
		&commission.ReferralRelationshipID,
		&commission.OrderID,
		&commission.CustomerID,
		&amountJSON,
		&commission.CommissionRate,
		&status,
		&orderTotalJSON,
		&notes,
		&paidAt,
		&reversedAt,
		&commission.IsActive,
		&commission.CreatedAt,
		&commission.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if status.Valid {
		commission.Status = domain.CommissionStatus(status.String)
	}
	if notes.Valid {
		commission.Notes = notes.String
	}
	if len(amountJSON) > 0 {
		var amount goodmoney.Money
		if err := json.Unmarshal(amountJSON, &amount); err == nil {
			commission.Amount = &amount
		}
	}
	if len(orderTotalJSON) > 0 {
		var orderTotal goodmoney.Money
		if err := json.Unmarshal(orderTotalJSON, &orderTotal); err == nil {
			commission.OrderTotal = &orderTotal
		}
	}
	if paidAt.Valid {
		commission.PaidAt = &paidAt.Time
	}
	if reversedAt.Valid {
		commission.ReversedAt = &reversedAt.Time
	}

	return &commission, nil
}


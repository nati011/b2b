package repository

import (
	"context"
	userservice "marketplace/internal/infra/user/service"
	"database/sql"
	"errors"
)

type phoneValidationRepository struct {
	db *sql.DB
}

// NewPhoneValidationRepository creates a new phone validation repository
func NewPhoneValidationRepository(db *sql.DB) userservice.PhoneValidationRepository {
	return &phoneValidationRepository{db: db}
}

// FindByCode retrieves a validation scheme by code
func (r *phoneValidationRepository) FindByCode(ctx context.Context, code string) (*userservice.PhoneValidationScheme, error) {
	query := `
		SELECT id, code, name, description, regex_pattern, is_default, active
		FROM phone_validation_schemes
		WHERE code = $1 AND active = true
	`

	var scheme userservice.PhoneValidationScheme
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&scheme.ID,
		&scheme.Code,
		&scheme.Name,
		&scheme.Description,
		&scheme.RegexPattern,
		&scheme.IsDefault,
		&scheme.Active,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("phone validation scheme not found")
		}
		return nil, err
	}

	return &scheme, nil
}

// FindDefault retrieves the default validation scheme
func (r *phoneValidationRepository) FindDefault(ctx context.Context) (*userservice.PhoneValidationScheme, error) {
	query := `
		SELECT id, code, name, description, regex_pattern, is_default, active
		FROM phone_validation_schemes
		WHERE is_default = true AND active = true
		ORDER BY created_at ASC
		LIMIT 1
	`

	var scheme userservice.PhoneValidationScheme
	err := r.db.QueryRowContext(ctx, query).Scan(
		&scheme.ID,
		&scheme.Code,
		&scheme.Name,
		&scheme.Description,
		&scheme.RegexPattern,
		&scheme.IsDefault,
		&scheme.Active,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no default phone validation scheme found")
		}
		return nil, err
	}

	return &scheme, nil
}

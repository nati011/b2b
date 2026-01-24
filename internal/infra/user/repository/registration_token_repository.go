package repository

import (
	"context"
	"marketplace/internal/infra/user/domain"
	"database/sql"
	"errors"
)

var (
	ErrRegistrationTokenNotFound = errors.New("registration token not found")
)

// registrationTokenRepository implements service.RegistrationTokenRepository for database storage
type registrationTokenRepository struct {
	db *sql.DB
}

// NewRegistrationTokenRepository creates a new registration token repository
func NewRegistrationTokenRepository(db *sql.DB) *registrationTokenRepository {
	return &registrationTokenRepository{db: db}
}

// Create inserts a new registration token
func (r *registrationTokenRepository) Create(ctx context.Context, token *domain.RegistrationToken) error {
	query := `
		INSERT INTO registration_tokens (id, user_id, token, used, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.Used,
		token.ExpiresAt,
		token.CreatedAt,
	)

	return err
}

// FindByToken retrieves a registration token by token value
func (r *registrationTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RegistrationToken, error) {
	query := `
		SELECT id, user_id, token, used, expires_at, created_at
		FROM registration_tokens
		WHERE token = $1
	`

	var regToken domain.RegistrationToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&regToken.ID,
		&regToken.UserID,
		&regToken.Token,
		&regToken.Used,
		&regToken.ExpiresAt,
		&regToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRegistrationTokenNotFound
		}
		return nil, err
	}

	return &regToken, nil
}

// MarkAsUsed marks a registration token as used
func (r *registrationTokenRepository) MarkAsUsed(ctx context.Context, tokenID string) error {
	query := `
		UPDATE registration_tokens
		SET used = true
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRegistrationTokenNotFound
	}

	return nil
}

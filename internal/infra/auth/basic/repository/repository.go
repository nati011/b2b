package repository

import (
	"context"
	"database/sql"
	"errors"

	"marketplace/internal/infra/auth/basic/domain"
)

// Repository defines the interface for basic auth credential storage
type Repository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Credential, error)
	Create(ctx context.Context, cred *domain.Credential) error
	Update(ctx context.Context, cred *domain.Credential) error
	Delete(ctx context.Context, id string) error
}

// CredentialRepository implements Repository for database storage
type CredentialRepository struct {
	db *sql.DB
}

// NewCredentialRepository creates a new credential repository
func NewCredentialRepository(db *sql.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

// FindByUsername retrieves a credential by username (case-insensitive)
func (r *CredentialRepository) FindByUsername(ctx context.Context, username string) (*domain.Credential, error) {
	query := `
		SELECT id, username, password, user_id, active, created_at, updated_at
		FROM basic_auth_credentials
		WHERE LOWER(username) = LOWER($1) AND active = true
	`

	var cred domain.Credential
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&cred.ID,
		&cred.Username,
		&cred.Password,
		&cred.UserID,
		&cred.Active,
		&cred.CreatedAt,
		&cred.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCredentialNotFound
		}
		return nil, err
	}

	return &cred, nil
}

// Create inserts a new credential
func (r *CredentialRepository) Create(ctx context.Context, cred *domain.Credential) error {
	query := `
		INSERT INTO basic_auth_credentials (id, username, password, user_id, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		cred.ID,
		cred.Username,
		cred.Password,
		cred.UserID,
		cred.Active,
		cred.CreatedAt,
		cred.UpdatedAt,
	)

	return err
}

// Update updates an existing credential
func (r *CredentialRepository) Update(ctx context.Context, cred *domain.Credential) error {
	query := `
		UPDATE basic_auth_credentials
		SET username = $2, password = $3, user_id = $4, active = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		cred.ID,
		cred.Username,
		cred.Password,
		cred.UserID,
		cred.Active,
		cred.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrCredentialNotFound
	}

	return nil
}

// Delete removes a credential
func (r *CredentialRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM basic_auth_credentials WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrCredentialNotFound
	}

	return nil
}


package permission

import (
	"context"
	"marketplace/internal/infra/authz/permission/domain"
	"marketplace/pkg/pagination"
	"database/sql"
	"errors"
	"strings"
)

var (
	// ErrPermissionNotFound indicates a permission lookup failure.
	ErrPermissionNotFound = errors.New("permission not found")
	// ErrPermissionAlreadyExists indicates duplicate IDs or resource/action combos.
	ErrPermissionAlreadyExists = errors.New("permission already exists")
)

// PermissionRepository persists permissions.
type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// Create inserts a new permission record.
func (r *PermissionRepository) Create(ctx context.Context, perm *domain.Permission) error {
	query := `
		INSERT INTO permissions (id, resource, action, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, perm.ID, perm.Resource.Code, perm.Action, perm.Description, perm.CreatedAt, perm.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrPermissionAlreadyExists
		}
		return err
	}
	return nil
}

// Update modifies an existing permission.
func (r *PermissionRepository) Update(ctx context.Context, perm *domain.Permission) error {
	query := `
		UPDATE permissions
		SET resource = $2, action = $3, description = $4, updated_at = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, perm.ID, perm.Resource.Code, perm.Action, perm.Description, perm.UpdatedAt)
	if err != nil {
		if isInvalidUUIDError(err) {
			return ErrPermissionNotFound
		}
		if isUniqueViolation(err) {
			return ErrPermissionAlreadyExists
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPermissionNotFound
	}

	return nil
}

// FindByID fetches a permission by identifier.
func (r *PermissionRepository) FindByID(ctx context.Context, id string) (*domain.Permission, error) {
	query := `
		SELECT p.id, p.action, p.description, p.created_at, p.updated_at,
		       res.id, res.code, res.service, res.description
		FROM permissions p
		JOIN resources res ON res.code = p.resource
		WHERE p.id = $1
	`

	var perm domain.Permission
	var resID, resCode, resService, resDescription string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&perm.ID,
		&perm.Action,
		&perm.Description,
		&perm.CreatedAt,
		&perm.UpdatedAt,
		&resID,
		&resCode,
		&resService,
		&resDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUIDError(err) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	perm.Resource = domain.ResourceRef{
		ID:          resID,
		Code:        resCode,
		Service:     resService,
		Description: resDescription,
	}
	return &perm, nil
}

// Delete removes a permission by identifier.
func (r *PermissionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM permissions WHERE id = $1`, id)
	if err != nil {
		if isInvalidUUIDError(err) {
			return ErrPermissionNotFound
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrPermissionNotFound
	}
	return nil
}

// FindAll returns a paginated list of permissions.
func (r *PermissionRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Permission], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM permissions`).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Permission]{}, err
	}

	query := `
		SELECT p.id, p.action, p.description, p.created_at, p.updated_at,
		       res.id, res.code, res.service, res.description
		FROM permissions p
		JOIN resources res ON res.code = p.resource
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.Permission]{}, err
	}
	defer rows.Close()

	var perms []*domain.Permission
	for rows.Next() {
		var perm domain.Permission
		var resID, resCode, resService, resDescription string
		if err := rows.Scan(
			&perm.ID,
			&perm.Action,
			&perm.Description,
			&perm.CreatedAt,
			&perm.UpdatedAt,
			&resID,
			&resCode,
			&resService,
			&resDescription,
		); err != nil {
			return pagination.PageResult[*domain.Permission]{}, err
		}
		perm.Resource = domain.ResourceRef{
			ID:          resID,
			Code:        resCode,
			Service:     resService,
			Description: resDescription,
		}
		perms = append(perms, &perm)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.Permission]{}, err
	}

	return pagination.NewPageResult(perms, total, pageReq), nil
}

// Exists checks if a permission exists by ID.
func (r *PermissionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM permissions WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		if isInvalidUUIDError(err) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

// FindByResourceAndAction fetches a permission by resource code and action.
func (r *PermissionRepository) FindByResourceAndAction(ctx context.Context, resourceCode, action string) (*domain.Permission, error) {
	query := `
		SELECT p.id, p.action, p.description, p.created_at, p.updated_at,
		       res.id, res.code, res.service, res.description
		FROM permissions p
		JOIN resources res ON res.code = p.resource
		WHERE p.resource = $1 AND p.action = $2
	`

	var perm domain.Permission
	var resID, resCode, resService, resDescription string
	err := r.db.QueryRowContext(ctx, query, resourceCode, action).Scan(
		&perm.ID,
		&perm.Action,
		&perm.Description,
		&perm.CreatedAt,
		&perm.UpdatedAt,
		&resID,
		&resCode,
		&resService,
		&resDescription,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	perm.Resource = domain.ResourceRef{
		ID:          resID,
		Code:        resCode,
		Service:     resService,
		Description: resDescription,
	}
	return &perm, nil
}

func isInvalidUUIDError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "22P02") || strings.Contains(errStr, "invalid input syntax for type uuid")
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "23505") || strings.Contains(errStr, "duplicate key value violates unique constraint")
}

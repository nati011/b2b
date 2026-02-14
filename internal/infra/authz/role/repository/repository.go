package role

import (
	"context"
	"marketplace/internal/infra/authz/role/domain"
	roleservice "marketplace/internal/infra/authz/role/service"
	dbinfra "marketplace/internal/infra/db"
	"marketplace/pkg/pagination"
	"database/sql"
	"errors"
	"strings"
)

var (
	// ErrRoleNotFound is imported from service package
	ErrRoleNotFound = roleservice.ErrRoleNotFound
)

// RoleRepository persists roles and their permissions.
type RoleRepository struct {
	db        *sql.DB
	txManager dbinfra.TxManager
}

func NewRoleRepository(db *sql.DB, txManager dbinfra.TxManager) *RoleRepository {
	return &RoleRepository{
		db:        db,
		txManager: txManager,
	}
}

// isInvalidUUIDError checks if the error is a PostgreSQL invalid UUID error.
func isInvalidUUIDError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "22P02") || strings.Contains(errStr, "invalid input syntax for type uuid")
}

// Create inserts a new role with its permissions.
func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	return r.txManager.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		insertRoleQuery := `
			INSERT INTO roles (id, name, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`

		if _, err := tx.ExecContext(ctx, insertRoleQuery, role.ID, role.Name, role.Description, role.CreatedAt, role.UpdatedAt); err != nil {
			return err
		}

		if err := r.replacePermissionIDs(ctx, tx, role.ID, role.PermissionIDs); err != nil {
			return err
		}

		return nil
	})
}

// Update modifies an existing role and refreshes its permissions.
func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) (err error) {
	return r.txManager.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		updateQuery := `
			UPDATE roles
			SET name = $2, description = $3, updated_at = $4
			WHERE id = $1
		`

		result, err := tx.ExecContext(ctx, updateQuery, role.ID, role.Name, role.Description, role.UpdatedAt)
		if err != nil {
			if isInvalidUUIDError(err) {
				return ErrRoleNotFound
			}
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return ErrRoleNotFound
		}

		if err := r.replacePermissionIDs(ctx, tx, role.ID, role.PermissionIDs); err != nil {
			return err
		}

		return nil
	})
}

// FindByID fetches a role and eagerly loads its permissions.
func (r *RoleRepository) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	var role domain.Role

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUIDError(err) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	permissionIDs, err := r.fetchPermissionIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	role.PermissionIDs = permissionIDs

	return &role, nil
}

// Delete removes a role and its permissions.
func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.txManager.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, id); err != nil {
			if isInvalidUUIDError(err) {
				return ErrRoleNotFound
			}
			return err
		}

		result, err := tx.ExecContext(ctx, `DELETE FROM roles WHERE id = $1`, id)
		if err != nil {
			if isInvalidUUIDError(err) {
				return ErrRoleNotFound
			}
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return ErrRoleNotFound
		}

		return nil
	})
}

// FindAll returns a paginated list of roles with their permissions.
func (r *RoleRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Role], error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM roles`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Role]{}, err
	}

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.Role]{}, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return pagination.PageResult[*domain.Role]{}, err
		}

		permissionIDs, err := r.fetchPermissionIDs(ctx, role.ID)
		if err != nil {
			return pagination.PageResult[*domain.Role]{}, err
		}
		role.PermissionIDs = permissionIDs
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.Role]{}, err
	}

	return pagination.NewPageResult(roles, total, pageReq), nil
}

// Exists confirms if a role exists by identifier.
func (r *RoleRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if isInvalidUUIDError(err) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *RoleRepository) fetchPermissionIDs(ctx context.Context, roleID string) ([]string, error) {
	query := `
		SELECT p.id
		FROM role_permissions rp
		JOIN permissions p ON p.resource = rp.resource AND p.action = rp.action
		WHERE rp.role_id = $1
		ORDER BY p.id
	`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissionIDs []string
	for rows.Next() {
		var permissionID string
		if err := rows.Scan(&permissionID); err != nil {
			return nil, err
		}
		permissionIDs = append(permissionIDs, permissionID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissionIDs, nil
}

func (r *RoleRepository) replacePermissionIDs(ctx context.Context, tx *sql.Tx, roleID string, permissionIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}

	if len(permissionIDs) == 0 {
		return nil
	}

	// Fetch permission details to get resource and action for the junction table
	// We need resource and action to insert into role_permissions table
	query := `
		SELECT p.resource, p.action
		FROM permissions p
		WHERE p.id = ANY($1)
	`
	rows, err := tx.QueryContext(ctx, query, permissionIDs)
	if err != nil {
		return err
	}
	defer rows.Close()

	type permRef struct {
		resource string
		action   string
	}
	var permRefs []permRef
	for rows.Next() {
		var ref permRef
		if err := rows.Scan(&ref.resource, &ref.action); err != nil {
			return err
		}
		permRefs = append(permRefs, ref)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(permRefs) != len(permissionIDs) {
		return errors.New("some permission IDs not found")
	}

	insertQuery := `
		INSERT INTO role_permissions (role_id, resource, action)
		VALUES ($1, $2, $3)
	`

	for _, ref := range permRefs {
		if _, err := tx.ExecContext(ctx, insertQuery, roleID, ref.resource, ref.action); err != nil {
			return err
		}
	}

	return nil
}

// FindByName returns a role by its case-insensitive name.
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE LOWER(name) = LOWER($1)
	`

	var role domain.Role
	if err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	permissionIDs, err := r.fetchPermissionIDs(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.PermissionIDs = permissionIDs
	return &role, nil
}

// GetByName returns a role by its case-insensitive name.
// This method implements the RoleReader interface for user service.
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	return r.FindByName(ctx, name)
}

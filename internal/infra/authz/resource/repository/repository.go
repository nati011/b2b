package resource

import (
	"context"
	"marketplace/internal/infra/authz/resource/domain"
	dbinfra "marketplace/internal/infra/db"
	"database/sql"
	"time"
)

// Repository persists resources and their allowed actions.
type Repository struct {
	db        *sql.DB
	txManager dbinfra.TxManager
}

func NewRepository(db *sql.DB, txManager dbinfra.TxManager) *Repository {
	return &Repository{
		db:        db,
		txManager: txManager,
	}
}

// Create inserts a new resource with the provided actions.
func (r *Repository) Create(ctx context.Context, resource *domain.Resource) error {
	return r.txManager.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		insertResource := `
			INSERT INTO resources (id, code, service, description, deprecated_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		if _, err := tx.ExecContext(ctx, insertResource, resource.ID, resource.Code, resource.Service, resource.Description, resource.DeprecatedAt, resource.CreatedAt, resource.UpdatedAt); err != nil {
			return err
		}

		if err := r.ensureActions(ctx, tx, resource.ID, resource.Actions); err != nil {
			return err
		}

		return nil
	})
}

// Update modifies resource metadata and ensures the provided actions exist.
func (r *Repository) Update(ctx context.Context, resource *domain.Resource) error {
	return r.txManager.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := `
			UPDATE resources
			SET service = $2, description = $3, deprecated_at = $4, updated_at = $5
			WHERE id = $1
		`

		result, err := tx.ExecContext(ctx, query, resource.ID, resource.Service, resource.Description, resource.DeprecatedAt, resource.UpdatedAt)
		if err != nil {
			return err
		}

		if _, err := result.RowsAffected(); err != nil {
			return err
		}

		if err := r.ensureActions(ctx, tx, resource.ID, resource.Actions); err != nil {
			return err
		}

		return nil
	})
}

// Deprecate marks a resource as deprecated.
func (r *Repository) Deprecate(ctx context.Context, id string, deprecatedAt *time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE resources SET deprecated_at = $2, updated_at = $2 WHERE id = $1
	`, id, deprecatedAt)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// FindByCode fetches a resource (with actions) by code.
func (r *Repository) FindByCode(ctx context.Context, code string) (*domain.Resource, error) {
	query := `
		SELECT id, code, service, description, deprecated_at, created_at, updated_at
		FROM resources
		WHERE LOWER(code) = LOWER($1)
	`
	var res domain.Resource
	if err := r.db.QueryRowContext(ctx, query, code).Scan(
		&res.ID,
		&res.Code,
		&res.Service,
		&res.Description,
		&res.DeprecatedAt,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, err
	}

	actions, err := r.fetchActions(ctx, res.ID)
	if err != nil {
		return nil, err
	}
	res.AttachActions(actions)

	return &res, nil
}

// List fetches all resources including their actions.
func (r *Repository) List(ctx context.Context) ([]*domain.Resource, error) {
	query := `
		SELECT id, code, service, description, deprecated_at, created_at, updated_at
		FROM resources
		ORDER BY code ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*domain.Resource
	for rows.Next() {
		var res domain.Resource
		if err := rows.Scan(
			&res.ID,
			&res.Code,
			&res.Service,
			&res.Description,
			&res.DeprecatedAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return nil, err
		}

		actions, err := r.fetchActions(ctx, res.ID)
		if err != nil {
			return nil, err
		}
		res.AttachActions(actions)
		resources = append(resources, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

// FindAction resolves a single action for the provided resource code.
func (r *Repository) FindAction(ctx context.Context, code, action string) (*domain.Action, error) {
	query := `
		SELECT ra.id, ra.resource_id, ra.action, ra.description, ra.deprecated_at, ra.created_at, ra.updated_at
		FROM resource_actions ra
		JOIN resources res ON res.id = ra.resource_id
		WHERE LOWER(res.code) = LOWER($1) AND LOWER(ra.action) = LOWER($2)
	`

	var act domain.Action
	if err := r.db.QueryRowContext(ctx, query, code, action).Scan(
		&act.ID,
		&act.ResourceID,
		&act.Name,
		&act.Description,
		&act.DeprecatedAt,
		&act.CreatedAt,
		&act.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &act, nil
}

// EnsureActions inserts action rows when they do not exist.
func (r *Repository) ensureActions(ctx context.Context, tx *sql.Tx, resourceID string, actions []domain.Action) error {
	if len(actions) == 0 {
		return nil
	}

	query := `
		INSERT INTO resource_actions (id, resource_id, action, description, deprecated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT ON CONSTRAINT resource_actions_resource_action_uc DO UPDATE
		SET description = EXCLUDED.description,
		    updated_at = EXCLUDED.updated_at
	`
	for _, action := range actions {
		if _, err := tx.ExecContext(ctx, query, action.ID, action.ResourceID, action.Name, action.Description, action.DeprecatedAt, action.CreatedAt, action.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

// fetchActions loads all actions for the given resource.
func (r *Repository) fetchActions(ctx context.Context, resourceID string) ([]domain.Action, error) {
	query := `
		SELECT id, resource_id, action, description, deprecated_at, created_at, updated_at
		FROM resource_actions
		WHERE resource_id = $1
		ORDER BY action ASC
	`

	rows, err := r.db.QueryContext(ctx, query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.Action
	for rows.Next() {
		var action domain.Action
		if err := rows.Scan(
			&action.ID,
			&action.ResourceID,
			&action.Name,
			&action.Description,
			&action.DeprecatedAt,
			&action.CreatedAt,
			&action.UpdatedAt,
		); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}

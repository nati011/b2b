package supplier

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"marketplace/internal/core/supplier/domain"
	supplierservice "marketplace/internal/core/supplier/service"
	"marketplace/pkg/pagination"
)

var (
	// ErrSupplierNotFound is imported from the service package.
	ErrSupplierNotFound = supplierservice.ErrSupplierNotFound
)

// SupplierRepository persists suppliers.
type SupplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

// Create inserts a new supplier record.
func (r *SupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	query := `
		INSERT INTO suppliers (
			business_name,
			status,
			support_email,
			support_phone,
			is_active,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, FALSE)
		RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		supplier.BusinessName,
		string(supplier.Status),
		nullableString(supplier.SupportEmail),
		nullableString(supplier.SupportPhone),
		supplier.IsActive,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	).Scan(&supplier.ID)
}

// Update modifies an existing supplier record.
func (r *SupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	query := `
		UPDATE suppliers
		SET business_name = $2,
		    status = $3,
		    support_email = $4,
		    support_phone = $5,
		    is_active = $6,
		    last_modified = $7
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		supplier.ID,
		supplier.BusinessName,
		string(supplier.Status),
		nullableString(supplier.SupportEmail),
		nullableString(supplier.SupportPhone),
		supplier.IsActive,
		supplier.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSupplierNotFound
	}

	return nil
}

// FindByID fetches a supplier by identifier.
func (r *SupplierRepository) FindByID(ctx context.Context, id int64) (*domain.Supplier, error) {
	query := `
		SELECT id, business_name, status, support_email, support_phone, is_active, created_date, last_modified
		FROM suppliers
		WHERE id = $1 AND is_deleted = FALSE
	`
	supplier, err := scanSupplier(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}
	return supplier, nil
}

// FindByEmail fetches a supplier by support email address.
func (r *SupplierRepository) FindByEmail(ctx context.Context, email string) (*domain.Supplier, error) {
	query := `
		SELECT id, business_name, status, support_email, support_phone, is_active, created_date, last_modified
		FROM suppliers
		WHERE support_email = $1 AND is_deleted = FALSE
	`
	supplier, err := scanSupplier(r.db.QueryRowContext(ctx, query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}
	return supplier, nil
}

// FindByPhone fetches a supplier by support phone number.
func (r *SupplierRepository) FindByPhone(ctx context.Context, phone string) (*domain.Supplier, error) {
	query := `
		SELECT id, business_name, status, support_email, support_phone, is_active, created_date, last_modified
		FROM suppliers
		WHERE support_phone = $1 AND is_deleted = FALSE
	`
	supplier, err := scanSupplier(r.db.QueryRowContext(ctx, query, phone))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}
	return supplier, nil
}

// Delete performs a soft delete on the supplier.
func (r *SupplierRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE suppliers
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
		return ErrSupplierNotFound
	}

	return nil
}

// FindAll returns a paginated list of suppliers.
func (r *SupplierRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Supplier], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM suppliers WHERE is_deleted = FALSE`).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Supplier]{}, err
	}

	query := `
		SELECT id, business_name, status, support_email, support_phone, is_active, created_date, last_modified
		FROM suppliers
		WHERE is_deleted = FALSE
		ORDER BY created_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.Supplier]{}, err
	}
	defer rows.Close()

	var suppliers []*domain.Supplier
	for rows.Next() {
		supplier, err := scanSupplier(rows)
		if err != nil {
			return pagination.PageResult[*domain.Supplier]{}, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.Supplier]{}, err
	}

	return pagination.NewPageResult(suppliers, total, pageReq), nil
}

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

type supplierScanner interface {
	Scan(dest ...interface{}) error
}

func scanSupplier(scanner supplierScanner) (*domain.Supplier, error) {
	var supplier domain.Supplier
	var status, supportEmail, supportPhone sql.NullString
	if err := scanner.Scan(
		&supplier.ID,
		&supplier.BusinessName,
		&status,
		&supportEmail,
		&supportPhone,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	); err != nil {
		return nil, err
	}

	applySupplierNulls(&supplier, status, supportEmail, supportPhone)
	return &supplier, nil
}

func applySupplierNulls(supplier *domain.Supplier, status, supportEmail, supportPhone sql.NullString) {
	if status.Valid {
		supplier.Status = domain.SupplierStatus(status.String)
	}
	if supportEmail.Valid {
		supplier.SupportEmail = supportEmail.String
	}
	if supportPhone.Valid {
		supplier.SupportPhone = supportPhone.String
	}
}


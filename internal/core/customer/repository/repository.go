package customer

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"marketplace/internal/core/customer/domain"
	customerservice "marketplace/internal/core/customer/service"
	"marketplace/pkg/pagination"
)

var (
	// ErrCustomerNotFound is imported from the service package.
	ErrCustomerNotFound = customerservice.ErrCustomerNotFound
)

// CustomerRepository persists customers.
type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create inserts a new customer record.
func (r *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (
			full_name,
			status,
			city,
			region,
			woreda,
			phone_number,
			email,
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
		customer.FullName,
		string(customer.Status),
		nullableString(customer.City),
		nullableString(customer.Region),
		nullableString(customer.Woreda),
		nullableString(customer.PhoneNumber),
		nullableString(customer.Email),
		customer.IsActive,
		customer.CreatedAt,
		customer.UpdatedAt,
	).Scan(&customer.ID)
}

// Update modifies an existing customer record.
func (r *CustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers
		SET full_name = $2,
		    status = $3,
		    city = $4,
		    region = $5,
		    woreda = $6,
		    phone_number = $7,
		    email = $8,
		    is_active = $9,
		    last_modified = $10
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		customer.ID,
		customer.FullName,
		string(customer.Status),
		nullableString(customer.City),
		nullableString(customer.Region),
		nullableString(customer.Woreda),
		nullableString(customer.PhoneNumber),
		nullableString(customer.Email),
		customer.IsActive,
		customer.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCustomerNotFound
	}

	return nil
}

// FindByID fetches a customer by identifier.
func (r *CustomerRepository) FindByID(ctx context.Context, id int64) (*domain.Customer, error) {
	query := `
		SELECT id, full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified
		FROM customers
		WHERE id = $1 AND is_deleted = FALSE
	`
	customer, err := scanCustomer(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return customer, nil
}

// FindByEmail fetches a customer by email address.
func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	query := `
		SELECT id, full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified
		FROM customers
		WHERE email = $1 AND is_deleted = FALSE
	`
	customer, err := scanCustomer(r.db.QueryRowContext(ctx, query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return customer, nil
}

// FindByPhoneNumber fetches a customer by phone number.
func (r *CustomerRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Customer, error) {
	query := `
		SELECT id, full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified
		FROM customers
		WHERE phone_number = $1 AND is_deleted = FALSE
	`
	customer, err := scanCustomer(r.db.QueryRowContext(ctx, query, phoneNumber))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return customer, nil
}

// Delete performs a soft delete on the customer.
func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE customers
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
		return ErrCustomerNotFound
	}

	return nil
}

// FindAll returns a paginated list of customers.
func (r *CustomerRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Customer], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE is_deleted = FALSE`).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Customer]{}, err
	}

	query := `
		SELECT id, full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified
		FROM customers
		WHERE is_deleted = FALSE
		ORDER BY created_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return pagination.PageResult[*domain.Customer]{}, err
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		customer, err := scanCustomer(rows)
		if err != nil {
			return pagination.PageResult[*domain.Customer]{}, err
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.Customer]{}, err
	}

	return pagination.NewPageResult(customers, total, pageReq), nil
}

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

type customerScanner interface {
	Scan(dest ...interface{}) error
}

func scanCustomer(scanner customerScanner) (*domain.Customer, error) {
	var customer domain.Customer
	var city, region, woreda, phoneNumber, email sql.NullString
	var status sql.NullString
	if err := scanner.Scan(
		&customer.ID,
		&customer.FullName,
		&status,
		&city,
		&region,
		&woreda,
		&phoneNumber,
		&email,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	); err != nil {
		return nil, err
	}

	applyCustomerNulls(&customer, status, city, region, woreda, phoneNumber, email)
	return &customer, nil
}

func applyCustomerNulls(customer *domain.Customer, status, city, region, woreda, phoneNumber, email sql.NullString) {
	if status.Valid {
		customer.Status = domain.CustomerStatus(status.String)
	}
	if city.Valid {
		customer.City = city.String
	}
	if region.Valid {
		customer.Region = region.String
	}
	if woreda.Valid {
		customer.Woreda = woreda.String
	}
	if phoneNumber.Valid {
		customer.PhoneNumber = phoneNumber.String
	}
	if email.Valid {
		customer.Email = email.String
	}
}

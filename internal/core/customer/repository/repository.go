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
	return r.FindAllWithSupplierFilter(ctx, pageReq, 0, "")
}

// FindAllWithSupplierFilter returns a paginated list of customers, optionally filtered by supplier.
// If supplierID > 0, only returns customers who have placed orders containing products from that supplier.
// If search is provided, filters customers by name, email, or phone number.
func (r *CustomerRepository) FindAllWithSupplierFilter(ctx context.Context, pageReq pagination.PageRequest, supplierID int64, search string) (pagination.PageResult[*domain.Customer], error) {
	var countQuery string
	var query string
	var countArgs []interface{}
	var queryArgs []interface{}

	// Build search pattern if search is provided
	searchPattern := ""
	if strings.TrimSpace(search) != "" {
		searchPattern = "%" + strings.ToLower(strings.TrimSpace(search)) + "%"
	}

	if supplierID > 0 {
		// Filter customers who have orders with products from this supplier
		if searchPattern != "" {
			// With search: supplier is $1, search is $2, limit is $3, offset is $4
			countQuery = `
				SELECT COUNT(DISTINCT c.id)
				FROM customers c
				INNER JOIN orders o ON o.customer_id = c.id AND o.is_deleted = FALSE
				WHERE c.is_deleted = FALSE
				  AND EXISTS (
					SELECT 1 
					FROM jsonb_array_elements(o.cart_snapshot) AS item
					JOIN products p ON (item->>'product_id')::int = p.id
					WHERE p.supplier_id = $1 AND p.is_deleted = FALSE
				  ) AND (
					LOWER(c.full_name) LIKE $2 OR
					LOWER(c.email) LIKE $2 OR
					LOWER(c.phone_number) LIKE $2
				  )`
			query = `
				SELECT DISTINCT c.id, c.full_name, c.status, c.city, c.region, c.woreda, c.phone_number, c.email, c.is_active, c.created_date, c.last_modified
				FROM customers c
				INNER JOIN orders o ON o.customer_id = c.id AND o.is_deleted = FALSE
				WHERE c.is_deleted = FALSE
				  AND EXISTS (
					SELECT 1 
					FROM jsonb_array_elements(o.cart_snapshot) AS item
					JOIN products p ON (item->>'product_id')::int = p.id
					WHERE p.supplier_id = $1 AND p.is_deleted = FALSE
				  ) AND (
					LOWER(c.full_name) LIKE $2 OR
					LOWER(c.email) LIKE $2 OR
					LOWER(c.phone_number) LIKE $2
				  )
				ORDER BY c.id ASC
				LIMIT $3 OFFSET $4`
			countArgs = []interface{}{supplierID, searchPattern}
			queryArgs = []interface{}{supplierID, searchPattern, pageReq.Limit, pageReq.Offset}
		} else {
			// Without search: supplier is $1, limit is $2, offset is $3
			countQuery = `
				SELECT COUNT(DISTINCT c.id)
				FROM customers c
				INNER JOIN orders o ON o.customer_id = c.id AND o.is_deleted = FALSE
				WHERE c.is_deleted = FALSE
				  AND EXISTS (
					SELECT 1 
					FROM jsonb_array_elements(o.cart_snapshot) AS item
					JOIN products p ON (item->>'product_id')::int = p.id
					WHERE p.supplier_id = $1 AND p.is_deleted = FALSE
				  )`
			query = `
				SELECT DISTINCT c.id, c.full_name, c.status, c.city, c.region, c.woreda, c.phone_number, c.email, c.is_active, c.created_date, c.last_modified
				FROM customers c
				INNER JOIN orders o ON o.customer_id = c.id AND o.is_deleted = FALSE
				WHERE c.is_deleted = FALSE
				  AND EXISTS (
					SELECT 1 
					FROM jsonb_array_elements(o.cart_snapshot) AS item
					JOIN products p ON (item->>'product_id')::int = p.id
					WHERE p.supplier_id = $1 AND p.is_deleted = FALSE
				  )
				ORDER BY c.id ASC
				LIMIT $2 OFFSET $3`
			countArgs = []interface{}{supplierID}
			queryArgs = []interface{}{supplierID, pageReq.Limit, pageReq.Offset}
		}
	} else {
		// No supplier filter - return all customers
		if searchPattern != "" {
			// With search: search is $1, limit is $2, offset is $3
			countQuery = `SELECT COUNT(*) FROM customers WHERE is_deleted = FALSE AND (
				LOWER(full_name) LIKE $1 OR
				LOWER(email) LIKE $1 OR
				LOWER(phone_number) LIKE $1
			)`
			query = `
				SELECT id, full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified
				FROM customers
				WHERE is_deleted = FALSE AND (
					LOWER(full_name) LIKE $1 OR
					LOWER(email) LIKE $1 OR
					LOWER(phone_number) LIKE $1
				)
				ORDER BY id ASC
				LIMIT $2 OFFSET $3`
			countArgs = []interface{}{searchPattern}
			queryArgs = []interface{}{searchPattern, pageReq.Limit, pageReq.Offset}
		} else {
			// Without search: limit is $1, offset is $2
			countQuery = `SELECT COUNT(*) FROM customers WHERE is_deleted = FALSE`
			query = `
				SELECT id, full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified
				FROM customers
				WHERE is_deleted = FALSE
				ORDER BY id ASC
				LIMIT $1 OFFSET $2`
			countArgs = []interface{}{}
			queryArgs = []interface{}{pageReq.Limit, pageReq.Offset}
		}
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return pagination.PageResult[*domain.Customer]{}, err
	}

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
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

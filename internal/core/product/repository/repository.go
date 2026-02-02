package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"marketplace/internal/core/product/domain"
	productservice "marketplace/internal/core/product/service"
)

var (
	// ErrProductNotFound is imported from the service package.
	ErrProductNotFound = productservice.ErrProductNotFound
)

// Repository persists products and category mappings.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a product and its category mappings.
func (r *Repository) Create(ctx context.Context, product *domain.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO products (
			name,
			description,
			external_id,
			attributes,
			unit,
			is_active,
			supplier_id,
			price,
			total_quantity,
			reserved_quantity,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, FALSE)
		RETURNING id, available_quantity
	`

	if err = tx.QueryRowContext(
		ctx,
		query,
		product.Name,
		nullableString(product.Description),
		nullableString(product.ExternalID),
		product.Attributes,
		nullableString(product.Unit),
		product.IsActive,
		product.SupplierID,
		product.Price,
		product.TotalQuantity,
		product.ReservedQuantity,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID, &product.AvailableQuantity); err != nil {
		return err
	}

	if err = r.replaceCategoriesTx(ctx, tx, product.ID, product.CategoryIDs); err != nil {
		return err
	}

	return tx.Commit()
}

// Update modifies a product and its category mappings.
func (r *Repository) Update(ctx context.Context, product *domain.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		UPDATE products
		SET name = $2,
		    description = $3,
		    external_id = $4,
		    attributes = $5,
		    unit = $6,
		    is_active = $7,
		    price = $8,
		    total_quantity = $9,
		    reserved_quantity = $10,
		    last_modified = $11
		WHERE id = $1 AND is_deleted = FALSE
		RETURNING available_quantity
	`

	var available int
	result := tx.QueryRowContext(
		ctx,
		query,
		product.ID,
		product.Name,
		nullableString(product.Description),
		nullableString(product.ExternalID),
		product.Attributes,
		nullableString(product.Unit),
		product.IsActive,
		product.Price,
		product.TotalQuantity,
		product.ReservedQuantity,
		product.UpdatedAt,
	)
	if err = result.Scan(&available); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	product.AvailableQuantity = available

	if err = r.replaceCategoriesTx(ctx, tx, product.ID, product.CategoryIDs); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete soft-deletes a product and its category mappings.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, `UPDATE products SET is_deleted = TRUE, last_modified = NOW() WHERE id = $1 AND is_deleted = FALSE`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrProductNotFound
	}

	if _, err = tx.ExecContext(ctx, `UPDATE product_categories SET is_deleted = TRUE, last_modified = NOW() WHERE product_id = $1 AND is_deleted = FALSE`, id); err != nil {
		return err
	}

	return tx.Commit()
}

// FindByID retrieves a product and its category IDs.
func (r *Repository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, description, external_id, attributes, unit, is_active, supplier_id,
		       price, total_quantity, reserved_quantity, available_quantity, created_date, last_modified
		FROM products
		WHERE id = $1 AND is_deleted = FALSE
	`

	product, err := scanProduct(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	categoryIDs, err := r.fetchCategoryIDs(ctx, product.ID)
	if err != nil {
		return nil, err
	}
	product.CategoryIDs = categoryIDs
	return product, nil
}

// List retrieves products with optional filters.
func (r *Repository) List(ctx context.Context, query productservice.ListQuery) ([]*domain.Product, int, error) {
	whereClause, args := buildListFilters(query)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products p %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit, offset := resolvePagination(query.Limit, query.Offset)

	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT p.id, p.name, p.description, p.external_id, p.attributes, p.unit, p.is_active,
		       p.supplier_id, p.price, p.total_quantity, p.reserved_quantity, p.available_quantity,
		       p.created_date, p.last_modified, pc.category_id
		FROM products p
		LEFT JOIN product_categories pc
		  ON pc.product_id = p.id AND pc.is_deleted = FALSE
		%s
		ORDER BY p.created_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products, err := collectProducts(rows)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *Repository) fetchCategoryIDs(ctx context.Context, productID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT category_id FROM product_categories WHERE product_id = $1 AND is_deleted = FALSE`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *Repository) replaceCategoriesTx(ctx context.Context, tx *sql.Tx, productID int64, categoryIDs []int64) error {
	if _, err := tx.ExecContext(ctx, `UPDATE product_categories SET is_deleted = TRUE, last_modified = NOW() WHERE product_id = $1 AND is_deleted = FALSE`, productID); err != nil {
		return err
	}

	if len(categoryIDs) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO product_categories (product_id, category_id, created_date, last_modified, is_deleted)
		VALUES ($1, $2, $3, $4, FALSE)
	`
	now := time.Now()
	for _, categoryID := range categoryIDs {
		if _, err := tx.ExecContext(ctx, insertQuery, productID, categoryID, now, now); err != nil {
			return err
		}
	}
	return nil
}

type productScanner interface {
	Scan(dest ...interface{}) error
}

func scanProduct(scanner productScanner) (*domain.Product, error) {
	var product domain.Product
	var description sql.NullString
	var externalID sql.NullString
	var attributes []byte
	var unit sql.NullString
	var price sql.NullFloat64
	var total sql.NullInt64
	var reserved sql.NullInt64
	var available sql.NullInt64

	if err := scanner.Scan(
		&product.ID,
		&product.Name,
		&description,
		&externalID,
		&attributes,
		&unit,
		&product.IsActive,
		&product.SupplierID,
		&price,
		&total,
		&reserved,
		&available,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if description.Valid {
		product.Description = description.String
	}
	if externalID.Valid {
		product.ExternalID = externalID.String
	}
	if len(attributes) > 0 {
		product.Attributes = attributes
	}
	if unit.Valid {
		product.Unit = unit.String
	}
	if price.Valid {
		product.Price = &price.Float64
	}
	if total.Valid {
		product.TotalQuantity = int(total.Int64)
	}
	if reserved.Valid {
		product.ReservedQuantity = int(reserved.Int64)
	}
	if available.Valid {
		product.AvailableQuantity = int(available.Int64)
	}

	return &product, nil
}

func scanProductWithCategory(scanner productScanner) (*domain.Product, *int64, error) {
	var product domain.Product
	var description sql.NullString
	var externalID sql.NullString
	var attributes []byte
	var unit sql.NullString
	var price sql.NullFloat64
	var total sql.NullInt64
	var reserved sql.NullInt64
	var available sql.NullInt64
	var categoryID sql.NullInt64

	if err := scanner.Scan(
		&product.ID,
		&product.Name,
		&description,
		&externalID,
		&attributes,
		&unit,
		&product.IsActive,
		&product.SupplierID,
		&price,
		&total,
		&reserved,
		&available,
		&product.CreatedAt,
		&product.UpdatedAt,
		&categoryID,
	); err != nil {
		return nil, nil, err
	}

	if description.Valid {
		product.Description = description.String
	}
	if externalID.Valid {
		product.ExternalID = externalID.String
	}
	if len(attributes) > 0 {
		product.Attributes = attributes
	}
	if unit.Valid {
		product.Unit = unit.String
	}
	if price.Valid {
		product.Price = &price.Float64
	}
	if total.Valid {
		product.TotalQuantity = int(total.Int64)
	}
	if reserved.Valid {
		product.ReservedQuantity = int(reserved.Int64)
	}
	if available.Valid {
		product.AvailableQuantity = int(available.Int64)
	}

	if categoryID.Valid {
		value := categoryID.Int64
		return &product, &value, nil
	}
	return &product, nil, nil
}

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

func nullableInt(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

func appendUnique(values []int64, value int64) []int64 {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func buildListFilters(query productservice.ListQuery) (string, []interface{}) {
	filters := []string{"p.is_deleted = FALSE"}
	args := []interface{}{}

	if query.SupplierID > 0 {
		args = append(args, query.SupplierID)
		filters = append(filters, fmt.Sprintf("p.supplier_id = $%d", len(args)))
	}

	if query.IsActive != nil {
		args = append(args, *query.IsActive)
		filters = append(filters, fmt.Sprintf("p.is_active = $%d", len(args)))
	}

	if query.CategoryID > 0 {
		args = append(args, query.CategoryID)
		filters = append(filters, fmt.Sprintf("EXISTS (SELECT 1 FROM product_categories pc WHERE pc.product_id = p.id AND pc.category_id = $%d AND pc.is_deleted = FALSE)", len(args)))
	}

	return "WHERE " + strings.Join(filters, " AND "), args
}

func resolvePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func collectProducts(rows *sql.Rows) ([]*domain.Product, error) {
	byID := make(map[int64]*domain.Product)
	orderIDs := make([]int64, 0)

	for rows.Next() {
		product, categoryID, err := scanProductWithCategory(rows)
		if err != nil {
			return nil, err
		}
		if existing, ok := byID[product.ID]; ok {
			if categoryID != nil {
				existing.CategoryIDs = appendUnique(existing.CategoryIDs, *categoryID)
			}
			continue
		}
		if categoryID != nil {
			product.CategoryIDs = []int64{*categoryID}
		} else {
			product.CategoryIDs = []int64{}
		}
		byID[product.ID] = product
		orderIDs = append(orderIDs, product.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	products := make([]*domain.Product, 0, len(orderIDs))
	for _, id := range orderIDs {
		if product, ok := byID[id]; ok {
			products = append(products, product)
		}
	}

	return products, nil
}

// RecordPriceChange records a price change in the price_history table.
func (r *Repository) RecordPriceChange(ctx context.Context, productID int64, oldPrice, newPrice *float64, userID *int64, userEmail, reason string) error {
	query := `
		INSERT INTO price_history (
			product_id,
			old_price,
			new_price,
			changed_by_user_id,
			changed_by_user_email,
			reason,
			created_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	var oldPriceVal interface{}
	if oldPrice != nil {
		oldPriceVal = *oldPrice
	}

	var newPriceVal interface{}
	if newPrice != nil {
		newPriceVal = *newPrice
	}

	var userIDVal interface{}
	if userID != nil {
		userIDVal = *userID
	}

	var userEmailVal interface{}
	if userEmail != "" {
		userEmailVal = userEmail
	}

	var reasonVal interface{}
	if reason != "" {
		reasonVal = reason
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		productID,
		oldPriceVal,
		newPriceVal,
		userIDVal,
		userEmailVal,
		reasonVal,
	)

	return err
}

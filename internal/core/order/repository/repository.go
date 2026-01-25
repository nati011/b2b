package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"marketplace/internal/core/order/domain"
	orderservice "marketplace/internal/core/order/service"

	goodmoney "github.com/the-nucleus-project/good_money"
)

var (
	// ErrOrderNotFound is imported from the service package.
	ErrOrderNotFound = orderservice.ErrOrderNotFound
)

// Repository persists orders and line items.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts an order and its line items.
func (r *Repository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var totalValue interface{}
	if order.Total != nil {
		var valErr error
		totalValue, valErr = order.Total.Value()
		if valErr != nil {
			return valErr
		}
	}

	query := `
		INSERT INTO orders (
			customer_id,
			status,
			payment_status,
			delivery_status,
			confirmation_status,
			total,
			customer_snapshot,
			shipping_address_snapshot,
			billing_address_snapshot,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, FALSE)
		RETURNING id
	`

	if err = tx.QueryRowContext(
		ctx,
		query,
		order.CustomerID,
		string(order.Status),
		nullableString(order.PaymentStatus),
		nullableString(order.DeliveryStatus),
		nullableString(order.ConfirmationStatus),
		totalValue,
		order.CustomerSnapshot,
		order.ShippingAddressSnapshot,
		order.BillingAddressSnapshot,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID); err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO o_items (
			order_id,
			product_id,
			quantity,
			price,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, FALSE)
	`

	for i := range order.Items {
		item := &order.Items[i]
		item.OrderID = order.ID
		item.CreatedAt = time.Now()
		item.UpdatedAt = item.CreatedAt
		var priceValue interface{}
		if item.Price != nil {
			var valErr error
			priceValue, valErr = item.Price.Value()
			if valErr != nil {
				return valErr
			}
		}
		if _, err = tx.ExecContext(
			ctx,
			itemQuery,
			order.ID,
			item.ProductID,
			item.Quantity,
			priceValue,
			item.CreatedAt,
			item.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Update modifies an existing order record.
func (r *Repository) Update(ctx context.Context, order *domain.Order) error {
	var totalValue interface{}
	if order.Total != nil {
		var valErr error
		totalValue, valErr = order.Total.Value()
		if valErr != nil {
			return valErr
		}
	}

	query := `
		UPDATE orders
		SET status = $2,
		    payment_status = $3,
		    delivery_status = $4,
		    confirmation_status = $5,
		    total = $6,
		    customer_snapshot = $7,
		    shipping_address_snapshot = $8,
		    billing_address_snapshot = $9,
		    last_modified = $10
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		order.ID,
		string(order.Status),
		nullableString(order.PaymentStatus),
		nullableString(order.DeliveryStatus),
		nullableString(order.ConfirmationStatus),
		totalValue,
		order.CustomerSnapshot,
		order.ShippingAddressSnapshot,
		order.BillingAddressSnapshot,
		order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrOrderNotFound
	}

	return nil
}

// UpdateStatus updates only the order status.
func (r *Repository) UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE orders SET status = $2, last_modified = NOW() WHERE id = $1 AND is_deleted = FALSE`,
		orderID,
		string(status),
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// FindByID retrieves an order and its line items.
func (r *Repository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `
		SELECT id, customer_id, status, payment_status, delivery_status, confirmation_status,
		       total, customer_snapshot, shipping_address_snapshot, billing_address_snapshot,
		       created_date, last_modified
		FROM orders
		WHERE id = $1 AND is_deleted = FALSE
	`

	order, err := scanOrder(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	items, err := r.findItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

// ListByCustomer retrieves customer orders with optional status filter.
func (r *Repository) ListByCustomer(ctx context.Context, query orderservice.CustomerOrderQuery) ([]*domain.Order, int, error) {
	filters := []string{"is_deleted = FALSE"}
	args := []interface{}{}

	if query.CustomerID > 0 {
		args = append(args, query.CustomerID)
		filters = append(filters, fmt.Sprintf("customer_id = $%d", len(args)))
	}

	status := strings.TrimSpace(query.Status)
	if status != "" && strings.ToUpper(status) != "ALL" {
		args = append(args, strings.ToLower(status))
		filters = append(filters, fmt.Sprintf("status = $%d", len(args)))
	}

	whereClause := "WHERE " + strings.Join(filters, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT id, customer_id, status, payment_status, delivery_status, confirmation_status,
		       total, customer_snapshot, shipping_address_snapshot, billing_address_snapshot,
		       created_date, last_modified
		FROM orders
		%s
		ORDER BY created_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *Repository) findItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	query := `
		SELECT order_id, product_id, quantity, price, created_date, last_modified
		FROM o_items
		WHERE order_id = $1 AND is_deleted = FALSE
		ORDER BY created_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var price goodmoney.Money
		if err := rows.Scan(
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&price,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if price.Currency() != "" {
			item.Price = &price
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type orderScanner interface {
	Scan(dest ...interface{}) error
}

func scanOrder(scanner orderScanner) (*domain.Order, error) {
	var order domain.Order
	var status sql.NullString
	var paymentStatus sql.NullString
	var deliveryStatus sql.NullString
	var confirmationStatus sql.NullString
	var total goodmoney.Money
	var customerSnapshot []byte
	var shippingSnapshot []byte
	var billingSnapshot []byte

	if err := scanner.Scan(
		&order.ID,
		&order.CustomerID,
		&status,
		&paymentStatus,
		&deliveryStatus,
		&confirmationStatus,
		&total,
		&customerSnapshot,
		&shippingSnapshot,
		&billingSnapshot,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if status.Valid {
		order.Status = domain.OrderStatus(status.String)
	} else {
		order.Status = domain.OrderStatusPending
	}
	if paymentStatus.Valid {
		order.PaymentStatus = paymentStatus.String
	}
	if deliveryStatus.Valid {
		order.DeliveryStatus = deliveryStatus.String
	}
	if confirmationStatus.Valid {
		order.ConfirmationStatus = confirmationStatus.String
	}
	if total.Currency() != "" {
		order.Total = &total
	}
	if len(customerSnapshot) > 0 {
		order.CustomerSnapshot = customerSnapshot
	}
	if len(shippingSnapshot) > 0 {
		order.ShippingAddressSnapshot = shippingSnapshot
	}
	if len(billingSnapshot) > 0 {
		order.BillingAddressSnapshot = billingSnapshot
	}

	return &order, nil
}

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"marketplace/internal/core/order/domain"
	orderservice "marketplace/internal/core/order/service"
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

	// Serialize items to JSON for cart_snapshot
	cartSnapshot, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to serialize cart items: %w", err)
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
			cart_snapshot,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, FALSE)
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
		order.Total,
		order.CustomerSnapshot,
		cartSnapshot,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID); err != nil {
		return err
	}

	return tx.Commit()
}

// Update modifies an existing order record.
func (r *Repository) Update(ctx context.Context, order *domain.Order) error {
	// Serialize items to JSON for cart_snapshot
	cartSnapshot, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to serialize cart items: %w", err)
	}

	query := `
		UPDATE orders
		SET status = $2,
		    payment_status = $3,
		    delivery_status = $4,
		    confirmation_status = $5,
		    total = $6,
		    customer_snapshot = $7,
		    cart_snapshot = $8,
		    last_modified = $9
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
		order.Total,
		order.CustomerSnapshot,
		cartSnapshot,
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
		       total, customer_snapshot, cart_snapshot,
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
		       total, customer_snapshot, cart_snapshot,
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


type orderScanner interface {
	Scan(dest ...interface{}) error
}

func scanOrder(scanner orderScanner) (*domain.Order, error) {
	var order domain.Order
	var status sql.NullString
	var paymentStatus sql.NullString
	var deliveryStatus sql.NullString
	var confirmationStatus sql.NullString
	var total sql.NullFloat64
	var customerSnapshot []byte
	var cartSnapshot []byte

	if err := scanner.Scan(
		&order.ID,
		&order.CustomerID,
		&status,
		&paymentStatus,
		&deliveryStatus,
		&confirmationStatus,
		&total,
		&customerSnapshot,
		&cartSnapshot,
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
	if total.Valid {
		order.Total = &total.Float64
	}
	if len(customerSnapshot) > 0 {
		order.CustomerSnapshot = customerSnapshot
	}
	if len(cartSnapshot) > 0 {
		order.CartSnapshot = cartSnapshot
		// Deserialize cart_snapshot to Items
		var items []domain.OrderItem
		if err := json.Unmarshal(cartSnapshot, &items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cart snapshot: %w", err)
		}
		order.Items = items
	} else {
		order.Items = []domain.OrderItem{}
	}

	return &order, nil
}

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

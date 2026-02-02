package repository

import (
	"context"
	roleDomain "marketplace/internal/infra/authz/role/domain"
	"marketplace/internal/infra/user/domain"
	userservice "marketplace/internal/infra/user/service"
	"marketplace/pkg/pagination"
	"database/sql"
	"errors"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// isInvalidUUIDError checks if the error is a PostgreSQL invalid UUID error
func isInvalidUUIDError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "22P02") || strings.Contains(errStr, "invalid input syntax for type uuid")
}

// Create inserts a new user into the database
func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, external_id, email, phone_number, name, user_type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	var email, phoneNumber *string
	if !user.Email.IsEmpty() {
		emailStr := user.Email.String()
		email = &emailStr
	}
	if !user.PhoneNumber.IsEmpty() {
		phoneStr := user.PhoneNumber.String()
		phoneNumber = &phoneStr
	}

	var externalID *string
	if user.ExternalID != "" {
		externalID = &user.ExternalID
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		externalID,
		email,
		phoneNumber,
		user.Name,
		string(user.UserType),
		string(user.Status),
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

// FindByID retrieves a user by ID
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, external_id, email, phone_number, name, user_type, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	var email, phoneNumber, externalID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&externalID,
		&email,
		&phoneNumber,
		&user.Name,
		&user.UserType,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUIDError(err) {
			return nil, userservice.ErrUserNotFound
		}
		return nil, err
	}

	if externalID.Valid {
		user.ExternalID = externalID.String
	}

	if email.Valid && email.String != "" {
		emailObj, err := domain.NewEmail(email.String)
		if err != nil {
			return nil, err
		}
		user.Email = emailObj
	}

	if phoneNumber.Valid && phoneNumber.String != "" {
		phoneObj, err := domain.NewPhoneNumber(phoneNumber.String, nil)
		if err != nil {
			return nil, err
		}
		user.PhoneNumber = phoneObj
	}

	// Fetch role IDs for the user
	roleIDs, err := r.fetchUserRoleIDs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.RoleIDs = roleIDs

	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, external_id, email, phone_number, name, user_type, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	var emailVal, phoneNumber, externalID sql.NullString
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&externalID,
		&emailVal,
		&phoneNumber,
		&user.Name,
		&user.UserType,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userservice.ErrUserNotFound
		}
		return nil, err
	}

	if externalID.Valid {
		user.ExternalID = externalID.String
	}

	if emailVal.Valid && emailVal.String != "" {
		emailObj, err := domain.NewEmail(emailVal.String)
		if err != nil {
			return nil, err
		}
		user.Email = emailObj
	}

	if phoneNumber.Valid && phoneNumber.String != "" {
		phoneObj, err := domain.NewPhoneNumber(phoneNumber.String, nil)
		if err != nil {
			return nil, err
		}
		user.PhoneNumber = phoneObj
	}

	// Fetch role IDs for the user
	roleIDs, err := r.fetchUserRoleIDs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.RoleIDs = roleIDs

	return &user, nil
}

// FindByPhoneNumber retrieves a user by phone number
func (r *Repository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	query := `
		SELECT id, external_id, email, phone_number, name, user_type, status, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`

	var user domain.User
	var email, phoneNumberVal, externalID sql.NullString
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&externalID,
		&email,
		&phoneNumberVal,
		&user.Name,
		&user.UserType,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userservice.ErrUserNotFound
		}
		return nil, err
	}

	if externalID.Valid {
		user.ExternalID = externalID.String
	}

	if email.Valid && email.String != "" {
		emailObj, err := domain.NewEmail(email.String)
		if err != nil {
			return nil, err
		}
		user.Email = emailObj
	}

	if phoneNumberVal.Valid && phoneNumberVal.String != "" {
		phoneObj, err := domain.NewPhoneNumber(phoneNumberVal.String, nil)
		if err != nil {
			return nil, err
		}
		user.PhoneNumber = phoneObj
	}

	// Fetch role IDs for the user
	roleIDs, err := r.fetchUserRoleIDs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.RoleIDs = roleIDs

	return &user, nil
}

// Update updates an existing user in the database
func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET external_id = $2, email = $3, phone_number = $4, name = $5, user_type = $6, status = $7, updated_at = $8
		WHERE id = $1
	`

	var email, phoneNumber *string
	if !user.Email.IsEmpty() {
		emailStr := user.Email.String()
		email = &emailStr
	}
	if !user.PhoneNumber.IsEmpty() {
		phoneStr := user.PhoneNumber.String()
		phoneNumber = &phoneStr
	}

	var externalID *string
	if user.ExternalID != "" {
		externalID = &user.ExternalID
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		externalID,
		email,
		phoneNumber,
		user.Name,
		string(user.UserType),
		string(user.Status),
		user.UpdatedAt,
	)
	return err
}

// Delete removes a user from the database
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if isInvalidUUIDError(err) {
			return userservice.ErrUserNotFound
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return userservice.ErrUserNotFound
	}

	return nil
}

// FindPage retrieves a paginated list of users
func (r *Repository) FindPage(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.User], error) {
	return r.FindPageWithSupplierEmail(ctx, pageReq, "")
}

// FindPageWithSupplierEmail retrieves a paginated list of users, optionally filtered by supplier email.
// If supplierEmail is empty, returns all users. Otherwise, filters users whose email matches the supplier's support_email.
func (r *Repository) FindPageWithSupplierEmail(ctx context.Context, pageReq pagination.PageRequest, supplierEmail string) (pagination.PageResult[*domain.User], error) {
	var total int
	var countQuery string
	var countArgs []interface{}
	
	if supplierEmail != "" {
		// Filter users whose email matches a supplier's support_email
		countQuery = `
			SELECT COUNT(*) 
			FROM users u
			INNER JOIN suppliers s ON u.email = s.support_email
			WHERE s.support_email = $1 AND s.is_deleted = FALSE
		`
		countArgs = []interface{}{supplierEmail}
	} else {
		countQuery = `SELECT COUNT(*) FROM users`
		countArgs = []interface{}{}
	}
	
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return pagination.PageResult[*domain.User]{}, err
	}

	var query string
	var queryArgs []interface{}
	
	if supplierEmail != "" {
		query = `
			SELECT u.id, u.external_id, u.email, u.phone_number, u.name, u.user_type, u.status, u.created_at, u.updated_at
			FROM users u
			INNER JOIN suppliers s ON u.email = s.support_email
			WHERE s.support_email = $1 AND s.is_deleted = FALSE
			ORDER BY u.created_at DESC
			LIMIT $2 OFFSET $3
		`
		queryArgs = []interface{}{supplierEmail, pageReq.Limit, pageReq.Offset}
	} else {
		query = `
			SELECT id, external_id, email, phone_number, name, user_type, status, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		queryArgs = []interface{}{pageReq.Limit, pageReq.Offset}
	}

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return pagination.PageResult[*domain.User]{}, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user, err := r.scanUserFromRow(ctx, rows)
		if err != nil {
			return pagination.PageResult[*domain.User]{}, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return pagination.PageResult[*domain.User]{}, err
	}

	return pagination.NewPageResult(users, total, pageReq), nil
}

// scanUserFromRow scans a user from a database row
func (r *Repository) scanUserFromRow(ctx context.Context, rows *sql.Rows) (*domain.User, error) {
	var user domain.User
	var email, phoneNumber, externalID sql.NullString
	if err := rows.Scan(
		&user.ID,
		&externalID,
		&email,
		&phoneNumber,
		&user.Name,
		&user.UserType,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if externalID.Valid {
		user.ExternalID = externalID.String
	}

	if email.Valid && email.String != "" {
		emailObj, err := domain.NewEmail(email.String)
		if err != nil {
			return nil, err
		}
		user.Email = emailObj
	}

	if phoneNumber.Valid && phoneNumber.String != "" {
		phoneObj, err := domain.NewPhoneNumber(phoneNumber.String, nil)
		if err != nil {
			return nil, err
		}
		user.PhoneNumber = phoneObj
	}

	// Fetch role IDs for the user
	roleIDs, err := r.fetchUserRoleIDs(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.RoleIDs = roleIDs

	return &user, nil
}

// Exists checks if a user exists by ID
func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

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

// AssignRole assigns a role to a user
func (r *Repository) AssignRole(ctx context.Context, userID, roleID string) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// RevokeRole revokes a role from a user
func (r *Repository) RevokeRole(ctx context.Context, userID, roleID string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// ListRoles retrieves all roles assigned to a user.
// Note: This returns full Role objects for API responses, not for domain storage.
func (r *Repository) ListRoles(ctx context.Context, userID string) ([]roleDomain.Role, error) {
	return r.fetchUserRoles(ctx, userID)
}

// fetchUserRoleIDs retrieves only role IDs for a user (for domain storage)
func (r *Repository) fetchUserRoleIDs(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT r.id
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roleIDs []string
	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		roleIDs = append(roleIDs, roleID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roleIDs, nil
}

func (r *Repository) fetchUserRoles(ctx context.Context, userID string) ([]roleDomain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []roleDomain.Role
	for rows.Next() {
		var role roleDomain.Role
		var descriptionNull sql.NullString
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&descriptionNull,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if descriptionNull.Valid {
			role.Description = descriptionNull.String
		}
		permissionIDs, err := r.fetchRolePermissionIDs(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.PermissionIDs = permissionIDs
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repository) fetchRolePermissionIDs(ctx context.Context, roleID string) ([]string, error) {
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

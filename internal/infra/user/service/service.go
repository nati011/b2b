package service

import (
	"context"
	"errors"
	roleDomain "marketplace/internal/infra/authz/role/domain"
	"marketplace/internal/infra/user/domain"
	"marketplace/pkg/logger"
	"marketplace/pkg/pagination"
)

// Business logic errors
var (
	// ErrEmailAlreadyExists indicates that a user with the given email already exists
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrPhoneNumberAlreadyExists indicates that a user with the given phone number already exists
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")

	// ErrUserNotFound indicates that a user was not found
	ErrUserNotFound = errors.New("user not found")
)

// Repository defines the interface for user repository operations.
type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	FindPage(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.User], error)
	Exists(ctx context.Context, id string) (bool, error)
	AssignRole(ctx context.Context, userID, roleID string) error
	RevokeRole(ctx context.Context, userID, roleID string) error
	ListRoles(ctx context.Context, userID string) ([]roleDomain.Role, error)
}

// RoleReader defines the interface for role lookup operations.
type RoleReader interface {
	FindByID(ctx context.Context, id string) (*roleDomain.Role, error)
	GetByName(ctx context.Context, name string) (*roleDomain.Role, error)
}

// RolePermissionChecker defines the interface for checking role permissions.
type RolePermissionChecker interface {
	HasPermission(ctx context.Context, roleID, resource, action string) (bool, error)
}

// CreateUserResult contains the result of creating a user, including optional registration token
type CreateUserResult struct {
	User              *domain.User
	RegistrationToken string // Only populated for self-service users
}

// Service handles business logic for user operations
type Service struct {
	repository               Repository
	roles                    RoleReader
	phoneValidation          *PhoneValidationService
	registrationTokenService *RegistrationTokenService
	userTypeRoleMapping      map[string]string // Maps user_type (as stored in DB) to role name
}

// NewService creates a new user service with the given repository
func NewService(repo Repository, roles RoleReader, phoneValidation *PhoneValidationService, registrationTokenService *RegistrationTokenService) *Service {
	return &Service{
		repository:               repo,
		roles:                    roles,
		phoneValidation:          phoneValidation,
		registrationTokenService: registrationTokenService,
		userTypeRoleMapping:      make(map[string]string),
	}
}

// SetUserTypeRoleMapping sets the mapping from user_type to role name
func (s *Service) SetUserTypeRoleMapping(mapping map[string]string) {
	if mapping != nil {
		s.userTypeRoleMapping = mapping
	}
}

// ValidatePhoneNumber validates a phone number using the phone validation service
func (s *Service) ValidatePhoneNumber(ctx context.Context, phoneNumber string) error {
	if s.phoneValidation == nil {
		// Fallback to basic validation if service not available
		return nil
	}
	return s.phoneValidation.ValidatePhoneNumber(ctx, phoneNumber)
}

// Create creates a new user and returns the result including optional registration token
func (s *Service) Create(ctx context.Context, externalID string, email domain.Email, phoneNumber domain.PhoneNumber, name string, userType domain.UserType) (*CreateUserResult, error) {
	// Check uniqueness: email if provided
	if !email.IsEmpty() {
		existingUser, err := s.repository.FindByEmail(ctx, email.String())
		if err == nil && existingUser != nil {
			logger.Warn("User creation failed: email already exists", "email", email.String())
			return nil, ErrEmailAlreadyExists
		}
	}

	// Check uniqueness: phone number if provided
	if !phoneNumber.IsEmpty() {
		existingUser, err := s.repository.FindByPhoneNumber(ctx, phoneNumber.String())
		if err == nil && existingUser != nil {
			logger.Warn("User creation failed: phone number already exists", "phone_number", phoneNumber.String())
			return nil, ErrPhoneNumberAlreadyExists
		}
	}

	user, err := domain.NewUser(externalID, email, phoneNumber, name, userType)
	if err != nil {
		logger.Error("User creation failed: validation error", "email", email.String(), "error", err)
		return nil, err
	}

	if err := s.repository.Create(ctx, user); err != nil {
		logger.Error("User creation failed: repository error", "email", email.String(), "error", err)
		return nil, err
	}

	// Automatically assign role based on user type if mapping is configured
	userTypeStr := string(userType)
	if roleName, ok := s.userTypeRoleMapping[userTypeStr]; ok {
		if err := s.assignRoleByName(ctx, user.ID, roleName); err != nil {
			// Log error but don't fail user creation
			logger.Warn("Failed to assign role to user", "user_id", user.ID, "user_type", userTypeStr, "role_name", roleName, "error", err)
		} else {
			logger.Info("Role assigned to user", "user_id", user.ID, "user_type", userTypeStr, "role_name", roleName)
		}
	}

	// Generate registration token
	var registrationToken string
	token, err := s.registrationTokenService.GenerateToken(ctx, user.ID)
	if err != nil {
		// Log error but don't fail user creation
		logger.Warn("Failed to generate registration token", "user_id", user.ID, "error", err)
	} else {
		registrationToken = token.Token
	}

	logger.Info("User created successfully", "user_id", user.ID, "email", email.String(), "user_type", string(userType))
	return &CreateUserResult{
		User:              user,
		RegistrationToken: registrationToken,
	}, nil
}

// assignRoleByName assigns a role to a user by looking up the role by name
func (s *Service) assignRoleByName(ctx context.Context, userID, roleName string) error {
	role, err := s.roles.GetByName(ctx, roleName)
	if err != nil {
		return err
	}
	return s.repository.AssignRole(ctx, userID, role.ID)
}

// Get retrieves a user by ID
func (s *Service) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			logger.Debug("User retrieval failed: user not found", "user_id", id)
		} else {
			logger.Error("User retrieval failed: repository error", "user_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("User retrieved successfully", "user_id", id)
	return user, nil
}

// FindByEmail retrieves a user by email address
func (s *Service) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			logger.Debug("User retrieval by email failed: user not found", "email", email)
		} else {
			logger.Error("User retrieval by email failed: repository error", "email", email, "error", err)
		}
		return nil, err
	}
	logger.Debug("User retrieved by email successfully", "user_id", user.ID, "email", email)
	return user, nil
}

// FindByPhoneNumber retrieves a user by phone number
func (s *Service) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	user, err := s.repository.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			logger.Debug("User retrieval by phone number failed: user not found", "phone_number", phoneNumber)
		} else {
			logger.Error("User retrieval by phone number failed: repository error", "phone_number", phoneNumber, "error", err)
		}
		return nil, err
	}
	logger.Debug("User retrieved by phone number successfully", "user_id", user.ID, "phone_number", phoneNumber)
	return user, nil
}

// Update updates an existing user
// externalID can be nil to preserve existing value, or a pointer to a string (including empty string to clear)
func (s *Service) Update(ctx context.Context, id string, externalID *string, email domain.Email, name string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("User update failed: user not found", "user_id", id)
		return nil, err
	}

	if user.Email.String() != email.String() {
		existingUser, err := s.repository.FindByEmail(ctx, email.String())
		if err == nil && existingUser != nil && existingUser.ID != id {
			logger.Warn("User update failed: email already exists", "user_id", id, "email", email.String())
			return nil, ErrEmailAlreadyExists
		}
	}

	// Preserve existing externalID if not provided
	externalIDValue := user.ExternalID
	if externalID != nil {
		externalIDValue = *externalID
	}

	if err := user.Update(externalIDValue, email, name); err != nil {
		logger.Error("User update failed: validation error", "user_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, user); err != nil {
		logger.Error("User update failed: repository error", "user_id", id, "error", err)
		return nil, err
	}

	logger.Info("User updated successfully", "user_id", id, "email", email.String())
	return user, nil
}

// Delete removes a user
func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("User deletion failed: user not found", "user_id", id)
		return err
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		logger.Error("User deletion failed: repository error", "user_id", id, "error", err)
		return err
	}

	logger.Info("User deleted successfully", "user_id", id)
	return nil
}

// Activate activates a user
func (s *Service) Activate(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("User activation failed: user not found", "user_id", id)
		return nil, err
	}

	if err := user.Activate(); err != nil {
		logger.Error("User activation failed: validation error", "user_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, user); err != nil {
		logger.Error("User activation failed: repository error", "user_id", id, "error", err)
		return nil, err
	}

	logger.Info("User activated successfully", "user_id", id)
	return user, nil
}

// Deactivate deactivates a user
func (s *Service) Deactivate(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("User deactivation failed: user not found", "user_id", id)
		return nil, err
	}

	if err := user.Deactivate(); err != nil {
		logger.Error("User deactivation failed: validation error", "user_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, user); err != nil {
		logger.Error("User deactivation failed: repository error", "user_id", id, "error", err)
		return nil, err
	}

	logger.Info("User deactivated successfully", "user_id", id)
	return user, nil
}

// Suspend suspends a user
func (s *Service) Suspend(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("User suspension failed: user not found", "user_id", id)
		return nil, err
	}

	if err := user.Suspend(); err != nil {
		logger.Error("User suspension failed: validation error", "user_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, user); err != nil {
		logger.Error("User suspension failed: repository error", "user_id", id, "error", err)
		return nil, err
	}

	logger.Info("User suspended successfully", "user_id", id)
	return user, nil
}

// FindPage retrieves a paginated list of users
func (s *Service) FindPage(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.User], error) {
	result, err := s.repository.FindPage(ctx, pageReq)
	if err != nil {
		logger.Error("User pagination failed: repository error", "page", pageReq.Page, "limit", pageReq.Limit, "error", err)
		return pagination.PageResult[*domain.User]{}, err
	}
	logger.Debug("User pagination completed", "page", result.Page, "total", result.Total, "items", len(result.Items))
	return result, nil
}

// AssignRole assigns a role to a user and returns the updated user.
func (s *Service) AssignRole(ctx context.Context, userID, roleID string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if _, err := s.roles.FindByID(ctx, roleID); err != nil {
		return nil, err
	}

	if err := s.repository.AssignRole(ctx, userID, roleID); err != nil {
		return nil, err
	}

	// Update domain with new role ID
	user.AssignRole(roleID)

	// Reload user to get updated role IDs
	user, err = s.repository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// RevokeRole removes a role assignment from a user.
func (s *Service) RevokeRole(ctx context.Context, userID, roleID string) (*domain.User, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if _, err := s.roles.FindByID(ctx, roleID); err != nil {
		return nil, err
	}

	if err := s.repository.RevokeRole(ctx, userID, roleID); err != nil {
		return nil, err
	}

	// Update domain by removing role ID
	user.RevokeRole(roleID)

	// Reload user to get updated role IDs
	user, err = s.repository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ListRolesForUser lists all roles assigned to a user.
func (s *Service) ListRolesForUser(ctx context.Context, userID string) ([]roleDomain.Role, error) {
	if _, err := s.repository.FindByID(ctx, userID); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			logger.Debug("List roles failed: user not found", "user_id", userID)
		} else {
			logger.Error("List roles failed: user lookup error", "user_id", userID, "error", err)
		}
		return nil, err
	}
	roles, err := s.repository.ListRoles(ctx, userID)
	if err != nil {
		logger.Error("List roles failed: repository error", "user_id", userID, "error", err)
		return nil, err
	}
	logger.Debug("Roles listed successfully", "user_id", userID, "role_count", len(roles))
	return roles, nil
}

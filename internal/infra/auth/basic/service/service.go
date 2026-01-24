package service

import (
	"context"
	"errors"

	"marketplace/internal/infra/auth/basic/domain"
	userDomain "marketplace/internal/infra/user/domain"
	httputil "marketplace/pkg/http"

	"golang.org/x/crypto/bcrypt"
)

// Permission constants - these match user module constants but are defined here to avoid circular dependencies
const (
	ResourceUsers = "users"
	ActionCreate  = "create"
)

var (
	ErrTargetUserNotFound                  = errors.New("target user not found")
	ErrSelfServiceOwnCredentialsOnly       = errors.New("self-service users can only manage their own credentials")
	ErrInsufficientPermissions             = errors.New("insufficient permissions to manage credentials for other users")
	ErrRegistrationTokenRequired           = errors.New("registration_token is required for self-service users")
	ErrUserIDRequired                      = errors.New("user_id is required")
	ErrRegistrationTokenForSelfServiceOnly = errors.New("registration token can only be used for self-service users")
	ErrCredentialNotActive                 = errors.New("credential is not active")
	ErrInvalidPassword                     = errors.New("invalid password")
	ErrPasswordMismatch                    = errors.New("password mismatch")
	ErrUsernameRequired                    = errors.New("username is required")
	ErrPasswordRequired                    = errors.New("password is required")
	ErrAuthenticationRequired              = errors.New("authentication required")
)

// UserLoader defines the interface for loading users by ID
type UserLoader interface {
	Get(ctx context.Context, id string) (*userDomain.User, error)
}

// PermissionChecker defines the interface for checking user permissions
type PermissionChecker interface {
	HasPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

// Repository defines the interface for basic auth credential storage
type Repository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Credential, error)
	Create(ctx context.Context, cred *domain.Credential) error
	Update(ctx context.Context, cred *domain.Credential) error
	Delete(ctx context.Context, id string) error
}

// Service handles business logic for basic auth credentials
type Service struct {
	repository        Repository
	userLoader        UserLoader
	permissionChecker PermissionChecker
}

// Repository exposes the repository for handler authorization checks
func (s *Service) Repository() Repository {
	return s.repository
}

// NewService creates a new basic auth service
func NewService(repository Repository, userLoader UserLoader, permissionChecker PermissionChecker) *Service {
	return &Service{
		repository:        repository,
		userLoader:        userLoader,
		permissionChecker: permissionChecker,
	}
}

// ValidateCredentials validates username and password
func (s *Service) ValidateCredentials(ctx context.Context, username, password string) (*domain.Credential, error) {
	cred, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrCredentialNotFound) {
			return nil, domain.ErrCredentialNotFound
		}
		return nil, err
	}

	if !cred.IsActive() {
		return nil, ErrCredentialNotActive
	}

	// Verify password using bcrypt
	if err := s.verifyPassword(cred.Password, password); err != nil {
		return nil, ErrInvalidPassword
	}

	return cred, nil
}

// verifyPassword verifies a password against a stored bcrypt hash
func (s *Service) verifyPassword(storedPasswordHash, providedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(providedPassword)); err != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// HashPassword hashes a password using bcrypt
func (s *Service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CreateCredentialParams contains the parameters for creating a credential
type CreateCredentialParams struct {
	Username string
	Password string
	UserID   string // Optional: for authenticated officer users creating credentials for others
}

// CreateCredentialWithAuthorization creates a new credential with authorization checks
// It resolves the userID from context (registration token) or request, validates authorization,
// and creates the credential.
func (s *Service) CreateCredentialWithAuthorization(ctx context.Context, req CreateCredentialParams) (*domain.Credential, error) {
	// Validate required fields
	if err := s.validateCredentialRequest(req); err != nil {
		return nil, err
	}

	authenticatedUser := httputil.UserFromContext(ctx)
	userID, _, err := s.resolveUserID(ctx, authenticatedUser, req)
	if err != nil {
		return nil, err
	}

	// Create the credential
	return s.CreateCredential(ctx, req.Username, req.Password, userID)
}

// validateCredentialRequest validates the credential request parameters.
func (s *Service) validateCredentialRequest(req CreateCredentialParams) error {
	if req.Username == "" {
		return ErrUsernameRequired
	}
	if req.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// resolveUserID resolves the user ID and validates authorization for credential creation.
func (s *Service) resolveUserID(ctx context.Context, authenticatedUser *userDomain.User, req CreateCredentialParams) (string, *userDomain.User, error) {
	if authenticatedUser == nil {
		return s.resolveUserIDForUnauthenticated(ctx)
	}
	return s.resolveUserIDForAuthenticated(ctx, authenticatedUser, req)
}

// resolveUserIDForUnauthenticated handles user ID resolution for unauthenticated requests (self-service via registration token).
func (s *Service) resolveUserIDForUnauthenticated(ctx context.Context) (string, *userDomain.User, error) {
	userID := httputil.RegistrationUserIDFromContext(ctx)
	if userID == "" {
		return "", nil, ErrRegistrationTokenRequired
	}

	targetUser, err := s.userLoader.Get(ctx, userID)
	if err != nil {
		return "", nil, ErrTargetUserNotFound
	}

	if targetUser.UserType != userDomain.UserTypeSelfService {
		return "", nil, ErrRegistrationTokenForSelfServiceOnly
	}

	return userID, targetUser, nil
}

// resolveUserIDForAuthenticated handles user ID resolution for authenticated requests (officer-created credentials).
func (s *Service) resolveUserIDForAuthenticated(ctx context.Context, authenticatedUser *userDomain.User, req CreateCredentialParams) (string, *userDomain.User, error) {
	if req.UserID == "" {
		return "", nil, ErrUserIDRequired
	}

	userID := req.UserID
	targetUser, err := s.userLoader.Get(ctx, userID)
	if err != nil {
		return "", nil, ErrTargetUserNotFound
	}

	// Authorization: Check if user can create credential for target user
	if err := s.checkCreateCredentialPermission(ctx, authenticatedUser, targetUser, userID); err != nil {
		return "", nil, err
	}

	return userID, targetUser, nil
}

// CreateCredential creates a new credential with hashed password
// userID is required and must link to an existing user entity
func (s *Service) CreateCredential(ctx context.Context, username, password string, userID string) (*domain.Credential, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if username == "" {
		return nil, ErrUsernameRequired
	}

	if password == "" {
		return nil, ErrPasswordRequired
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	cred, err := domain.NewCredential(username, hashedPassword, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, cred); err != nil {
		return nil, err
	}

	return cred, nil
}

// checkCreateCredentialPermission checks if the authenticated user can create a credential for the target user
func (s *Service) checkCreateCredentialPermission(ctx context.Context, authenticatedUser *userDomain.User, targetUser *userDomain.User, targetUserID string) error {
	// Allow unauthenticated users to create credentials for self-service users only
	// This enables users to create their own login credentials after self-registration
	if authenticatedUser == nil {
		if targetUser.UserType == userDomain.UserTypeSelfService {
			return nil // Allow unauthenticated users to create credentials for self-service users
		}
		return ErrInsufficientPermissions
	}

	// Self-service users can only create credentials for themselves
	if authenticatedUser.UserType == userDomain.UserTypeSelfService {
		return s.checkSelfServiceAccess(authenticatedUser.ID, targetUserID)
	}

	// Officer users can create for themselves or others with permission
	if authenticatedUser.UserType == userDomain.UserTypeOfficer {
		return s.checkOfficerCreateAccess(ctx, authenticatedUser, targetUser, targetUserID)
	}

	return ErrInsufficientPermissions
}

// checkSelfServiceAccess verifies self-service users can only access their own resources
func (s *Service) checkSelfServiceAccess(authenticatedUserID, targetUserID string) error {
	if authenticatedUserID != targetUserID {
		return ErrSelfServiceOwnCredentialsOnly
	}
	return nil
}

// checkOfficerCreateAccess verifies officer users can create credentials
func (s *Service) checkOfficerCreateAccess(ctx context.Context, authenticatedUser *userDomain.User, targetUser *userDomain.User, targetUserID string) error {
	// Officers can always create for themselves
	if authenticatedUser.ID == targetUserID {
		return nil
	}

	// For others, require permission
	hasPerm, err := s.permissionChecker.HasPermission(ctx, authenticatedUser.ID, ResourceUsers, ActionCreate)
	if err != nil || !hasPerm {
		return ErrInsufficientPermissions
	}

	// If target is an officer, also require permission
	if targetUser.UserType == userDomain.UserTypeOfficer {
		hasPerm, err := s.permissionChecker.HasPermission(ctx, authenticatedUser.ID, ResourceUsers, ActionCreate)
		if err != nil || !hasPerm {
			return ErrInsufficientPermissions
		}
	}

	return nil
}

// UpdatePasswordWithAuthorization updates the password for a credential with authorization checks
func (s *Service) UpdatePasswordWithAuthorization(ctx context.Context, username, newPassword string) error {
	// Get the credential to find the user ID
	cred, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	// Get target user to check authorization
	targetUser, err := s.userLoader.Get(ctx, cred.UserID)
	if err != nil {
		return ErrTargetUserNotFound
	}

	// Check authorization
	authenticatedUser := httputil.UserFromContext(ctx)
	if authenticatedUser == nil {
		return ErrAuthenticationRequired
	}

	if err := s.checkManageCredentialPermission(ctx, authenticatedUser, targetUser, cred.UserID); err != nil {
		return err
	}

	// Update password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := cred.UpdatePassword(hashedPassword); err != nil {
		return err
	}

	return s.repository.Update(ctx, cred)
}

// DeactivateCredentialWithAuthorization deactivates a credential with authorization checks
func (s *Service) DeactivateCredentialWithAuthorization(ctx context.Context, username string) error {
	// Get the credential to find the user ID
	cred, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	// Get target user to check authorization
	targetUser, err := s.userLoader.Get(ctx, cred.UserID)
	if err != nil {
		return ErrTargetUserNotFound
	}

	// Check authorization
	authenticatedUser := httputil.UserFromContext(ctx)
	if authenticatedUser == nil {
		return ErrAuthenticationRequired
	}

	if err := s.checkManageCredentialPermission(ctx, authenticatedUser, targetUser, cred.UserID); err != nil {
		return err
	}

	// Deactivate credential
	cred.Deactivate()
	return s.repository.Update(ctx, cred)
}

// UpdatePassword updates the password for a credential
func (s *Service) UpdatePassword(ctx context.Context, username, newPassword string) error {
	cred, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := cred.UpdatePassword(hashedPassword); err != nil {
		return err
	}

	return s.repository.Update(ctx, cred)
}

// DeactivateCredential deactivates a credential
func (s *Service) DeactivateCredential(ctx context.Context, username string) error {
	cred, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	cred.Deactivate()
	return s.repository.Update(ctx, cred)
}

// checkManageCredentialPermission checks if the authenticated user can manage (update/deactivate) a credential for the target user
func (s *Service) checkManageCredentialPermission(ctx context.Context, authenticatedUser *userDomain.User, targetUser *userDomain.User, targetUserID string) error {
	// Self-service users can only manage their own credentials
	if authenticatedUser.UserType == userDomain.UserTypeSelfService {
		if authenticatedUser.ID != targetUserID {
			return ErrSelfServiceOwnCredentialsOnly
		}
		return nil
	}

	// For officer users, check if they have permission or are managing their own
	if authenticatedUser.UserType == userDomain.UserTypeOfficer {
		if authenticatedUser.ID != targetUserID {
			// Managing another user's credential - require permission
			hasPerm, err := s.permissionChecker.HasPermission(ctx, authenticatedUser.ID, ResourceUsers, ActionCreate)
			if err != nil || !hasPerm {
				return ErrInsufficientPermissions
			}
		}
	}

	// If target user is an officer and we're not the same user, require permission
	if targetUser.UserType == userDomain.UserTypeOfficer {
		if authenticatedUser.ID != targetUserID {
			hasPerm, err := s.permissionChecker.HasPermission(ctx, authenticatedUser.ID, ResourceUsers, ActionCreate)
			if err != nil || !hasPerm {
				return ErrInsufficientPermissions
			}
		}
	}

	return nil
}

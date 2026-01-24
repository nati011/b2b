package service

import (
	"context"
	roleDomain "marketplace/internal/infra/authz/role/domain"
	"marketplace/internal/infra/user/domain"
	"marketplace/pkg/pagination"
	"errors"
	"testing"
	"time"
)

// Test constants to avoid duplication
const (
	testEmail          = "test@example.com"
	testUserName       = "Test User"
	existingEmail      = "existing@example.com"
	updatedName        = "Updated Name"
	errExpectedError   = "Expected error but got none"
	errUnexpectedError = "Unexpected error: %v"
	errExpectedUser    = "Expected user but got nil"
	errExpectedEmail   = "Expected Email %s, got %s"
	errExpectedStatus  = "Expected Status %s, got %s"
	errUserNotFound    = "user not found"
)

// mockUserRepository is a mock implementation of UserRepository for testing
type mockUserRepository struct {
	users                 map[string]*domain.User
	usersByEmail          map[string]*domain.User
	usersByPhoneNumber    map[string]*domain.User
	createFunc            func(context.Context, *domain.User) error
	findByIDFunc          func(context.Context, string) (*domain.User, error)
	findByEmailFunc       func(context.Context, string) (*domain.User, error)
	findByPhoneNumberFunc func(context.Context, string) (*domain.User, error)
	updateFunc            func(context.Context, *domain.User) error
	deleteFunc            func(context.Context, string) error
	findPageFunc          func(context.Context, pagination.PageRequest) (pagination.PageResult[*domain.User], error)
	existsFunc            func(context.Context, string) (bool, error)
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:              make(map[string]*domain.User),
		usersByEmail:       make(map[string]*domain.User),
		usersByPhoneNumber: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	if !user.Email.IsEmpty() {
		if _, exists := m.usersByEmail[user.Email.String()]; exists {
			return ErrEmailAlreadyExists
		}
	}
	if !user.PhoneNumber.IsEmpty() {
		if _, exists := m.usersByPhoneNumber[user.PhoneNumber.String()]; exists {
			// Phone number already exists - using email error as placeholder until separate error is defined
			return ErrEmailAlreadyExists
		}
	}
	m.users[user.ID] = user
	if !user.Email.IsEmpty() {
		m.usersByEmail[user.Email.String()] = user
	}
	if !user.PhoneNumber.IsEmpty() {
		m.usersByPhoneNumber[user.PhoneNumber.String()] = user
	}
	return nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	user, exists := m.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, email)
	}
	user, exists := m.usersByEmail[email]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	if m.findByPhoneNumberFunc != nil {
		return m.findByPhoneNumberFunc(ctx, phoneNumber)
	}
	user, exists := m.usersByPhoneNumber[phoneNumber]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	if _, exists := m.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	m.users[user.ID] = user
	m.usersByEmail[user.Email.String()] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	user, exists := m.users[id]
	if !exists {
		return ErrUserNotFound
	}
	delete(m.users, id)
	delete(m.usersByEmail, user.Email.String())
	return nil
}

func (m *mockUserRepository) FindPage(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.User], error) {
	if m.findPageFunc != nil {
		return m.findPageFunc(ctx, pageReq)
	}
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	total := len(users)

	if pageReq.Offset >= total {
		return pagination.NewPageResult([]*domain.User{}, total, pageReq), nil
	}

	end := pageReq.Offset + pageReq.Limit
	if end > total {
		end = total
	}

	paginatedUsers := users[pageReq.Offset:end]
	return pagination.NewPageResult(paginatedUsers, total, pageReq), nil
}

func (m *mockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	_, exists := m.users[id]
	return exists, nil
}

func (m *mockUserRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	user, exists := m.users[userID]
	if !exists {
		return ErrUserNotFound
	}
	for _, rid := range user.RoleIDs {
		if rid == roleID {
			return nil
		}
	}
	user.AssignRole(roleID)
	return nil
}

func (m *mockUserRepository) RevokeRole(ctx context.Context, userID, roleID string) error {
	user, exists := m.users[userID]
	if !exists {
		return ErrUserNotFound
	}
	user.RevokeRole(roleID)
	return nil
}

func (m *mockUserRepository) ListRoles(ctx context.Context, userID string) ([]roleDomain.Role, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, ErrUserNotFound
	}
	// Convert RoleIDs to Role objects for the mock
	roles := make([]roleDomain.Role, len(user.RoleIDs))
	for i, roleID := range user.RoleIDs {
		roles[i] = roleDomain.Role{ID: roleID}
	}
	return roles, nil
}

type mockRegistrationTokenRepository struct {
	tokens map[string]*domain.RegistrationToken
}

func (m *mockRegistrationTokenRepository) Create(ctx context.Context, token *domain.RegistrationToken) error {
	if m.tokens == nil {
		m.tokens = make(map[string]*domain.RegistrationToken)
	}
	m.tokens[token.ID] = token
	return nil
}

func (m *mockRegistrationTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RegistrationToken, error) {
	if m.tokens == nil {
		return nil, errors.New("token not found")
	}
	for _, t := range m.tokens {
		if t.Token == token {
			return t, nil
		}
	}
	return nil, errors.New("token not found")
}

func (m *mockRegistrationTokenRepository) MarkAsUsed(ctx context.Context, tokenID string) error {
	if m.tokens == nil {
		return errors.New("token not found")
	}
	if token, ok := m.tokens[tokenID]; ok {
		token.Used = true
		return nil
	}
	return errors.New("token not found")
}

func createTestUser(id, externalID, emailStr, name string, userType domain.UserType) *domain.User {
	email, _ := domain.NewEmail(emailStr)
	var phoneNumber domain.PhoneNumber
	user := &domain.User{
		ID:          id,
		ExternalID:  externalID,
		Email:       email,
		PhoneNumber: phoneNumber,
		Name:        name,
		Status:      domain.UserStatusInactive,
		UserType:    userType,
		RoleIDs:     []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return user
}

// assertError checks if an error matches expectations
func assertError(t *testing.T, err error, wantErr bool, errContains string) bool {
	if !wantErr {
		return false
	}
	if err == nil {
		t.Error(errExpectedError)
		return true
	}
	if errContains != "" && !contains(err.Error(), errContains) {
		t.Errorf("Expected error to contain '%s', got '%s'", errContains, err.Error())
	}
	return true
}

// assertNoError checks that no error occurred
func assertNoError(t *testing.T, err error) bool {
	if err != nil {
		t.Errorf(errUnexpectedError, err)
		return false
	}
	return true
}

// assertUserNotNull checks that user is not nil
func assertUserNotNull(t *testing.T, user *domain.User) bool {
	if user == nil {
		t.Error(errExpectedUser)
		return false
	}
	return true
}

// assertUserFields validates user fields match expected values
func assertUserFields(t *testing.T, user *domain.User, expectedExternalID, expectedEmail, expectedName string, expectedUserType domain.UserType) {
	if user.ExternalID != expectedExternalID {
		t.Errorf("Expected ExternalID %s, got %s", expectedExternalID, user.ExternalID)
	}
	if user.Email.String() != expectedEmail {
		t.Errorf(errExpectedEmail, expectedEmail, user.Email.String())
	}
	if user.Name != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, user.Name)
	}
	if user.UserType != expectedUserType {
		t.Errorf("Expected UserType %s, got %s", expectedUserType, user.UserType)
	}
	if user.Status != domain.UserStatusInactive {
		t.Errorf(errExpectedStatus, domain.UserStatusInactive, user.Status)
	}
}

// assertUpdatedUserFields validates updated user fields
func assertUpdatedUserFields(t *testing.T, user *domain.User, expectedEmail, expectedName string) {
	if user.Email.String() != expectedEmail {
		t.Errorf(errExpectedEmail, expectedEmail, user.Email.String())
	}
	if user.Name != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, user.Name)
	}
}

func TestUserServiceCreate(t *testing.T) {
	tests := []struct {
		testName    string
		externalID  string
		email       string
		userName    string
		userType    domain.UserType
		setupMock   func(*mockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			testName:   "successful creation",
			externalID: "ext123",
			email:      testEmail,
			userName:   testUserName,
			userType:   domain.UserTypeSelfService,
			setupMock: func(m *mockUserRepository) {
				// No setup needed for successful creation
			},
			wantErr: false,
		},
		{
			testName:   "email already exists",
			externalID: "ext123",
			email:      existingEmail,
			userName:   testUserName,
			userType:   domain.UserTypeSelfService,
			setupMock: func(m *mockUserRepository) {
				existingUser := createTestUser("id1", "ext1", existingEmail, "Existing User", domain.UserTypeSelfService)
				m.users["id1"] = existingUser
				m.usersByEmail[existingEmail] = existingUser
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			// Create email - assume valid email for business logic tests
			email, err := domain.NewEmail(tt.email)
			if err != nil {
				t.Fatalf("Failed to create email: %v", err)
			}

			ctx := context.Background()
			var phoneNumber domain.PhoneNumber
			result, err := service.Create(ctx, tt.externalID, email, phoneNumber, tt.userName, tt.userType)

			if assertError(t, err, tt.wantErr, tt.errContains) {
				return
			}

			if !assertNoError(t, err) || result == nil || !assertUserNotNull(t, result.User) {
				return
			}

			assertUserFields(t, result.User, tt.externalID, tt.email, tt.userName, tt.userType)
		})
	}
}

func TestUserServiceGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "successful get",
			id:   "user1",
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				m.users["user1"] = user
			},
			wantErr: false,
		},
		{
			name: errUserNotFound,
			id:   "nonexistent",
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			ctx := context.Background()
			user, err := service.Get(ctx, tt.id)

			if assertError(t, err, tt.wantErr, "") {
				return
			}

			if !assertNoError(t, err) || !assertUserNotNull(t, user) {
				return
			}

			if user.ID != tt.id {
				t.Errorf("Expected ID %s, got %s", tt.id, user.ID)
			}
		})
	}
}

func TestUserServiceFindByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name:  "successful find",
			email: testEmail,
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				m.users["user1"] = user
				m.usersByEmail[testEmail] = user
			},
			wantErr: false,
		},
		{
			name:  errUserNotFound,
			email: "nonexistent@example.com",
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			ctx := context.Background()
			user, err := service.FindByEmail(ctx, tt.email)

			if assertError(t, err, tt.wantErr, "") {
				return
			}

			if !assertNoError(t, err) || !assertUserNotNull(t, user) {
				return
			}

			if user.Email.String() != tt.email {
				t.Errorf(errExpectedEmail, tt.email, user.Email.String())
			}
		})
	}
}

func TestUserServiceUpdate(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		newEmail    string
		newName     string
		setupMock   func(*mockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:     "successful update",
			id:       "user1",
			newEmail: "updated@example.com",
			newName:  updatedName,
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				m.users["user1"] = user
				m.usersByEmail[testEmail] = user
			},
			wantErr: false,
		},
		{
			name:     errUserNotFound,
			id:       "nonexistent",
			newEmail: "updated@example.com",
			newName:  updatedName,
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
		{
			name:     "email already exists for another user",
			id:       "user1",
			newEmail: existingEmail,
			newName:  updatedName,
			setupMock: func(m *mockUserRepository) {
				user1 := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				user2 := createTestUser("user2", "ext2", existingEmail, "Existing User", domain.UserTypeSelfService)
				m.users["user1"] = user1
				m.users["user2"] = user2
				m.usersByEmail[testEmail] = user1
				m.usersByEmail[existingEmail] = user2
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:     "update with same email",
			id:       "user1",
			newEmail: testEmail,
			newName:  updatedName,
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				m.users["user1"] = user
				m.usersByEmail[testEmail] = user
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			// Create email - assume valid email for business logic tests
			email, err := domain.NewEmail(tt.newEmail)
			if err != nil {
				t.Fatalf("Failed to create email: %v", err)
			}

			ctx := context.Background()
			emptyStr := ""
			user, err := service.Update(ctx, tt.id, &emptyStr, email, tt.newName)

			if assertError(t, err, tt.wantErr, tt.errContains) {
				return
			}

			if !assertNoError(t, err) || !assertUserNotNull(t, user) {
				return
			}

			assertUpdatedUserFields(t, user, tt.newEmail, tt.newName)
		})
	}
}

func TestUserServiceDelete(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "successful delete",
			id:   "user1",
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				m.users["user1"] = user
				m.usersByEmail[testEmail] = user
			},
			wantErr: false,
		},
		{
			name: errUserNotFound,
			id:   "nonexistent",
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			ctx := context.Background()
			err := service.Delete(ctx, tt.id)

			if assertError(t, err, tt.wantErr, "") {
				return
			}

			if !assertNoError(t, err) {
				return
			}

			// Verify user was deleted
			ctx2 := context.Background()
			exists, _ := mockRepo.Exists(ctx2, tt.id)
			if exists {
				t.Errorf("User should have been deleted but still exists")
			}
		})
	}
}

func TestUserServiceActivate(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "successful activation",
			id:   "user1",
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				user.Status = domain.UserStatusInactive
				m.users["user1"] = user
			},
			wantErr: false,
		},
		{
			name: errUserNotFound,
			id:   "nonexistent",
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			ctx := context.Background()
			user, err := service.Activate(ctx, tt.id)

			if assertError(t, err, tt.wantErr, "") {
				return
			}

			if !assertNoError(t, err) {
				return
			}

			if user.Status != domain.UserStatusActive {
				t.Errorf(errExpectedStatus, domain.UserStatusActive, user.Status)
			}
		})
	}
}

func TestUserServiceDeactivate(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "successful deactivation",
			id:   "user1",
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				user.Status = domain.UserStatusActive
				m.users["user1"] = user
			},
			wantErr: false,
		},
		{
			name: errUserNotFound,
			id:   "nonexistent",
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			ctx := context.Background()
			user, err := service.Deactivate(ctx, tt.id)

			if assertError(t, err, tt.wantErr, "") {
				return
			}

			if !assertNoError(t, err) {
				return
			}

			if user.Status != domain.UserStatusInactive {
				t.Errorf(errExpectedStatus, domain.UserStatusInactive, user.Status)
			}
		})
	}
}

func TestUserServiceSuspend(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "successful suspension",
			id:   "user1",
			setupMock: func(m *mockUserRepository) {
				user := createTestUser("user1", "ext1", testEmail, testUserName, domain.UserTypeSelfService)
				user.Status = domain.UserStatusActive
				m.users["user1"] = user
			},
			wantErr: false,
		},
		{
			name: errUserNotFound,
			id:   "nonexistent",
			setupMock: func(m *mockUserRepository) {
				// No setup needed for not found case
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			mockTokenService := NewRegistrationTokenService(&mockRegistrationTokenRepository{})
			service := NewService(mockRepo, &stubRoleReader{}, nil, mockTokenService)

			ctx := context.Background()
			user, err := service.Suspend(ctx, tt.id)

			if assertError(t, err, tt.wantErr, "") {
				return
			}

			if !assertNoError(t, err) {
				return
			}

			if user.Status != domain.UserStatusSuspended {
				t.Errorf(errExpectedStatus, domain.UserStatusSuspended, user.Status)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

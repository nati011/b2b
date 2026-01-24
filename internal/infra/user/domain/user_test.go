package domain

import (
	"strings"
	"testing"
	"time"

	roleDomain "marketplace/internal/infra/authz/role/domain"
)

const (
	testEmail               = "test@example.com"
	testUserName            = "Test User"
	testRoleID              = "role-1"
	errMsgExpectedErrorNone = "Expected error but got none"
	errMsgUnexpectedError   = "Unexpected error: %v"
	errMsgExpectedStatus    = "Expected Status %s, got %s"
	errMsgExpectedUpdatedAt = "Expected UpdatedAt to be updated"
)

// assertUserCreationError is a helper function to reduce cognitive complexity
func assertUserCreationError(t *testing.T, err error, wantErr bool, errContains string) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Errorf(errMsgExpectedErrorNone)
			return
		}
		if errContains != "" && !strings.Contains(err.Error(), errContains) {
			t.Errorf("Expected error to contain '%s', got '%s'", errContains, err.Error())
		}
		return
	}
	if err != nil {
		t.Errorf(errMsgUnexpectedError, err)
	}
}

// assertUserProperties is a helper function to reduce cognitive complexity
func assertUserProperties(t *testing.T, user *User, externalID string, userName string, userType UserType) {
	t.Helper()
	if user == nil {
		t.Errorf("Expected user but got nil")
		return
	}
	if user.ExternalID != externalID {
		t.Errorf("Expected ExternalID %s, got %s", externalID, user.ExternalID)
	}
	if user.Name != userName {
		t.Errorf("Expected Name %s, got %s", userName, user.Name)
	}
	if user.UserType != userType {
		t.Errorf("Expected UserType %s, got %s", userType, user.UserType)
	}
	if user.Status != UserStatusInactive {
		t.Errorf("Expected Status %s, got %s", UserStatusInactive, user.Status)
	}
	if len(user.RoleIDs) != 0 {
		t.Errorf("Expected empty roles, got %d", len(user.RoleIDs))
	}
}

// userUpdateExpectations holds expected values for user update assertions
type userUpdateExpectations struct {
	externalID string
	email      Email
	name       string
}

// assertUserUpdateResult is a helper function to reduce cognitive complexity
func assertUserUpdateResult(t *testing.T, user *User, err error, wantErr bool, errContains string, expected userUpdateExpectations) {
	t.Helper()
	assertUserCreationError(t, err, wantErr, errContains)
	if !wantErr {
		if user.ExternalID != expected.externalID {
			t.Errorf("Expected ExternalID %s, got %s", expected.externalID, user.ExternalID)
		}
		if user.Email.String() != expected.email.String() {
			t.Errorf("Expected Email %s, got %s", expected.email.String(), user.Email.String())
		}
		if user.Name != expected.name {
			t.Errorf("Expected Name %s, got %s", expected.name, user.Name)
		}
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		externalID  string
		email       Email
		phoneNumber PhoneNumber
		userName    string
		userType    UserType
		wantErr     bool
		errContains string
	}{
		{
			name:        "successful creation with email",
			externalID:  "ext-123",
			email:       mustNewEmail(t, testEmail),
			phoneNumber: PhoneNumber{},
			userName:    testUserName,
			userType:    UserTypeSelfService,
			wantErr:     false,
		},
		{
			name:        "successful creation with phone number",
			externalID:  "ext-124",
			email:       Email{},
			phoneNumber: mustNewPhoneNumber(t, "1234567890"),
			userName:    testUserName,
			userType:    UserTypeOfficer,
			wantErr:     false,
		},
		{
			name:        "successful creation with both email and phone",
			externalID:  "ext-125",
			email:       mustNewEmail(t, testEmail),
			phoneNumber: mustNewPhoneNumber(t, "1234567890"),
			userName:    testUserName,
			userType:    UserTypeSelfService,
			wantErr:     false,
		},
		{
			name:        "fails when both email and phone are empty",
			externalID:  "ext-126",
			email:       Email{},
			phoneNumber: PhoneNumber{},
			userName:    testUserName,
			userType:    UserTypeSelfService,
			wantErr:     true,
			errContains: "either email or phone number must be provided",
		},
		{
			name:        "fails with empty external ID when provided",
			externalID:  "",
			email:       mustNewEmail(t, testEmail),
			phoneNumber: PhoneNumber{},
			userName:    testUserName,
			userType:    UserTypeSelfService,
			wantErr:     false, // External ID is optional
		},
		{
			name:        "fails with empty name",
			externalID:  "ext-127",
			email:       mustNewEmail(t, testEmail),
			phoneNumber: PhoneNumber{},
			userName:    "",
			userType:    UserTypeSelfService,
			wantErr:     true,
			errContains: "empty",
		},
		{
			name:        "fails with invalid user type",
			externalID:  "ext-128",
			email:       mustNewEmail(t, testEmail),
			phoneNumber: PhoneNumber{},
			userName:    testUserName,
			userType:    UserType("invalid"),
			wantErr:     true,
			errContains: "contain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.externalID, tt.email, tt.phoneNumber, tt.userName, tt.userType)
			assertUserCreationError(t, err, tt.wantErr, tt.errContains)
			if !tt.wantErr {
				assertUserProperties(t, user, tt.externalID, tt.userName, tt.userType)
			}
		})
	}
}

func TestUserUpdate(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		externalID  string
		email       Email
		newName     string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful update",
			user:       createTestUser(t, "ext-1", testEmail, "Old Name", UserTypeSelfService),
			externalID: "ext-1-updated",
			email:      mustNewEmail(t, "updated@example.com"),
			newName:    "New Name",
			wantErr:    false,
		},
		{
			name:        "fails with empty name",
			user:        createTestUser(t, "ext-2", testEmail, "Old Name", UserTypeSelfService),
			externalID:  "ext-2",
			email:       mustNewEmail(t, testEmail),
			newName:     "",
			wantErr:     true,
			errContains: "empty",
		},
		{
			name:        "fails with invalid external ID length",
			user:        createTestUser(t, "ext-3", testEmail, "Test Name", UserTypeSelfService),
			externalID:  string(make([]byte, 300)), // Too long
			email:       mustNewEmail(t, testEmail),
			newName:     "Test Name",
			wantErr:     true,
			errContains: "long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.user.UpdatedAt
			err := tt.user.Update(tt.externalID, tt.email, tt.newName)
			expected := userUpdateExpectations{
				externalID: tt.externalID,
				email:      tt.email,
				name:       tt.newName,
			}
			assertUserUpdateResult(t, tt.user, err, tt.wantErr, tt.errContains, expected)
			if !tt.wantErr && !tt.user.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf(errMsgExpectedUpdatedAt)
			}
		})
	}
}

// assertStatusChange is a helper function to reduce cognitive complexity for status change tests
func assertStatusChange(t *testing.T, user *User, err error, wantErr bool, expectedStatus UserStatus, originalUpdatedAt time.Time, checkUpdatedAt bool) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Errorf(errMsgExpectedErrorNone)
		}
		return
	}
	if err != nil {
		t.Errorf(errMsgUnexpectedError, err)
		return
	}
	if user.Status != expectedStatus {
		t.Errorf(errMsgExpectedStatus, expectedStatus, user.Status)
	}
	if checkUpdatedAt && !user.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf(errMsgExpectedUpdatedAt)
	}
}

func TestUserActivate(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		expectedStatus UserStatus
		wantErr        bool
	}{
		{
			name:           "successful activation from inactive",
			user:           createTestUserWithStatus(t, UserStatusInactive),
			expectedStatus: UserStatusActive,
			wantErr:        false,
		},
		{
			name:           "no change when already active",
			user:           createTestUserWithStatus(t, UserStatusActive),
			expectedStatus: UserStatusActive,
			wantErr:        false,
		},
		{
			name:           "activation from suspended",
			user:           createTestUserWithStatus(t, UserStatusSuspended),
			expectedStatus: UserStatusActive,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.user.UpdatedAt
			originalStatus := tt.user.Status
			err := tt.user.Activate()
			checkUpdatedAt := originalStatus != tt.expectedStatus
			assertStatusChange(t, tt.user, err, tt.wantErr, tt.expectedStatus, originalUpdatedAt, checkUpdatedAt)
		})
	}
}

func TestUserDeactivate(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		expectedStatus UserStatus
		wantErr        bool
	}{
		{
			name:           "successful deactivation from active",
			user:           createTestUserWithStatus(t, UserStatusActive),
			expectedStatus: UserStatusInactive,
			wantErr:        false,
		},
		{
			name:           "no change when already inactive",
			user:           createTestUserWithStatus(t, UserStatusInactive),
			expectedStatus: UserStatusInactive,
			wantErr:        false,
		},
		{
			name:           "deactivation from suspended",
			user:           createTestUserWithStatus(t, UserStatusSuspended),
			expectedStatus: UserStatusInactive,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.user.UpdatedAt
			originalStatus := tt.user.Status
			err := tt.user.Deactivate()
			checkUpdatedAt := originalStatus != tt.expectedStatus
			assertStatusChange(t, tt.user, err, tt.wantErr, tt.expectedStatus, originalUpdatedAt, checkUpdatedAt)
		})
	}
}

func TestUserSuspend(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		expectedStatus UserStatus
		wantErr        bool
	}{
		{
			name:           "successful suspension from active",
			user:           createTestUserWithStatus(t, UserStatusActive),
			expectedStatus: UserStatusSuspended,
			wantErr:        false,
		},
		{
			name:           "successful suspension from inactive",
			user:           createTestUserWithStatus(t, UserStatusInactive),
			expectedStatus: UserStatusSuspended,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.user.UpdatedAt
			err := tt.user.Suspend()
			assertStatusChange(t, tt.user, err, tt.wantErr, tt.expectedStatus, originalUpdatedAt, true)
		})
	}
}

// assertRoleAssignment is a helper function to reduce cognitive complexity
func assertRoleAssignment(t *testing.T, user *User, role roleDomain.Role, expectedCount int, originalUpdatedAt time.Time, initialCount int) {
	t.Helper()
	if len(user.RoleIDs) != expectedCount {
		t.Errorf("Expected %d roles, got %d", expectedCount, len(user.RoleIDs))
	}
	found := false
	for _, roleID := range user.RoleIDs {
		if roleID == role.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected role %s to be assigned", role.ID)
	}
	// UpdatedAt should be updated when a new role is added, but not when role already exists
	if expectedCount > initialCount && !user.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Expected UpdatedAt to be updated when role was added")
	}
}

// assertRoleRevocation is a helper function to reduce cognitive complexity
func assertRoleRevocation(t *testing.T, user *User, roleID string, expectedCount int, originalUpdatedAt time.Time, initialCount int) {
	t.Helper()
	if len(user.RoleIDs) != expectedCount {
		t.Errorf("Expected %d roles, got %d", expectedCount, len(user.RoleIDs))
	}
	found := false
	for _, rid := range user.RoleIDs {
		if rid == roleID {
			found = true
			break
		}
	}
	if found {
		t.Errorf("Expected role %s to be revoked", roleID)
	}
	if expectedCount < initialCount && !user.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Expected UpdatedAt to be updated when role was revoked")
	}
}

func TestUserAssignRole(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		role          roleDomain.Role
		expectedCount int
	}{
		{
			name:          "assigns role when user has no roles",
			user:          createTestUser(t, "ext-1", testEmail, testUserName, UserTypeSelfService),
			role:          roleDomain.Role{ID: testRoleID, Name: "Admin"},
			expectedCount: 1,
		},
		{
			name:          "assigns role when user has other roles",
			user:          createTestUserWithRoles(t, []roleDomain.Role{{ID: testRoleID, Name: "Admin"}}),
			role:          roleDomain.Role{ID: "role-2", Name: "User"},
			expectedCount: 2,
		},
		{
			name:          "does not duplicate when role already assigned",
			user:          createTestUserWithRoles(t, []roleDomain.Role{{ID: testRoleID, Name: "Admin"}}),
			role:          roleDomain.Role{ID: testRoleID, Name: "Admin"},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.user.UpdatedAt
			initialCount := len(tt.user.RoleIDs)
			tt.user.AssignRole(tt.role.ID)
			assertRoleAssignment(t, tt.user, tt.role, tt.expectedCount, originalUpdatedAt, initialCount)
		})
	}
}

func TestUserRevokeRole(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		roleID        string
		expectedCount int
	}{
		{
			name:          "revokes role when role exists",
			user:          createTestUserWithRoles(t, []roleDomain.Role{{ID: "role-1", Name: "Admin"}, {ID: "role-2", Name: "User"}}),
			roleID:        testRoleID,
			expectedCount: 1,
		},
		{
			name:          "no change when role does not exist",
			user:          createTestUserWithRoles(t, []roleDomain.Role{{ID: testRoleID, Name: "Admin"}}),
			roleID:        "role-999",
			expectedCount: 1,
		},
		{
			name:          "no change when user has no roles",
			user:          createTestUser(t, "ext-1", testEmail, testUserName, UserTypeSelfService),
			roleID:        testRoleID,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.user.UpdatedAt
			initialCount := len(tt.user.RoleIDs)
			tt.user.RevokeRole(tt.roleID)
			assertRoleRevocation(t, tt.user, tt.roleID, tt.expectedCount, originalUpdatedAt, initialCount)
		})
	}
}

// TestUserHasPermission removed - HasPermission method moved to PermissionChecker service
// Permission checking is now handled at the service layer, not the domain layer

func TestUserCanLogin(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		expectedResult bool
	}{
		{
			name:           "returns true when user is active",
			user:           createTestUserWithStatus(t, UserStatusActive),
			expectedResult: true,
		},
		{
			name:           "returns false when user is inactive",
			user:           createTestUserWithStatus(t, UserStatusInactive),
			expectedResult: false,
		},
		{
			name:           "returns false when user is suspended",
			user:           createTestUserWithStatus(t, UserStatusSuspended),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.CanLogin()
			if result != tt.expectedResult {
				t.Errorf("Expected CanLogin() = %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestUserTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected UserType
	}{
		{
			name:     "converts selfservice",
			input:    "selfservice",
			expected: UserTypeSelfService,
		},
		{
			name:     "converts self_service",
			input:    "self_service",
			expected: UserTypeSelfService,
		},
		{
			name:     "converts officer",
			input:    "officer",
			expected: UserTypeOfficer,
		},
		{
			name:     "defaults to selfservice for invalid input",
			input:    "invalid",
			expected: UserTypeSelfService,
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: UserTypeSelfService,
		},
		{
			name:     "handles whitespace",
			input:    "  officer  ",
			expected: UserTypeOfficer,
		},
		{
			name:     "case insensitive",
			input:    "OFFICER",
			expected: UserTypeOfficer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UserTypeFromString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected UserTypeFromString(%s) = %s, got %s", tt.input, tt.expected, result)
			}
		})
	}
}

// Helper functions

func mustNewEmail(t *testing.T, value string) Email {
	email, err := NewEmail(value)
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	return email
}

func mustNewPhoneNumber(t *testing.T, value string) PhoneNumber {
	phone, err := NewPhoneNumber(value, nil)
	if err != nil {
		t.Fatalf("Failed to create phone number: %v", err)
	}
	return phone
}

func createTestUser(t *testing.T, externalID, emailStr, name string, userType UserType) *User {
	email := mustNewEmail(t, emailStr)
	user, err := NewUser(externalID, email, PhoneNumber{}, name, userType)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestUserWithStatus(t *testing.T, status UserStatus) *User {
	user := createTestUser(t, "ext-1", testEmail, testUserName, UserTypeSelfService)
	user.Status = status
	return user
}

func createTestUserWithRoles(t *testing.T, roles []roleDomain.Role) *User {
	user := createTestUser(t, "ext-1", testEmail, testUserName, UserTypeSelfService)
	// Assign role IDs from the roles
	for _, role := range roles {
		user.AssignRole(role.ID)
	}
	return user
}

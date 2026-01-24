//go:build integration
// +build integration

package repository

import (
	"context"
	roleDomain "marketplace/internal/infra/authz/role/domain"
	"marketplace/internal/infra/user/domain"
	userservice "marketplace/internal/infra/user/service"
	"marketplace/pkg/pagination"
	"marketplace/test"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	errMsgFailedCreateEmail       = "Failed to create email: %v"
	errMsgFailedCreateUser        = "Failed to create user: %v"
	errMsgExpectedUserNil         = "Expected user but got nil"
	errMsgExpectedStatus          = "Expected Status '%s', got '%s'"
	errMsgExpectedPersistedStatus = "Expected persisted Status '%s', got '%s'"
	errMsgExpectedErrUserNotFound = "Expected ErrUserNotFound, got %v"
	errMsgFailedActivateUser      = "Failed to activate user: %v"
	testNonExistentID             = "nonexistent-id"
	testEmailFindByEmail          = "findbyemail@example.com"
	testEmailUpdated              = "updated@example.com"
	testNameUpdated               = "Updated Name"
)

type stubRoleReader struct{}

func (stubRoleReader) FindByID(ctx context.Context, id string) (*roleDomain.Role, error) {
	return &roleDomain.Role{ID: id}, nil
}

var (
	service   *userservice.Service
	repo      *Repository
	ctx       context.Context
	db        *sql.DB
	cleanup   func()
	setupOnce sync.Once
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupTestDB(t *testing.T) {
	setupOnce.Do(func() {
		// Setup test DB with panic recovery
		var setupDB *sql.DB
		var setupCleanup func()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Skipf("Docker is not accessible (panic: %v). Skipping integration tests.", r)
				}
			}()
			setupDB, setupCleanup = test.SetupTestDB(t)
		}()

		if setupDB == nil {
			return // Test was skipped
		}

		db = setupDB
		cleanup = setupCleanup

		repo = NewRepository(db)
		service = userservice.NewService(repo, &stubRoleReader{}, nil, nil)
		ctx = context.Background()
	})
}

func TestUserRepositoryIntegration(t *testing.T) {
	setupTestDB(t)
	if db == nil {
		t.Skip("Database setup failed, skipping integration tests")
	}

	t.Run("Create", func(t *testing.T) {
		email, err := domain.NewEmail("test@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext123", email, "Test User", domain.UserTypeSelfService)
		if err == nil {
			user = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		if user == nil {
			t.Fatal(errMsgExpectedUserNil)
		}

		if user.ID == "" {
			t.Error("Expected user ID to be set")
		}

		if user.ExternalID != "ext123" {
			t.Errorf("Expected ExternalID 'ext123', got '%s'", user.ExternalID)
		}

		if user.Email.String() != "test@example.com" {
			t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email.String())
		}

		if user.Name != "Test User" {
			t.Errorf("Expected Name 'Test User', got '%s'", user.Name)
		}

		if user.UserType != domain.UserTypeSelfService {
			t.Errorf("Expected UserType '%s', got '%s'", domain.UserTypeSelfService, user.UserType)
		}

		if user.Status != domain.UserStatusInactive {
			t.Errorf(errMsgExpectedStatus, domain.UserStatusInactive, user.Status)
		}

		if user.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if user.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Create_DuplicateEmail", func(t *testing.T) {
		email, err := domain.NewEmail("duplicate@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		// Create first user
		_, err = service.Create(ctx, "ext1", email, "First User", domain.UserTypeSelfService)
		if err != nil {
			t.Fatalf("Failed to create first user: %v", err)
		}

		// Try to create second user with same email
		_, err = service.Create(ctx, "ext2", email, "Second User", domain.UserTypeSelfService)
		if err == nil {
			t.Fatal("Expected error when creating user with duplicate email, got nil")
		}

		if err != userservice.ErrEmailAlreadyExists {
			t.Errorf("Expected ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		email, err := domain.NewEmail("findbyid@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext456", email, "Find By ID User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		foundUser, err := service.Get(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to find user by ID: %v", err)
		}

		if foundUser == nil {
			t.Fatal(errMsgExpectedUserNil)
		}

		if foundUser.ID != createdUser.ID {
			t.Errorf("Expected ID '%s', got '%s'", createdUser.ID, foundUser.ID)
		}

		if foundUser.Email.String() != "findbyid@example.com" {
			t.Errorf("Expected Email 'findbyid@example.com', got '%s'", foundUser.Email.String())
		}

		if foundUser.Name != "Find By ID User" {
			t.Errorf("Expected Name 'Find By ID User', got '%s'", foundUser.Name)
		}
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		_, err := service.Get(ctx, testNonExistentID)
		if err == nil {
			t.Fatal("Expected error when finding non-existent user, got nil")
		}

		if err != userservice.ErrUserNotFound {
			t.Errorf(errMsgExpectedErrUserNotFound, err)
		}
	})

	t.Run("FindByEmail", func(t *testing.T) {
		email, err := domain.NewEmail(testEmailFindByEmail)
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext789", email, "Find By Email User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		foundUser, err := service.FindByEmail(ctx, testEmailFindByEmail)
		if err != nil {
			t.Fatalf("Failed to find user by email: %v", err)
		}

		if foundUser == nil {
			t.Fatal(errMsgExpectedUserNil)
		}

		if foundUser.ID != createdUser.ID {
			t.Errorf("Expected ID '%s', got '%s'", createdUser.ID, foundUser.ID)
		}

		if foundUser.Email.String() != testEmailFindByEmail {
			t.Errorf("Expected Email '%s', got '%s'", testEmailFindByEmail, foundUser.Email.String())
		}
	})

	t.Run("FindByEmail_NotFound", func(t *testing.T) {
		_, err := service.FindByEmail(ctx, "nonexistent@example.com")
		if err == nil {
			t.Fatal("Expected error when finding non-existent user, got nil")
		}

		if err != userservice.ErrUserNotFound {
			t.Errorf(errMsgExpectedErrUserNotFound, err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		email, err := domain.NewEmail("update@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext999", email, "Original Name", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		newEmail, err := domain.NewEmail(testEmailUpdated)
		if err != nil {
			t.Fatalf("Failed to create new email: %v", err)
		}

		emptyStr := ""
		updatedUser, err := service.Update(ctx, createdUser.ID, &emptyStr, newEmail, testNameUpdated)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		if updatedUser.Email.String() != testEmailUpdated {
			t.Errorf("Expected Email 'updated@example.com', got '%s'", updatedUser.Email.String())
		}

		if updatedUser.Name != testNameUpdated {
			t.Errorf("Expected Name 'Updated Name', got '%s'", updatedUser.Name)
		}

		if !updatedUser.UpdatedAt.After(createdUser.UpdatedAt) {
			t.Error("Expected UpdatedAt to be after creation time")
		}

		// Verify changes persisted
		foundUser, err := service.Get(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to find updated user: %v", err)
		}

		if foundUser.Email.String() != testEmailUpdated {
			t.Errorf("Expected persisted Email 'updated@example.com', got '%s'", foundUser.Email.String())
		}

		if foundUser.Name != testNameUpdated {
			t.Errorf("Expected persisted Name 'Updated Name', got '%s'", foundUser.Name)
		}
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		email, err := domain.NewEmail("notfound@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		emptyStr := ""
		_, err = service.Update(ctx, testNonExistentID, &emptyStr, email, "Name")
		if err == nil {
			t.Fatal("Expected error when updating non-existent user, got nil")
		}

		if err != userservice.ErrUserNotFound {
			t.Errorf(errMsgExpectedErrUserNotFound, err)
		}
	})

	t.Run("Update_DuplicateEmail", func(t *testing.T) {
		email1, err := domain.NewEmail("user1@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		email2, err := domain.NewEmail("user2@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result1, err := service.Create(ctx, "ext1", email1, "User 1", domain.UserTypeSelfService)
		if err == nil {
			user1 = result1.User
		}
		if err != nil {
			t.Fatalf("Failed to create user1: %v", err)
		}

		_, err = service.Create(ctx, "ext2", email2, "User 2", domain.UserTypeSelfService)
		if err != nil {
			t.Fatalf("Failed to create user2: %v", err)
		}

		// Try to update user1 with user2's email
		emptyStr := ""
		_, err = service.Update(ctx, user1.ID, &emptyStr, email2, "User 1")
		if err == nil {
			t.Fatal("Expected error when updating with duplicate email, got nil")
		}

		if err != userservice.ErrEmailAlreadyExists {
			t.Errorf("Expected ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		email, err := domain.NewEmail("delete@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext-delete", email, "Delete User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		err = service.Delete(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify user is deleted
		_, err = service.Get(ctx, createdUser.ID)
		if err == nil {
			t.Fatal("Expected error when finding deleted user, got nil")
		}

		if err != userservice.ErrUserNotFound {
			t.Errorf(errMsgExpectedErrUserNotFound, err)
		}
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		err := service.Delete(ctx, "nonexistent-id")
		if err == nil {
			t.Fatal("Expected error when deleting non-existent user, got nil")
		}

		if err != userservice.ErrUserNotFound {
			t.Errorf(errMsgExpectedErrUserNotFound, err)
		}
	})

	t.Run("Activate", func(t *testing.T) {
		email, err := domain.NewEmail("activate@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext-activate", email, "Activate User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		if createdUser.Status != domain.UserStatusInactive {
			t.Errorf("Expected initial status '%s', got '%s'", domain.UserStatusInactive, createdUser.Status)
		}

		activatedUser, err := service.Activate(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf(errMsgFailedActivateUser, err)
		}

		if activatedUser.Status != domain.UserStatusActive {
			t.Errorf(errMsgExpectedStatus, domain.UserStatusActive, activatedUser.Status)
		}

		// Verify activation persisted
		foundUser, err := service.Get(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to find activated user: %v", err)
		}

		if foundUser.Status != domain.UserStatusActive {
			t.Errorf(errMsgExpectedPersistedStatus, domain.UserStatusActive, foundUser.Status)
		}
	})

	t.Run("Deactivate", func(t *testing.T) {
		email, err := domain.NewEmail("deactivate@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext-deactivate", email, "Deactivate User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		// Activate first
		_, err = service.Activate(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf(errMsgFailedActivateUser, err)
		}

		// Then deactivate
		deactivatedUser, err := service.Deactivate(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to deactivate user: %v", err)
		}

		if deactivatedUser.Status != domain.UserStatusInactive {
			t.Errorf(errMsgExpectedStatus, domain.UserStatusInactive, deactivatedUser.Status)
		}

		// Verify deactivation persisted
		foundUser, err := service.Get(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to find deactivated user: %v", err)
		}

		if foundUser.Status != domain.UserStatusInactive {
			t.Errorf(errMsgExpectedPersistedStatus, domain.UserStatusInactive, foundUser.Status)
		}
	})

	t.Run("Suspend", func(t *testing.T) {
		email, err := domain.NewEmail("suspend@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext-suspend", email, "Suspend User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		// Activate first
		_, err = service.Activate(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf(errMsgFailedActivateUser, err)
		}

		// Then suspend
		suspendedUser, err := service.Suspend(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to suspend user: %v", err)
		}

		if suspendedUser.Status != domain.UserStatusSuspended {
			t.Errorf(errMsgExpectedStatus, domain.UserStatusSuspended, suspendedUser.Status)
		}

		// Verify suspension persisted
		foundUser, err := service.Get(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to find suspended user: %v", err)
		}

		if foundUser.Status != domain.UserStatusSuspended {
			t.Errorf(errMsgExpectedPersistedStatus, domain.UserStatusSuspended, foundUser.Status)
		}
	})

	t.Run("FindPage", func(t *testing.T) {
		// Create multiple users
		for i := 0; i < 5; i++ {
			email, err := domain.NewEmail(fmt.Sprintf("page%d@example.com", i))
			if err != nil {
				t.Fatalf(errMsgFailedCreateEmail, err)
			}
			_, err = service.Create(ctx, fmt.Sprintf("ext-page%d", i), email, fmt.Sprintf("User %d", i), domain.UserTypeSelfService)
			if err != nil {
				t.Fatalf(errMsgFailedCreateUser, err)
			}
		}

		pageReq := pagination.NewPageRequest(1, 3)

		result, err := service.FindPage(ctx, pageReq)
		if err != nil {
			t.Fatalf("Failed to find page: %v", err)
		}

		if result.Total < 5 {
			t.Errorf("Expected at least 5 total users, got %d", result.Total)
		}

		if len(result.Items) > pageReq.Limit {
			t.Errorf("Expected at most %d items, got %d", pageReq.Limit, len(result.Items))
		}

		if result.Page != pageReq.Page {
			t.Errorf("Expected page %d, got %d", pageReq.Page, result.Page)
		}

		if result.Limit != pageReq.Limit {
			t.Errorf("Expected limit %d, got %d", pageReq.Limit, result.Limit)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		email, err := domain.NewEmail("exists@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		result, err := service.Create(ctx, "ext-exists", email, "Exists User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}

		exists, err := repo.Exists(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("Failed to check existence: %v", err)
		}

		if !exists {
			t.Error("Expected user to exist, got false")
		}

		exists, err = repo.Exists(ctx, testNonExistentID)
		if err != nil {
			t.Fatalf("Failed to check existence: %v", err)
		}

		if exists {
			t.Error("Expected user to not exist, got true")
		}
	})

	t.Run("Timestamps", func(t *testing.T) {
		email, err := domain.NewEmail("timestamps@example.com")
		if err != nil {
			t.Fatalf(errMsgFailedCreateEmail, err)
		}

		beforeCreate := time.Now()
		result, err := service.Create(ctx, "ext-timestamps", email, "Timestamps User", domain.UserTypeSelfService)
		if err == nil {
			createdUser = result.User
		}
		if err != nil {
			t.Fatalf(errMsgFailedCreateUser, err)
		}
		afterCreate := time.Now()

		if createdUser.CreatedAt.Before(beforeCreate) || createdUser.CreatedAt.After(afterCreate) {
			t.Error("CreatedAt should be between beforeCreate and afterCreate")
		}

		if createdUser.UpdatedAt.Before(beforeCreate) || createdUser.UpdatedAt.After(afterCreate) {
			t.Error("UpdatedAt should be between beforeCreate and afterCreate")
		}

		// Wait a bit and update
		time.Sleep(10 * time.Millisecond)
		beforeUpdate := time.Now()
		newEmail, err := domain.NewEmail("updated-timestamps@example.com")
		if err != nil {
			t.Fatalf("Failed to create new email: %v", err)
		}
		emptyStr := ""
		updatedUser, err := service.Update(ctx, createdUser.ID, &emptyStr, newEmail, testNameUpdated)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}
		afterUpdate := time.Now()

		if updatedUser.UpdatedAt.Before(beforeUpdate) || updatedUser.UpdatedAt.After(afterUpdate) {
			t.Error("UpdatedAt should be updated after update operation")
		}

		// Allow small difference due to database timestamp precision
		timeDiff := updatedUser.CreatedAt.Sub(createdUser.CreatedAt)
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff > time.Second {
			t.Errorf("CreatedAt should not change after update. Difference: %v", timeDiff)
		}
	})
}

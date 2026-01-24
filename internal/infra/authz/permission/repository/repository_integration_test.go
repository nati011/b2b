//go:build integration
// +build integration

package permission

import (
	"context"
	resourcepkg "marketplace/internal/infra/authz/resource"
	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	"marketplace/pkg/pagination"
	"marketplace/test"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	service      *PermissionService
	repo         *PermissionRepository
	resourceRepo *resourcepkg.Repository
	resourceSvc  *resourcepkg.Service
	ctx          context.Context
	db           *sql.DB
	cleanup      func()
	setupOnce    sync.Once
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupTestDB(t *testing.T) {
	setupOnce.Do(func() {
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
			return
		}

		db = setupDB
		cleanup = setupCleanup

		repo = NewPermissionRepository(db)
		resourceRepo = resourcepkg.NewRepository(db, nil)
		resourceSvc = resourcepkg.NewService(resourceRepo)
		service = NewPermissionService(repo, resourceSvc)
		ctx = context.Background()
	})
}

func newPermissionID() string {
	return uuid.New().String()
}

func createTestResource(t *testing.T, code, serviceName string) *resourceDomain.Resource {
	t.Helper()
	resource, err := resourceDomain.NewResource(code, serviceName, "test description")
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}
	action, err := resourceDomain.NewAction(resource.ID, "read", "read action")
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}
	resource.AttachActions([]resourceDomain.Action{*action})
	return resource
}

func truncatePermissionTables(t *testing.T) {
	t.Helper()
	if db == nil {
		t.Skip("Database not initialized, skipping truncate")
		return
	}

	if _, err := db.ExecContext(ctx, `TRUNCATE TABLE permissions, resource_actions, resources RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("failed to truncate permission tables: %v", err)
	}
}

func TestPermissionRepositoryIntegration(t *testing.T) {
	setupTestDB(t)
	if db == nil {
		t.Skip("Database setup failed, skipping integration tests")
	}
	t.Cleanup(func() {
		if cleanup != nil {
			cleanup()
			cleanup = nil
		}
	})

	t.Run("Create", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Create permission
		permID := newPermissionID()
		perm, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create permission: %v", err)
		}

		if perm == nil {
			t.Fatal("Expected permission but got nil")
		}

		if perm.ID != permID {
			t.Errorf("Expected ID '%s', got '%s'", permID, perm.ID)
		}

		if perm.Resource.Code != resourceCode {
			t.Errorf("Expected Resource.Code '%s', got '%s'", resourceCode, perm.Resource.Code)
		}

		if perm.Action != "read" {
			t.Errorf("Expected Action 'read', got '%s'", perm.Action)
		}

		// Verify permission was created
		found, err := repo.FindByID(ctx, permID)
		if err != nil {
			t.Fatalf("Failed to find created permission: %v", err)
		}

		if found == nil {
			t.Fatal("Expected permission but got nil")
		}

		if found.ID != permID {
			t.Errorf("Expected ID '%s', got '%s'", permID, found.ID)
		}
	})

	t.Run("Create_Duplicate", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Create first permission
		permID := newPermissionID()
		_, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create first permission: %v", err)
		}

		// Try to create duplicate
		_, err = service.Create(ctx, permID, resourceCode, "read", "Duplicate")
		if err == nil {
			t.Fatal("Expected error for duplicate permission, got nil")
		}

		if err != ErrPermissionAlreadyExists {
			t.Errorf("Expected ErrPermissionAlreadyExists, got %v", err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Create permission
		permID := newPermissionID()
		created, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create permission: %v", err)
		}

		// Find permission
		found, err := repo.FindByID(ctx, permID)
		if err != nil {
			t.Fatalf("Failed to find permission: %v", err)
		}

		if found == nil {
			t.Fatal("Expected permission but got nil")
		}

		if found.ID != created.ID {
			t.Errorf("Expected ID '%s', got '%s'", created.ID, found.ID)
		}

		if found.Resource.Code != resourceCode {
			t.Errorf("Expected Resource.Code '%s', got '%s'", resourceCode, found.Resource.Code)
		}

		if found.Action != "read" {
			t.Errorf("Expected Action 'read', got '%s'", found.Action)
		}
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		_, err := repo.FindByID(ctx, "nonexistent-id")
		if err == nil {
			t.Fatal("Expected error when finding non-existent permission, got nil")
		}

		if err != ErrPermissionNotFound {
			t.Errorf("Expected ErrPermissionNotFound, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Add write action to resource
		writeAction, err := resourceDomain.NewAction(resource.ID, "write", "write action")
		if err != nil {
			t.Fatalf("Failed to create write action: %v", err)
		}
		resource.AttachActions([]resourceDomain.Action{*writeAction})
		if err := resourceRepo.Update(ctx, resource); err != nil {
			t.Fatalf("Failed to update resource with write action: %v", err)
		}

		// Create permission
		permID := newPermissionID()
		created, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create permission: %v", err)
		}

		// Update permission
		updated, err := service.Update(ctx, permID, resourceCode, "write", "Write test resource")
		if err != nil {
			t.Fatalf("Failed to update permission: %v", err)
		}

		if updated.Action != "write" {
			t.Errorf("Expected Action 'write', got '%s'", updated.Action)
		}

		if updated.Description != "Write test resource" {
			t.Errorf("Expected Description 'Write test resource', got '%s'", updated.Description)
		}

		if !updated.UpdatedAt.After(created.UpdatedAt) {
			t.Error("Expected UpdatedAt to be after creation time")
		}

		// Verify changes persisted
		found, err := repo.FindByID(ctx, permID)
		if err != nil {
			t.Fatalf("Failed to find updated permission: %v", err)
		}

		if found.Action != "write" {
			t.Errorf("Expected persisted Action 'write', got '%s'", found.Action)
		}

		if found.Description != "Write test resource" {
			t.Errorf("Expected persisted Description 'Write test resource', got '%s'", found.Description)
		}
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		_, err := service.Update(ctx, "nonexistent-id", resourceCode, "read", "Read test resource")
		if err == nil {
			t.Fatal("Expected error when updating non-existent permission, got nil")
		}

		if err != ErrPermissionNotFound {
			t.Errorf("Expected ErrPermissionNotFound, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Create permission
		permID := newPermissionID()
		created, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create permission: %v", err)
		}

		// Delete permission
		err = service.Delete(ctx, created.ID)
		if err != nil {
			t.Fatalf("Failed to delete permission: %v", err)
		}

		// Verify permission is deleted
		_, err = repo.FindByID(ctx, created.ID)
		if err == nil {
			t.Fatal("Expected error when finding deleted permission, got nil")
		}

		if err != ErrPermissionNotFound {
			t.Errorf("Expected ErrPermissionNotFound, got %v", err)
		}
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		err := service.Delete(ctx, "nonexistent-id")
		if err == nil {
			t.Fatal("Expected error when deleting non-existent permission, got nil")
		}

		if err != ErrPermissionNotFound {
			t.Errorf("Expected ErrPermissionNotFound, got %v", err)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Create multiple permissions
		for i := 0; i < 5; i++ {
			permID := newPermissionID()
			_, err := service.Create(ctx, permID, resourceCode, "read", fmt.Sprintf("Permission %d", i))
			if err != nil {
				t.Fatalf("Failed to create permission %d: %v", i, err)
			}
		}

		pageReq := pagination.NewPageRequest(1, 3)

		result, err := service.List(ctx, pageReq)
		if err != nil {
			t.Fatalf("Failed to list permissions: %v", err)
		}

		if result.Total < 5 {
			t.Errorf("Expected at least 5 total permissions, got %d", result.Total)
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
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Create permission
		permID := newPermissionID()
		created, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create permission: %v", err)
		}

		exists, err := repo.Exists(ctx, created.ID)
		if err != nil {
			t.Fatalf("Failed to check existence: %v", err)
		}

		if !exists {
			t.Error("Expected permission to exist, got false")
		}

		exists, err = repo.Exists(ctx, "nonexistent-id")
		if err != nil {
			t.Fatalf("Failed to check existence: %v", err)
		}

		if exists {
			t.Error("Expected permission to not exist, got true")
		}
	})

	t.Run("Timestamps", func(t *testing.T) {
		t.Cleanup(func() { truncatePermissionTables(t) })

		// Create resource first
		resourceCode := "test-" + uuid.New().String()
		resource := createTestResource(t, resourceCode, "test-service")
		if err := resourceRepo.Create(ctx, resource); err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		beforeCreate := time.Now()
		permID := newPermissionID()
		created, err := service.Create(ctx, permID, resourceCode, "read", "Read test resource")
		if err != nil {
			t.Fatalf("Failed to create permission: %v", err)
		}
		afterCreate := time.Now()

		if created.CreatedAt.Before(beforeCreate) || created.CreatedAt.After(afterCreate) {
			t.Error("CreatedAt should be between beforeCreate and afterCreate")
		}

		if created.UpdatedAt.Before(beforeCreate) || created.UpdatedAt.After(afterCreate) {
			t.Error("UpdatedAt should be between beforeCreate and afterCreate")
		}

		// Wait a bit and update
		time.Sleep(10 * time.Millisecond)
		beforeUpdate := time.Now()
		updated, err := service.Update(ctx, permID, resourceCode, "read", "Updated description")
		if err != nil {
			t.Fatalf("Failed to update permission: %v", err)
		}
		afterUpdate := time.Now()

		if updated.UpdatedAt.Before(beforeUpdate) || updated.UpdatedAt.After(afterUpdate) {
			t.Error("UpdatedAt should be updated after update operation")
		}

		// CreatedAt should not change
		timeDiff := updated.CreatedAt.Sub(created.CreatedAt)
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff > time.Second {
			t.Errorf("CreatedAt should not change after update. Difference: %v", timeDiff)
		}
	})
}


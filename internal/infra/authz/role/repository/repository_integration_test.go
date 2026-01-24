//go:build integration
// +build integration

package role

import (
	"context"
	permissionrepo "marketplace/internal/infra/authz/permission/repository"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourcerepo "marketplace/internal/infra/authz/resource/repository"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	dbinfra "marketplace/internal/infra/db"
	"marketplace/pkg/pagination"
	"marketplace/test"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
)

var (
	roleService     *RoleService
	roleRepo        *RoleRepository
	permService     *permissionservice.PermissionService
	permRepo        *permissionrepo.PermissionRepository
	resourceService *resourceservice.Service
	ctx             context.Context
	db              *sql.DB
	cleanup         func()
	setupOnce       sync.Once
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupRoleTestDB(t *testing.T) {
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

		txManager := dbinfra.NewTxManager(db)
		roleRepo = NewRoleRepository(db, txManager)

		// Setup resource and permission services for role service
		resourceRepo := resourcerepo.NewRepository(db, txManager)
		resourceService = resourceservice.NewService(resourceRepo)
		permRepo = permissionrepo.NewPermissionRepository(db)
		permService = permissionservice.NewPermissionService(permRepo, resourceService)

		roleService = NewRoleService(roleRepo, permRepo)
		ctx = context.Background()
	})
}

func newRoleID() string {
	return uuid.New().String()
}

func createTestPermission(t *testing.T, resourceCode, action string) string {
	t.Helper()
	// Create resource first if it doesn't exist
	actionSpecs := []resourceservice.ActionSpec{
		{Name: action, Description: action + " action"},
	}
	_, err := resourceService.Register(ctx, resourceCode, "test-service", "test resource", actionSpecs)
	if err != nil && err != resourceservice.ErrResourceAlreadyExists {
		// If resource exists, try to update it with the action
		_, updateErr := resourceService.Update(ctx, resourceCode, "test-service", "test resource", actionSpecs)
		if updateErr != nil {
			t.Fatalf("Failed to register/update resource: %v", err)
		}
	}

	// Create permission
	permID := uuid.New().String()
	perm, err := permService.Create(ctx, permID, resourceCode, action, fmt.Sprintf("%s %s", action, resourceCode))
	if err != nil {
		t.Fatalf("Failed to create permission: %v", err)
	}
	return perm.ID
}

func TestRoleRepositoryIntegration(t *testing.T) {
	setupRoleTestDB(t)
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
		t.Cleanup(func() { truncateRoleTables(t) })
		permID1 := createTestPermission(t, "users", "read")
		permID2 := createTestPermission(t, "users", "write")

		role, err := roleService.Create(ctx, newRoleID(), "Create Role", "role used for create test", []string{permID1, permID2})
		if err != nil {
			t.Fatalf("failed to create role: %v", err)
		}
		if role == nil {
			t.Fatal("expected role but got nil")
		}
		if role.ID == "" {
			t.Fatal("expected role to have ID")
		}
		if len(role.PermissionIDs) != 2 {
			t.Fatalf("expected 2 permission IDs, got %d", len(role.PermissionIDs))
		}
	})

	t.Run("Create_Duplicate", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		id := newRoleID()
		_, err := roleService.Create(ctx, id, "Duplicate", "duplicate role", []string{})
		if err != nil {
			t.Fatalf("failed to seed duplicate role: %v", err)
		}
		_, err = roleService.Create(ctx, id, "Duplicate", "duplicate role", []string{})
		if err == nil {
			t.Fatal("expected error for duplicate role, got nil")
		}
		if err != ErrRoleAlreadyExists {
			t.Fatalf("expected ErrRoleAlreadyExists, got %v", err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		id := newRoleID()
		permID := createTestPermission(t, "roles", "view")

		_, err := roleService.Create(ctx, id, "Find Role", "role for find", []string{permID})
		if err != nil {
			t.Fatalf("failed to create role: %v", err)
		}

		role, err := roleService.Get(ctx, id)
		if err != nil {
			t.Fatalf("failed to fetch role: %v", err)
		}
		if role.ID != id {
			t.Fatalf("expected ID %s, got %s", id, role.ID)
		}
		if len(role.PermissionIDs) != 1 {
			t.Fatalf("expected 1 permission ID, got %d", len(role.PermissionIDs))
		}
		if role.PermissionIDs[0] != permID {
			t.Fatalf("expected permission ID %s, got %s", permID, role.PermissionIDs[0])
		}
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		_, err := roleService.Get(ctx, newRoleID())
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if err != ErrRoleNotFound {
			t.Fatalf("expected ErrRoleNotFound, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		id := newRoleID()
		permID1 := createTestPermission(t, "roles", "read")

		_, err := roleService.Create(ctx, id, "Role Update", "role before update", []string{permID1})
		if err != nil {
			t.Fatalf("failed to create role: %v", err)
		}

		permID2 := createTestPermission(t, "roles", "edit")
		updatedPermIDs := []string{permID1, permID2}

		updatedRole, err := roleService.Update(ctx, id, "Role Updated", "role after update", updatedPermIDs)
		if err != nil {
			t.Fatalf("failed to update role: %v", err)
		}
		if updatedRole.Name != "Role Updated" || updatedRole.Description != "role after update" {
			t.Fatal("role fields were not updated")
		}
		if len(updatedRole.PermissionIDs) != len(updatedPermIDs) {
			t.Fatalf("expected %d permission IDs, got %d", len(updatedPermIDs), len(updatedRole.PermissionIDs))
		}

		storedRole, err := roleService.Get(ctx, id)
		if err != nil {
			t.Fatalf("failed to read back updated role: %v", err)
		}
		if len(storedRole.PermissionIDs) != len(updatedPermIDs) {
			t.Fatalf("expected persisted permission IDs length %d, got %d", len(updatedPermIDs), len(storedRole.PermissionIDs))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		id := newRoleID()
		permID := createTestPermission(t, "roles", "delete")

		_, err := roleService.Create(ctx, id, "Role Delete", "role for delete", []string{permID})
		if err != nil {
			t.Fatalf("failed to create role: %v", err)
		}

		if err := roleService.Delete(ctx, id); err != nil {
			t.Fatalf("failed to delete role: %v", err)
		}

		if _, err := roleService.Get(ctx, id); err == nil {
			t.Fatal("expected error when fetching deleted role")
		}
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		err := roleService.Delete(ctx, newRoleID())
		if err == nil {
			t.Fatal("expected error for non-existent role")
		}
		if err != ErrRoleNotFound {
			t.Fatalf("expected ErrRoleNotFound, got %v", err)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		for i := 0; i < 5; i++ {
			id := newRoleID()
			permID := createTestPermission(t, "roles", fmt.Sprintf("action-%d", i))

			_, err := roleService.Create(ctx, id, fmt.Sprintf("Role %d", i), fmt.Sprintf("role %d desc", i), []string{permID})
			if err != nil && err != ErrRoleAlreadyExists {
				t.Fatalf("failed to seed roles: %v", err)
			}
		}

		pageReq := pagination.NewPageRequest(1, 3)
		result, err := roleService.List(ctx, pageReq)
		if err != nil {
			t.Fatalf("failed to list roles: %v", err)
		}
		if len(result.Items) > pageReq.Limit {
			t.Fatalf("expected at most %d items, got %d", pageReq.Limit, len(result.Items))
		}
		if result.Total < 5 {
			t.Fatalf("expected at least 5 total roles, got %d", result.Total)
		}
		for _, role := range result.Items {
			if len(role.PermissionIDs) == 0 {
				t.Fatal("expected paged roles to include permission IDs")
			}
		}
	})

	t.Run("Exists", func(t *testing.T) {
		t.Cleanup(func() { truncateRoleTables(t) })
		id := newRoleID()
		_, err := roleService.Create(ctx, id, "Role Exists", "role exists desc", []string{})
		if err != nil && err != ErrRoleAlreadyExists {
			t.Fatalf("failed to create role: %v", err)
		}

		exists, err := roleRepo.Exists(ctx, id)
		if err != nil {
			t.Fatalf("exists returned error: %v", err)
		}
		if !exists {
			t.Fatal("expected role to exist")
		}

		exists, err = roleRepo.Exists(ctx, "role-missing")
		if err != nil {
			t.Fatalf("exists returned error: %v", err)
		}
		if exists {
			t.Fatal("expected role to not exist")
		}
	})
}

func truncateRoleTables(t *testing.T) {
	t.Helper()
	if db == nil {
		t.Skip("Database not initialized, skipping truncate")
		return
	}

	if _, err := db.ExecContext(ctx, `TRUNCATE TABLE role_permissions, roles RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("failed to truncate role tables: %v", err)
	}
}

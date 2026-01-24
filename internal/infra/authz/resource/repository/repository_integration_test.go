//go:build integration
// +build integration

package resource

import (
	"context"
	"marketplace/internal/infra/authz/resource/domain"
	dbinfra "marketplace/internal/infra/db"
	"marketplace/test"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	resourceService *Service
	resourceRepo    *Repository
	ctx             context.Context
	db              *sql.DB
	cleanup         func()
	setupOnce       sync.Once
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupResourceTestDB(t *testing.T) {
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
		resourceRepo = NewRepository(db, txManager)
		resourceService = NewService(resourceRepo)
		ctx = context.Background()
	})
}

func newResourceCode() string {
	return "test-" + uuid.New().String()
}

func TestResourceRepositoryIntegration(t *testing.T) {
	setupResourceTestDB(t)
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
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")
		actions := []domain.Action{
			domain.NewAction(resource.ID, "read", "read action"),
			domain.NewAction(resource.ID, "write", "write action"),
		}
		resource.AttachActions(actions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		// Verify resource was created
		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find created resource: %v", err)
		}
		if found == nil {
			t.Fatal("expected resource but got nil")
		}
		if found.Code != code {
			t.Errorf("expected code '%s', got '%s'", code, found.Code)
		}
		if found.Service != "test-service" {
			t.Errorf("expected service 'test-service', got '%s'", found.Service)
		}
		if len(found.Actions) != 2 {
			t.Fatalf("expected 2 actions, got %d", len(found.Actions))
		}
	})

	t.Run("Create_Duplicate", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create first resource: %v", err)
		}

		// Try to create duplicate - repository returns database error, not business logic error
		duplicate := domain.NewResource(code, "another-service", "another description")
		err = resourceRepo.Create(ctx, duplicate)
		if err == nil {
			t.Fatal("expected database error for duplicate resource, got nil")
		}
		// Repository returns raw database error, validation is handled at service level
	})

	t.Run("FindByCode", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "find-service", "find description")
		actions := []domain.Action{
			domain.NewAction(resource.ID, "read", "read action"),
		}
		resource.AttachActions(actions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find resource: %v", err)
		}
		if found == nil {
			t.Fatal("expected resource but got nil")
		}
		if found.Code != code {
			t.Errorf("expected code '%s', got '%s'", code, found.Code)
		}
		if found.Service != "find-service" {
			t.Errorf("expected service 'find-service', got '%s'", found.Service)
		}
		if len(found.Actions) != 1 {
			t.Fatalf("expected 1 action, got %d", len(found.Actions))
		}
		if found.Actions[0].Name != "read" {
			t.Errorf("expected action name 'read', got '%s'", found.Actions[0].Name)
		}
	})

	t.Run("FindByCode_CaseInsensitive", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := "TestResource"
		resource := domain.NewResource(code, "test-service", "test description")

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		// Try finding with different case
		found, err := resourceRepo.FindByCode(ctx, "testresource")
		if err != nil {
			t.Fatalf("failed to find resource with different case: %v", err)
		}
		if found == nil {
			t.Fatal("expected resource but got nil")
		}
		if found.Code != "testresource" {
			t.Errorf("expected normalized code 'testresource', got '%s'", found.Code)
		}
	})

	t.Run("FindByCode_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		_, err := resourceRepo.FindByCode(ctx, "nonexistent")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "original-service", "original description")
		actions := []domain.Action{
			domain.NewAction(resource.ID, "read", "read action"),
		}
		resource.AttachActions(actions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		// Update resource
		resource.Service = "updated-service"
		resource.Description = "updated description"
		resource.UpdatedAt = time.Now()
		newActions := []domain.Action{
			domain.NewAction(resource.ID, "read", "updated read action"),
			domain.NewAction(resource.ID, "write", "write action"),
		}
		resource.AttachActions(newActions)

		err = resourceRepo.Update(ctx, resource)
		if err != nil {
			t.Fatalf("failed to update resource: %v", err)
		}

		// Verify update
		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find updated resource: %v", err)
		}
		if found.Service != "updated-service" {
			t.Errorf("expected service 'updated-service', got '%s'", found.Service)
		}
		if found.Description != "updated description" {
			t.Errorf("expected description 'updated description', got '%s'", found.Description)
		}
		if len(found.Actions) != 2 {
			t.Fatalf("expected 2 actions, got %d", len(found.Actions))
		}
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		resource := domain.NewResource("nonexistent", "service", "description")
		// Repository doesn't validate existence - that's handled at service level
		err := resourceRepo.Update(ctx, resource)
		if err != nil {
			t.Fatalf("repository should not return error for non-existent resource (validation is at service level), got %v", err)
		}
	})

	t.Run("Deprecate", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		now := time.Now()
		err = resourceRepo.Deprecate(ctx, resource.ID, &now)
		if err != nil {
			t.Fatalf("failed to deprecate resource: %v", err)
		}

		// Verify deprecation
		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find deprecated resource: %v", err)
		}
		if found.DeprecatedAt == nil {
			t.Fatal("expected DeprecatedAt to be set")
		}
		if !found.IsDeprecated() {
			t.Error("expected resource to be deprecated")
		}
	})

	t.Run("Deprecate_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		now := time.Now()
		// Repository doesn't validate existence - that's handled at service level
		err := resourceRepo.Deprecate(ctx, "nonexistent-id", &now)
		if err != nil {
			t.Fatalf("repository should not return error for non-existent resource (validation is at service level), got %v", err)
		}
	})

	t.Run("List", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		// Create multiple resources
		for i := 0; i < 3; i++ {
			code := newResourceCode()
			resource := domain.NewResource(code, fmt.Sprintf("service-%d", i), fmt.Sprintf("description %d", i))
			actions := []domain.Action{
				domain.NewAction(resource.ID, fmt.Sprintf("action-%d", i), fmt.Sprintf("action %d", i)),
			}
			resource.AttachActions(actions)
			if err := resourceRepo.Create(ctx, resource); err != nil {
				t.Fatalf("failed to create resource %d: %v", i, err)
			}
		}

		resources, err := resourceRepo.List(ctx)
		if err != nil {
			t.Fatalf("failed to list resources: %v", err)
		}
		if len(resources) < 3 {
			t.Errorf("expected at least 3 resources, got %d", len(resources))
		}

		// Verify resources are sorted by code
		for i := 1; i < len(resources); i++ {
			if resources[i].Code < resources[i-1].Code {
				t.Errorf("resources not sorted by code: %s should come before %s", resources[i-1].Code, resources[i].Code)
			}
		}

		// Verify actions are included
		for _, resource := range resources {
			if len(resource.Actions) == 0 {
				t.Errorf("expected resource '%s' to have actions", resource.Code)
			}
		}
	})

	t.Run("FindAction", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")
		actions := []domain.Action{
			domain.NewAction(resource.ID, "read", "read action"),
			domain.NewAction(resource.ID, "write", "write action"),
		}
		resource.AttachActions(actions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		action, err := resourceRepo.FindAction(ctx, code, "read")
		if err != nil {
			t.Fatalf("failed to find action: %v", err)
		}
		if action == nil {
			t.Fatal("expected action but got nil")
		}
		if action.Name != "read" {
			t.Errorf("expected action name 'read', got '%s'", action.Name)
		}
		if action.ResourceID != resource.ID {
			t.Errorf("expected ResourceID '%s', got '%s'", resource.ID, action.ResourceID)
		}
	})

	t.Run("FindAction_CaseInsensitive", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")
		actions := []domain.Action{
			domain.NewAction(resource.ID, "read", "read action"),
		}
		resource.AttachActions(actions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		// Try finding with different case
		action, err := resourceRepo.FindAction(ctx, code, "READ")
		if err != nil {
			t.Fatalf("failed to find action with different case: %v", err)
		}
		if action == nil {
			t.Fatal("expected action but got nil")
		}
		if action.Name != "read" {
			t.Errorf("expected normalized action name 'read', got '%s'", action.Name)
		}
	})

	t.Run("FindAction_NotFound", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		_, err = resourceRepo.FindAction(ctx, code, "nonexistent")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("Create_WithActions", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")
		actions := []domain.Action{
			domain.NewAction(resource.ID, "create", "create action"),
			domain.NewAction(resource.ID, "read", "read action"),
			domain.NewAction(resource.ID, "update", "update action"),
			domain.NewAction(resource.ID, "delete", "delete action"),
		}
		resource.AttachActions(actions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource with actions: %v", err)
		}

		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find resource: %v", err)
		}
		if len(found.Actions) != 4 {
			t.Fatalf("expected 4 actions, got %d", len(found.Actions))
		}

		// Verify actions are sorted
		expectedOrder := []string{"create", "delete", "read", "update"}
		for i, expected := range expectedOrder {
			if found.Actions[i].Name != expected {
				t.Errorf("expected action at index %d to be '%s', got '%s'", i, expected, found.Actions[i].Name)
			}
		}
	})

	t.Run("Update_WithNewActions", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		resource := domain.NewResource(code, "test-service", "test description")
		initialActions := []domain.Action{
			domain.NewAction(resource.ID, "read", "read action"),
		}
		resource.AttachActions(initialActions)

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		// Update with new actions
		resource.UpdatedAt = time.Now()
		newActions := []domain.Action{
			domain.NewAction(resource.ID, "read", "updated read action"),
			domain.NewAction(resource.ID, "write", "write action"),
			domain.NewAction(resource.ID, "delete", "delete action"),
		}
		resource.AttachActions(newActions)

		err = resourceRepo.Update(ctx, resource)
		if err != nil {
			t.Fatalf("failed to update resource: %v", err)
		}

		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find updated resource: %v", err)
		}
		if len(found.Actions) != 3 {
			t.Fatalf("expected 3 actions, got %d", len(found.Actions))
		}
		if found.Actions[0].Description != "updated read action" {
			t.Errorf("expected updated description for read action, got '%s'", found.Actions[0].Description)
		}
	})

	t.Run("Timestamps", func(t *testing.T) {
		t.Cleanup(func() { truncateResourceTables(t) })
		code := newResourceCode()
		beforeCreate := time.Now()
		resource := domain.NewResource(code, "test-service", "test description")

		err := resourceRepo.Create(ctx, resource)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}
		afterCreate := time.Now()

		found, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find resource: %v", err)
		}

		if found.CreatedAt.Before(beforeCreate) || found.CreatedAt.After(afterCreate) {
			t.Error("CreatedAt should be between beforeCreate and afterCreate")
		}

		// Wait a bit and update
		time.Sleep(10 * time.Millisecond)
		beforeUpdate := time.Now()
		resource.Service = "updated-service"
		resource.UpdatedAt = time.Now()
		err = resourceRepo.Update(ctx, resource)
		if err != nil {
			t.Fatalf("failed to update resource: %v", err)
		}
		afterUpdate := time.Now()

		updated, err := resourceRepo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("failed to find updated resource: %v", err)
		}

		if updated.UpdatedAt.Before(beforeUpdate) || updated.UpdatedAt.After(afterUpdate) {
			t.Error("UpdatedAt should be updated after update operation")
		}

		// CreatedAt should not change
		if !updated.CreatedAt.Equal(found.CreatedAt) {
			t.Error("CreatedAt should not change after update")
		}
	})
}

func truncateResourceTables(t *testing.T) {
	t.Helper()
	if db == nil {
		t.Skip("Database not initialized, skipping truncate")
		return
	}

	if _, err := db.ExecContext(ctx, `TRUNCATE TABLE resource_actions, resources RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("failed to truncate resource tables: %v", err)
	}
}

package role

import (
	"context"
	permissionDomain "marketplace/internal/infra/authz/permission/domain"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	"marketplace/internal/infra/authz/role/domain"
	"marketplace/pkg/pagination"
	"errors"
	"strings"
	"testing"
	"time"
)

var testResourceRef = permissionDomain.ResourceRef{
	ID:      "res-users",
	Code:    "users",
	Service: "user-management",
}

func newPermission(id, action string) *permissionDomain.Permission {
	perm, err := permissionDomain.NewPermissionEntity(id, testResourceRef, action, action+" permission")
	if err != nil {
		panic(err)
	}
	return perm
}

type mockRoleRepository struct {
	roles        map[string]*domain.Role
	createFunc   func(context.Context, *domain.Role) error
	updateFunc   func(context.Context, *domain.Role) error
	findByIDFunc func(context.Context, string) (*domain.Role, error)
	deleteFunc   func(context.Context, string) error
	findAllFunc  func(context.Context, pagination.PageRequest) (pagination.PageResult[*domain.Role], error)
	existsFunc   func(context.Context, string) (bool, error)
}

func newMockRoleRepository() *mockRoleRepository {
	return &mockRoleRepository{
		roles: make(map[string]*domain.Role),
	}
}

type mockPermissionLookup struct {
	perms map[string]*permissionDomain.Permission
}

func newMockPermissionLookup(perms ...*permissionDomain.Permission) *mockPermissionLookup {
	m := &mockPermissionLookup{perms: make(map[string]*permissionDomain.Permission)}
	for _, perm := range perms {
		m.perms[perm.ID] = perm
	}
	return m
}

func (m *mockPermissionLookup) FindByID(ctx context.Context, id string) (*permissionDomain.Permission, error) {
	if perm, ok := m.perms[id]; ok {
		return perm, nil
	}
	return nil, permissionservice.ErrPermissionNotFound
}

func (m *mockRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, role)
	}
	if _, ok := m.roles[role.ID]; ok {
		return ErrRoleAlreadyExists
	}
	m.roles[role.ID] = role
	return nil
}

func (m *mockRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, role)
	}
	if _, ok := m.roles[role.ID]; !ok {
		return ErrRoleNotFound
	}
	m.roles[role.ID] = role
	return nil
}

func (m *mockRoleRepository) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	role, ok := m.roles[id]
	if !ok {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (m *mockRoleRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	if _, ok := m.roles[id]; !ok {
		return ErrRoleNotFound
	}
	delete(m.roles, id)
	return nil
}

func (m *mockRoleRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Role], error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, pageReq)
	}
	roles := make([]*domain.Role, 0, len(m.roles))
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	total := len(roles)
	if total == 0 || pageReq.Offset >= total {
		return pagination.NewPageResult([]*domain.Role{}, total, pageReq), nil
	}

	end := pageReq.Offset + pageReq.Limit
	if end > total {
		end = total
	}

	return pagination.NewPageResult(roles[pageReq.Offset:end], total, pageReq), nil
}

func (m *mockRoleRepository) Exists(ctx context.Context, id string) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	_, ok := m.roles[id]
	return ok, nil
}

func (m *mockRoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	for _, role := range m.roles {
		if strings.EqualFold(role.Name, name) {
			return role, nil
		}
	}
	return nil, ErrRoleNotFound
}

func createTestRole(id, name string) *domain.Role {
	return &domain.Role{
		ID:            id,
		Name:          name,
		Description:   name + " description",
		PermissionIDs: []string{"perm-" + id},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func TestRoleServiceCreate(t *testing.T) {
	tests := []struct {
		name        string
		roleID      string
		setup       func(*mockRoleRepository)
		wantErr     bool
		errExpected error
	}{
		{
			name:   "successful creation",
			roleID: "role-1",
			setup:  func(m *mockRoleRepository) {},
		},
		{
			name:   "duplicate role",
			roleID: "role-1",
			setup: func(m *mockRoleRepository) {
				m.roles["role-1"] = createTestRole("role-1", "Admin")
			},
			wantErr:     true,
			errExpected: ErrRoleAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRoleRepository()
			tt.setup(repo)
			perm := newPermission("perm-"+tt.roleID, "create")
			lookup := newMockPermissionLookup(perm)
			service := NewRoleService(repo, lookup)

			ctx := context.Background()
			role, err := service.Create(ctx, tt.roleID, "Admin", "Admin role", []string{perm.ID})

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if !errors.Is(err, tt.errExpected) {
					t.Fatalf("expected error %v, got %v", tt.errExpected, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if role == nil || role.ID != tt.roleID {
				t.Fatalf("expected role with ID %s", tt.roleID)
			}

			if _, exists := repo.roles[tt.roleID]; !exists {
				t.Fatalf("role was not persisted in repository")
			}
		})
	}
}

func TestRoleServiceGet(t *testing.T) {
	tests := []struct {
		name    string
		roleID  string
		setup   func(*mockRoleRepository)
		wantErr bool
	}{
		{
			name:   "role found",
			roleID: "role-1",
			setup: func(m *mockRoleRepository) {
				m.roles["role-1"] = createTestRole("role-1", "Admin")
			},
		},
		{
			name:    "role missing",
			roleID:  "role-2",
			setup:   func(m *mockRoleRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRoleRepository()
			tt.setup(repo)
			service := NewRoleService(repo, newMockPermissionLookup())

			ctx := context.Background()
			role, err := service.Get(ctx, tt.roleID)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if role.ID != tt.roleID {
				t.Fatalf("expected ID %s, got %s", tt.roleID, role.ID)
			}
		})
	}
}

func TestRoleServiceUpdate(t *testing.T) {
	tests := []struct {
		name            string
		roleID          string
		setup           func(*mockRoleRepository)
		newName         string
		newDesc         string
		newPermIDs      []string
		expectedActions []string
		wantErr         bool
		errExpected     error
	}{
		{
			name:            "update success",
			roleID:          "role-1",
			newName:         "Updated",
			newDesc:         "Updated desc",
			newPermIDs:      []string{"perm-update"},
			expectedActions: []string{"view"},
			setup: func(m *mockRoleRepository) {
				m.roles["role-1"] = createTestRole("role-1", "Admin")
			},
		},
		{
			name:    "role missing",
			roleID:  "role-2",
			setup:   func(m *mockRoleRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRoleRepository()
			tt.setup(repo)
			var lookup *mockPermissionLookup
			if len(tt.newPermIDs) > 0 {
				perms := make([]*permissionDomain.Permission, len(tt.newPermIDs))
				for i, id := range tt.newPermIDs {
					perms[i] = newPermission(id, tt.expectedActions[i])
				}
				lookup = newMockPermissionLookup(perms...)
			} else {
				lookup = newMockPermissionLookup()
			}
			service := NewRoleService(repo, lookup)

			ctx := context.Background()
			role, err := service.Update(ctx, tt.roleID, tt.newName, tt.newDesc, tt.newPermIDs)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if role.Name != tt.newName {
				t.Fatalf("expected name %s, got %s", tt.newName, role.Name)
			}

			if role.Description != tt.newDesc {
				t.Fatalf("expected description %s, got %s", tt.newDesc, role.Description)
			}

			if len(role.PermissionIDs) != len(tt.newPermIDs) {
				t.Fatalf("expected %d permission IDs, got %d", len(tt.newPermIDs), len(role.PermissionIDs))
			}
			// Verify permission IDs match
			permIDMap := make(map[string]bool)
			for _, id := range role.PermissionIDs {
				permIDMap[id] = true
			}
			for _, expectedID := range tt.newPermIDs {
				if !permIDMap[expectedID] {
					t.Fatalf("expected permission ID %s not found in role", expectedID)
				}
			}
		})
	}
}

func TestRoleServiceDelete(t *testing.T) {
	tests := []struct {
		name    string
		roleID  string
		setup   func(*mockRoleRepository)
		wantErr bool
	}{
		{
			name:   "delete success",
			roleID: "role-1",
			setup: func(m *mockRoleRepository) {
				m.roles["role-1"] = createTestRole("role-1", "Admin")
			},
		},
		{
			name:    "role missing",
			roleID:  "role-2",
			setup:   func(m *mockRoleRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRoleRepository()
			tt.setup(repo)
			service := NewRoleService(repo, newMockPermissionLookup())

			ctx := context.Background()
			err := service.Delete(ctx, tt.roleID)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if _, exists := repo.roles[tt.roleID]; exists {
				t.Fatalf("role should have been deleted")
			}
		})
	}
}

func TestRoleServiceList(t *testing.T) {
	repo := newMockRoleRepository()
	repo.roles["role-1"] = createTestRole("role-1", "Role 1")
	repo.roles["role-2"] = createTestRole("role-2", "Role 2")
	repo.roles["role-3"] = createTestRole("role-3", "Role 3")

	service := NewRoleService(repo, newMockPermissionLookup())
	ctx := context.Background()

	pageReq := pagination.NewPageRequest(1, 2)
	result, err := service.List(ctx, pageReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(result.Items))
	}

	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}

	if !result.HasNext {
		t.Fatalf("expected HasNext to be true")
	}

	// request page beyond range
	pageReq = pagination.NewPageRequest(3, 2)
	result, err = service.List(ctx, pageReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 0 {
		t.Fatalf("expected empty result for out-of-range page")
	}
}

package permission

import (
	"context"
	"marketplace/internal/infra/authz/permission/domain"
	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	"marketplace/pkg/pagination"
	"errors"
	"strings"
	"testing"
	"time"
)

// mockPermissionRepository is a mock implementation of permissionRepository for testing
type mockPermissionRepository struct {
	permissions  map[string]*domain.Permission
	createFunc   func(context.Context, *domain.Permission) error
	updateFunc   func(context.Context, *domain.Permission) error
	findByIDFunc func(context.Context, string) (*domain.Permission, error)
	deleteFunc   func(context.Context, string) error
	findAllFunc  func(context.Context, pagination.PageRequest) (pagination.PageResult[*domain.Permission], error)
	existsFunc   func(context.Context, string) (bool, error)
}

func newMockPermissionRepository() *mockPermissionRepository {
	return &mockPermissionRepository{
		permissions: make(map[string]*domain.Permission),
	}
}

func (m *mockPermissionRepository) Create(ctx context.Context, perm *domain.Permission) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, perm)
	}
	if _, exists := m.permissions[perm.ID]; exists {
		return ErrPermissionAlreadyExists
	}
	m.permissions[perm.ID] = perm
	return nil
}

func (m *mockPermissionRepository) Update(ctx context.Context, perm *domain.Permission) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, perm)
	}
	if _, exists := m.permissions[perm.ID]; !exists {
		return ErrPermissionNotFound
	}
	m.permissions[perm.ID] = perm
	return nil
}

func (m *mockPermissionRepository) FindByID(ctx context.Context, id string) (*domain.Permission, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	perm, exists := m.permissions[id]
	if !exists {
		return nil, ErrPermissionNotFound
	}
	return perm, nil
}

func (m *mockPermissionRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	if _, exists := m.permissions[id]; !exists {
		return ErrPermissionNotFound
	}
	delete(m.permissions, id)
	return nil
}

func (m *mockPermissionRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Permission], error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, pageReq)
	}
	perms := make([]*domain.Permission, 0, len(m.permissions))
	for _, perm := range m.permissions {
		perms = append(perms, perm)
	}
	total := len(perms)

	if pageReq.Offset >= total {
		return pagination.NewPageResult([]*domain.Permission{}, total, pageReq), nil
	}

	end := pageReq.Offset + pageReq.Limit
	if end > total {
		end = total
	}

	paginatedPerms := perms[pageReq.Offset:end]
	return pagination.NewPageResult(paginatedPerms, total, pageReq), nil
}

func (m *mockPermissionRepository) Exists(ctx context.Context, id string) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	_, exists := m.permissions[id]
	return exists, nil
}

func (m *mockPermissionRepository) FindByResourceAndAction(ctx context.Context, resourceCode, action string) (*domain.Permission, error) {
	for _, perm := range m.permissions {
		if perm.Resource.Code == resourceCode && perm.Action == action {
			return perm, nil
		}
	}
	return nil, ErrPermissionNotFound
}

// mockResourceCatalog is a mock implementation of resourceCatalog for testing
type mockResourceCatalog struct {
	resolveActionFunc func(context.Context, string, string) (*resourceDomain.Resource, *resourceDomain.Action, error)
}

func newMockResourceCatalog() *mockResourceCatalog {
	return &mockResourceCatalog{}
}

func (m *mockResourceCatalog) ResolveAction(ctx context.Context, code, action string) (*resourceDomain.Resource, *resourceDomain.Action, error) {
	if m.resolveActionFunc != nil {
		return m.resolveActionFunc(ctx, code, action)
	}
	// Default: return a valid resource and action
	resource, err := resourceDomain.NewResource(code, "test-service", "test description")
	if err != nil {
		return nil, nil, err
	}
	act, err := resourceDomain.NewAction(resource.ID, action, action+" description")
	if err != nil {
		return nil, nil, err
	}
	resource.AttachActions([]resourceDomain.Action{act})
	return resource, &act, nil
}

func createTestPermission(id, resourceCode, action, description string) *domain.Permission {
	now := time.Now()
	return &domain.Permission{
		ID: id,
		Resource: domain.ResourceRef{
			ID:          "resource-id",
			Code:        resourceCode,
			Service:     "test-service",
			Description: "test resource",
		},
		Action:      action,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestPermissionServiceCreate(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		resourceCode string
		action       string
		description  string
		setupMock    func(*mockPermissionRepository, *mockResourceCatalog)
		wantErr      bool
		errContains  string
	}{
		{
			testName:     "successful creation",
			id:           "perm1",
			resourceCode: "users",
			action:       "read",
			description:  "Read users",
			setupMock: func(m *mockPermissionRepository, rc *mockResourceCatalog) {
				// No setup needed - permission doesn't exist and resource resolves successfully
			},
			wantErr: false,
		},
		{
			testName:     "permission already exists",
			id:           "perm1",
			resourceCode: "users",
			action:       "read",
			description:  "Read users",
			setupMock: func(m *mockPermissionRepository, rc *mockResourceCatalog) {
				existing := createTestPermission("perm1", "users", "read", "Existing")
				m.permissions["perm1"] = existing
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			testName:     "resource not found",
			id:           "perm1",
			resourceCode: "nonexistent",
			action:       "read",
			description:  "Read users",
			setupMock: func(m *mockPermissionRepository, rc *mockResourceCatalog) {
				rc.resolveActionFunc = func(ctx context.Context, code, action string) (*resourceDomain.Resource, *resourceDomain.Action, error) {
					return nil, nil, errors.New("resource not found")
				}
			},
			wantErr:     true,
			errContains: "resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockRepo := newMockPermissionRepository()
			mockResourceCat := newMockResourceCatalog()
			tt.setupMock(mockRepo, mockResourceCat)
			service := NewPermissionService(mockRepo, mockResourceCat)

			ctx := context.Background()
			perm, err := service.Create(ctx, tt.id, tt.resourceCode, tt.action, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if perm == nil {
				t.Errorf("Expected permission but got nil")
				return
			}

			if perm.ID != tt.id {
				t.Errorf("Expected ID %s, got %s", tt.id, perm.ID)
			}

			if perm.Resource.Code != tt.resourceCode {
				t.Errorf("Expected Resource.Code %s, got %s", tt.resourceCode, perm.Resource.Code)
			}

			if perm.Action != tt.action {
				t.Errorf("Expected Action %s, got %s", tt.action, perm.Action)
			}

			if perm.Description != tt.description {
				t.Errorf("Expected Description %s, got %s", tt.description, perm.Description)
			}
		})
	}
}

func TestPermissionServiceGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockPermissionRepository)
		wantErr   bool
	}{
		{
			name: "successful get",
			id:   "perm1",
			setupMock: func(m *mockPermissionRepository) {
				perm := createTestPermission("perm1", "users", "read", "Read users")
				m.permissions["perm1"] = perm
			},
			wantErr: false,
		},
		{
			name:      "permission not found",
			id:        "nonexistent",
			setupMock: func(m *mockPermissionRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockPermissionRepository()
			tt.setupMock(mockRepo)
			mockResourceCat := newMockResourceCatalog()
			service := NewPermissionService(mockRepo, mockResourceCat)

			ctx := context.Background()
			perm, err := service.Get(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if perm == nil {
				t.Errorf("Expected permission but got nil")
				return
			}

			if perm.ID != tt.id {
				t.Errorf("Expected ID %s, got %s", tt.id, perm.ID)
			}
		})
	}
}

func TestPermissionServiceUpdate(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		resourceCode string
		action       string
		description  string
		setupMock    func(*mockPermissionRepository, *mockResourceCatalog)
		wantErr      bool
		errContains  string
	}{
		{
			name:         "successful update",
			id:           "perm1",
			resourceCode: "users",
			action:       "write",
			description:  "Write users",
			setupMock: func(m *mockPermissionRepository, rc *mockResourceCatalog) {
				perm := createTestPermission("perm1", "users", "read", "Read users")
				m.permissions["perm1"] = perm
			},
			wantErr: false,
		},
		{
			name:         "permission not found",
			id:           "nonexistent",
			resourceCode: "users",
			action:       "read",
			description:  "Read users",
			setupMock:    func(m *mockPermissionRepository, rc *mockResourceCatalog) {},
			wantErr:      true,
		},
		{
			name:         "resource not found",
			id:           "perm1",
			resourceCode: "nonexistent",
			action:       "read",
			description:  "Read users",
			setupMock: func(m *mockPermissionRepository, rc *mockResourceCatalog) {
				perm := createTestPermission("perm1", "users", "read", "Read users")
				m.permissions["perm1"] = perm
				rc.resolveActionFunc = func(ctx context.Context, code, action string) (*resourceDomain.Resource, *resourceDomain.Action, error) {
					return nil, nil, errors.New("resource not found")
				}
			},
			wantErr:     true,
			errContains: "resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockPermissionRepository()
			mockResourceCat := newMockResourceCatalog()
			tt.setupMock(mockRepo, mockResourceCat)
			service := NewPermissionService(mockRepo, mockResourceCat)

			ctx := context.Background()
			perm, err := service.Update(ctx, tt.id, tt.resourceCode, tt.action, tt.description)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if perm == nil {
				t.Errorf("Expected permission but got nil")
				return
			}

			if perm.Action != tt.action {
				t.Errorf("Expected Action %s, got %s", tt.action, perm.Action)
			}

			if perm.Description != tt.description {
				t.Errorf("Expected Description %s, got %s", tt.description, perm.Description)
			}
		})
	}
}

func TestPermissionServiceDelete(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockPermissionRepository)
		wantErr   bool
	}{
		{
			name: "successful delete",
			id:   "perm1",
			setupMock: func(m *mockPermissionRepository) {
				perm := createTestPermission("perm1", "users", "read", "Read users")
				m.permissions["perm1"] = perm
			},
			wantErr: false,
		},
		{
			name:      "permission not found",
			id:        "nonexistent",
			setupMock: func(m *mockPermissionRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockPermissionRepository()
			tt.setupMock(mockRepo)
			mockResourceCat := newMockResourceCatalog()
			service := NewPermissionService(mockRepo, mockResourceCat)

			ctx := context.Background()
			err := service.Delete(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify permission was deleted
			ctx2 := context.Background()
			exists, _ := mockRepo.Exists(ctx2, tt.id)
			if exists {
				t.Errorf("Permission should have been deleted but still exists")
			}
		})
	}
}

func TestPermissionServiceList(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		limit     int
		setupMock func(*mockPermissionRepository)
		wantErr   bool
		validate  func(*testing.T, pagination.PageResult[*domain.Permission])
	}{
		{
			name:  "successful list",
			page:  1,
			limit: 10,
			setupMock: func(m *mockPermissionRepository) {
				for i := 0; i < 5; i++ {
					perm := createTestPermission("perm"+string(rune(i)), "users", "read", "Read users")
					m.permissions["perm"+string(rune(i))] = perm
				}
			},
			wantErr: false,
			validate: func(t *testing.T, result pagination.PageResult[*domain.Permission]) {
				if len(result.Items) != 5 {
					t.Errorf("Expected 5 items, got %d", len(result.Items))
				}
				if result.Total != 5 {
					t.Errorf("Expected total 5, got %d", result.Total)
				}
			},
		},
		{
			name:  "empty list",
			page:  1,
			limit: 10,
			setupMock: func(m *mockPermissionRepository) {
				// No permissions
			},
			wantErr: false,
			validate: func(t *testing.T, result pagination.PageResult[*domain.Permission]) {
				if len(result.Items) != 0 {
					t.Errorf("Expected 0 items, got %d", len(result.Items))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockPermissionRepository()
			tt.setupMock(mockRepo)
			mockResourceCat := newMockResourceCatalog()
			service := NewPermissionService(mockRepo, mockResourceCat)

			ctx := context.Background()
			pageReq := pagination.NewPageRequest(tt.page, tt.limit)
			result, err := service.List(ctx, pageReq)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

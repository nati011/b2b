package resource

import (
	"context"
	"marketplace/internal/infra/authz/resource/domain"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"
)

// mockResourceRepository is a mock implementation of resourceRepository for testing
type mockResourceRepository struct {
	resources      map[string]*domain.Resource
	actions        map[string]*domain.Action // key: "code:action"
	createFunc     func(context.Context, *domain.Resource) error
	updateFunc     func(context.Context, *domain.Resource) error
	deprecateFunc  func(context.Context, string, *time.Time) error
	findByCodeFunc func(context.Context, string) (*domain.Resource, error)
	listFunc       func(context.Context) ([]*domain.Resource, error)
	findActionFunc func(context.Context, string, string) (*domain.Action, error)
}

func newMockResourceRepository() *mockResourceRepository {
	return &mockResourceRepository{
		resources: make(map[string]*domain.Resource),
		actions:   make(map[string]*domain.Action),
	}
}

func (m *mockResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, resource)
	}
	if _, exists := m.resources[resource.Code]; exists {
		return ErrResourceAlreadyExists
	}
	m.resources[resource.Code] = resource
	for _, action := range resource.Actions {
		key := resource.Code + ":" + action.Name
		m.actions[key] = &action
	}
	return nil
}

func (m *mockResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, resource)
	}
	if _, exists := m.resources[resource.Code]; !exists {
		return ErrResourceNotFound
	}
	m.resources[resource.Code] = resource
	for _, action := range resource.Actions {
		key := resource.Code + ":" + action.Name
		m.actions[key] = &action
	}
	return nil
}

func (m *mockResourceRepository) Deprecate(ctx context.Context, id string, deprecatedAt *time.Time) error {
	if m.deprecateFunc != nil {
		return m.deprecateFunc(ctx, id, deprecatedAt)
	}
	for _, resource := range m.resources {
		if resource.ID == id {
			resource.DeprecatedAt = deprecatedAt
			return nil
		}
	}
	return ErrResourceNotFound
}

func (m *mockResourceRepository) FindByCode(ctx context.Context, code string) (*domain.Resource, error) {
	if m.findByCodeFunc != nil {
		return m.findByCodeFunc(ctx, code)
	}
	resource, exists := m.resources[code]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return resource, nil
}

func (m *mockResourceRepository) List(ctx context.Context) ([]*domain.Resource, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	resources := make([]*domain.Resource, 0, len(m.resources))
	for _, resource := range m.resources {
		resources = append(resources, resource)
	}
	return resources, nil
}

func (m *mockResourceRepository) FindAction(ctx context.Context, code, action string) (*domain.Action, error) {
	if m.findActionFunc != nil {
		return m.findActionFunc(ctx, code, action)
	}
	key := code + ":" + action
	act, exists := m.actions[key]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return act, nil
}

func TestServiceRegister(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		service     string
		description string
		actions     []ActionSpec
		setupMock   func(*mockResourceRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *domain.Resource)
	}{
		{
			name:        "successfully registers new resource",
			code:        "users",
			service:     "user-service",
			description: "User management resource",
			actions: []ActionSpec{
				{Name: "read", Description: "read users"},
				{Name: "write", Description: "write users"},
			},
			setupMock: func(m *mockResourceRepository) {},
			wantErr:   false,
			validate: func(t *testing.T, r *domain.Resource) {
				if r == nil {
					t.Fatal("Expected resource but got nil")
				}
				if r.Code != "users" {
					t.Errorf("Expected code 'users', got '%s'", r.Code)
				}
				if r.Service != "user-service" {
					t.Errorf("Expected service 'user-service', got '%s'", r.Service)
				}
				if len(r.Actions) != 2 {
					t.Errorf("Expected 2 actions, got %d", len(r.Actions))
				}
			},
		},
		{
			name:        "fails when resource already exists",
			code:        "users",
			service:     "user-service",
			description: "User management resource",
			actions:     []ActionSpec{},
			setupMock: func(m *mockResourceRepository) {
				existing, err := domain.NewResource("users", "existing", "existing")
				if err != nil {
					t.Fatalf("failed to create existing resource: %v", err)
				}
				m.resources["users"] = existing
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:        "handles repository error",
			code:        "users",
			service:     "user-service",
			description: "User management resource",
			actions:     []ActionSpec{},
			setupMock: func(m *mockResourceRepository) {
				m.findByCodeFunc = func(ctx context.Context, code string) (*domain.Resource, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:     true,
			errContains: "database error",
		},
		{
			name:        "normalizes code and action names",
			code:        "  USER_MANAGEMENT  ",
			service:     "user-service",
			description: "User management",
			actions: []ActionSpec{
				{Name: "  CREATE  ", Description: "create action"},
			},
			setupMock: func(m *mockResourceRepository) {},
			wantErr:   false,
			validate: func(t *testing.T, r *domain.Resource) {
				if r.Code != "user_management" {
					t.Errorf("Expected normalized code 'user_management', got '%s'", r.Code)
				}
				if len(r.Actions) > 0 && r.Actions[0].Name != "create" {
					t.Errorf("Expected normalized action name 'create', got '%s'", r.Actions[0].Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockResourceRepository()
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			ctx := context.Background()
			resource, err := service.Register(ctx, tt.code, tt.service, tt.description, tt.actions)

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

			if tt.validate != nil {
				tt.validate(t, resource)
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		service     string
		description string
		actions     []ActionSpec
		setupMock   func(*mockResourceRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *domain.Resource)
	}{
		{
			name:        "successfully updates existing resource",
			code:        "users",
			service:     "updated-service",
			description: "Updated description",
			actions: []ActionSpec{
				{Name: "read", Description: "read users"},
			},
			setupMock: func(m *mockResourceRepository) {
				existing, err := domain.NewResource("users", "original", "original")
				if err != nil {
					t.Fatalf("failed to create existing resource: %v", err)
				}
				m.resources["users"] = existing
			},
			wantErr: false,
			validate: func(t *testing.T, r *domain.Resource) {
				if r.Service != "updated-service" {
					t.Errorf("Expected service 'updated-service', got '%s'", r.Service)
				}
				if r.Description != "Updated description" {
					t.Errorf("Expected description 'Updated description', got '%s'", r.Description)
				}
			},
		},
		{
			name:        "fails when resource not found",
			code:        "nonexistent",
			service:     "service",
			description: "description",
			actions:     []ActionSpec{},
			setupMock:   func(m *mockResourceRepository) {},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:        "updates without changing actions when empty",
			code:        "users",
			service:     "updated-service",
			description: "Updated description",
			actions:     []ActionSpec{},
			setupMock: func(m *mockResourceRepository) {
				existing, err := domain.NewResource("users", "original", "original")
				if err != nil {
					t.Fatalf("failed to create existing resource: %v", err)
				}
				existing.AttachActions([]domain.Action{
					func() domain.Action {
						a, err := domain.NewAction(existing.ID, "read", "read action")
						if err != nil {
							t.Fatalf("failed to create action: %v", err)
						}
						return a
					}(),
				})
				m.resources["users"] = existing
			},
			wantErr: false,
			validate: func(t *testing.T, r *domain.Resource) {
				if len(r.Actions) != 1 {
					t.Errorf("Expected existing actions to be preserved, got %d actions", len(r.Actions))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockResourceRepository()
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			ctx := context.Background()
			resource, err := service.Update(ctx, tt.code, tt.service, tt.description, tt.actions)

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

			if tt.validate != nil {
				tt.validate(t, resource)
			}
		})
	}
}

func TestServiceDeprecate(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		setupMock   func(*mockResourceRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully deprecates resource",
			code: "users",
			setupMock: func(m *mockResourceRepository) {
				existing, err := domain.NewResource("users", "service", "description")
				if err != nil {
					t.Fatalf("failed to create existing resource: %v", err)
				}
				m.resources["users"] = existing
			},
			wantErr: false,
		},
		{
			name:        "fails when resource not found",
			code:        "nonexistent",
			setupMock:   func(m *mockResourceRepository) {},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockResourceRepository()
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			ctx := context.Background()
			err := service.Deprecate(ctx, tt.code)

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
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		setupMock   func(*mockResourceRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *domain.Resource)
	}{
		{
			name: "successfully gets resource",
			code: "users",
			setupMock: func(m *mockResourceRepository) {
				existing, err := domain.NewResource("users", "service", "description")
				if err != nil {
					t.Fatalf("failed to create existing resource: %v", err)
				}
				m.resources["users"] = existing
			},
			wantErr: false,
			validate: func(t *testing.T, r *domain.Resource) {
				if r == nil {
					t.Fatal("Expected resource but got nil")
				}
				if r.Code != "users" {
					t.Errorf("Expected code 'users', got '%s'", r.Code)
				}
			},
		},
		{
			name:        "fails when resource not found",
			code:        "nonexistent",
			setupMock:   func(m *mockResourceRepository) {},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockResourceRepository()
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			ctx := context.Background()
			resource, err := service.Get(ctx, tt.code)

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

			if tt.validate != nil {
				tt.validate(t, resource)
			}
		})
	}
}

func TestServiceList(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mockResourceRepository)
		wantErr   bool
		validate  func(*testing.T, []*domain.Resource)
	}{
		{
			name: "successfully lists all resources",
			setupMock: func(m *mockResourceRepository) {
				r1, err := domain.NewResource("users", "service1", "desc1")
				if err != nil {
					t.Fatalf("failed to create resource: %v", err)
				}
				r2, err := domain.NewResource("roles", "service2", "desc2")
				if err != nil {
					t.Fatalf("failed to create resource: %v", err)
				}
				m.resources["users"] = r1
				m.resources["roles"] = r2
			},
			wantErr: false,
			validate: func(t *testing.T, resources []*domain.Resource) {
				if len(resources) != 2 {
					t.Errorf("Expected 2 resources, got %d", len(resources))
				}
			},
		},
		{
			name:      "returns empty list when no resources",
			setupMock: func(m *mockResourceRepository) {},
			wantErr:   false,
			validate: func(t *testing.T, resources []*domain.Resource) {
				if len(resources) != 0 {
					t.Errorf("Expected 0 resources, got %d", len(resources))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockResourceRepository()
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			ctx := context.Background()
			resources, err := service.List(ctx)

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
				tt.validate(t, resources)
			}
		})
	}
}

func TestServiceResolveAction(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		action      string
		setupMock   func(*mockResourceRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *domain.Resource, *domain.Action)
	}{
		{
			name:   "successfully resolves action",
			code:   "users",
			action: "read",
			setupMock: func(m *mockResourceRepository) {
				resource, err := domain.NewResource("users", "service", "description")
				if err != nil {
					t.Fatalf("failed to create resource: %v", err)
				}
				actions := []domain.Action{
					func() domain.Action {
						a, err := domain.NewAction(resource.ID, "read", "read action")
						if err != nil {
							t.Fatalf("failed to create action: %v", err)
						}
						return a
					}(),
				}
				resource.AttachActions(actions)
				m.resources["users"] = resource
			},
			wantErr: false,
			validate: func(t *testing.T, r *domain.Resource, a *domain.Action) {
				if r == nil {
					t.Fatal("Expected resource but got nil")
				}
				if a == nil {
					t.Fatal("Expected action but got nil")
				}
				if a.Name != "read" {
					t.Errorf("Expected action name 'read', got '%s'", a.Name)
				}
			},
		},
		{
			name:        "fails when resource not found",
			code:        "nonexistent",
			action:      "read",
			setupMock:   func(m *mockResourceRepository) {},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "fails when resource is deprecated",
			code:   "users",
			action: "read",
			setupMock: func(m *mockResourceRepository) {
				resource, err := domain.NewResource("users", "service", "description")
				if err != nil {
					t.Fatalf("failed to create resource: %v", err)
				}
				now := time.Now()
				resource.DeprecatedAt = &now
				m.resources["users"] = resource
			},
			wantErr:     true,
			errContains: "deprecated",
		},
		{
			name:   "fails when action not found",
			code:   "users",
			action: "nonexistent",
			setupMock: func(m *mockResourceRepository) {
				resource, err := domain.NewResource("users", "service", "description")
				if err != nil {
					t.Fatalf("failed to create resource: %v", err)
				}
				m.resources["users"] = resource
			},
			wantErr:     true,
			errContains: "action not found",
		},
		{
			name:   "fails when action is deprecated",
			code:   "users",
			action: "read",
			setupMock: func(m *mockResourceRepository) {
				resource, err := domain.NewResource("users", "service", "description")
				if err != nil {
					t.Fatalf("failed to create resource: %v", err)
				}
				action, err := domain.NewAction(resource.ID, "read", "read action")
				if err != nil {
					t.Fatalf("failed to create action: %v", err)
				}
				now := time.Now()
				action.DeprecatedAt = &now
				resource.AttachActions([]domain.Action{action})
				m.resources["users"] = resource
			},
			wantErr:     true,
			errContains: "action is deprecated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockResourceRepository()
			tt.setupMock(mockRepo)
			service := NewService(mockRepo)

			ctx := context.Background()
			resource, action, err := service.ResolveAction(ctx, tt.code, tt.action)

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

			if tt.validate != nil {
				tt.validate(t, resource, action)
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

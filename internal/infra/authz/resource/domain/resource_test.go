package domain

import (
	"testing"
	"time"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		service     string
		description string
		validate    func(*testing.T, *Resource)
	}{
		{
			name:        "creates resource with normalized code",
			code:        "  USER_MANAGEMENT  ",
			service:     "user-service",
			description: "User management resource",
			validate: func(t *testing.T, r *Resource) {
				if r.Code != "user_management" {
					t.Errorf("Expected code 'user_management', got '%s'", r.Code)
				}
			},
		},
		{
			name:        "trims whitespace from service and description",
			code:        "test",
			service:     "  service-name  ",
			description: "  description text  ",
			validate: func(t *testing.T, r *Resource) {
				if r.Service != "service-name" {
					t.Errorf("Expected service 'service-name', got '%s'", r.Service)
				}
				if r.Description != "description text" {
					t.Errorf("Expected description 'description text', got '%s'", r.Description)
				}
			},
		},
		{
			name:        "generates UUID for ID",
			code:        "test",
			service:     "test-service",
			description: "test description",
			validate: func(t *testing.T, r *Resource) {
				if r.ID == "" {
					t.Error("Expected ID to be set")
				}
				if len(r.ID) != 36 { // UUID v4 format
					t.Errorf("Expected ID to be UUID format (36 chars), got length %d", len(r.ID))
				}
			},
		},
		{
			name:        "initializes empty actions slice",
			code:        "test",
			service:     "test-service",
			description: "test description",
			validate: func(t *testing.T, r *Resource) {
				if r.Actions == nil {
					t.Error("Expected Actions to be initialized, got nil")
				}
				if len(r.Actions) != 0 {
					t.Errorf("Expected empty Actions slice, got %d items", len(r.Actions))
				}
			},
		},
		{
			name:        "sets CreatedAt and UpdatedAt",
			code:        "test",
			service:     "test-service",
			description: "test description",
			validate: func(t *testing.T, r *Resource) {
				if r.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set")
				}
				if r.UpdatedAt.IsZero() {
					t.Error("Expected UpdatedAt to be set")
				}
				if !r.CreatedAt.Equal(r.UpdatedAt) {
					t.Error("Expected CreatedAt and UpdatedAt to be equal initially")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := NewResource(tt.code, tt.service, tt.description)
			if err != nil {
				t.Fatalf("Expected resource but got error: %v", err)
			}
			if resource == nil {
				t.Fatal("Expected resource but got nil")
			}
			tt.validate(t, resource)
		})
	}
}

func TestResourceUpdate(t *testing.T) {
	resource, err := NewResource("test", "original-service", "original description")
	if err != nil {
		t.Fatalf("Expected resource but got error: %v", err)
	}
	originalUpdatedAt := resource.UpdatedAt

	// Wait a bit to ensure UpdatedAt changes
	time.Sleep(10 * time.Millisecond)

	if err := resource.Update("new-service", "new description"); err != nil {
		t.Fatalf("Expected update to succeed but got error: %v", err)
	}

	if resource.Service != "new-service" {
		t.Errorf("Expected service 'new-service', got '%s'", resource.Service)
	}
	if resource.Description != "new description" {
		t.Errorf("Expected description 'new description', got '%s'", resource.Description)
	}
	if !resource.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
	if resource.CreatedAt != originalUpdatedAt {
		t.Error("Expected CreatedAt to remain unchanged")
	}
}

func TestResourceIsDeprecated(t *testing.T) {
	tests := []struct {
		name           string
		deprecatedAt   *time.Time
		expectedResult bool
	}{
		{
			name:           "not deprecated when DeprecatedAt is nil",
			deprecatedAt:   nil,
			expectedResult: false,
		},
		{
			name:           "deprecated when DeprecatedAt is set",
			deprecatedAt:   timePtr(time.Now()),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := NewResource("test", "service", "description")
			if err != nil {
				t.Fatalf("Expected resource but got error: %v", err)
			}
			resource.DeprecatedAt = tt.deprecatedAt

			result := resource.IsDeprecated()
			if result != tt.expectedResult {
				t.Errorf("Expected IsDeprecated() to return %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestResourceAttachActions(t *testing.T) {
	resource, err := NewResource("test", "service", "description")
	if err != nil {
		t.Fatalf("Expected resource but got error: %v", err)
	}
	actions := []Action{
		func() Action {
			a, err := NewAction(resource.ID, "read", "read action")
			if err != nil {
				t.Fatalf("Expected action but got error: %v", err)
			}
			return a
		}(),
		func() Action {
			a, err := NewAction(resource.ID, "write", "write action")
			if err != nil {
				t.Fatalf("Expected action but got error: %v", err)
			}
			return a
		}(),
	}

	resource.AttachActions(actions)

	if len(resource.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(resource.Actions))
	}
	if resource.Actions[0].Name != "read" {
		t.Errorf("Expected first action name 'read', got '%s'", resource.Actions[0].Name)
	}
	if resource.Actions[1].Name != "write" {
		t.Errorf("Expected second action name 'write', got '%s'", resource.Actions[1].Name)
	}
}

func TestResourceFindAction(t *testing.T) {
	resource, err := NewResource("test", "service", "description")
	if err != nil {
		t.Fatalf("Expected resource but got error: %v", err)
	}
	actions := []Action{
		func() Action {
			a, err := NewAction(resource.ID, "read", "read action")
			if err != nil {
				t.Fatalf("Expected action but got error: %v", err)
			}
			return a
		}(),
		func() Action {
			a, err := NewAction(resource.ID, "write", "write action")
			if err != nil {
				t.Fatalf("Expected action but got error: %v", err)
			}
			return a
		}(),
		func() Action {
			a, err := NewAction(resource.ID, "delete", "delete action")
			if err != nil {
				t.Fatalf("Expected action but got error: %v", err)
			}
			return a
		}(),
	}
	resource.AttachActions(actions)

	tests := []struct {
		name          string
		actionName    string
		expectedFound bool
		expectedName  string
	}{
		{
			name:          "finds existing action with exact case",
			actionName:    "read",
			expectedFound: true,
			expectedName:  "read",
		},
		{
			name:          "finds existing action with different case",
			actionName:    "READ",
			expectedFound: true,
			expectedName:  "read",
		},
		{
			name:          "finds existing action with mixed case",
			actionName:    "Write",
			expectedFound: true,
			expectedName:  "write",
		},
		{
			name:          "finds existing action with whitespace",
			actionName:    "  delete  ",
			expectedFound: true,
			expectedName:  "delete",
		},
		{
			name:          "does not find non-existent action",
			actionName:    "nonexistent",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, found := resource.FindAction(tt.actionName)

			if found != tt.expectedFound {
				t.Errorf("Expected found=%v, got %v", tt.expectedFound, found)
			}

			if tt.expectedFound {
				if action == nil {
					t.Fatal("Expected action but got nil")
				}
				if action.Name != tt.expectedName {
					t.Errorf("Expected action name '%s', got '%s'", tt.expectedName, action.Name)
				}
			} else {
				if action != nil {
					t.Errorf("Expected nil action, got %v", action)
				}
			}
		})
	}
}

func TestNewAction(t *testing.T) {
	resourceID := "resource-123"
	action, err := NewAction(resourceID, "  CREATE  ", "  Create action  ")
	if err != nil {
		t.Fatalf("Expected action but got error: %v", err)
	}

	if action.ID == "" {
		t.Error("Expected ID to be set")
	}
	if len(action.ID) != 36 {
		t.Errorf("Expected ID to be UUID format (36 chars), got length %d", len(action.ID))
	}
	if action.ResourceID != resourceID {
		t.Errorf("Expected ResourceID '%s', got '%s'", resourceID, action.ResourceID)
	}
	if action.Name != "create" {
		t.Errorf("Expected normalized name 'create', got '%s'", action.Name)
	}
	if action.Description != "Create action" {
		t.Errorf("Expected trimmed description 'Create action', got '%s'", action.Description)
	}
	if action.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if action.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
	if !action.CreatedAt.Equal(action.UpdatedAt) {
		t.Error("Expected CreatedAt and UpdatedAt to be equal initially")
	}
}

func TestActionIsDeprecated(t *testing.T) {
	tests := []struct {
		name           string
		deprecatedAt   *time.Time
		expectedResult bool
	}{
		{
			name:           "not deprecated when DeprecatedAt is nil",
			deprecatedAt:   nil,
			expectedResult: false,
		},
		{
			name:           "deprecated when DeprecatedAt is set",
			deprecatedAt:   timePtr(time.Now()),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := NewAction("resource-123", "test", "test action")
			if err != nil {
				t.Fatalf("Expected action but got error: %v", err)
			}
			action.DeprecatedAt = tt.deprecatedAt

			result := action.IsDeprecated()
			if result != tt.expectedResult {
				t.Errorf("Expected IsDeprecated() to return %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestNormalizeCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercases and trims",
			input:    "  USER_MANAGEMENT  ",
			expected: "user_management",
		},
		{
			name:     "handles already normalized",
			input:    "user_management",
			expected: "user_management",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles whitespace only",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeCode(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestNormalizeAction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercases and trims",
			input:    "  CREATE_USER  ",
			expected: "create_user",
		},
		{
			name:     "handles already normalized",
			input:    "create_user",
			expected: "create_user",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles whitespace only",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAction(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

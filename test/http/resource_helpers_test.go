package http

import (
	"context"
	"errors"
	"testing"

	resourceservice "marketplace/internal/infra/authz/resource/service"
)

func ensureDefaultResources(t *testing.T, service *resourceservice.Service) {
	t.Helper()

	actions := []resourceservice.ActionSpec{
		{Name: "view", Description: "view"},
		{Name: "create", Description: "create"},
		{Name: "update", Description: "update"},
		{Name: "delete", Description: "delete"},
		{Name: "manage", Description: "manage"},
		{Name: "assign_roles", Description: "assign roles"},
	}

	ensureResource(t, service, "users", "user-management", "user resource", actions)
}

func ensureResource(t *testing.T, service *resourceservice.Service, code, serviceName, description string, actions []resourceservice.ActionSpec) {
	t.Helper()

	ctx := context.Background()
	_, err := service.Register(ctx, code, serviceName, description, actions)
	if err == nil {
		return
	}
	if errors.Is(err, resourceservice.ErrResourceAlreadyExists) {
		if _, updateErr := service.Update(ctx, code, serviceName, description, actions); updateErr != nil {
			t.Fatalf("update resource %s: %v", code, updateErr)
		}
		return
	}
	t.Fatalf("register resource %s: %v", code, err)
}


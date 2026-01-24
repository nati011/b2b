package e2e

import (
	"context"
	"testing"

	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	resourcerepo "marketplace/internal/infra/authz/resource/repository"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	"marketplace/internal/infra/db"
	testutil "marketplace/test"
)

func TestResourceServiceBootstrap(t *testing.T) {
	service, cleanup := newResourceService(t)
	defer cleanup()

	manifest := loadManifest(t)
	bootstrapResources(t, service, manifest)
	assertManifestResources(t, service, manifest)
}

func TestResourceServiceBootstrapIdempotent(t *testing.T) {
	service, cleanup := newResourceService(t)
	defer cleanup()

	manifest := loadManifest(t)
	bootstrapResources(t, service, manifest)
	bootstrapResources(t, service, manifest)
	assertManifestResources(t, service, manifest)
}

func newResourceService(t *testing.T) (*resourceservice.Service, func()) {
	t.Helper()

	dbConn, cleanupDB := testutil.SetupTestDB(t)

	txManager := db.NewTxManager(dbConn)
	repo := resourcerepo.NewRepository(dbConn, txManager)
	service := resourceservice.NewService(repo)

	cleanup := func() {
		cleanupDB()
	}

	return service, cleanup
}

func loadManifest(t *testing.T) []resourceDomain.ManifestResource {
	t.Helper()

	manifestPath := testutil.GetResourceManifestPath(t)
	manifest, err := resourceDomain.LoadManifestFromFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}
	return manifest
}

func bootstrapResources(t *testing.T, service *resourceservice.Service, manifest []resourceDomain.ManifestResource) {
	t.Helper()

	ctx := context.Background()
	if err := service.Bootstrap(ctx, manifest); err != nil {
		t.Fatalf("failed to bootstrap resources: %v", err)
	}
}

func assertManifestResources(t *testing.T, service *resourceservice.Service, manifest []resourceDomain.ManifestResource) {
	t.Helper()

	for _, spec := range manifest {
		assertResourceMatchesSpec(t, service, spec)
	}
}

func assertResourceMatchesSpec(t *testing.T, service *resourceservice.Service, spec resourceDomain.ManifestResource) {
	t.Helper()

	resourceEntity, err := service.Get(context.Background(), spec.Code)
	if err != nil {
		t.Fatalf("failed to get bootstrapped resource %s: %v", spec.Code, err)
	}

	if resourceEntity.Code != spec.Code {
		t.Errorf("expected code %s, got %s", spec.Code, resourceEntity.Code)
	}
	if resourceEntity.Service != spec.Service {
		t.Errorf("expected service %s, got %s", spec.Service, resourceEntity.Service)
	}
	if resourceEntity.Description != spec.Description {
		t.Errorf("expected description %s, got %s", spec.Description, resourceEntity.Description)
	}

	if len(resourceEntity.Actions) != len(spec.Actions) {
		t.Errorf("expected %d actions for resource %s, got %d",
			len(spec.Actions), spec.Code, len(resourceEntity.Actions))
	}

	for _, expectedAction := range spec.Actions {
		action, found := resourceEntity.FindAction(expectedAction.Name)
		if !found {
			t.Errorf("expected action %s not found in resource %s", expectedAction.Name, spec.Code)
			continue
		}

		if action.Description != expectedAction.Description {
			t.Errorf("expected action description %s, got %s for action %s in resource %s",
				expectedAction.Description, action.Description, expectedAction.Name, spec.Code)
		}
	}
}

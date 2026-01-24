package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	permissionhttp "marketplace/internal/infra/authz/permission/api/http"
	permissionrepo "marketplace/internal/infra/authz/permission/repository"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	resourcerepo "marketplace/internal/infra/authz/resource/repository"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	dbinfra "marketplace/internal/infra/db"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/google/uuid"
)

func TestPermissionHTTPCreateAndGet(t *testing.T) {
	server, cleanup := newPermissionHTTPServer(t)
	defer cleanup()

	createPayload := permissionhttp.CreatePermissionRequest{
		ID:          uuid.NewString(),
		Resource:    "users",
		Action:      "view",
		Description: "view accounts",
	}

	created := doCreatePermission(t, server, createPayload, http.StatusCreated)

	getURL := fmt.Sprintf("%s/permissions/%s", server.URL, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var fetched permissionhttp.PermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	if fetched.ID != created.ID {
		t.Fatalf("expected fetched ID %s, got %s", created.ID, fetched.ID)
	}
	if fetched.Resource != createPayload.Resource {
		t.Fatalf("expected resource %s, got %s", createPayload.Resource, fetched.Resource)
	}
}

func TestPermissionHTTPList(t *testing.T) {
	server, cleanup := newPermissionHTTPServer(t)
	defer cleanup()

	actions := []string{"view", "create", "update"}
	for i := 0; i < len(actions); i++ {
		payload := permissionhttp.CreatePermissionRequest{
			ID:          uuid.NewString(),
			Resource:    "users",
			Action:      actions[i],
			Description: "test perm",
		}
		doCreatePermission(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s/permissions?page=1&limit=2", server.URL)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var list permissionhttp.PermissionListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if list.Total != 3 {
		t.Fatalf("expected total 3, got %d", list.Total)
	}
	if len(list.Items) == 0 {
		t.Fatalf("expected list items")
	}
}

func TestPermissionHTTPCreateDuplicate(t *testing.T) {
	server, cleanup := newPermissionHTTPServer(t)
	defer cleanup()

	payload := permissionhttp.CreatePermissionRequest{
		ID:          uuid.NewString(),
		Resource:    "users",
		Action:      "view",
		Description: "duplicate check",
	}

	doCreatePermission(t, server, payload, http.StatusCreated)

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+"/permissions", payload)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

func TestPermissionHTTPCreateIdempotent(t *testing.T) {
	server, cleanup := newPermissionHTTPServer(t)
	defer cleanup()

	payload := permissionhttp.CreatePermissionRequest{
		ID:          uuid.NewString(),
		Resource:    "users",
		Action:      "create",
		Description: "idempotent check",
	}
	key := uuid.NewString()

	resp := doJSONRequestWithHeaders(t, server.Client(), http.MethodPost, server.URL+"/permissions", payload, map[string]string{
		"Idempotency-Key": key,
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var first permissionhttp.PermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&first); err != nil {
		t.Fatalf("decode first response: %v", err)
	}

	resp2 := doJSONRequestWithHeaders(t, server.Client(), http.MethodPost, server.URL+"/permissions", payload, map[string]string{
		"Idempotency-Key": key,
	})
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp2.StatusCode)
	}

	var second permissionhttp.PermissionResponse
	if err := json.NewDecoder(resp2.Body).Decode(&second); err != nil {
		t.Fatalf("decode second response: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf("expected identical permission IDs, got %s vs %s", first.ID, second.ID)
	}
}

func newPermissionHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	db, cleanupDB := testutil.SetupTestDB(t)

	txManager := dbinfra.NewTxManager(db)
	resourceRepo := resourcerepo.NewRepository(db, txManager)
	resourceService := resourceservice.NewService(resourceRepo)

	manifestPath := testutil.GetResourceManifestPath(t)
	manifest, err := resourceDomain.LoadManifestFromFile(manifestPath)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if err := resourceService.Bootstrap(context.Background(), manifest); err != nil {
		t.Fatalf("bootstrap resources: %v", err)
	}
	ensureDefaultResources(t, resourceService)

	repo := permissionrepo.NewPermissionRepository(db)
	service := permissionservice.NewPermissionService(repo, resourceService)
	handler := permissionhttp.NewPermissionHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	permissionhttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, cleanup
}

func doCreatePermission(t *testing.T, server *httptest.Server, payload permissionhttp.CreatePermissionRequest, expectedStatus int) permissionhttp.PermissionResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+"/permissions", payload)
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	var created permissionhttp.PermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected created permission ID to be set")
	}
	return created
}

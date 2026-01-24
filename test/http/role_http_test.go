package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	permissionDomain "marketplace/internal/infra/authz/permission/domain"
	permissionrepo "marketplace/internal/infra/authz/permission/repository"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	resourcerepo "marketplace/internal/infra/authz/resource/repository"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	rolehttp "marketplace/internal/infra/authz/role/api/http"
	rolerepo "marketplace/internal/infra/authz/role/repository"
	roleservice "marketplace/internal/infra/authz/role/service"
	dbinfra "marketplace/internal/infra/db"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/google/uuid"
)

const (
	rolesPath = "/roles"
)

func TestRoleHTTPCreateAndGet(t *testing.T) {
	server, permService, cleanup := newRoleHTTPServer(t)
	defer cleanup()

	perm := seedPermission(t, permService, "users", "view")

	createPayload := rolehttp.CreateRoleRequest{
		ID:            uuid.NewString(),
		Name:          "Role Create 1",
		Description:   "role used for create-get flow",
		PermissionIDs: []string{perm.ID},
	}

	createdRole := doCreateRole(t, server, createPayload, http.StatusCreated)

	getURL := fmt.Sprintf("%s%s/%s", server.URL, rolesPath, createdRole.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(expectedStatusFmt, http.StatusOK, resp.StatusCode)
	}

	var fetched rolehttp.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode role get response: %v", err)
	}

	if fetched.ID != createdRole.ID {
		t.Fatalf("expected fetched ID %s, got %s", createdRole.ID, fetched.ID)
	}
	if fetched.Name != createPayload.Name {
		t.Fatalf("expected fetched name %s, got %s", createPayload.Name, fetched.Name)
	}
	if len(fetched.Permissions) != len(createPayload.PermissionIDs) {
		t.Fatalf("expected %d permissions, got %d", len(createPayload.PermissionIDs), len(fetched.Permissions))
	}
}

func TestRoleHTTPListRoles(t *testing.T) {
	server, permService, cleanup := newRoleHTTPServer(t)
	defer cleanup()

	actions := []string{"view", "create", "update"}
	for i := 0; i < len(actions); i++ {
		perm := seedPermission(t, permService, "users", actions[i])
		payload := rolehttp.CreateRoleRequest{
			ID:            uuid.NewString(),
			Name:          fmt.Sprintf("Role %d", i),
			Description:   "list role",
			PermissionIDs: []string{perm.ID},
		}
		doCreateRole(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s%s?page=1&limit=2", server.URL, rolesPath)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(expectedStatusFmt, http.StatusOK, resp.StatusCode)
	}

	var list rolehttp.RoleListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if list.Total != 3 {
		t.Fatalf("expected total 3, got %d", list.Total)
	}
	if len(list.Items) == 0 {
		t.Fatalf("expected paged roles to contain items")
	}
}

func TestRoleHTTPCreateDuplicate(t *testing.T) {
	server, permService, cleanup := newRoleHTTPServer(t)
	defer cleanup()

	perm := seedPermission(t, permService, "users", "assign_roles")
	payload := rolehttp.CreateRoleRequest{
		ID:            uuid.NewString(),
		Name:          "Role Duplicate",
		Description:   "first role",
		PermissionIDs: []string{perm.ID},
	}

	doCreateRole(t, server, payload, http.StatusCreated)

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+rolesPath, payload)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected conflict status, got %d", resp.StatusCode)
	}
}

func TestRoleHTTPCreateIdempotent(t *testing.T) {
	server, permService, cleanup := newRoleHTTPServer(t)
	defer cleanup()

	perm := seedPermission(t, permService, "users", "update")
	payload := rolehttp.CreateRoleRequest{
		ID:            uuid.NewString(),
		Name:          "Role Idempotent",
		Description:   "verify idempotent behavior",
		PermissionIDs: []string{perm.ID},
	}
	key := uuid.NewString()

	resp := doJSONRequestWithHeaders(t, server.Client(), http.MethodPost, server.URL+rolesPath, payload, map[string]string{
		"Idempotency-Key": key,
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf(expectedStatusFmt, http.StatusCreated, resp.StatusCode)
	}

	var first rolehttp.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&first); err != nil {
		t.Fatalf("decode first response: %v", err)
	}

	resp2 := doJSONRequestWithHeaders(t, server.Client(), http.MethodPost, server.URL+rolesPath, payload, map[string]string{
		"Idempotency-Key": key,
	})
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp2.StatusCode)
	}

	var second rolehttp.RoleResponse
	if err := json.NewDecoder(resp2.Body).Decode(&second); err != nil {
		t.Fatalf("decode second response: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf("expected identical role IDs, got %s vs %s", first.ID, second.ID)
	}
}

func newRoleHTTPServer(t *testing.T) (*httptest.Server, *permissionservice.PermissionService, func()) {
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

	permRepo := permissionrepo.NewPermissionRepository(db)
	permService := permissionservice.NewPermissionService(permRepo, resourceService)

	repo := rolerepo.NewRoleRepository(db, txManager)
	service := roleservice.NewRoleService(repo, permRepo)
	handler := rolehttp.NewRoleHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	rolehttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, permService, cleanup
}

func doCreateRole(t *testing.T, server *httptest.Server, payload rolehttp.CreateRoleRequest, expectedStatus int) rolehttp.RoleResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+rolesPath, payload)
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Fatalf(expectedStatusFmt, expectedStatus, resp.StatusCode)
	}

	var created rolehttp.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected created role ID to be set")
	}
	return created
}

func seedPermission(t *testing.T, service *permissionservice.PermissionService, resourceCode, action string) *permissionDomain.Permission {
	t.Helper()
	perm, err := service.Create(context.Background(), uuid.NewString(), resourceCode, action, fmt.Sprintf("%s:%s", resourceCode, action))
	if err != nil {
		t.Fatalf("create permission: %v", err)
	}
	return perm
}

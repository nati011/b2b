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
	roleDomain "marketplace/internal/infra/authz/role/domain"
	rolerepo "marketplace/internal/infra/authz/role/repository"
	roleservice "marketplace/internal/infra/authz/role/service"
	dbinfra "marketplace/internal/infra/db"
	"marketplace/internal/infra/idempotency"
	userhttp "marketplace/internal/infra/user/api/http"
	userDomain "marketplace/internal/infra/user/domain"
	userrepository "marketplace/internal/infra/user/repository"
	userservice "marketplace/internal/infra/user/service"
	httputil "marketplace/pkg/http"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/google/uuid"
)

const (
	usersPath = "/users"
)

func TestUserHTTPCreateAndGet(t *testing.T) {
	server, _, _, cleanup := newUserHTTPServer(t)
	defer cleanup()

	createPayload := userhttp.CreateUserRequest{
		ExternalID: "ext-create-1",
		Email:      "create.and.get@example.com",
		Name:       "Create And Get",
		UserType:   "selfservice",
	}

	createdUser := doCreateUser(t, server, createPayload, http.StatusCreated)

	getURL := fmt.Sprintf("%s%s/%s", server.URL, usersPath, createdUser.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(expectedStatusFmt, http.StatusOK, resp.StatusCode)
	}

	var fetched userhttp.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	if fetched.ID != createdUser.ID {
		t.Fatalf("expected fetched ID %s, got %s", createdUser.ID, fetched.ID)
	}
	if fetched.Email != createPayload.Email {
		t.Fatalf("expected fetched email %s, got %s", createPayload.Email, fetched.Email)
	}
	if fetched.Status != "inactive" {
		t.Fatalf("expected default status inactive, got %s", fetched.Status)
	}
}

func TestUserHTTPListUsers(t *testing.T) {
	server, _, _, cleanup := newUserHTTPServer(t)
	defer cleanup()

	payloads := []userhttp.CreateUserRequest{
		{
			ExternalID: "ext-list-1",
			Email:      "list.one@example.com",
			Name:       "List One",
			UserType:   "selfservice",
		},
		{
			ExternalID: "ext-list-2",
			Email:      "list.two@example.com",
			Name:       "List Two",
			UserType:   "officer",
		},
	}

	for _, payload := range payloads {
		doCreateUser(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s%s?page=1&limit=10", server.URL, usersPath)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(expectedStatusFmt, http.StatusOK, resp.StatusCode)
	}

	var list userhttp.UserListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if list.Total != len(payloads) {
		t.Fatalf("expected total %d, got %d", len(payloads), list.Total)
	}
	if len(list.Items) != len(payloads) {
		t.Fatalf("expected %d items, got %d", len(payloads), len(list.Items))
	}
}

func TestUserHTTPCreateValidationError(t *testing.T) {
	server, _, _, cleanup := newUserHTTPServer(t)
	defer cleanup()

	createPayload := userhttp.CreateUserRequest{
		ExternalID: "ext-invalid-email",
		Email:      "invalid-email",
		Name:       "Invalid Email",
		UserType:   "selfservice",
	}

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+usersPath, createPayload)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf(expectedStatusFmt, http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUserHTTPCreateIdempotent(t *testing.T) {
	server, _, _, cleanup := newUserHTTPServer(t)
	defer cleanup()

	createPayload := userhttp.CreateUserRequest{
		ExternalID: "ext-idem-1",
		Email:      "idem-user@example.com",
		Name:       "Idem User",
		UserType:   "selfservice",
	}
	key := uuid.NewString()

	resp := doJSONRequestWithHeaders(t, server.Client(), http.MethodPost, server.URL+usersPath, createPayload, map[string]string{
		"Idempotency-Key": key,
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf(expectedStatusFmt, http.StatusCreated, resp.StatusCode)
	}

	var first userhttp.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&first); err != nil {
		t.Fatalf("decode first response: %v", err)
	}

	resp2 := doJSONRequestWithHeaders(t, server.Client(), http.MethodPost, server.URL+"/users", createPayload, map[string]string{
		"Idempotency-Key": key,
	})
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp2.StatusCode)
	}

	var second userhttp.UserResponse
	if err := json.NewDecoder(resp2.Body).Decode(&second); err != nil {
		t.Fatalf("decode second response: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf("expected same user ID for idempotent requests, got %s vs %s", first.ID, second.ID)
	}
}

func newUserHTTPServer(t *testing.T) (*httptest.Server, *permissionservice.PermissionService, *roleservice.RoleService, func()) {
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

	roleRepo := rolerepo.NewRoleRepository(db, txManager)
	roleService := roleservice.NewRoleService(roleRepo, permRepo)

	repo := userrepository.NewRepository(db)
	registrationTokenRepo := userrepository.NewRegistrationTokenRepository(db)
	registrationTokenService := userservice.NewRegistrationTokenService(registrationTokenRepo)
	service := userservice.NewService(repo, roleRepo, nil, registrationTokenService)
	userPermissionChecker := allowAllPermissionChecker{}
	handler := userhttp.NewUserHandler(service, userPermissionChecker)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	authUser := &userDomain.User{ID: uuid.NewString()}

	mux := http.NewServeMux()
	userhttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(withAuthenticatedUser(authUser, mux)))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, permService, roleService, cleanup
}

func doCreateUser(t *testing.T, server *httptest.Server, payload userhttp.CreateUserRequest, expectedStatus int) userhttp.UserResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+usersPath, payload)
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Fatalf(expectedStatusFmt, expectedStatus, resp.StatusCode)
	}

	var created userhttp.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected created user ID to be set")
	}
	return created
}

func TestUserHTTPAssignListRevokeRoles(t *testing.T) {
	server, permService, roleService, cleanup := newUserHTTPServer(t)
	defer cleanup()

	perm := createTestPermission(t, permService, "users", "view")
	roleEntity := createTestRole(t, roleService, "Support", []string{perm.ID})

	userPayload := userhttp.CreateUserRequest{
		ExternalID: "role-assign",
		Email:      "role.assign@example.com",
		Name:       "Role Assign",
		UserType:   "officer",
	}
	createdUser := doCreateUser(t, server, userPayload, http.StatusCreated)

	assignPayload := userhttp.AssignRoleRequest{RoleID: roleEntity.ID}
	assignURL := fmt.Sprintf("%s%s/%s/roles", server.URL, usersPath, createdUser.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodPost, assignURL, assignPayload)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("assign role status %d", resp.StatusCode)
	}

	var updated userhttp.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("decode assign response: %v", err)
	}
	if len(updated.Roles) != 1 || updated.Roles[0].ID != roleEntity.ID {
		t.Fatalf("expected role assignment to reflect in user response")
	}

	listURL := fmt.Sprintf("%s%s/%s/roles", server.URL, usersPath, createdUser.ID)
	listResp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("list roles status %d", listResp.StatusCode)
	}

	var list userhttp.UserRoleResponseEnvelope
	if err := json.NewDecoder(listResp.Body).Decode(&list); err != nil {
		t.Fatalf("decode role list: %v", err)
	}
	if len(list.Roles) != 1 {
		t.Fatalf("expected 1 assigned role, got %d", len(list.Roles))
	}

	deleteURL := fmt.Sprintf("%s%s/%s/roles/%s", server.URL, usersPath, createdUser.ID, roleEntity.ID)
	revokeResp := doJSONRequest(t, server.Client(), http.MethodDelete, deleteURL, nil)
	defer revokeResp.Body.Close()
	if revokeResp.StatusCode != http.StatusOK {
		t.Fatalf("revoke role status %d", revokeResp.StatusCode)
	}
	var revoked userhttp.UserResponse
	if err := json.NewDecoder(revokeResp.Body).Decode(&revoked); err != nil {
		t.Fatalf("decode revoke response: %v", err)
	}
	if len(revoked.Roles) != 0 {
		t.Fatalf("expected no roles after revoke, got %d", len(revoked.Roles))
	}
}

func createTestPermission(t *testing.T, service *permissionservice.PermissionService, resourceCode, action string) *permissionDomain.Permission {
	t.Helper()
	perm, err := service.Create(context.Background(), uuid.NewString(), resourceCode, action, action)
	if err != nil {
		t.Fatalf("create permission: %v", err)
	}
	return perm
}

func createTestRole(t *testing.T, service *roleservice.RoleService, name string, permissionIDs []string) *roleDomain.Role {
	t.Helper()
	roleID := uuid.NewString()
	roleEntity, err := service.Create(context.Background(), roleID, name, "test role", permissionIDs)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	return roleEntity
}

func withAuthenticatedUser(user *userDomain.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := httputil.WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type allowAllPermissionChecker struct{}

func (allowAllPermissionChecker) HasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	return true, nil
}

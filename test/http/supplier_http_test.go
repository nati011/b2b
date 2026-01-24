package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	supplierhttp "marketplace/internal/core/supplier/api/http"
	supplierrepo "marketplace/internal/core/supplier/repository"
	supplierservice "marketplace/internal/core/supplier/service"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/stretchr/testify/require"
)

const supplierPath = "/supplier"

func TestSupplierHTTPCreateAndGet(t *testing.T) {
	server, cleanup := newSupplierHTTPServer(t)
	defer cleanup()

	createPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Acme Corporation",
		SupportEmail: "support@acme.com",
		SupportPhone: "0912345678",
		Status:       "active",
	}

	createdSupplier := doCreateSupplier(t, server, createPayload, http.StatusCreated)

	getURL := fmt.Sprintf("%s%s/%d", server.URL, supplierPath, createdSupplier.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched supplierhttp.SupplierResponse
	err := json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)

	require.Equal(t, createdSupplier.ID, fetched.ID)
	require.Equal(t, createPayload.BusinessName, fetched.BusinessName)
	require.Equal(t, createPayload.SupportEmail, fetched.SupportEmail)
	require.Equal(t, "active", fetched.Status)
	require.True(t, fetched.IsActive)
}

func TestSupplierHTTPListSuppliers(t *testing.T) {
	server, cleanup := newSupplierHTTPServer(t)
	defer cleanup()

	payloads := []supplierhttp.CreateSupplierRequest{
		{
			BusinessName: "Supplier One",
			SupportEmail: "one@example.com",
			SupportPhone: "0911111111",
			Status:       "active",
		},
		{
			BusinessName: "Supplier Two",
			SupportEmail: "two@example.com",
			SupportPhone: "0922222222",
			Status:       "inactive",
		},
	}

	for _, payload := range payloads {
		doCreateSupplier(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s%s?page=1&limit=10", server.URL, supplierPath)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse supplierhttp.SupplierListResponse
	err := json.NewDecoder(resp.Body).Decode(&listResponse)
	require.NoError(t, err)

	require.GreaterOrEqual(t, listResponse.Total, 2)
	require.GreaterOrEqual(t, len(listResponse.Items), 2)
}

func TestSupplierHTTPUpdate(t *testing.T) {
	server, cleanup := newSupplierHTTPServer(t)
	defer cleanup()

	createPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Original Corp",
		SupportEmail: "original@example.com",
		SupportPhone: "0910000000",
		Status:       "active",
	}

	created := doCreateSupplier(t, server, createPayload, http.StatusCreated)

	updatePayload := supplierhttp.UpdateSupplierRequest{
		BusinessName: "Updated Corp",
		SupportEmail: "updated@example.com",
		SupportPhone: "0920000000",
		Status:       "active",
	}

	updateURL := fmt.Sprintf("%s%s/%d", server.URL, supplierPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodPut, updateURL, updatePayload)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updated supplierhttp.SupplierResponse
	err := json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(t, err)

	require.Equal(t, updatePayload.BusinessName, updated.BusinessName)
	require.Equal(t, updatePayload.SupportEmail, updated.SupportEmail)
	require.Equal(t, updatePayload.SupportPhone, updated.SupportPhone)
}

func TestSupplierHTTPDelete(t *testing.T) {
	server, cleanup := newSupplierHTTPServer(t)
	defer cleanup()

	createPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "To Delete",
		SupportEmail: "delete@example.com",
		SupportPhone: "0930000000",
		Status:       "active",
	}

	created := doCreateSupplier(t, server, createPayload, http.StatusCreated)

	deleteURL := fmt.Sprintf("%s%s/%d", server.URL, supplierPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodDelete, deleteURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	getURL := fmt.Sprintf("%s%s/%d", server.URL, supplierPath, created.ID)
	getResp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer getResp.Body.Close()

	require.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestSupplierHTTPCreateEmailConflict(t *testing.T) {
	server, cleanup := newSupplierHTTPServer(t)
	defer cleanup()

	payload := supplierhttp.CreateSupplierRequest{
		BusinessName: "First Supplier",
		SupportEmail: "duplicate@example.com",
		SupportPhone: "0911111111",
		Status:       "active",
	}

	doCreateSupplier(t, server, payload, http.StatusCreated)

	// Try to create another with same email
	payload2 := supplierhttp.CreateSupplierRequest{
		BusinessName: "Second Supplier",
		SupportEmail: "duplicate@example.com",
		SupportPhone: "0922222222",
		Status:       "active",
	}

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+supplierPath, payload2)
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func newSupplierHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	db, cleanupDB := testutil.SetupTestDB(t)

	repo := supplierrepo.NewSupplierRepository(db)
	service := supplierservice.NewSupplierService(repo)
	handler := supplierhttp.NewSupplierHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	supplierhttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, cleanup
}

func doCreateSupplier(t *testing.T, server *httptest.Server, payload supplierhttp.CreateSupplierRequest, expectedStatus int) supplierhttp.SupplierResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+supplierPath, payload)
	defer resp.Body.Close()

	require.Equal(t, expectedStatus, resp.StatusCode, "expected status %d, got %d", expectedStatus, resp.StatusCode)

	var created supplierhttp.CreateSupplierResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.Supplier.ID)

	return created.Supplier
}

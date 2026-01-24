package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customerhttp "marketplace/internal/core/customer/api/http"
	customerrepo "marketplace/internal/core/customer/repository"
	customerservice "marketplace/internal/core/customer/service"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/stretchr/testify/require"
)

const customerPath = "/customer"

func TestCustomerHTTPCreateAndGet(t *testing.T) {
	server, cleanup := newCustomerHTTPServer(t)
	defer cleanup()

	createPayload := customerhttp.CreateCustomerRequest{
		FullName:    "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "0912345678",
		City:        "Addis Ababa",
		Region:      "Addis Ababa",
		Woreda:      "01",
		Status:      "active",
	}

	createdCustomer := doCreateCustomer(t, server, createPayload, http.StatusCreated)

	getURL := fmt.Sprintf("%s%s/%d", server.URL, customerPath, createdCustomer.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched customerhttp.CustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)

	require.Equal(t, createdCustomer.ID, fetched.ID)
	require.Equal(t, createPayload.FullName, fetched.FullName)
	require.Equal(t, createPayload.Email, fetched.Email)
	require.Equal(t, "active", fetched.Status)
	require.True(t, fetched.IsActive)
}

func TestCustomerHTTPListCustomers(t *testing.T) {
	server, cleanup := newCustomerHTTPServer(t)
	defer cleanup()

	payloads := []customerhttp.CreateCustomerRequest{
		{
			FullName:    "Customer One",
			Email:       "one@example.com",
			PhoneNumber: "0911111111",
			Status:      "active",
		},
		{
			FullName:    "Customer Two",
			Email:       "two@example.com",
			PhoneNumber: "0922222222",
			Status:      "inactive",
		},
	}

	for _, payload := range payloads {
		doCreateCustomer(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s%s?page=1&limit=10", server.URL, customerPath)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse customerhttp.CustomerListResponse
	err := json.NewDecoder(resp.Body).Decode(&listResponse)
	require.NoError(t, err)

	require.GreaterOrEqual(t, listResponse.Total, 2)
	require.GreaterOrEqual(t, len(listResponse.Items), 2)
}

func TestCustomerHTTPUpdate(t *testing.T) {
	server, cleanup := newCustomerHTTPServer(t)
	defer cleanup()

	createPayload := customerhttp.CreateCustomerRequest{
		FullName:    "Original Name",
		Email:       "original@example.com",
		PhoneNumber: "0910000000",
		Status:      "active",
	}

	created := doCreateCustomer(t, server, createPayload, http.StatusCreated)

	updatePayload := customerhttp.UpdateCustomerRequest{
		FullName:    "Updated Name",
		Email:       "updated@example.com",
		PhoneNumber: "0920000000",
		Status:      "active",
	}

	updateURL := fmt.Sprintf("%s%s/%d", server.URL, customerPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodPut, updateURL, updatePayload)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updated customerhttp.CustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(t, err)

	require.Equal(t, updatePayload.FullName, updated.FullName)
	require.Equal(t, updatePayload.Email, updated.Email)
	require.Equal(t, updatePayload.PhoneNumber, updated.PhoneNumber)
}

func TestCustomerHTTPDelete(t *testing.T) {
	server, cleanup := newCustomerHTTPServer(t)
	defer cleanup()

	createPayload := customerhttp.CreateCustomerRequest{
		FullName:    "To Delete",
		Email:       "delete@example.com",
		PhoneNumber: "0930000000",
		Status:      "active",
	}

	created := doCreateCustomer(t, server, createPayload, http.StatusCreated)

	deleteURL := fmt.Sprintf("%s%s/%d", server.URL, customerPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodDelete, deleteURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	getURL := fmt.Sprintf("%s%s/%d", server.URL, customerPath, created.ID)
	getResp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer getResp.Body.Close()

	require.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestCustomerHTTPCreateEmailConflict(t *testing.T) {
	server, cleanup := newCustomerHTTPServer(t)
	defer cleanup()

	payload := customerhttp.CreateCustomerRequest{
		FullName:    "First Customer",
		Email:       "duplicate@example.com",
		PhoneNumber: "0911111111",
		Status:      "active",
	}

	doCreateCustomer(t, server, payload, http.StatusCreated)

	// Try to create another with same email
	payload2 := customerhttp.CreateCustomerRequest{
		FullName:    "Second Customer",
		Email:       "duplicate@example.com",
		PhoneNumber: "0922222222",
		Status:      "active",
	}

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+customerPath, payload2)
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func newCustomerHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	db, cleanupDB := testutil.SetupTestDB(t)

	repo := customerrepo.NewCustomerRepository(db)
	service := customerservice.NewCustomerService(repo)
	handler := customerhttp.NewCustomerHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	customerhttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, cleanup
}

func doCreateCustomer(t *testing.T, server *httptest.Server, payload customerhttp.CreateCustomerRequest, expectedStatus int) customerhttp.CustomerResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+customerPath, payload)
	defer resp.Body.Close()

	require.Equal(t, expectedStatus, resp.StatusCode, "expected status %d, got %d", expectedStatus, resp.StatusCode)

	var created customerhttp.CreateCustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.Customer.ID)

	return created.Customer
}

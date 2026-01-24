package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customerhttp "marketplace/internal/core/customer/api/http"
	orderhttp "marketplace/internal/core/order/api/http"
	orderrepo "marketplace/internal/core/order/repository"
	orderservice "marketplace/internal/core/order/service"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/stretchr/testify/require"
)

const orderPath = "/order"

func TestOrderHTTPCreateAndGet(t *testing.T) {
	server, cleanup := newOrderHTTPServer(t)
	defer cleanup()

	// Create customer and supplier first
	customerServer, customerCleanup := newCustomerHTTPServer(t)
	defer customerCleanup()

	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "Order Customer",
		Email:       "order@example.com",
		PhoneNumber: "0912345678",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create order
	total := 199.99
	createPayload := orderhttp.CreateOrderRequest{
		CustomerID: customer.ID,
		Status:     "pending",
		Total:      &total,
		Items: []orderhttp.OrderItemRequest{
			{
				Quantity: 2,
				Price:    &total,
			},
		},
	}

	createdOrder := doCreateOrder(t, server, createPayload, http.StatusCreated)

	// Get order
	getURL := fmt.Sprintf("%s%s?id=%d", server.URL, orderPath, createdOrder.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched orderhttp.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)

	require.Equal(t, createdOrder.ID, fetched.ID)
	require.Equal(t, customer.ID, fetched.CustomerID)
	require.Equal(t, "pending", fetched.Status)
	require.NotNil(t, fetched.Total)
	require.Equal(t, total, *fetched.Total)
}

func TestOrderHTTPListCustomerOrders(t *testing.T) {
	server, cleanup := newOrderHTTPServer(t)
	defer cleanup()

	// Create customer
	customerServer, customerCleanup := newCustomerHTTPServer(t)
	defer customerCleanup()

	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "List Customer",
		Email:       "list@example.com",
		PhoneNumber: "0911111111",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create multiple orders
	total1 := 100.0
	total2 := 200.0
	payloads := []orderhttp.CreateOrderRequest{
		{
			CustomerID: customer.ID,
			Status:     "pending",
			Total:      &total1,
			Items: []orderhttp.OrderItemRequest{
				{Quantity: 1, Price: &total1},
			},
		},
		{
			CustomerID: customer.ID,
			Status:     "confirmed",
			Total:      &total2,
			Items: []orderhttp.OrderItemRequest{
				{Quantity: 2, Price: &total2},
			},
		},
	}

	for _, payload := range payloads {
		doCreateOrder(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s/orders/customer?customer_id=%d&limit=10&offset=0", server.URL, customer.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse orderhttp.OrderListResponse
	err := json.NewDecoder(resp.Body).Decode(&listResponse)
	require.NoError(t, err)

	require.GreaterOrEqual(t, listResponse.Total, 2)
	require.GreaterOrEqual(t, len(listResponse.Orders), 2)
}

func TestOrderHTTPUpdateStatus(t *testing.T) {
	server, cleanup := newOrderHTTPServer(t)
	defer cleanup()

	// Create customer
	customerServer, customerCleanup := newCustomerHTTPServer(t)
	defer customerCleanup()

	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "Update Customer",
		Email:       "update@example.com",
		PhoneNumber: "0920000000",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create order
	total := 150.0
	createPayload := orderhttp.CreateOrderRequest{
		CustomerID: customer.ID,
		Status:     "pending",
		Total:      &total,
		Items: []orderhttp.OrderItemRequest{
			{Quantity: 1, Price: &total},
		},
	}

	created := doCreateOrder(t, server, createPayload, http.StatusCreated)
	require.Equal(t, "pending", created.Status)

	// Update status
	updateURL := fmt.Sprintf("%s%s?id=%d&command=confirmed", server.URL, orderPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodPatch, updateURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updated orderhttp.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(t, err)

	require.Equal(t, "confirmed", updated.Status)
}

func newOrderHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	db, cleanupDB := testutil.SetupTestDB(t)

	repo := orderrepo.NewRepository(db)
	service := orderservice.NewService(repo)
	handler := orderhttp.NewOrderHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	orderhttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, cleanup
}

func doCreateOrder(t *testing.T, server *httptest.Server, payload orderhttp.CreateOrderRequest, expectedStatus int) orderhttp.OrderResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+orderPath, payload)
	defer resp.Body.Close()

	require.Equal(t, expectedStatus, resp.StatusCode, "expected status %d, got %d", expectedStatus, resp.StatusCode)

	var created orderhttp.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	return created
}


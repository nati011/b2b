package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customerhttp "marketplace/internal/core/customer/api/http"
	customerrepo "marketplace/internal/core/customer/repository"
	customerservice "marketplace/internal/core/customer/service"
	orderhttp "marketplace/internal/core/order/api/http"
	orderrepo "marketplace/internal/core/order/repository"
	orderservice "marketplace/internal/core/order/service"
	producthttp "marketplace/internal/core/product/api/http"
	supplierhttp "marketplace/internal/core/supplier/api/http"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/stretchr/testify/require"
)

const orderPath = "/order"

func TestOrderHTTPCreateAndGet(t *testing.T) {
	db, cleanupDB := testutil.SetupTestDB(t)
	defer cleanupDB()

	// Create both servers with the same database
	orderServer := newOrderHTTPServerWithDB(t, db)
	customerServer := newCustomerHTTPServerWithDB(t, db)

	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "Order Customer",
		Email:       "order@example.com",
		PhoneNumber: "0912345678",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create supplier and product
	supplierServer := newSupplierHTTPServerWithDB(t, db)
	productServer := newProductHTTPServerWithDB(t, db)

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Order Supplier",
		SupportEmail: "ordersupplier@example.com",
		SupportPhone: "0912345678",
		Status:       "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	price := 99.99
	productPayload := producthttp.CreateProductRequest{
		Name:          "Order Product",
		SupplierID:    supplier.ID,
		Price:         &price,
		TotalQuantity: 100,
		IsActive:      true,
	}
	product := doCreateProduct(t, productServer, productPayload, http.StatusCreated)

	// Create order
	total := 199.99
	createPayload := orderhttp.CreateOrderRequest{
		CustomerID: customer.ID,
		Status:     "pending",
		Total:      &total,
		Items: []orderhttp.OrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  2,
				Price:     &total,
			},
		},
	}

	createdOrder := doCreateOrder(t, orderServer, createPayload, http.StatusCreated)

	// Get order
	getURL := fmt.Sprintf("%s%s?id=%d", orderServer.URL, orderPath, createdOrder.ID)
	resp := doJSONRequest(t, orderServer.Client(), http.MethodGet, getURL, nil)
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
	db, cleanupDB := testutil.SetupTestDB(t)
	defer cleanupDB()

	// Create both servers with the same database
	orderServer := newOrderHTTPServerWithDB(t, db)
	customerServer := newCustomerHTTPServerWithDB(t, db)

	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "List Customer",
		Email:       "list@example.com",
		PhoneNumber: "0911111111",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create supplier and products
	supplierServer := newSupplierHTTPServerWithDB(t, db)
	productServer := newProductHTTPServerWithDB(t, db)

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "List Supplier",
		SupportEmail: "listsupplier@example.com",
		SupportPhone: "0911111111",
		Status:       "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	price1 := 50.0
	price2 := 100.0
	product1 := doCreateProduct(t, productServer, producthttp.CreateProductRequest{
		Name:          "Product 1",
		SupplierID:    supplier.ID,
		Price:         &price1,
		TotalQuantity: 50,
		IsActive:      true,
	}, http.StatusCreated)
	product2 := doCreateProduct(t, productServer, producthttp.CreateProductRequest{
		Name:          "Product 2",
		SupplierID:    supplier.ID,
		Price:         &price2,
		TotalQuantity: 75,
		IsActive:      true,
	}, http.StatusCreated)

	// Create multiple orders
	total1 := 100.0
	total2 := 200.0
	payloads := []orderhttp.CreateOrderRequest{
		{
			CustomerID: customer.ID,
			Status:     "pending",
			Total:      &total1,
			Items: []orderhttp.OrderItemRequest{
				{ProductID: product1.ID, Quantity: 1, Price: &total1},
			},
		},
		{
			CustomerID: customer.ID,
			Status:     "confirmed",
			Total:      &total2,
			Items: []orderhttp.OrderItemRequest{
				{ProductID: product2.ID, Quantity: 2, Price: &total2},
			},
		},
	}

	for _, payload := range payloads {
		doCreateOrder(t, orderServer, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s/orders/customer?customer_id=%d&limit=10&offset=0", orderServer.URL, customer.ID)
	resp := doJSONRequest(t, orderServer.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse orderhttp.OrderListResponse
	err := json.NewDecoder(resp.Body).Decode(&listResponse)
	require.NoError(t, err)

	require.GreaterOrEqual(t, listResponse.Total, 2)
	require.GreaterOrEqual(t, len(listResponse.Orders), 2)
}

func TestOrderHTTPUpdateStatus(t *testing.T) {
	db, cleanupDB := testutil.SetupTestDB(t)
	defer cleanupDB()

	// Create both servers with the same database
	orderServer := newOrderHTTPServerWithDB(t, db)
	customerServer := newCustomerHTTPServerWithDB(t, db)

	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "Update Customer",
		Email:       "update@example.com",
		PhoneNumber: "0920000000",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create supplier and product
	supplierServer := newSupplierHTTPServerWithDB(t, db)
	productServer := newProductHTTPServerWithDB(t, db)

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Update Supplier",
		SupportEmail: "updatesupplier@example.com",
		SupportPhone: "0920000000",
		Status:       "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	price := 150.0
	product := doCreateProduct(t, productServer, producthttp.CreateProductRequest{
		Name:          "Update Product",
		SupplierID:    supplier.ID,
		Price:         &price,
		TotalQuantity: 100,
		IsActive:      true,
	}, http.StatusCreated)

	// Create order
	total := 150.0
	createPayload := orderhttp.CreateOrderRequest{
		CustomerID: customer.ID,
		Status:     "pending",
		Total:      &total,
		Items: []orderhttp.OrderItemRequest{
			{ProductID: product.ID, Quantity: 1, Price: &total},
		},
	}

	created := doCreateOrder(t, orderServer, createPayload, http.StatusCreated)
	require.Equal(t, "pending", created.Status)

	// Update status
	updateURL := fmt.Sprintf("%s%s?id=%d&command=confirmed", orderServer.URL, orderPath, created.ID)
	resp := doJSONRequest(t, orderServer.Client(), http.MethodPatch, updateURL, nil)
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
	return newOrderHTTPServerWithDB(t, db), cleanupDB
}

func newOrderHTTPServerWithDB(t *testing.T, db *sql.DB) *httptest.Server {
	t.Helper()

	repo := orderrepo.NewRepository(db)
	service := orderservice.NewService(repo)
	handler := orderhttp.NewOrderHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	orderhttp.RegisterHTTPRoutes(mux, handler)

	return httptest.NewServer(idempoMiddleware.Handle(mux))
}

func newCustomerHTTPServerWithDB(t *testing.T, db *sql.DB) *httptest.Server {
	t.Helper()

	repo := customerrepo.NewCustomerRepository(db)
	service := customerservice.NewCustomerService(repo)
	handler := customerhttp.NewCustomerHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	customerhttp.RegisterHTTPRoutes(mux, handler)

	return httptest.NewServer(idempoMiddleware.Handle(mux))
}

func doCreateOrder(t *testing.T, server *httptest.Server, payload orderhttp.CreateOrderRequest, expectedStatus int) orderhttp.OrderResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+orderPath, payload)
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		bodyBytes := make([]byte, 1024)
		n, _ := resp.Body.Read(bodyBytes)
		t.Logf("Response body: %s", string(bodyBytes[:n]))
	}
	require.Equal(t, expectedStatus, resp.StatusCode, "expected status %d, got %d", expectedStatus, resp.StatusCode)

	var created orderhttp.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	return created
}

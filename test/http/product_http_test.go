package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	producthttp "marketplace/internal/core/product/api/http"
	productrepo "marketplace/internal/core/product/repository"
	productservice "marketplace/internal/core/product/service"
	supplierhttp "marketplace/internal/core/supplier/api/http"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/stretchr/testify/require"
)

const productPath = "/product"

func TestProductHTTPCreateAndGet(t *testing.T) {
	server, cleanup := newProductHTTPServer(t)
	defer cleanup()

	// First create a supplier for the product
	supplierServer, supplierCleanup := newSupplierHTTPServer(t)
	defer supplierCleanup()

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Test Supplier",
		SupportEmail: "supplier@example.com",
		SupportPhone: "0912345678",
		Status:      "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	// Create product
	price := 99.99
	createPayload := producthttp.CreateProductRequest{
		Name:          "Test Product",
		Description:   "A test product",
		SupplierID:     supplier.ID,
		Price:         &price,
		TotalQuantity: 100,
		IsActive:      true,
	}

	createdProduct := doCreateProduct(t, server, createPayload, http.StatusCreated)

	// Get product
	getURL := fmt.Sprintf("%s%s?id=%d", server.URL, productPath, createdProduct.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched producthttp.ProductResponse
	err := json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)

	require.Equal(t, createdProduct.ID, fetched.ID)
	require.Equal(t, createPayload.Name, fetched.Name)
	require.Equal(t, createPayload.Description, fetched.Description)
	require.Equal(t, supplier.ID, fetched.SupplierID)
	require.True(t, fetched.IsActive)
}

func TestProductHTTPListProducts(t *testing.T) {
	server, cleanup := newProductHTTPServer(t)
	defer cleanup()

	// Create supplier
	supplierServer, supplierCleanup := newSupplierHTTPServer(t)
	defer supplierCleanup()

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "List Supplier",
		SupportEmail: "list@example.com",
		SupportPhone: "0911111111",
		Status:      "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	price1 := 10.0
	price2 := 20.0
	payloads := []producthttp.CreateProductRequest{
		{
			Name:          "Product One",
			SupplierID:    supplier.ID,
			Price:         &price1,
			TotalQuantity: 50,
			IsActive:      true,
		},
		{
			Name:          "Product Two",
			SupplierID:    supplier.ID,
			Price:         &price2,
			TotalQuantity: 75,
			IsActive:      true,
		},
	}

	for _, payload := range payloads {
		doCreateProduct(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s/products?limit=10&offset=0", server.URL)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse producthttp.ProductListResponse
	err := json.NewDecoder(resp.Body).Decode(&listResponse)
	require.NoError(t, err)

	require.GreaterOrEqual(t, listResponse.Total, 2)
	require.GreaterOrEqual(t, len(listResponse.Products), 2)
}

func TestProductHTTPUpdate(t *testing.T) {
	server, cleanup := newProductHTTPServer(t)
	defer cleanup()

	// Create supplier
	supplierServer, supplierCleanup := newSupplierHTTPServer(t)
	defer supplierCleanup()

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Update Supplier",
		SupportEmail: "update@example.com",
		SupportPhone: "0920000000",
		Status:      "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	price := 50.0
	createPayload := producthttp.CreateProductRequest{
		Name:          "Original Product",
		SupplierID:    supplier.ID,
		Price:         &price,
		TotalQuantity: 100,
		IsActive:      true,
	}

	created := doCreateProduct(t, server, createPayload, http.StatusCreated)

	updatedPrice := 75.0
	updatePayload := producthttp.UpdateProductRequest{
		Name:          "Updated Product",
		Description:   "Updated description",
		Price:         &updatedPrice,
		TotalQuantity: 150,
		IsActive:      true,
	}

	updateURL := fmt.Sprintf("%s%s?id=%d", server.URL, productPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodPut, updateURL, updatePayload)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updated producthttp.ProductResponse
	err := json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(t, err)

	require.Equal(t, updatePayload.Name, updated.Name)
	require.Equal(t, updatePayload.Description, updated.Description)
}

func TestProductHTTPDelete(t *testing.T) {
	server, cleanup := newProductHTTPServer(t)
	defer cleanup()

	// Create supplier
	supplierServer, supplierCleanup := newSupplierHTTPServer(t)
	defer supplierCleanup()

	supplierPayload := supplierhttp.CreateSupplierRequest{
		BusinessName: "Delete Supplier",
		SupportEmail: "delete@example.com",
		SupportPhone: "0930000000",
		Status:      "active",
	}
	supplier := doCreateSupplier(t, supplierServer, supplierPayload, http.StatusCreated)

	price := 30.0
	createPayload := producthttp.CreateProductRequest{
		Name:          "To Delete",
		SupplierID:    supplier.ID,
		Price:         &price,
		TotalQuantity: 50,
		IsActive:      true,
	}

	created := doCreateProduct(t, server, createPayload, http.StatusCreated)

	deleteURL := fmt.Sprintf("%s%s?id=%d", server.URL, productPath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodDelete, deleteURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	getURL := fmt.Sprintf("%s%s?id=%d", server.URL, productPath, created.ID)
	getResp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer getResp.Body.Close()

	require.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func newProductHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	db, cleanupDB := testutil.SetupTestDB(t)

	repo := productrepo.NewRepository(db)
	service := productservice.NewService(repo)
	handler := producthttp.NewProductHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	producthttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, cleanup
}

func doCreateProduct(t *testing.T, server *httptest.Server, payload producthttp.CreateProductRequest, expectedStatus int) producthttp.ProductResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+productPath, payload)
	defer resp.Body.Close()

	require.Equal(t, expectedStatus, resp.StatusCode, "expected status %d, got %d", expectedStatus, resp.StatusCode)

	var created producthttp.ProductResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	return created
}


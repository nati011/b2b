package domain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"b2b.nati011.github.com/config"
	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var (
	benchmarkCatalogueHandler *Catalogue
	benchmarkServer           *httptest.Server
	benchmarkOnce             sync.Once
)

func setupBenchmark() {
	benchmarkOnce.Do(func() {
		// Setup test database
		db := db_test_container.Setup()

		// Setup configuration
		cfg := config.Config{}

		// Setup pagination
		pagination := config.DefaultPaginationBuilder().Build()

		// Initialize containers
		applicationContainer := application_core.NewContainer(
			db,
			&cfg,
			pagination,
		)

		domainContainer := domain_core.NewContainer(
			*applicationContainer,
			cfg.BaseUrl,
			cfg.FrontendUrl,
			db,
		)

		// Initialize catalogue handler
		InitCatalogue()
		handlers := handler.GetHandlers()
		for _, h := range handlers {
			if catalogueHandler, ok := h.(*Catalogue); ok {
				benchmarkCatalogueHandler = catalogueHandler
				if err := catalogueHandler.Init(
					applicationContainer.AuthMiddleware,
					applicationContainer,
					domainContainer,
				); err != nil {
					panic("failed to initialize catalogue handler: " + err.Error())
				}
				break
			}
		}

		if benchmarkCatalogueHandler == nil {
			panic("catalogue handler not found")
		}

		// Setup HTTP mux and test server
		mux := http.NewServeMux()
		benchmarkCatalogueHandler.Routes(mux)
		benchmarkServer = httptest.NewServer(mux)
	})
}

// BenchmarkGetCatalogueHandler benchmarks the GET /api/v1/catalogue endpoint
func BenchmarkGetCatalogueHandler(b *testing.B) {
	setupBenchmark()
	url := benchmarkServer.URL + "/api/v1/catalogue"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			b.Fatalf("failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatalf("failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			b.Fatalf("unexpected status code: %d", resp.StatusCode)
		}

		// Read response body to ensure complete request processing
		resp.Body.Close()
	}
}

// BenchmarkGetCatalogueHandlerParallel benchmarks with parallel requests
func BenchmarkGetCatalogueHandlerParallel(b *testing.B) {
	setupBenchmark()
	url := benchmarkServer.URL + "/api/v1/catalogue"

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				b.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				b.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
				b.Fatalf("unexpected status code: %d", resp.StatusCode)
			}

			resp.Body.Close()
		}
	})
}

// BenchmarkCatalogueServiceGetAll benchmarks the service layer directly
func BenchmarkCatalogueServiceGetAll(b *testing.B) {
	setupBenchmark()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := benchmarkCatalogueHandler.service.GetAll(ctx)
		if err != nil {
			b.Fatalf("failed to get catalogue: %v", err)
		}
	}
}

// BenchmarkCatalogueServiceGetAllParallel benchmarks the service layer with parallel calls
func BenchmarkCatalogueServiceGetAllParallel(b *testing.B) {
	setupBenchmark()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := benchmarkCatalogueHandler.service.GetAll(ctx)
			if err != nil {
				b.Fatalf("failed to get catalogue: %v", err)
			}
		}
	})
}

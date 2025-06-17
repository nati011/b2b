package application

import (
	"net/http"
	"time"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthResponse struct {
	Status       string            `json:"status"`
	Timestamp    string            `json:"timestamp"`
	Service      string            `json:"service"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type HealthHandler struct {
}

func InitHealth() {
	handler.Register(new(HealthHandler))

	handler.RegisterResource("/health")
	handler.RegisterResource("/metrics/promethus")
}

func (h *HealthHandler) Init(authMiddleWare *middleware.Auth, services *application_core.Container, domainService *domain_core.Container) error {
	return nil
}

func (h *HealthHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /metrics/promethus", promhttp.Handler())
	mux.Handle("GET /health", http.HandlerFunc(h.healthHandler))
}

func (h *HealthHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement proper check health !!!
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "b2b-backend",
		Version:   "1.0.0",
		Dependencies: map[string]string{
			"postgres": "connected",
			"keycloak": "connected",
		},
	}
	util.OperationSuccessMessageResponse(w, response)
}

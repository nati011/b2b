package application

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/health"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthReportResponse struct {
	Status       string            `json:"status"`
	Timestamp    string            `json:"timestamp"`
	Service      string            `json:"service"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type HealthHandler struct {
	service health.Provider
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
	response, err := h.service.GetSystemHealth()
	if err != nil {
		switch err {
		case health.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessMessageResponse(w, HealthReportResponse(response))
}

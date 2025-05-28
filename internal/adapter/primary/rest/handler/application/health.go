package application

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthHandler struct {
}

func InitHealth() {
	handler.Register(new(HealthHandler))
}

func (d *HealthHandler) Init(authMiddleWare *middleware.Auth, services *application_core.Container, domainService *domain_core.Container) error {
	return nil
}

func (d *HealthHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /metric", promhttp.Handler())

}

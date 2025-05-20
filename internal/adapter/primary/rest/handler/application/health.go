package application

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type HealthHandler struct {
}

func InitHealth() {
	handler.Register(new(HealthHandler))
}

func (d *HealthHandler) Init(authMiddleWare *util.AuthMiddleware, services *application_core.Container, domainService *domain_core.Container) error {
	return nil
}

func (d *HealthHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", d.CheckHealth)

}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

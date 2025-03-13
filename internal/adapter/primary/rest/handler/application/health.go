package handler

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core"
)

type HealthHandler struct {
}

func init() {
	handler.Register(new(HealthHandler))
}

func (d *HealthHandler) Init(services *core.MasterContainer) error {
	return nil
}

func (d *HealthHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", d.CheckHealth)

}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

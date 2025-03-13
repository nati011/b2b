package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core"

	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
	port "b2b.nati011.github.com/internal/port/distributor"
	// util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
)

type DistributorHandler struct {
	distributorService service.Provider
}

func InitDistributor() {
	handler.Register(new(DistributorHandler))
}

func (d *DistributorHandler) Init(services *core.MasterContainer) error {
	// d.service = services.ResourceService
	return nil
}

func (d *DistributorHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/distributor/register", d.RegisterDistributor)

}

func (h *DistributorHandler) RegisterDistributor(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req port.RegisterDistributorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := h.distributorService.Create(ctx, &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (h *DistributorHandler) AddBusinessInformattion(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req port.CreateBusinessInformation

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := h.distributorService.AddBusinessInformattion(ctx, &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

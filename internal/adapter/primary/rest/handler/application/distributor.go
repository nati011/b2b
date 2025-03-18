package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	port "b2b.nati011.github.com/internal/port/application/distributor"
	// util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
)

type DistributorHandler struct {
	distributorService distributor.Provider
}

func InitDistributor() {
	handler.Register(new(DistributorHandler))
}

func (d *DistributorHandler) Init(services *application_core.Container, domainService *domain_core.Container) error {
	d.distributorService = services.DistributorService
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

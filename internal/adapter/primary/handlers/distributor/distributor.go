package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
	port "b2b.nati011.github.com/internal/port/distributor"
)

type DistributorHandler struct {
	distributorService service.Provider
}

func NewDistributorHandler(distributorService service.Provider) DistributorHandler {
	return DistributorHandler{
		distributorService: distributorService,
	}
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

	registerResponse, err := h.distributorService.CreateBusinessInformation(ctx, &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

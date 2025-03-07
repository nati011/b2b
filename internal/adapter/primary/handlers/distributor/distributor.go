package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	distributorDTO "b2b.nati011.github.com/internal/core/domain/distributor/model/dto"
	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
)

type DistributorHandler struct {
	distributorContainer *service.Container
}

func NewDistributorHandler(distributorContainer *service.Container) DistributorHandler {
	return DistributorHandler{
		distributorContainer: distributorContainer,
	}
}

func (h *DistributorHandler) RegisterDistributor(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req distributorDTO.RegisterDistributorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := h.distributorContainer.DistributorProvider.Create(ctx, &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

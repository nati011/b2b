package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	distributor "b2b.nati011.github.com/internal/core/domain/distributor"
)

type DistributorHandler struct {
	distributorService *distributor.DistributorService
}

func NewDistributorHandler(distributorService *distributor.DistributorService) *DistributorHandler {
	return &DistributorHandler{
		distributorService: distributorService,
	}
}

func (h *DistributorHandler) RegisterDistributor(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req distributor.RegisterDistributorRequest
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

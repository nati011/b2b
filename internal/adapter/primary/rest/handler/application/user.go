package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core/application/user"

	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type UserHandler struct {
	service user.Provider
}

func InitUser() {
	handler.Register(new(UserHandler))
}

func (a *UserHandler) Init(services *application_core.Container, domainService *domain_core.Container) error {
	a.service = services.UserService
	return nil
}

func (a *UserHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("PATCH /api/v1/user", a.UpdateProfile)
	mux.HandleFunc("GET /api/v1/user", a.FetchProfile)

}

func (a *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req user.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := a.service.Update(ctx, &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (a *UserHandler) FetchProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("oopsy, not implemented (yet?)"))
}

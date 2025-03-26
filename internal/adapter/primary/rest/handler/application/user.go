package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core/application/user"

	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type GetUserResponse struct {
	Id         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Username   string    `json:"username"`
	DOB        time.Time `json:"dob"`
	IsActive   bool      `json:"is_active"`
	ExternalId string    `json:"external_id"`
}

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
	var req user.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := a.service.Update(r.Context(), &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (a *UserHandler) FetchProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("oopsy, not implemented (yet?)"))
}

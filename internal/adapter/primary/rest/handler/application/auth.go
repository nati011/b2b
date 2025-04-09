package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/auth"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type AuthHandler struct {
	service    auth.Provider
	middleware util.AuthMiddleware
}

func InitAuth() {
	handler.Register(new(AuthHandler))
}

func (a *AuthHandler) Init(services *application_core.Container, domainService *domain_core.Container) error {
	a.service = services.AuthService
	a.middleware = *services.AuthMiddleware
	return nil
}

func (a *AuthHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", a.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", a.RefreshToken)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()
	var requestBody auth.LoginUserRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	loginResponse, err := h.service.ClientLogin(r.Context(), requestBody)
	if err != nil {
		switch err {
		case auth.ErrUnknown:
			util.ServerErrorResponse(w, err)
		default:
			util.UnauthorizedResponse(w)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"body": loginResponse})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()
	var requestBody auth.RefreshTokenRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	refreshResponse, err := h.service.RefreshToken(r.Context(), requestBody)
	if err != nil {
		switch err {
		case auth.ErrUnknown:
			util.ServerErrorResponse(w, err)
		default:
			util.UnauthorizedResponse(w)
			return
		}
	}

	util.OperationSuccessResponse(w, util.Envelope{"body": refreshResponse})
}

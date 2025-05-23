package application

import (
	"encoding/json"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	auth "b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/middleware"

	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type AuthHandler struct {
	authMiddleware middleware.Auth
	service        auth.Provider
}

type ResetPasswordRequest struct {
	NewPassword string `json:"password"`
}

func InitAuth() {
	handler.Register(new(AuthHandler))
}

func (a *AuthHandler) Init(authMiddleWare *middleware.Auth, services *application_core.Container, domainService *domain_core.Container) error {
	a.service = services.AuthService
	a.authMiddleware = *services.AuthMiddleware
	return nil
}

func (a *AuthHandler) Routes(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/v1/auth/reset/{token}", func(w http.ResponseWriter, r *http.Request) {
		a.authMiddleware.RequireNoAuthentication(http.HandlerFunc(a.ResetCredentialsHandler)).ServeHTTP(w, r)
	})
}

func (h *AuthHandler) ResetCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody ResetPasswordRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	restToken, err := util.GetStringPathParam(r, 5)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	err = h.service.ResetClientCredentials(r.Context(), &auth.ResetCredentialsRequest{
		NewPassword: requestBody.NewPassword,
		ResetToken:  restToken,
	})
	if err != nil {
		util.ServerErrorResponse(w, err)
		return
	}
	util.OperationSuccessMessageResponse(w, "password reset successfully")
}

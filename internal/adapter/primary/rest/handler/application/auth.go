package application

import (
	"encoding/json"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	auth "b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/middleware"
	oauth "b2b.nati011.github.com/internal/core/application/oauth"
	"b2b.nati011.github.com/internal/core/application/resource"

	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type AuthHandler struct {
	authMiddleware middleware.Auth
	service        auth.Provider
	oauthService   oauth.Provider
}

type ResetPasswordRequest struct {
	NewPassword string `json:"password"`
}

func InitAuth() {
	handler.Register(new(AuthHandler))
	initRegisterResources()
}

func initRegisterResources() {
	handler.RegisterResource(resource.CreateRequest{
		Name:     "auth_login",
		Action:   "ALL",
		Resource: "/api/v1/auth/login",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "auth_logout",
		Action:   "ALL",
		Resource: "/api/v1/auth/logout",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "auth_refresh",
		Action:   "ALL",
		Resource: "/api/v1/auth/refresh",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "auth_reset",
		Action:   "ALL",
		Resource: "/api/v1/auth/reset/{param}",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "sso",
		Action:   "ALL",
		Resource: "/api/v1/auth/sso",
	})
}

func (a *AuthHandler) Init(authMiddleWare *middleware.Auth, services *application_core.Container, domainService *domain_core.Container) error {
	a.service = services.AuthService
	a.authMiddleware = *services.AuthMiddleware
	return nil
}

func (a *AuthHandler) Routes(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		a.authMiddleware.RequireNoAuthentication(http.HandlerFunc(a.LoginHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/auth/sso", func(w http.ResponseWriter, r *http.Request) {
		a.authMiddleware.RequireNoAuthentication(http.HandlerFunc(a.SSOHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		a.authMiddleware.RequireNoAuthentication(http.HandlerFunc(a.LogoutHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		a.authMiddleware.RequireNoAuthentication(http.HandlerFunc(a.RefreshTokenHandler)).ServeHTTP(w, r)
	})

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

func (a *AuthHandler) SSOHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	loginResponse, err := a.service.ClientLogin(r.Context(), &requestBody)
	if err != nil {
		switch err {
		case auth.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.UnauthorizedResponse(w)
			return
		}
	}
	util.OperationSuccessResponse(w, loginResponse.JWT)
}

func (a *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refreshResponse, err := a.service.RefreshToken(r.Context(), &req)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	util.OperationSuccessResponse(w, refreshResponse.JWT)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.ClientLogout(r.Context(), req)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	util.OperationSuccessMessageResponse(w, "Logged out successfully")
}

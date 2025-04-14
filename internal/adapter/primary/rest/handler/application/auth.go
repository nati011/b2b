package application

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/auth"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type RegisterUserRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	BirthDate   time.Time `json:"birth_date"`
	PhoneNumber string    `json:"phone_number"`
	ExternalId  string    `json:"external_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
}

type RegisterUserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginAuthResponse struct {
	JWT JWT
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not_before_policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

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
	util.OperationSuccessResponse(w, util.Envelope{"body": LoginAuthResponse{
		JWT: JWT(loginResponse.JWT),
	}})
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

	util.OperationSuccessResponse(w, util.Envelope{"body": LoginAuthResponse{
		JWT: JWT(refreshResponse.JWT),
	}})
}

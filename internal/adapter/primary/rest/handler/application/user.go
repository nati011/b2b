package application

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/user"

	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrUnknownUserCommand = errors.New("unknown command")
)

type CreateUserRequest struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Username   string    `json:"username"`
	DOB        time.Time `json:"dob"`
	ExternalId string    `json:"external_id"`
	Password   string    `json:"password"`
}

type GetUserAssignedRoleResponse struct {
	Id int `json:"id"`
}

type GetAllUserAssignedRoleResponse struct {
	List []GetUserAssignedRoleResponse `json:"list"`
}

type UpdateUserRequest struct {
	Id         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Username   string    `json:"username"`
	DOB        time.Time `json:"dob"`
	ExternalId string    `json:"external_id"`
}

type GetUserByParamRequest struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Username   string `json:"username"`
	IsActive   bool   `json:"is_active"`
	ExternalId string `json:"external_id"`
}

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

type GetAllUserResponse struct {
	List []GetUserResponse `json:"list"`
}

type UserHandler struct {
	authMiddleware middleware.Auth
	service        user.Provider
}

type InitResetPasswordRequest struct {
	Email string `json:"email"`
}

func InitUser() {
	handler.Register(new(UserHandler))

	handler.RegisterResource("/api/v1/user")
	handler.RegisterResource("/api/v1/user/{id}/status")
	handler.RegisterResource("/api/v1/user/{id}/role/{role_id}")
	handler.RegisterResource("/api/v1/user/init_auth_reset")
}

func (a *UserHandler) Init(authMiddleWare *middleware.Auth, services *application_core.Container, domainService *domain_core.Container) error {
	a.service = services.UserService
	a.authMiddleware = *services.AuthMiddleware
	return nil
}

func (u *UserHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireNoAuthentication(http.HandlerFunc(u.GetUser)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireAuthentication(http.HandlerFunc(u.CreateUser)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireAuthentication(http.HandlerFunc(u.UpdateProfile)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/user/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireAuthentication(http.HandlerFunc(u.StatusHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/user/{id}/role/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireNoAuthentication(http.HandlerFunc(u.RoleHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/user/{id}/role", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireNoAuthentication(http.HandlerFunc(u.GetRoleHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/user/init_reset", func(w http.ResponseWriter, r *http.Request) {
		u.authMiddleware.RequireNoAuthentication(http.HandlerFunc(u.InitResetTokenHandler)).ServeHTTP(w, r)
	})
}

const (
	ACTIVATE_USER_COMMAND   = "activate"
	DEACTIVATE_USER_COMMAND = "deactivate"
)

func (p *UserHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	const ParamCommand = "command"

	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	if paramCommandValue != "" {
		typedParamId, err := util.GetPathParam(r, 4)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case ACTIVATE_USER_COMMAND:
			err = p.service.Activate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case user.ErrIdNotFound:

					util.RequestErrorResponse(w, err)
					return
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"user": typedParamId})
		case DEACTIVATE_USER_COMMAND:
			err = p.service.Deactivate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case user.ErrIdNotFound:
					util.RequestErrorResponse(w, err)
					return
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"user": typedParamId})
		default:
			util.RequestErrorResponse(w, ErrUnknownUserCommand)
		}
	}
}

const (
	ASSIGN_ROLE_COMMAND          = "assign_role"
	REMOVE_ASSIGNED_ROLE_COMMAND = "remove_role"
)

func (a *UserHandler) RoleHandler(w http.ResponseWriter, r *http.Request) {
	const ParamCommand = "command"

	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)

	if paramCommandValue != "" {
		typedParamRoleId, err := util.GetPathParam(r, 6)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		typedParamId, err := util.GetPathParam(r, 4)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}

		switch paramCommandValue {
		case ASSIGN_ROLE_COMMAND:
			err = a.service.AssignRole(r.Context(), typedParamId, typedParamRoleId)
			if err != nil {
				switch err {
				case user.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"role": typedParamRoleId})
		case REMOVE_ASSIGNED_ROLE_COMMAND:
			err = a.service.RemoveAssignedRole(r.Context(), typedParamId, typedParamRoleId)
			if err != nil {
				switch err {
				case user.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"role": typedParamRoleId})
		default:
			util.RequestErrorResponse(w, ErrUnknownUserCommand)
		}
	}
}

func (a *UserHandler) GetRoleHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	resp, err := a.service.GetAllAssignedRoles(r.Context(), typedParamId)
	if err != nil {
		switch err {
		case user.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"roles": resp})

}

func (a *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateUserRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	id, err := a.service.Create(r.Context(), (*user.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case user.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"user": id})
}

func (a *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	const ParamEmail = "email"
	const ParamPhone = "phone"
	const ParamUsername = "username"
	const ParamIsActive = "is_active"
	const ParamExternalId = "external_id"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	paramNameValue := paramValues.Get(ParamName)
	paramEmailValue := paramValues.Get(ParamEmail)
	paramPhoneValue := paramValues.Get(ParamPhone)
	paramUsernameValue := paramValues.Get(ParamUsername)
	paramIsActiveValue := paramValues.Get(ParamIsActive)
	ParamExternalIdValue := paramValues.Get(ParamExternalId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := a.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case user.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"user": GetUserResponse(resp)})
	} else if paramNameValue != "" || paramEmailValue != "" || paramPhoneValue != "" ||
		paramUsernameValue != "" || paramIsActiveValue != "" || ParamExternalIdValue != "" {

		var typedparamIsActiveValue bool
		switch paramNameValue {
		case "true":
			typedparamIsActiveValue = true
		default:
			typedparamIsActiveValue = false
		}

		resp, err := a.service.GetByParam(r.Context(), &user.GetByParam{
			Email:      strings.Trim(paramEmailValue, `"`),
			Phone:      strings.Trim(paramPhoneValue, `"`),
			Username:   strings.Trim(paramUsernameValue, `"`),
			ExternalId: strings.Trim(ParamExternalIdValue, `"`),
			IsActive:   typedparamIsActiveValue,
		})

		if err != nil {
			switch err {
			case user.ErrEmptyGetContent:
			case user.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"users": resp})
	} else {
		resp, err := a.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case user.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		var response GetAllUserResponse
		for _, i := range resp.List {
			response.List = append(response.List, GetUserResponse(i))
		}
		util.OperationSuccessResponse(w, util.Envelope{"users": response})
	}
}

func (u *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()
	var requestBody UpdateUserRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := u.service.Update(r.Context(), (*user.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case user.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"user": id})
}

func (u *UserHandler) InitResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()
	var requestBody InitResetPasswordRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	err = u.service.ResetPassword(r.Context(), requestBody.Email)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	util.OperationSuccessMessageResponse(w, "password reset init successfully")
}

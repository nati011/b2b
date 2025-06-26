package application

import (
	"errors"
	"net/http"
	"time"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type GetIdentityResponse struct {
	Id          int                          `json:"id"`
	FirstName   string                       `json:"first_name"`
	LastName    string                       `json:"last_name"`
	Email       string                       `json:"email"`
	Phone       string                       `json:"phone"`
	Username    string                       `json:"username"`
	DOB         time.Time                    `json:"dob"`
	IsActive    bool                         `json:"is_active"`
	ExternalId  string                       `json:"external_id"`
	Permissions role.GetAllResourcesResponse `json:"permissions"`
}

type IdentityHandler struct {
	authMiddleware middleware.Auth
	userService    user.Provider
	roleService    role.Provider
}

func InitIdentity() {
	handler.Register(new(IdentityHandler))

	handler.RegisterResource(resource.CreateRequest{
		Name:     "identity_user",
		Action:   "ALL",
		Resource: "/api/v1/identity/user",
	})
}

func (i *IdentityHandler) Init(authMiddleWare *middleware.Auth, services *application_core.Container, domainService *domain_core.Container) error {
	i.userService = services.UserService
	i.roleService = services.RoleService
	i.authMiddleware = *services.AuthMiddleware
	return nil
}

func (i *IdentityHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/identity/user", func(w http.ResponseWriter, r *http.Request) {
		i.authMiddleware.RequireAuthentication(http.HandlerFunc(i.GetIdentity)).ServeHTTP(w, r)
	})
}

func (a *IdentityHandler) GetIdentity(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		util.ServerErrorResponse(w, errors.New("userId not found in context or is not an integer"))
		return
	}

	user_resp, err := a.userService.Get(r.Context(), userId)
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

	role_resp, err := a.userService.GetAllAssignedRoles(r.Context(), userId)
	if err != nil {
		switch err {
		case user.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		case user.ErrNoRoleAssigned:
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}

	allowed_resources := role.GetAllResourcesResponse{}
	for _, i := range role_resp.List {
		resource_resp, err := a.roleService.GetAllResources(r.Context(), i.Id)
		if err != nil {
			switch err {
			case resource.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		allowed_resources.List = append(allowed_resources.List, resource_resp.List...)
	}
	response := GetIdentityResponse{
		Id:          user_resp.Id,
		FirstName:   user_resp.FirstName,
		LastName:    user_resp.LastName,
		Email:       user_resp.Email,
		Phone:       user_resp.Phone,
		Username:    user_resp.Username,
		DOB:         user_resp.DOB,
		IsActive:    user_resp.IsActive,
		ExternalId:  user_resp.ExternalId,
		Permissions: allowed_resources,
	}
	util.OperationSuccessResponse(w, util.Envelope{"user": response})
}

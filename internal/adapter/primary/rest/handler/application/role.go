package application

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/role"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type CreateRoleRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type UpdateRoleRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type GetRoleRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GetRoleResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type GetAllRoleResponse struct {
	List []GetRoleResponse `json:"roles"`
}

type AddResourceToRoleRequest struct {
	ResourceId int `json:"resource_id"`
	RoleId     int `json:"role_id"`
}

type RemoveResourceFromRoleRequest struct {
	ResourceId int `json:"resource_id"`
	RoleId     int `json:"role_id"`
}

type HasResourceInRoleRequest struct {
	ResourceId int `json:"resource_id"`
	RoleId     int `json:"role_id"`
}

type GetAllResourcesResponse struct {
	List []GetResourceResponse `json:"resources"`
}

var (
	ErrUnknownCommand = errors.New("unknown command")
)

type Role struct {
	authMiddleware middleware.Auth
	service        role.Provider
}

func InitRole() {
	handler.Register(new(Role))

	handler.RegisterResource("/api/v1/role")
	handler.RegisterResource("/api/v1/role/resource")
	handler.RegisterResource("/api/v1/role/{id}")
}

func (r *Role) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.ApplicationServices.RoleService
	r.authMiddleware = *authMiddleWare
	return nil
}

func (ro *Role) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/role", func(w http.ResponseWriter, r *http.Request) {
		ro.authMiddleware.RequireAuthentication(http.HandlerFunc(ro.GetHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/role/resource", func(w http.ResponseWriter, r *http.Request) {
		ro.authMiddleware.RequireAuthentication(http.HandlerFunc(ro.GetAllResourcesHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/role", func(w http.ResponseWriter, r *http.Request) {
		ro.authMiddleware.RequireAuthentication(http.HandlerFunc(ro.CreateHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/role/{id}", func(w http.ResponseWriter, r *http.Request) {
		ro.authMiddleware.RequireAuthentication(http.HandlerFunc(ro.CommandHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PUT /api/v1/role", func(w http.ResponseWriter, r *http.Request) {
		ro.authMiddleware.RequireAuthentication(http.HandlerFunc(ro.UpdateHandler)).ServeHTTP(w, r)
	})
}

const (
	ADD_RESOURCE_COMMAND    = "add_resource"
	REMOVE_RESOURCE_COMMAND = "remove_resource"
)

func (ro *Role) CommandHandler(w http.ResponseWriter, r *http.Request) {
	typedParamRoleId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	const ParamResourceId = "resource_id"
	const ParamCommand = "command"
	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	paramResourceIdValue := paramValues.Get(ParamResourceId)
	if paramCommandValue != "" {
		typedParamResourceId, err := strconv.Atoi(paramResourceIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case ADD_RESOURCE_COMMAND:
			err = ro.service.AddResource(r.Context(), &role.AddResourceRequest{
				ResourceId: typedParamResourceId,
				RoleId:     typedParamRoleId,
			})
			if err != nil {
				switch err {
				case role.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"resource": typedParamResourceId})
		case REMOVE_RESOURCE_COMMAND:
			err = ro.service.RemoveResource(r.Context(), &role.RemoveResourceRequest{
				ResourceId: typedParamResourceId,
				RoleId:     typedParamRoleId,
			})
			if err != nil {
				switch err {
				case role.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"resource": typedParamResourceId})
		default:
			util.RequestErrorResponse(w, ErrUnknownCommand)
			return
		}
	}
}

func (ro *Role) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateRoleRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := ro.service.Update(r.Context(), (*role.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case role.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"role": id})
}

func (ro *Role) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateRoleRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	id, err := ro.service.Create(r.Context(), (*role.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		default:
			util.RequestErrorResponse(w, err)
			return
		case role.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"role": id})
}

func (ro *Role) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := ro.service.Get(r.Context(), &role.GetRequest{
			Id: typedParamId,
		})
		if err != nil {
			switch err {
			case role.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.WriteJSON(w, util.Envelope{"role": resp}, http.StatusAccepted)
	} else {
		resp, err := ro.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case role.ErrEmptyGetContent:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		var response GetAllRoleResponse
		for _, i := range resp.List {
			response.List = append(response.List, (GetRoleResponse)(i))
		}
		util.OperationSuccessResponse(w, response)
	}
}

func (ro *Role) GetAllResourcesHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := ro.service.GetAllResources(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case role.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		var response GetAllResourcesResponse
		for _, i := range resp.List {
			response.List = append(response.List, GetResourceResponse{
				Id:     i.Id,
				Name:   i.Name,
				Action: i.Action,
			})
		}
		util.WriteJSON(w, util.Envelope{"resources": response}, http.StatusAccepted)
	} else {
		resp, err := ro.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case role.ErrEmptyGetContent:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, resp)
	}
}

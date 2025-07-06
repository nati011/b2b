package application

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type CreateResourceRequest struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Name     string `json:"name"`
	Scope    string `json:"scope"`
}

type UpdateResourceRequest struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

type GetResourceResponse struct {
	Id       int    `json:"id"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Name     string `json:"name"`
	Scope    string `json:"scope"`
}

type GetAllResourceResponse struct {
	List []GetResourceResponse `json:"resources"`
}

type Resource struct {
	authMiddleware middleware.Auth
	service        resource.Provider
}

func InitResource() {
	handler.Register(new(Resource))

	handler.RegisterResource(resource.CreateRequest{
		Name:     "resource",
		Action:   "ALL",
		Resource: "/api/v1/resource",
	})
}

func (r *Resource) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = applicationServices.ResourceService
	r.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (re *Resource) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/resource", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.GetResourceHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/resource", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.CreateResourceHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/resource", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.UpdateResourceHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("DELETE /api/v1/resource", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.DeleteResourceHandler)).ServeHTTP(w, r)
	})
}

func (rs *Resource) GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := rs.service.Get(r.Context(), typedParamId)
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
		util.OperationSuccessResponse(w, resp)
	} else {
		var response GetAllResourceResponse
		resp, err := rs.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case resource.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, response)
				return
			case resource.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		for _, i := range resp.List {
			response.List = append(response.List, (GetResourceResponse)(i))
		}
		util.OperationSuccessResponse(w, response)
	}
}

func (rs *Resource) CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateResourceRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	id, err := rs.service.Create(r.Context(), (*resource.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		default:
			util.RequestErrorResponse(w, err)
			return
		case resource.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"resource": id})
}

func (rs *Resource) UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateResourceRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := rs.service.Update(r.Context(), &resource.UpdateRequest{
		Id:       requestBody.Id,
		Resource: requestBody.Resource,
		Name:     requestBody.Name,
		Action:   requestBody.Action,
	})
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
	util.OperationSuccessResponse(w, util.Envelope{"resource": id})
}

func (rs *Resource) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		err = rs.service.Delete(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case resource.ErrIdNotFound:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"resource": paramIdValue})
}

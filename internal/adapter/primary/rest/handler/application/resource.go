package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type CreateResourceRequest struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

type UpdateResourceRequest struct {
	Id     int    `json:"id"`
	Action string `json:"action"`
	Name   string `json:"name"`
}

type GetResponse struct {
	Id     int    `json:"id"`
	Action string `json:"action"`
	Name   string `json:"name"`
}

type GetAllResponse struct {
	List []GetResponse `json:"resources"`
}

type Resource struct {
	service resource.Provider
}

func InitResource() {
	handler.Register(new(Resource))
}

func (r *Resource) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = applicationServices.ResourceService
	return nil
}

func (r *Resource) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/resource", r.GetResourceHandler)
	mux.HandleFunc("POST /api/v1/resource", r.CreateResourceHandler)
	mux.HandleFunc("PATCH /api/v1/resource", r.UpdateResourceHandler)
	mux.HandleFunc("DELETE /api/v1/resource", r.DeleteResourceHandler)
}

func (rs *Resource) GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, r, err)
			return
		}
		resp, err := rs.service.Get(r.Context(), &resource.GetRequest{
			Id: typedParamId,
		})
		if err != nil {
			switch err {
			case resource.ErrEmptyGetContent,
				resource.ErrEmptyName,
				resource.ErrEmptyAction:

				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}

		util.WriteJSON(w, util.Envelope{"resource": resp}, http.StatusAccepted)
	} else {
		resp, err := rs.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case resource.ErrEmptyGetContent:
				util.RequestErrorResponse(w, r, err)
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}

		util.WriteJSON(w, util.Envelope{"resources": resp}, http.StatusAccepted)
	}
}

func (rs *Resource) CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateResourceRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}

	id, err := rs.service.Create(r.Context(), (*resource.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case resource.ErrDuplicateName,
			resource.ErrEmptyAction,
			resource.ErrEmptyName,
			resource.ErrIdNotFound,
			resource.ErrEmptyUpdateContent,
			resource.ErrEmptyGetContent:

			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"resource": id}, http.StatusAccepted)
}

func (rs *Resource) UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateResourceRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	id, err := rs.service.Update(r.Context(), (*resource.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case resource.ErrDuplicateName,
			resource.ErrEmptyAction,
			resource.ErrEmptyName,
			resource.ErrIdNotFound,
			resource.ErrEmptyUpdateContent,
			resource.ErrEmptyGetContent:

			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"resource": id}, http.StatusAccepted)
}

func (rs *Resource) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, r, err)
			return
		}
		err = rs.service.Delete(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case resource.ErrIdNotFound:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}
	}
	util.WriteJSON(w, util.Envelope{"resource": paramIdValue}, http.StatusAccepted)
}

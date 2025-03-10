package handler

import (
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	"b2b.nati011.github.com/internal/core"
	"b2b.nati011.github.com/internal/core/application/resource"
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

func init() {
	handler.Register(new(Resource))
}

func (r *Resource) Init(services *core.MasterContainer) error {
	// r.service = services.ResourceService
	return nil
}

func (r *Resource) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/resource", r.GetResourceHandler)
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
			case resource.ErrEmptyGetContent:
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}

		util.WriteJSON(w, util.Envelope{"resource": resp}, nil)
	} else {
		resp, err := rs.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case resource.ErrEmptyGetContent:
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}

		util.WriteJSON(w, util.Envelope{"resources": resp}, nil)
	}
}

func (rs *Resource) CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
}

func (rs *Resource) UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
}

func (rs *Resource) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
}

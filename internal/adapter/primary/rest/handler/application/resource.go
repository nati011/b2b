package handler

import (
	"net/http"
	"strconv"

	handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler"
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

type ResourceHandler struct {
	service resource.Provider
}

func NewResourceHandler(rp resource.Provider) ResourceHandler {
	return ResourceHandler{
		service: rp,
	}
}

func CreateResourceHandler(w http.ResponseWriter, r *http.Request) {}
func UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {}
func DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {}

func (rh *ResourceHandler) GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			handler.RequestErrorResponse(w, r, err)
			return
		}
		resp, err := rh.service.Get(r.Context(), &resource.GetRequest{
			Id: typedParamId,
		})
		if err != nil {
			switch err {
			case resource.ErrEmptyGetContent:
				handler.NotFoundResponse(w, r)
			default:
				handler.ServerErrorResponse(w, r, err)
			}
		}

		handler.WriteJSON(w, handler.Envelope{"resource": resp}, nil)
	} else {
		resp, err := rh.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case resource.ErrEmptyGetContent:
				handler.NotFoundResponse(w, r)
			default:
				handler.ServerErrorResponse(w, r, err)
			}
		}

		handler.WriteJSON(w, handler.Envelope{"resources": resp}, nil)
	}
}

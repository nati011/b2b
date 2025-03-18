package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/category"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type Category struct {
	service category.Provider
}

func InitCategory() {
	handler.Register(new(Category))
}

func (r *Category) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.CategoryService
	return nil
}

func (c *Category) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/category", c.GetHandler)
	mux.HandleFunc("POST /api/v1/category", c.CreateHandler)
	mux.HandleFunc("DELETE /api/v1/category", c.DeleteHandler)
}

func (c *Category) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, r, err)
			return
		}
		resp, err := c.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case category.ErrIdNotFound:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}
		util.WriteJSON(w, util.Envelope{"category": resp}, http.StatusAccepted)
	} else {
		resp, err := c.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case category.ErrEmptyGetContent:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}
		util.WriteJSON(w, util.Envelope{"categories": resp}, http.StatusAccepted)
	}
}

func (c *Category) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateCategoryRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	id, err := c.service.Create(r.Context(), (*category.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case category.ErrDescIsNotSupplied,
			category.ErrDuplicateName,
			category.ErrEmptyGetContent:

			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"category": id}, http.StatusAccepted)
}

func (c *Category) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, r, err)
			return
		}
		err = c.service.Remove(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case category.ErrIdNotFound:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}
	}
	util.WriteJSON(w, util.Envelope{"category": paramIdValue}, http.StatusAccepted)
}

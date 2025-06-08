package domain

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/category"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

type GetCategoryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GetAllCategoryResponse struct {
	List []GetCategoryResponse `json:"categories"`
}

type Category struct {
	authMiddleware middleware.Auth
	service        category.Provider
}

func InitCategory() {
	handler.Register(new(Category))

	handler.RegisterResource("/api/v1/category")
	handler.RegisterResource("/api/v1/category/{id}")
}

func (r *Category) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.CategoryService
	return nil
}

func (c *Category) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/category", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireNoAuthentication(http.HandlerFunc(c.GetHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/category", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireAuthentication(http.HandlerFunc(c.CreateHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/category/{id}", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireAuthentication(http.HandlerFunc(c.UpdateHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("DELETE /api/v1/category", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireAuthentication(http.HandlerFunc(c.DeleteHandler)).ServeHTTP(w, r)
	})
}

func (c *Category) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := c.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case category.ErrIdNotFound:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, GetCategoryResponse(resp))
	} else {
		var get_all_response GetAllCategoryResponse
		resp, err := c.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case category.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, get_all_response)
				return
			case category.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		for _, i := range resp.List {
			get_all_response.List = append(get_all_response.List, GetCategoryResponse(i))
		}
		util.OperationSuccessResponse(w, get_all_response)
		return
	}
}

func (c *Category) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateCategoryRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := c.service.Create(r.Context(), (*category.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case category.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"category": id})
}

func (c *Category) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateCategoryRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	err = c.service.Update(r.Context(), &category.UpdateRequest{
		Id:   typedParamId,
		Name: requestBody.Name,
	})
	if err != nil {
		switch err {
		case category.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessMessageResponse(w, "category updated successfully")
}

func (c *Category) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		err = c.service.Remove(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case category.ErrIdNotFound:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"category": paramIdValue})
}

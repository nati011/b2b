package handler

import (
	"encoding/json"
	// "errors"
	"io"
	"net/http"

	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
)

type CreateConfigurableProductRequest struct {
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	ExternalId    string   `json:"external_id"`
	AttributeKeys []string `json:"attributes"`
	Products      []int    `json:"products"`
	Images        []string `json:"images"`
}

type UpdateRequest struct {
	Id                int      `json:"id"`
	Name              string   `json:"name"`
	Desc              string   `json:"desc"`
	ExternalId        string   `json:"external_id"`
	Product           []int    `json:"product"`
	IsAvailableStatus bool     `json:"is_available"`
	Images            []string `json:"images"`
	AttributeKeys     []string `json:"attribute_keys"`
}

type ConfigurableProduct struct {
	service configurable_product.Provider
}

func InitConfigurableProduct() {
	handler.Register(new(ConfigurableProduct))
}

func (r *ConfigurableProduct) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.ConfigurableProductService
	return nil
}

func (p *ConfigurableProduct) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/configurable_product", p.CreateConfigurableProductHandler)
	// mux.HandleFunc("GET /api/v1/configurable_product", p.GetAllConfigurableProductsHandler)
	mux.HandleFunc("PUT /api/v1/configurable_product", p.UpdateHandler)
	mux.HandleFunc("PUT /api/v1/configurable_product/available", p.UpdateAvailabilityHandler)
	mux.HandleFunc("PUT /api/v1/configurable_product/disable", p.DisableConfigurableProductHandler)
}

func (p *ConfigurableProduct) CreateConfigurableProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateConfigurableProductRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := p.service.Create(r.Context(), (*configurable_product.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case configurable_product.ErrUnknown:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"product": id})
}

func (p *ConfigurableProduct) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	err = p.service.Update(r.Context(), (*configurable_product.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case configurable_product.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}

}

func (p *ConfigurableProduct) UpdateAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		err = p.service.Avail(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case configurable_product.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
	} else {
		util.RequestErrorResponse(w, util.ErrIdRequired)
	}

}

func (p *ConfigurableProduct) DisableConfigurableProductHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		err = p.service.Disable(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case configurable_product.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
	} else {
		util.RequestErrorResponse(w, util.ErrIdRequired)
	}

}

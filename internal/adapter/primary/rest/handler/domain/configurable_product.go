package domain

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
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
)

var (
	ErrUnknownCPCommand = errors.New("unknown command")
)

type CreateConfigurableProductRequest struct {
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	ExternalId    string   `json:"external_id"`
	AttributeKeys []string `json:"attribute_keys"`
	Products      []int    `json:"products"`
	Images        []string `json:"images"`
}

type PriceRangeResponse struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type GetConfigurableProductResponse struct {
	Id            int                 `json:"id"`
	Name          string              `json:"name"`
	Desc          string              `json:"desc"`
	ExternalId    string              `json:"external_id"`
	Attributes    []map[string]string `json:"attributes"`
	Products      []int               `json:"products"`
	IsAvailable   bool                `json:"is_available"`
	PriceRange    PriceRangeResponse  `json:"price_range"`
	CategoryId    []int               `json:"category_id"`
	DistributorId int                 `json:"distributor_id"`
	Images        []string            `json:"images"`
}

type GetAllConfigurableProductsResponse struct {
	List []GetConfigurableProductResponse `json:"configurable_products"`
}

type GetConfigurableProductsByParamRequest struct {
	Name       string `json:"name"`
	ExternalId string `json:"external_id"`
}

type UpdateConfigurableProductRequest struct {
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
	authMiddleware middleware.Auth
	service        configurable_product.Provider
}

func InitConfigurableProduct() {
	handler.Register(new(ConfigurableProduct))

	handler.RegisterResource("/api/v1/configurable_product")
	handler.RegisterResource("/api/v1/configurable_product/{id}/status")
}

func (c *ConfigurableProduct) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	c.service = domainService.ConfigurableProductService
	c.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (cp *ConfigurableProduct) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/configurable_product", func(w http.ResponseWriter, r *http.Request) {
		cp.authMiddleware.RequireNoAuthentication(http.HandlerFunc(cp.GetConfigurableProductsHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/configurable_product", func(w http.ResponseWriter, r *http.Request) {
		cp.authMiddleware.RequireAuthentication(http.HandlerFunc(cp.CreateConfigurableProductHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PUT /api/v1/configurable_product", func(w http.ResponseWriter, r *http.Request) {
		cp.authMiddleware.RequireAuthentication(http.HandlerFunc(cp.UpdateHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/configurable_product/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		cp.authMiddleware.RequireAuthentication(http.HandlerFunc(cp.StatusCommandHandler)).ServeHTTP(w, r)
	})
}

func (cp *ConfigurableProduct) GetConfigurableProductsHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	const ParamExternalId = "external_id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	paramNameValue := paramValues.Get(ParamName)
	ParamExternalIdValue := paramValues.Get(ParamExternalId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := cp.service.Get(r.Context(), typedParamId)
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
		util.OperationSuccessResponse(w, util.Envelope{"configurable_product": resp})
	} else if paramNameValue != "" {
		cp_name_resp, err := cp.service.GetByParam(r.Context(),
			&configurable_product.GetByParamRequest{
				Name: paramNameValue,
			})
		if err != nil {
			switch err {
			case configurable_product.ErrEmptyGetContent:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"configurable_products": cp_name_resp})
	} else if ParamExternalIdValue != "" {
		cp_extId_resp, err := cp.service.GetByParam(r.Context(),
			&configurable_product.GetByParamRequest{
				ExternalId: ParamExternalIdValue,
			})
		if err != nil {
			switch err {
			case configurable_product.ErrEmptyGetContent:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"configurable_products": cp_extId_resp})
	} else {
		configurable_products, err := cp.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case configurable_product.ErrEmptyGetContent:
			case configurable_product.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"configurable_products": configurable_products})
	}
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
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"configurable_product": id})
}

func (p *ConfigurableProduct) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateConfigurableProductRequest
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
	util.OperationSuccessResponse(w, util.Envelope{"configurable_product": requestBody.Id})

}

const (
	AVAIL_CP_COMMAND   = "activate"
	DISABLE_CP_COMMAND = "deactivate"
)

func (p *ConfigurableProduct) StatusCommandHandler(w http.ResponseWriter, r *http.Request) {
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
		case AVAIL_CP_COMMAND:
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
			util.OperationSuccessResponse(w, util.Envelope{"detail": "configurable product successfully activated"})
		case DISABLE_CP_COMMAND:
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
			util.OperationSuccessResponse(w, util.Envelope{"detail": "configurable product successfully deactivated"})
		default:
			util.RequestErrorResponse(w, ErrUnknownCPCommand)
		}
	}
}

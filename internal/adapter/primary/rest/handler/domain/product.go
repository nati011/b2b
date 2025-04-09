package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrUnknownProductCommand = errors.New("unknown command")
)

type CreateProductRequest struct {
	Name          string            `json:"name"`
	Desc          string            `json:"desc"`
	ExternalID    string            `json:"external_id"`
	Images        []string          `json:"images"`
	Price         float64           `json:"price"`
	Attributes    map[string]string `json:"attributes"`
	DistributorId int               `json:"distributor_id"`
	CategoryId    []int             `json:"category_id"`
}

type CreateConfigurableProductRequest struct {
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	ExternalId    string   `json:"external_id"`
	AttributeKeys []string `json:"attributes"`
	Products      []int    `json:"products"`
	Images        []string `json:"images"`
}

type ProductResponse struct {
	Id            int               `json:"id"`
	Name          string            `json:"name"`
	Desc          string            `json:"desc"`
	Price         float64           `json:"price"`
	Stock         int               `json:"stock"`
	ExternalID    string            `json:"external_id"`
	Attributes    map[string]string `json:"attributes"`
	Images        []string          `json:"images"`
	DistributorId int               `json:"distributor_id"`
	CategoryId    []int             `json:"categories"`
	IsActive      bool              `json:"is_active"`
}

type ConfigurableAttributesResponse struct {
	ProductId      int    `json:"product_id"`
	AttributeValue string `json:"attribute_value"`
}

type GetProductResponse struct {
	Name                   string                                      `json:"name"`
	Desc                   string                                      `json:"desc"`
	IsActive               bool                                        `json:"is_active"`
	Images                 []string                                    `json:"images"`
	ConfigurableAttributes map[string][]ConfigurableAttributesResponse `json:"configurable_attributes"`
	Configurables          []ProductResponse                           `json:"configurables"`
}

type GetAllProductResponse struct {
	List []GetProductResponse `json:"products"`
}

type GetProductByParamRequest struct {
	Name          string `json:"name"`
	ExternalID    string `json:"external_id"`
	DistributorId int    `json:"distributor_id"`
	CategoryId    []int  `json:"categories"`
	PriceMin      int    `json:"price_min"`
	PriceMax      int    `json:"price_max"`
}

type UpdateProductRequest struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	ExternalID string   `json:"extenal_id"`
	Price      float64  `json:"price"`
	Desc       string   `json:"desc"`
	Images     []string `json:"images"`
	CategoryId []int    `json:"categories"`
}

type GoodsReceivingRequest struct {
	Id     int `json:"id"`
	Amount int `json:"amount"`
}

type DispatchRequest struct {
	Id     int `json:"id"`
	Amount int `json:"amount"`
}

type GetProductsWithCategoriesRequest struct {
	List []int `json:"products"`
}

type Product struct {
	service                    product.Provider
	configurableProductservice configurable_product.Provider
}

func InitProduct() {
	handler.Register(new(Product))
}

func (r *Product) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.ProductService
	r.configurableProductservice = domainService.ConfigurableProductService
	return nil
}

func (p *Product) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/product", p.GetHandler)
	mux.HandleFunc("POST /api/v1/product", p.CreateHandler)
	mux.HandleFunc("PUT /api/v1/product", p.UpdateHandler)
	mux.HandleFunc("PATCH /api/v1/product/status", p.StatusHandler)
	mux.HandleFunc("PATCH /api/v1/product/stock", p.StockHandler)

	mux.HandleFunc("POST /api/v1/configurable_product", p.CreateConfigurableProductHandler)
	mux.HandleFunc("GET /api/v1/configurable_product", p.GetConfigurableProductHandler)
}

func (p *Product) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	const ParamExternalId = "external_id"
	const ParamCategoryId = "category_id"
	const ParamPriceMin = "price_min"
	const ParamPriceMax = "price_max"

	paramValues := r.URL.Query()
	paramNameValue := paramValues.Get(ParamName)
	ParamExternalIdValue := paramValues.Get(ParamExternalId)
	ParamCategoryIdValue := paramValues.Get(ParamCategoryId)
	ParamPriceMinValue := paramValues.Get(ParamPriceMin)
	ParamPriceMaxValue := paramValues.Get(ParamPriceMax)

	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		resp, err := p.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case product.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"product": ProductResponse{
			Id:            resp.Id,
			Name:          resp.Name,
			Desc:          resp.Desc,
			ExternalID:    resp.ExternalID,
			Images:        resp.Images,
			Price:         resp.Price,
			Attributes:    resp.Attributes,
			DistributorId: resp.DistributorId,
			CategoryId:    resp.CategoryId,
			Stock:         resp.Stock,
			IsActive:      resp.IsActive,
		}})

	} else if ParamCategoryIdValue != "" || ParamPriceMinValue != "" || ParamPriceMaxValue != "" {
		var typedCategoryId int
		var err error
		if ParamCategoryIdValue != "" {
			typedCategoryId, err = strconv.Atoi(ParamCategoryIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}

		var typedPriceMin int
		if ParamPriceMinValue != "" {
			typedPriceMin, err = strconv.Atoi(ParamPriceMinValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}

		var typedPriceMax int
		if ParamPriceMaxValue != "" {
			typedPriceMax, err = strconv.Atoi(ParamPriceMaxValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}

		var resp product.GetAllResponse

		if typedCategoryId != 0 {
			resp, err = p.service.GetByParam(r.Context(), &product.GetByParamRequest{
				Name:       paramNameValue,
				ExternalID: ParamExternalIdValue,
				CategoryId: []int{typedCategoryId},
				PriceMin:   typedPriceMin,
				PriceMax:   typedPriceMax,
			})
		} else {
			resp, err = p.service.GetByParam(r.Context(), &product.GetByParamRequest{
				Name:       paramNameValue,
				ExternalID: ParamExternalIdValue,
				PriceMin:   typedPriceMin,
				PriceMax:   typedPriceMax,
			})
		}

		if err != nil {
			switch err {
			case product.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"products": resp})
	} else {
		// build
		var resp []GetProductResponse
		pr, err := p.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case product.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		var configurables []ProductResponse
		for _, i := range pr.List {
			configurables = append(configurables, ProductResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        i.Images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			})

			var configurableAttribute = make(map[string][]ConfigurableAttributesResponse)
			for attr_key, attr_val := range i.Attributes {
				configurableAttribute[attr_key] = []ConfigurableAttributesResponse{
					{
						ProductId:      i.Id,
						AttributeValue: attr_val,
					},
				}
			}

			resp = append(resp, GetProductResponse{
				Name:                   i.Name,
				Desc:                   i.Desc,
				IsActive:               i.IsActive,
				Images:                 i.Images,
				ConfigurableAttributes: configurableAttribute,
				Configurables:          configurables,
			})
		}

		cp, err := p.configurableProductservice.GetAll(r.Context())
		if err != nil {
			switch err {
			case configurable_product.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}

		for _, j := range cp.List {
			var configurables []ProductResponse
			var configurableAttribute = make(map[string][]ConfigurableAttributesResponse)
			for _, i := range j.Products {
				resp, err := p.service.Get(r.Context(), i)
				if err != nil {
					if err != nil {
						switch err {
						case product.ErrIdNotFound:
						default:
							util.ServerErrorResponse(w, err)
						}
					}
				}
				configurables = append(configurables, ProductResponse{
					Id:            resp.Id,
					Name:          resp.Name,
					Desc:          resp.Desc,
					ExternalID:    resp.ExternalID,
					Images:        resp.Images,
					Price:         resp.Price,
					Attributes:    resp.Attributes,
					DistributorId: resp.DistributorId,
					CategoryId:    resp.CategoryId,
					Stock:         resp.Stock,
					IsActive:      resp.IsActive,
				})
				for attr_key, _ := range j.Attributes {
					configurableAttribute[attr_key] = append(configurableAttribute[attr_key], ConfigurableAttributesResponse{
						ProductId:      resp.Id,
						AttributeValue: resp.Attributes[attr_key],
					})
				}
			}

			resp = append(resp, GetProductResponse{
				Name:                   j.Name,
				Desc:                   j.Desc,
				IsActive:               j.IsAvailable,
				Images:                 j.Images,
				ConfigurableAttributes: configurableAttribute,
				Configurables:          configurables,
			})
		}
		util.OperationSuccessResponse(w, util.Envelope{"products": resp})
	}
}

func (p *Product) GetConfigurableProductHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	paramValues := r.URL.Query()
	paramNameValue := paramValues.Get(ParamName)
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		cp_resp, err := p.configurableProductservice.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case configurable_product.ErrIdNotFound:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
			}
		}

		var configurables []ProductResponse
		var configurableAttribute = make(map[string][]ConfigurableAttributesResponse)
		for _, i := range cp_resp.Products {
			resp, err := p.service.Get(r.Context(), i)
			if err != nil {
				if err != nil {
					switch err {
					case product.ErrIdNotFound:
					default:
						util.ServerErrorResponse(w, err)
					}
				}
			}
			configurables = append(configurables, ProductResponse{
				Id:            resp.Id,
				Name:          resp.Name,
				Desc:          resp.Desc,
				ExternalID:    resp.ExternalID,
				Images:        resp.Images,
				Price:         resp.Price,
				Attributes:    resp.Attributes,
				DistributorId: resp.DistributorId,
				CategoryId:    resp.CategoryId,
				Stock:         resp.Stock,
				IsActive:      resp.IsActive,
			})

			for attr_key, _ := range cp_resp.Attributes {
				configurableAttribute[attr_key] = []ConfigurableAttributesResponse{
					{
						ProductId:      resp.Id,
						AttributeValue: resp.Attributes[attr_key],
					},
				}
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"configurable_product": GetProductResponse{
			Name:                   cp_resp.Name,
			Desc:                   cp_resp.Desc,
			IsActive:               cp_resp.IsAvailable,
			Images:                 cp_resp.Images,
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		}})

	} else if paramNameValue != "" {
		cp_resp, err := p.configurableProductservice.GetByParam(r.Context(),
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
		var resp []GetProductResponse
		for _, j := range cp_resp.List {

			var configurables []ProductResponse
			var configurableAttribute = make(map[string][]ConfigurableAttributesResponse)
			for _, i := range j.Products {
				resp, err := p.service.Get(r.Context(), i)
				if err != nil {
					if err != nil {
						switch err {
						case product.ErrIdNotFound:
						default:
							util.ServerErrorResponse(w, err)
						}
					}
				}
				configurables = append(configurables, ProductResponse{
					Id:            resp.Id,
					Name:          resp.Name,
					Desc:          resp.Desc,
					ExternalID:    resp.ExternalID,
					Images:        resp.Images,
					Price:         resp.Price,
					Attributes:    resp.Attributes,
					DistributorId: resp.DistributorId,
					CategoryId:    resp.CategoryId,
					Stock:         resp.Stock,
					IsActive:      resp.IsActive,
				})
				for attr_key, _ := range resp.Attributes {
					configurableAttribute[attr_key] = []ConfigurableAttributesResponse{
						{
							ProductId:      resp.Id,
							AttributeValue: resp.Attributes[attr_key],
						},
					}
				}
			}

			resp = append(resp, GetProductResponse{
				Name:                   j.Name,
				Desc:                   j.Desc,
				IsActive:               j.IsAvailable,
				Images:                 j.Images,
				ConfigurableAttributes: configurableAttribute,
				Configurables:          configurables,
			})
		}
		util.OperationSuccessResponse(w, util.Envelope{"configurable_products": resp})
	}
}

func (p *Product) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateProductRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := p.service.Create(r.Context(), (*product.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case product.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"product": id})
}

func (p *Product) CreateConfigurableProductHandler(w http.ResponseWriter, r *http.Request) {
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
	id, err := p.configurableProductservice.Create(r.Context(), (*configurable_product.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case product.ErrUnknown:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"product": id})
}

func (p *Product) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateProductRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := p.service.Update(r.Context(), (*product.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case product.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"product": id})
}

const (
	ACTIVATE_COMMAND   = "activate"
	DEACTIVATE_COMMAND = "deactivate"
)

func (p *Product) StatusHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamCommand = "command"
	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	paramIdValue := paramValues.Get(ParamId)

	if paramCommandValue != "" && paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case ACTIVATE_COMMAND:
			err = p.service.Activate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case product.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return

				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"id": typedParamId})
		case DEACTIVATE_COMMAND:
			err = p.service.Deactivate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case product.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return

				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"id": typedParamId})
		default:
			util.RequestErrorResponse(w, ErrUnknownProductCommand)
		}
	}
}

const (
	RECEIVEGOODS_COMMAND = "receive"
	DEPLETE_COMMAND      = "deplete"
)

func (p *Product) StockHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamCommand = "command"
	const ParamAmount = "amount"

	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	paramIdValue := paramValues.Get(ParamId)
	ParamAmountValue := paramValues.Get(ParamAmount)

	if paramCommandValue != "" && paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		typedParamAmount, err := strconv.Atoi(ParamAmountValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case RECEIVEGOODS_COMMAND:
			err = p.service.ReceiveGoods(r.Context(), &product.GoodsReceivingRequest{
				Id:     typedParamId,
				Amount: typedParamAmount,
			})
			if err != nil {
				switch err {
				case product.ErrUnknown:
					util.ServerErrorResponse(w, err)
				default:
					util.RequestErrorResponse(w, err)
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"id": typedParamId})
		case DEPLETE_COMMAND:
			err = p.service.Dispatch(r.Context(), &product.DispatchRequest{
				Id:     typedParamId,
				Amount: typedParamAmount,
			})
			if err != nil {
				switch err {
				case product.ErrUnknown:
					util.ServerErrorResponse(w, err)
				default:
					util.RequestErrorResponse(w, err)
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"id": typedParamId})
		default:
			util.RequestErrorResponse(w, ErrUnknownProductCommand)
		}
	}
}

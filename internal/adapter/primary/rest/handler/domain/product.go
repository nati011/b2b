package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core/domain/product"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
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

type GetProductResponse struct {
	Id            int               `json:"id"`
	Name          string            `json:"name"`
	Desc          string            `json:"desc"`
	ExternalID    string            `json:"external_id"`
	Images        []string          `json:"images"`
	Price         float64           `json:"price"`
	Attributes    map[string]string `json:"attributes"`
	DistributorId int               `json:"distributor_id"`
	CategoryId    []int             `json:"categories"`
	Stock         int               `json:"stock"`
	IsActive      bool              `json:"is_active"`
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
	service product.Provider
}

func InitProduct() {
	handler.Register(new(Product))
}

func (r *Product) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.ProductService
	return nil
}

func (p *Product) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/product", p.GetHandler)
	mux.HandleFunc("POST /api/v1/product", p.CreateHandler)
	mux.HandleFunc("PUT /api/v1/product", p.UpdateHandler)
	mux.HandleFunc("PATCH /api/v1/product/status", p.StatusHandler)
	mux.HandleFunc("PATCH /api/v1/product/stock", p.StockHandler)
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
			util.RequestErrorResponse(w, r, err)
			return
		}
		resp, err := p.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case product.ErrIdNotFound:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}

		util.WriteJSON(w, util.Envelope{"product": resp}, http.StatusAccepted)
	} else {
		var typedCategoryId int
		var err error
		if ParamCategoryIdValue != "" {
			typedCategoryId, err = strconv.Atoi(ParamCategoryIdValue)
			if err != nil {
				util.RequestErrorResponse(w, r, err)
				return
			}
		}

		var typedPriceMin int
		if ParamPriceMinValue != "" {
			typedPriceMin, err = strconv.Atoi(ParamPriceMinValue)
			if err != nil {
				util.RequestErrorResponse(w, r, err)
				return
			}
		}

		var typedPriceMax int
		if ParamPriceMaxValue != "" {
			typedPriceMax, err = strconv.Atoi(ParamPriceMaxValue)
			if err != nil {
				util.RequestErrorResponse(w, r, err)
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
			case product.ErrCategoryNotFound,
				product.ErrEmptyGetContent:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}
		util.WriteJSON(w, util.Envelope{"products": resp}, http.StatusAccepted)
	}
}

func (p *Product) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateProductRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	id, err := p.service.Create(r.Context(), (*product.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case product.ErrIdNotFound,
			product.ErrNameNotSupplied,
			product.ErrNameDuplicate,
			product.ErrImagesMustBeAtleastTwo,
			product.ErrPriceNotSupplied,
			product.ErrAttributeValuesCannotBeEmpty,
			product.ErrCategoryNotFound:
			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"product": id}, http.StatusAccepted)
}

func (p *Product) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateProductRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	id, err := p.service.Update(r.Context(), (*product.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case product.ErrIdNotFound:
			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"product": id}, http.StatusAccepted)
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
			util.RequestErrorResponse(w, r, err)
			return
		}
		switch paramCommandValue {
		case ACTIVATE_COMMAND:
			err = p.service.Activate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case product.ErrIdNotFound,
					product.ErrAlreadyActive:

					util.RequestErrorResponse(w, r, err)
					return
				default:
					util.ServerErrorResponse(w, r, err)
					return
				}
			}
			return
		case DEACTIVATE_COMMAND:
			err = p.service.Deactivate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case product.ErrIdNotFound,
					product.ErrAlreadyInactive:

					util.RequestErrorResponse(w, r, err)
					return
				default:
					util.ServerErrorResponse(w, r, err)
					return
				}
			}
			return
		default:
			util.RequestErrorResponse(w, r, errors.New("unknown command"))
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
			util.RequestErrorResponse(w, r, err)
			return
		}
		typedParamAmount, err := strconv.Atoi(ParamAmountValue)
		if err != nil {
			util.RequestErrorResponse(w, r, err)
			return
		}
		switch paramCommandValue {
		case RECEIVEGOODS_COMMAND:
			p.service.ReceiveGoods(r.Context(), &product.GoodsReceivingRequest{
				Id:     typedParamId,
				Amount: typedParamAmount,
			})
			return
		case DEPLETE_COMMAND:
			p.service.Dispatch(r.Context(), &product.DispatchRequest{
				Id:     typedParamId,
				Amount: typedParamAmount,
			})
			return
		default:
			util.RequestErrorResponse(w, r, errors.New("unknown command"))
		}
	}
}

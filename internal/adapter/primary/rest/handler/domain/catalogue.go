package domain

import (
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/catalogue"
)

type Image struct {
	ImageUrl string
	BlurHash string
}

type CatalogueResponse struct {
	Id            int               `json:"id"`
	Name          string            `json:"name"`
	Desc          string            `json:"desc"`
	Price         float64           `json:"price"`
	Stock         int               `json:"stock"`
	ExternalID    string            `json:"external_id"`
	Attributes    map[string]string `json:"attributes"`
	Images        []Image           `json:"images"`
	DistributorId int               `json:"distributor_id"`
	CategoryId    []int             `json:"categories"`
	IsActive      bool              `json:"is_active"`
}

type CatalogueConfigurableAttributesResponse struct {
	ProductId      int    `json:"product_id"`
	AttributeValue string `json:"attribute_value"`
}

type GetCatalogueResponse struct {
	Name                   string                                               `json:"name"`
	Desc                   string                                               `json:"desc"`
	IsActive               bool                                                 `json:"is_active"`
	Images                 []Image                                              `json:"images"`
	ConfigurableAttributes map[string][]CatalogueConfigurableAttributesResponse `json:"configurable_attributes"`
	Configurables          []CatalogueResponse                                  `json:"configurables"`
}

type ProductResponse struct {
	Id             int               `json:"id"`
	Name           string            `json:"name"`
	Desc           string            `json:"desc"`
	Price          float64           `json:"price"`
	Stock          int               `json:"stock"`
	AvailableStock int               `json:"available_stock"`
	ReservedStock  int               `json:"reserved_stock"`
	ExternalID     string            `json:"external_id"`
	Attributes     map[string]string `json:"attributes"`
	Images         []Image           `json:"images"`
	DistributorId  int               `json:"distributor_id"`
	CategoryId     []int             `json:"categories"`
	IsActive       bool              `json:"is_active"`
}

type GetProductResponse struct {
	Name                   string                                      `json:"name"`
	Desc                   string                                      `json:"desc"`
	IsActive               bool                                        `json:"is_active"`
	Images                 []string                                    `json:"images"`
	ConfigurableAttributes map[string][]ConfigurableAttributesResponse `json:"configurable_attributes"`
	Configurables          []ProductResponse                           `json:"configurables"`
}

type GetAllCatalogueResponse struct {
	List []GetProductResponse `json:"products"`
}

type Catalogue struct {
	authMiddleware middleware.Auth
	service        catalogue.Provider
}

func InitCatalogue() {
	handler.Register(new(Catalogue))

	handler.RegisterResource(resource.CreateRequest{
		Name:     "catalogue",
		Action:   "ALL",
		Resource: "/api/v1/catalogue",
	})

	handler.RegisterResource(resource.CreateRequest{
		Name:     "catalogue_search",
		Action:   "ALL",
		Resource: "/api/v1/catalogue/search",
	})
}

func (c *Catalogue) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	c.service = domainService.CatalogueService
	c.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (c *Catalogue) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/catalogue", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireNoAuthentication(http.HandlerFunc(c.GetCatalogueHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/catalogue/search", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireNoAuthentication(http.HandlerFunc(c.SearchCatalogueHandler)).ServeHTTP(w, r)
	})
}

func (c *Catalogue) GetCatalogueHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := c.service.GetAll(r.Context())
	if err != nil {
		switch err {
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, resp)
}

func (c *Catalogue) SearchCatalogueHandler(w http.ResponseWriter, r *http.Request) {
	const ParamName = "name"
	const ParamPriceMin = "price_min"
	const ParamPriceMax = "price_max"
	paramValues := r.URL.Query()
	paramNameValue := paramValues.Get(ParamName)
	paramPriceMaxValue := paramValues.Get(ParamPriceMax)
	paramPriceMinValue := paramValues.Get(ParamPriceMin)

	if paramPriceMinValue != "" && paramPriceMaxValue != "" {
		typedParamPriceMinValue, err := strconv.Atoi(paramPriceMinValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}

		typedParamPriceMaxValue, err := strconv.Atoi(paramPriceMaxValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := c.service.Search(r.Context(), &catalogue.SearchCatalogueRequest{
			Name:     paramNameValue,
			PriceMin: typedParamPriceMinValue,
			PriceMax: typedParamPriceMaxValue,
		})
		if err != nil {
			switch err {
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, resp)
		return
	}

	resp, err := c.service.Search(r.Context(), &catalogue.SearchCatalogueRequest{
		Name: paramNameValue,
	})
	if err != nil {
		switch err {
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, resp)
}

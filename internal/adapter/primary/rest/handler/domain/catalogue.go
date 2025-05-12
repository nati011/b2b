package domain

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
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
	service catalogue.Provider
}

func InitCatalogue() {
	handler.Register(new(Catalogue))
}

func (c *Catalogue) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	c.service = domainService.CatalogueService
	return nil
}

func (c *Catalogue) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/catalogue", c.GetCatalogueHandler)
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

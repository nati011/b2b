package domain

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type CatalogueResponse struct {
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

type CatalogueConfigurableAttributesResponse struct {
	ProductId      int    `json:"product_id"`
	AttributeValue string `json:"attribute_value"`
}

type GetCatalogueResponse struct {
	Name                   string                                               `json:"name"`
	Desc                   string                                               `json:"desc"`
	IsActive               bool                                                 `json:"is_active"`
	Images                 []string                                             `json:"images"`
	ConfigurableAttributes map[string][]CatalogueConfigurableAttributesResponse `json:"configurable_attributes"`
	Configurables          []CatalogueResponse                                  `json:"configurables"`
}

type GetAllCatalogueResponse struct {
	List []GetProductResponse `json:"products"`
}

type Catalogue struct {
	service                    product.Provider
	configurableProductservice configurable_product.Provider
}

func InitCatalogue() {
	handler.Register(new(Catalogue))
}

func (r *Catalogue) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.ProductService
	r.configurableProductservice = domainService.ConfigurableProductService
	return nil
}

func (p *Catalogue) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/catalogue", p.GetCatalogueHandler)
}

func (p *Catalogue) GetCatalogueHandler(w http.ResponseWriter, r *http.Request) {
	// configurable products
	var resp []GetCatalogueResponse
	products_belonging_to_cps := []int{}
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
		var configurables []CatalogueResponse
		var configurableAttribute = make(map[string][]CatalogueConfigurableAttributesResponse)
		for _, i := range j.Products {
			//add product to cps belonging to products
			products_belonging_to_cps = append(products_belonging_to_cps, i)
			resp, err := p.service.Get(r.Context(), i)
			if err != nil {
				switch err {
				case product.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, err)
				}
			}
			configurables = append(configurables, CatalogueResponse{
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
			for _, key := range j.Attributes {
				configurableAttribute[key] = append(configurableAttribute[key], CatalogueConfigurableAttributesResponse{
					ProductId:      resp.Id,
					AttributeValue: resp.Attributes[key],
				})

			}
		}

		resp = append(resp, GetCatalogueResponse{
			Name:                   j.Name,
			Desc:                   j.Desc,
			IsActive:               j.IsAvailable,
			Images:                 j.Images,
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	// standalone products
	pr, err := p.service.GetAll(r.Context())
	if err != nil {
		switch err {
		case product.ErrEmptyGetContent:
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	var configurables []CatalogueResponse
	for _, i := range pr.List {
		if !productInSlice(products_belonging_to_cps, i.Id) {
			configurables = append(configurables, CatalogueResponse{
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

			var configurableAttribute = make(map[string][]CatalogueConfigurableAttributesResponse)
			for attr_key, attr_val := range i.Attributes {
				configurableAttribute[attr_key] = []CatalogueConfigurableAttributesResponse{
					{
						ProductId:      i.Id,
						AttributeValue: attr_val,
					},
				}
			}

			resp = append(resp, GetCatalogueResponse{
				Name:                   i.Name,
				Desc:                   i.Desc,
				IsActive:               i.IsActive,
				Images:                 i.Images,
				ConfigurableAttributes: configurableAttribute,
				Configurables:          configurables,
			})
		}
	}

	util.OperationSuccessResponse(w, util.Envelope{"products": resp})
}

func productInSlice(productsSlice []int, value int) bool {
	for _, v := range productsSlice {
		if v == value {
			return true
		}
	}
	return false
}

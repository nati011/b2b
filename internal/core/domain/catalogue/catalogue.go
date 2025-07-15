package catalogue

import (
	"context"

	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"
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
	Name                   string                                               `json:"name"`
	Desc                   string                                               `json:"desc"`
	IsActive               bool                                                 `json:"is_active"`
	Images                 []string                                             `json:"images"`
	ConfigurableAttributes map[string][]CatalogueConfigurableAttributesResponse `json:"configurable_attributes"`
	Configurables          []ProductResponse                                    `json:"configurables"`
}

type GetRequest struct {
	minPrice int32
	maxPrice int32
}

type GetAllCatalogueResponse struct {
	List []GetCatalogueResponse `json:"products"`
}

type Provider interface {
	GetAll(ctx context.Context) (GetAllCatalogueResponse, error)
	Get(ctx context.Context, req *GetRequest) (GetAllCatalogueResponse, error)
}

type CatalogueService struct {
	ProductService             product.Provider
	ConfigurableProductservice configurable_product.Provider
}

func NewCatalogueService(productService product.Provider, configurableProductservice configurable_product.Provider) Provider {
	return &CatalogueService{
		ProductService:             productService,
		ConfigurableProductservice: configurableProductservice,
	}
}

func (c *CatalogueService) GetAll(ctx context.Context) (GetAllCatalogueResponse, error) {
	// configurable products
	var resp []GetCatalogueResponse
	products_belonging_to_cps := []int{}
	cp, err := c.ConfigurableProductservice.GetAll(ctx)
	if err != nil {
		switch err {
		case configurable_product.ErrEmptyGetContent:
		default:
			return GetAllCatalogueResponse{}, err
		}
	}

	for _, j := range cp.List {
		var configurables []CatalogueResponse
		var configurableAttribute = make(map[string][]CatalogueConfigurableAttributesResponse)
		for _, i := range j.Products {
			//add product to cps belonging to products
			products_belonging_to_cps = append(products_belonging_to_cps, i)
			resp, err := c.ProductService.Get(ctx, i)
			if err != nil {
				switch err {
				case product.ErrIdNotFound:
				default:
					return GetAllCatalogueResponse{}, err
				}
			}
			var images []Image
			for _, value := range resp.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}

			configurables = append(configurables, CatalogueResponse{
				Id:            resp.Id,
				Name:          resp.Name,
				Desc:          resp.Desc,
				ExternalID:    resp.ExternalID,
				Images:        images,
				Price:         resp.Price,
				Attributes:    resp.Attributes,
				DistributorId: resp.DistributorId,
				CategoryId:    resp.CategoryId,
				Stock:         resp.AvailableStock,
				IsActive:      resp.IsActive,
			})
			for _, key := range j.Attributes {
				configurableAttribute[key] = append(configurableAttribute[key], CatalogueConfigurableAttributesResponse{
					ProductId:      resp.Id,
					AttributeValue: resp.Attributes[key],
				})

			}
		}
		var images []Image
		for _, value := range j.Images {
			image := Image{
				ImageUrl: value.ImageUrl,
				BlurHash: value.BlurHash,
			}
			images = append(images, image)
		}

		resp = append(resp, GetCatalogueResponse{
			Name:                   j.Name,
			Desc:                   j.Desc,
			IsActive:               j.IsAvailable,
			Images:                 images,
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	// standalone products
	pr, err := c.ProductService.GetAll(ctx)
	if err != nil {
		switch err {
		case product.ErrEmptyGetContent:
		default:
			return GetAllCatalogueResponse{}, err
		}
	}

	for _, i := range pr.List {
		if !productInSlice(products_belonging_to_cps, i.Id) {
			var images []Image
			var configurables []CatalogueResponse
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			configurables = append(configurables, CatalogueResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.AvailableStock,
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
			var imageCatalogue []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				imageCatalogue = append(imageCatalogue, image)
			}

			resp = append(resp, GetCatalogueResponse{
				Name:                   i.Name,
				Desc:                   i.Desc,
				IsActive:               i.IsActive,
				Images:                 imageCatalogue,
				ConfigurableAttributes: configurableAttribute,
				Configurables:          configurables,
			})
		}
	}
	return GetAllCatalogueResponse{
		List: resp,
	}, nil
}

func (c *CatalogueService) Get(ctx context.Context, req *GetRequest) (GetAllCatalogueResponse, error) {
	// configurable products
	var resp []GetCatalogueResponse
	products_belonging_to_cps := []int{}
	cp, err := c.ConfigurableProductservice.GetAll(ctx)
	if err != nil {
		switch err {
		case configurable_product.ErrEmptyGetContent:
		default:
			return GetAllCatalogueResponse{}, err
		}
	}

	for _, j := range cp.List {
		var configurables []CatalogueResponse
		var configurableAttribute = make(map[string][]CatalogueConfigurableAttributesResponse)
		for _, i := range j.Products {
			//add product to cps belonging to products
			products_belonging_to_cps = append(products_belonging_to_cps, i)
			resp, err := c.ProductService.Get(ctx, i)
			if err != nil {
				switch err {
				case product.ErrIdNotFound:
				default:
					return GetAllCatalogueResponse{}, err
				}
			}
			var images []Image
			for _, value := range resp.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}

			configurables = append(configurables, CatalogueResponse{
				Id:            resp.Id,
				Name:          resp.Name,
				Desc:          resp.Desc,
				ExternalID:    resp.ExternalID,
				Images:        images,
				Price:         resp.Price,
				Attributes:    resp.Attributes,
				DistributorId: resp.DistributorId,
				CategoryId:    resp.CategoryId,
				Stock:         resp.AvailableStock,
				IsActive:      resp.IsActive,
			})
			for _, key := range j.Attributes {
				configurableAttribute[key] = append(configurableAttribute[key], CatalogueConfigurableAttributesResponse{
					ProductId:      resp.Id,
					AttributeValue: resp.Attributes[key],
				})

			}
		}
		var images []Image
		for _, value := range j.Images {
			image := Image{
				ImageUrl: value.ImageUrl,
				BlurHash: value.BlurHash,
			}
			images = append(images, image)
		}

		resp = append(resp, GetCatalogueResponse{
			Name:                   j.Name,
			Desc:                   j.Desc,
			IsActive:               j.IsAvailable,
			Images:                 images,
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	// standalone products
	pr := product.GetAllResponse{}
	if req.minPrice == 0 && req.maxPrice == 0 {
		//if filter is used
		pr, err = c.ProductService.GetAll(ctx)
	} else {
		//if filter is not used
		pr, err = c.ProductService.GetByParam(ctx, &product.GetByParamRequest{
			PriceMin: int(req.minPrice),
			PriceMax: int(req.maxPrice),
		})
	}
	if err != nil {
		switch err {
		case product.ErrEmptyGetContent:
		default:
			return GetAllCatalogueResponse{}, err
		}
	}

	for _, i := range pr.List {
		if !productInSlice(products_belonging_to_cps, i.Id) {
			var images []Image
			var configurables []CatalogueResponse
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			configurables = append(configurables, CatalogueResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.AvailableStock,
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
			var imageCatalogue []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				imageCatalogue = append(imageCatalogue, image)
			}

			resp = append(resp, GetCatalogueResponse{
				Name:                   i.Name,
				Desc:                   i.Desc,
				IsActive:               i.IsActive,
				Images:                 imageCatalogue,
				ConfigurableAttributes: configurableAttribute,
				Configurables:          configurables,
			})
		}
	}
	return GetAllCatalogueResponse{
		List: resp,
	}, nil
}

func productInSlice(productsSlice []int, value int) bool {
	for _, v := range productsSlice {
		if v == value {
			return true
		}
	}
	return false
}

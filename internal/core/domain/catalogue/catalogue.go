package catalogue

import (
	"context"
	"fmt"
	"sync"
	"time"

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

type SearchCatalogueRequest struct {
	Name     string
	PriceMin int
	PriceMax int
}

type GetAllCatalogueResponse struct {
	List []GetCatalogueResponse `json:"products"`
}

type Provider interface {
	GetAll(ctx context.Context) (GetAllCatalogueResponse, error)
	Search(ctx context.Context, req *SearchCatalogueRequest) (GetAllCatalogueResponse, error)
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

// cacheEntry represents a cached catalogue response with expiration time
type cacheEntry struct {
	data      GetAllCatalogueResponse
	expiresAt time.Time
}

// CachedCatalogueService wraps a Provider with in-memory caching
type CachedCatalogueService struct {
	service      Provider
	cache        sync.Map
	cacheTTL     time.Duration
	searchCache  sync.Map
	searchTTL    time.Duration
}

// NewCachedCatalogueService creates a new cached catalogue service
// cacheTTL: TTL for GetAll cache (default: 5 minutes)
// searchTTL: TTL for Search cache (default: 2 minutes, shorter because search results may vary)
func NewCachedCatalogueService(service Provider, cacheTTL, searchTTL time.Duration) Provider {
	if cacheTTL == 0 {
		cacheTTL = 5 * time.Minute
	}
	if searchTTL == 0 {
		searchTTL = 2 * time.Minute
	}
	return &CachedCatalogueService{
		service:     service,
		cacheTTL:    cacheTTL,
		searchTTL:   searchTTL,
	}
}

// GetAll retrieves the catalogue, using cache if available
func (c *CachedCatalogueService) GetAll(ctx context.Context) (GetAllCatalogueResponse, error) {
	cacheKey := "catalogue:all"
	
	// Check cache
	if cached, found := c.cache.Load(cacheKey); found {
		entry := cached.(*cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.data, nil
		}
		// Cache expired, remove it
		c.cache.Delete(cacheKey)
	}
	
	// Fetch from service
	resp, err := c.service.GetAll(ctx)
	if err != nil {
		return GetAllCatalogueResponse{}, err
	}
	
	// Cache the response
	c.cache.Store(cacheKey, &cacheEntry{
		data:      resp,
		expiresAt: time.Now().Add(c.cacheTTL),
	})
	
	return resp, nil
}

// Search retrieves search results, using cache if available
func (c *CachedCatalogueService) Search(ctx context.Context, req *SearchCatalogueRequest) (GetAllCatalogueResponse, error) {
	// Create cache key from search parameters
	cacheKey := c.buildSearchCacheKey(req)
	
	// Check cache
	if cached, found := c.searchCache.Load(cacheKey); found {
		entry := cached.(*cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.data, nil
		}
		// Cache expired, remove it
		c.searchCache.Delete(cacheKey)
	}
	
	// Fetch from service
	resp, err := c.service.Search(ctx, req)
	if err != nil {
		return GetAllCatalogueResponse{}, err
	}
	
	// Cache the response
	c.searchCache.Store(cacheKey, &cacheEntry{
		data:      resp,
		expiresAt: time.Now().Add(c.searchTTL),
	})
	
	return resp, nil
}

// buildSearchCacheKey creates a cache key from search request parameters
func (c *CachedCatalogueService) buildSearchCacheKey(req *SearchCatalogueRequest) string {
	// Simple key construction - in production, consider using a hash for long strings
	return fmt.Sprintf("catalogue:search:%s:%d:%d", req.Name, req.PriceMin, req.PriceMax)
}

// InvalidateCache clears the cache (useful when products are updated)
func (c *CachedCatalogueService) InvalidateCache() {
	c.cache.Range(func(key, value interface{}) bool {
		c.cache.Delete(key)
		return true
	})
	c.searchCache.Range(func(key, value interface{}) bool {
		c.searchCache.Delete(key)
		return true
	})
}

func (c *CatalogueService) Search(ctx context.Context, req *SearchCatalogueRequest) (GetAllCatalogueResponse, error) {
	// Use parallel processing for independent operations
	var wg sync.WaitGroup
	var cpErr, prErr error
	var cp configurable_product.GetAllResponse
	var pr product.GetAllResponse

	// Fetch configurable products and standalone products in parallel
	wg.Add(2)

	go func() {
		defer wg.Done()
		cp, cpErr = c.ConfigurableProductservice.GetByParam(ctx, &configurable_product.GetByParamRequest{
			Name: req.Name,
		})
		if cpErr != nil {
			switch cpErr {
			case configurable_product.ErrEmptyGetContent:
				cpErr = nil
			}
		}
	}()

	go func() {
		defer wg.Done()
		pr, prErr = c.ProductService.Search(ctx, &product.SearchRequest{
			Name:     req.Name,
			PriceMin: req.PriceMin,
			PriceMax: req.PriceMax,
		})
		if prErr != nil {
			switch prErr {
			case product.ErrEmptyGetContent:
				prErr = nil
			}
		}
	}()

	wg.Wait()

	if cpErr != nil {
		return GetAllCatalogueResponse{}, cpErr
	}
	if prErr != nil {
		return GetAllCatalogueResponse{}, prErr
	}

	// Pre-allocate response slice with estimated capacity
	estimatedCapacity := len(cp.List) + len(pr.List)
	resp := make([]GetCatalogueResponse, 0, estimatedCapacity)
	
	// Use map for O(1) lookup instead of O(n) slice search
	products_belonging_to_cps := make(map[int]bool, len(cp.List)*2)

	// Process configurable products
	for _, j := range cp.List {
		// Pre-allocate configurables slice
		configurables := make([]CatalogueResponse, 0, len(j.Products))
		configurableAttribute := make(map[string][]CatalogueConfigurableAttributesResponse, len(j.Attributes))

		// Collect all product IDs first
		productIds := make([]int, 0, len(j.Products))
		for _, productId := range j.Products {
			products_belonging_to_cps[productId] = true
			productIds = append(productIds, productId)
		}

		// Fetch products (N+1 problem - will be optimized in Phase 2 with batch fetching)
		for _, productId := range productIds {
			productResp, err := c.ProductService.Get(ctx, productId)
			if err != nil {
				switch err {
				case product.ErrIdNotFound:
					continue
				default:
					return GetAllCatalogueResponse{}, err
				}
			}

			// Use optimized image conversion
			images := convertProductImages(productResp.Images)

			configurables = append(configurables, CatalogueResponse{
				Id:            productResp.Id,
				Name:          productResp.Name,
				Desc:          productResp.Desc,
				ExternalID:    productResp.ExternalID,
				Images:        images,
				Price:         productResp.Price,
				Attributes:    productResp.Attributes,
				DistributorId: productResp.DistributorId,
				CategoryId:    productResp.CategoryId,
				Stock:         productResp.AvailableStock,
				IsActive:      productResp.IsActive,
			})

			// Build configurable attributes map
			for _, key := range j.Attributes {
				if attrVal, exists := productResp.Attributes[key]; exists {
					configurableAttribute[key] = append(configurableAttribute[key], CatalogueConfigurableAttributesResponse{
						ProductId:      productResp.Id,
						AttributeValue: attrVal,
					})
				}
			}
		}

		// Convert configurable product images
		cpImages := convertConfigurableProductImages(j.Images)

		resp = append(resp, GetCatalogueResponse{
			Name:                   j.Name,
			Desc:                   j.Desc,
			IsActive:               j.IsAvailable,
			Images:                 cpImages,
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	// Process standalone products (exclude those already in configurable products)
	for _, i := range pr.List {
		if productInMap(products_belonging_to_cps, i.Id) {
			continue
		}

		// Use optimized image conversion
		images := convertProductImages(i.Images)

		// Pre-allocate configurables slice (single item for standalone product)
		configurables := make([]CatalogueResponse, 0, 1)
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

		// Build configurable attributes map
		configurableAttribute := make(map[string][]CatalogueConfigurableAttributesResponse, len(i.Attributes))
		for attrKey, attrVal := range i.Attributes {
			configurableAttribute[attrKey] = []CatalogueConfigurableAttributesResponse{
				{
					ProductId:      i.Id,
					AttributeValue: attrVal,
				},
			}
		}

		resp = append(resp, GetCatalogueResponse{
			Name:                   i.Name,
			Desc:                   i.Desc,
			IsActive:               i.IsActive,
			Images:                 images, // Reuse same images
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	return GetAllCatalogueResponse{
		List: resp,
	}, nil
}

func (c *CatalogueService) GetAll(ctx context.Context) (GetAllCatalogueResponse, error) {
	// Use parallel processing for independent operations
	var wg sync.WaitGroup
	var cpErr, prErr error
	var cp configurable_product.GetAllResponse
	var pr product.GetAllResponse

	// Fetch configurable products and standalone products in parallel
	wg.Add(2)

	go func() {
		defer wg.Done()
		cp, cpErr = c.ConfigurableProductservice.Catalogue(ctx)
		if cpErr != nil {
			switch cpErr {
			case configurable_product.ErrEmptyGetContent:
				cpErr = nil
			}
		}
	}()

	go func() {
		defer wg.Done()
		pr, prErr = c.ProductService.Catalogue(ctx)
		if prErr != nil {
			switch prErr {
			case product.ErrEmptyGetContent:
				prErr = nil
			}
		}
	}()

	wg.Wait()

	if cpErr != nil {
		return GetAllCatalogueResponse{}, cpErr
	}
	if prErr != nil {
		return GetAllCatalogueResponse{}, prErr
	}

	// Pre-allocate response slice with estimated capacity
	estimatedCapacity := len(cp.List) + len(pr.List)
	resp := make([]GetCatalogueResponse, 0, estimatedCapacity)
	
	// Use map for O(1) lookup instead of O(n) slice search
	products_belonging_to_cps := make(map[int]bool, len(cp.List)*2)

	// Process configurable products
	for _, j := range cp.List {
		// Pre-allocate configurables slice
		configurables := make([]CatalogueResponse, 0, len(j.Products))
		configurableAttribute := make(map[string][]CatalogueConfigurableAttributesResponse, len(j.Attributes))

		// Collect all product IDs first for batch processing
		productIds := make([]int, 0, len(j.Products))
		for _, productId := range j.Products {
			products_belonging_to_cps[productId] = true
			productIds = append(productIds, productId)
		}

		// Fetch products (N+1 problem - will be optimized in Phase 2 with batch fetching)
		// For now, we still need to call Get individually, but we've optimized the data structures
		for _, productId := range productIds {
			productResp, err := c.ProductService.Get(ctx, productId)
			if err != nil {
				switch err {
				case product.ErrIdNotFound:
					continue
				default:
					return GetAllCatalogueResponse{}, err
				}
			}

			// Use optimized image conversion
			images := convertProductImages(productResp.Images)

			configurables = append(configurables, CatalogueResponse{
				Id:            productResp.Id,
				Name:          productResp.Name,
				Desc:          productResp.Desc,
				ExternalID:    productResp.ExternalID,
				Images:        images,
				Price:         productResp.Price,
				Attributes:    productResp.Attributes,
				DistributorId: productResp.DistributorId,
				CategoryId:    productResp.CategoryId,
				Stock:         productResp.AvailableStock,
				IsActive:      productResp.IsActive,
			})

			// Build configurable attributes map
			for _, key := range j.Attributes {
				if attrVal, exists := productResp.Attributes[key]; exists {
					configurableAttribute[key] = append(configurableAttribute[key], CatalogueConfigurableAttributesResponse{
						ProductId:      productResp.Id,
						AttributeValue: attrVal,
					})
				}
			}
		}

		// Convert configurable product images
		cpImages := convertConfigurableProductImages(j.Images)

		resp = append(resp, GetCatalogueResponse{
			Name:                   j.Name,
			Desc:                   j.Desc,
			IsActive:               j.IsAvailable,
			Images:                 cpImages,
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	// Process standalone products (exclude those already in configurable products)
	for _, i := range pr.List {
		if productInMap(products_belonging_to_cps, i.Id) {
			continue
		}

		// Use optimized image conversion
			images := convertProductImages(i.Images)

		// Pre-allocate configurables slice (single item for standalone product)
		configurables := make([]CatalogueResponse, 0, 1)
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

		// Build configurable attributes map
		configurableAttribute := make(map[string][]CatalogueConfigurableAttributesResponse, len(i.Attributes))
		for attrKey, attrVal := range i.Attributes {
			configurableAttribute[attrKey] = []CatalogueConfigurableAttributesResponse{
				{
					ProductId:      i.Id,
					AttributeValue: attrVal,
				},
			}
		}

		resp = append(resp, GetCatalogueResponse{
			Name:                   i.Name,
			Desc:                   i.Desc,
			IsActive:               i.IsActive,
			Images:                 images, // Reuse same images
			ConfigurableAttributes: configurableAttribute,
			Configurables:          configurables,
		})
	}

	return GetAllCatalogueResponse{
		List: resp,
	}, nil
}

// productInMap checks if a product ID exists in a map (O(1) lookup)
func productInMap(productsMap map[int]bool, value int) bool {
	return productsMap[value]
}

// convertProductImages converts product images to catalogue images efficiently
func convertProductImages(portImages []product.Image) []Image {
	if len(portImages) == 0 {
		return nil
	}
	images := make([]Image, 0, len(portImages))
	for _, value := range portImages {
		images = append(images, Image{
			ImageUrl: value.ImageUrl,
			BlurHash: value.BlurHash,
		})
	}
	return images
}

// convertConfigurableProductImages converts configurable product images to catalogue images efficiently
func convertConfigurableProductImages(portImages []configurable_product.Image) []Image {
	if len(portImages) == 0 {
		return nil
	}
	images := make([]Image, 0, len(portImages))
	for _, value := range portImages {
		images = append(images, Image{
			ImageUrl: value.ImageUrl,
			BlurHash: value.BlurHash,
		})
	}
	return images
}

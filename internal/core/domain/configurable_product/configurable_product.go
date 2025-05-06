package configurable_product

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/domain/product"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/configurable_product"
)

var (
	ErrNameIsNotSupplied                = errors.New("oopsy, name is mandatory")
	ErrNameDuplicate                    = errors.New("oopsy, name already in use")
	ErrEmptyGetContent                  = errors.New("oopsy, empty get content")
	ErrUnknown                          = errors.New("oopsy, error unknown")
	ErrDescIsNotSupplied                = errors.New("oopsy, desc is mandatory")
	ErrImagesMustBeAtleastTwo           = errors.New("oopsy, atleast two images mandatory")
	ErrAttributeKeysMustBeAtleastOne    = errors.New("oopsy, atleast one attribute key mandatory")
	ErrProductsMustBeAtleastOne         = errors.New("oopsy, atleast product mandatory")
	ErrAlreadyAvailable                 = errors.New("oopsy, product already available")
	ErrAlreadyUnavailable               = errors.New("oopsy, product already unavailable")
	ErrIdNotFound                       = errors.New("oopsy, id not found")
	ErrProductNotFound                  = errors.New("oopsy, product not found")
	ErrAttributeKeysDoNotExistInProduct = errors.New("oopsy, attributes not found in products")
)

type CreateRequest struct {
	Name          string
	Desc          string
	ExternalId    string
	AttributeKeys []string
	Products      []int
	Images        []string
}

type PriceRangeResponse struct {
	Min int
	Max int
}

type GetResponse struct {
	Id            int
	Name          string
	Desc          string
	ExternalId    string
	Attributes    []string
	Products      []int
	IsAvailable   bool
	PriceRange    PriceRangeResponse
	CategoryId    []int
	DistributorId int
	Images        []string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name       string
	ExternalId string
}

type UpdateRequest struct {
	Id                int
	Name              string
	Desc              string
	ExternalId        string
	Product           []int
	IsAvailableStatus bool
	Images            []string
	AttributeKeys     []string
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Avail(ctx context.Context, id int) error
	Disable(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	Update(ctx context.Context, req *UpdateRequest) error
}

type ConfigurableProductService struct {
	DB             port.DB
	ProductService product.Provider
}

func NewConfigurableProductService(db port.DB, ps product.Provider) Provider {
	return &ConfigurableProductService{
		DB:             db,
		ProductService: ps,
	}
}

func (c *ConfigurableProductService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := c.validateName(ctx, req.Name)
	if err != nil {
		return 0, err
	}
	err = validateDesc(req.Desc)
	if err != nil {
		return 0, err
	}
	err = validateImages(req.Images)
	if err != nil {
		return 0, err
	}

	/*
		validate products before validating attribute keys
		because attribute validation relies on the products..
		we fetch the attributes in our products and ensure that the
		attributes exist
	*/
	err = c.validateProducts(ctx, req.Products)
	if err != nil {
		return 0, err
	}
	err = c.validateAttributekeys(ctx, req.AttributeKeys, req.Products)
	if err != nil {
		return 0, err
	}

	id, err := c.DB.Create(ctx, &port.CreateRequest{
		Name:              req.Name,
		Desc:              req.Desc,
		ExternalId:        req.ExternalId,
		IsAvailableStatus: false,
		Products:          req.Products,
		AttributeKeys:     req.AttributeKeys,
		Images:            req.Images,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}
	return id, nil
}

func (c *ConfigurableProductService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := c.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}

	return GetResponse{
		Id:            resp.Id,
		Name:          resp.Name,
		Desc:          resp.Desc,
		ExternalId:    resp.ExternalId,
		Products:      resp.Products,
		IsAvailable:   resp.IsAvailable,
		PriceRange:    PriceRangeResponse(resp.PriceRange),
		CategoryId:    resp.CategoryId,
		DistributorId: resp.DistributorId,
		Images:        resp.Images,
		Attributes:    resp.Attributes,
	}, nil
}

func (c *ConfigurableProductService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := []GetResponse{}
	if req.Name != "" {
		get_by_name_resp, err := c.DB.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
				return GetAllResponse{}, ErrEmptyGetContent
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range get_by_name_resp.List {
			resp = append(resp, GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalId:    i.ExternalId,
				Products:      i.Products,
				IsAvailable:   i.IsAvailable,
				PriceRange:    PriceRangeResponse(i.PriceRange),
				CategoryId:    i.CategoryId,
				DistributorId: i.DistributorId,
				Images:        i.Images,
			})
		}
	}

	if req.ExternalId != "" {
		get_by_name_resp, err := c.DB.GetByExternalId(ctx, req.ExternalId)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
				return GetAllResponse{}, ErrEmptyGetContent
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range get_by_name_resp.List {
			resp = append(resp, GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalId:    i.ExternalId,
				Products:      i.Products,
				IsAvailable:   i.IsAvailable,
				PriceRange:    PriceRangeResponse(i.PriceRange),
				CategoryId:    i.CategoryId,
				DistributorId: i.DistributorId,
				Images:        i.Images,
			})
		}
	}
	if len(resp) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return GetAllResponse{
		List: resp,
	}, nil
}

func (c *ConfigurableProductService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp := []GetResponse{}
	get_by_name_resp, err := c.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	for _, i := range get_by_name_resp.List {
		resp = append(resp, GetResponse{
			Id:            i.Id,
			Name:          i.Name,
			Desc:          i.Desc,
			ExternalId:    i.ExternalId,
			Products:      i.Products,
			IsAvailable:   i.IsAvailable,
			PriceRange:    PriceRangeResponse(i.PriceRange),
			CategoryId:    i.CategoryId,
			DistributorId: i.DistributorId,
			Images:        i.Images,
			Attributes:    i.Attributes,
		})
	}

	return GetAllResponse{
		List: resp,
	}, nil
}

func (c *ConfigurableProductService) Avail(ctx context.Context, id int) error {
	//validate
	resp, err := c.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//check if present
	if resp.IsAvailable {
		return ErrAlreadyAvailable
	}

	err = c.DB.UpdateIsAvailableStatus(ctx, &port.UpdateIsAvailableStatusRequest{
		Id:     id,
		Status: true,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (c *ConfigurableProductService) Disable(ctx context.Context, id int) error {
	//validate
	resp, err := c.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//check if present
	if !resp.IsAvailable {
		return ErrAlreadyUnavailable
	}

	err = c.DB.UpdateIsAvailableStatus(ctx, &port.UpdateIsAvailableStatusRequest{
		Id:     id,
		Status: false,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (c *ConfigurableProductService) Update(ctx context.Context, req *UpdateRequest) error {
	//validate
	_, err := c.DB.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	/*
		validate products before validating attribute keys
		because attribute validation relies on the products..
		we fetch the attributes in our products and ensure that the
		attributes exist
	*/
	if req.Name != "" {
		err = c.validateName(ctx, req.Name)
		if err != nil {
			return err
		}
		err := c.DB.UpdateName(ctx, &port.UpdateNameRequest{
			Id:   req.Id,
			Name: req.Name,
		})
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}

	if req.Desc != "" {
		err = validateDesc(req.Desc)
		if err != nil {
			return err
		}
		err := c.DB.UpdateDesc(ctx, &port.UpdateDescRequest{
			Id:   req.Id,
			Desc: req.Desc,
		})
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}

	if req.Product != nil {
		err = c.validateProducts(ctx, req.Product)
		if err != nil {
			return err
		}
		err := c.DB.UpdateProducts(ctx, &port.UpdateProductRequest{
			Id:         req.Id,
			ProductIds: req.Product,
		})
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}

	if req.AttributeKeys != nil {
		err = c.validateAttributekeys(ctx, req.AttributeKeys, req.Product)
		if err != nil {
			return err
		}

		err := c.DB.UpdateAttributes(ctx, &port.UpdateAttributes{
			Id:            req.Id,
			AttributeKeys: req.AttributeKeys,
		})
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}

	if req.Images != nil {
		err = validateImages(req.Images)
		if err != nil {
			return err
		}
		err = c.DB.UpdateImages(ctx, &port.UpdateImagesRequest{
			Id:     req.Id,
			Images: req.Images,
		})
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}

	if req.ExternalId != "" {
		err := c.DB.UpdateExternalId(ctx, &port.UpdateExternalIdRequest{
			Id:         req.Id,
			ExternalId: req.ExternalId,
		})
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}

	return nil
}

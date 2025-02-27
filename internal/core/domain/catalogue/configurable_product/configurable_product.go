package configurable_product

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/catalogue/configurable_product"
)

var (
	ErrNameIsNotSupplied             = errors.New("oopsy, name is mandatory")
	ErrNameDuplicate                 = errors.New("oopsy, name already in use")
	ErrEmptyGetContent               = errors.New("oopsy, empty get content")
	ErrUnknown                       = errors.New("oopsy, error unknown")
	ErrDescIsNotSupplied             = errors.New("oopsy, desc is mandatory")
	ErrImagesMustBeAtleastTwo        = errors.New("oopsy, atleast two images mandatory")
	ErrAttributeKeysMustBeAtleastOne = errors.New("oopsy, atleast one attribute key mandatory")
	ErrProductsMustBeAtleastOne      = errors.New("oopsy, atleast product mandatory")
	ErrAlreadyAvailable              = errors.New("oopsy, product already available")
	ErrAlreadyUnavailable            = errors.New("oopsy, product already unavailable")
	ErrIdNotFound                    = errors.New("oopsy, id not found")
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
	Attributes    map[string]string
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
	db port.DB
}

func NewConfigurableProductService(db port.DB) Provider {
	return &ConfigurableProductService{
		db: db,
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
	err = c.validateAttributekeys(ctx, req.AttributeKeys)
	if err != nil {
		return 0, err
	}

	id, err := c.db.Create(ctx, &port.CreateRequest{
		Name:              req.Name,
		Desc:              req.Desc,
		ExternalId:        req.ExternalId,
		IsAvailableStatus: false,
		Products:          req.Products,
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
	resp, err := c.db.Get(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
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
	}, nil
}

func (c *ConfigurableProductService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := []GetResponse{}
	if req.Name != "" {
		get_by_name_resp, err := c.db.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
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
		get_by_name_resp, err := c.db.GetByExternalId(ctx, req.ExternalId)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
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
	get_by_name_resp, err := c.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
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
	if len(resp) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}

	return GetAllResponse{
		List: resp,
	}, nil
}

func (c *ConfigurableProductService) Avail(ctx context.Context, id int) error {
	//validate
	resp, err := c.db.Get(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//check if present
	if resp.IsAvailable {
		return ErrAlreadyAvailable
	}

	err = c.db.UpdateIsAvailableStatus(ctx, &port.UpdateIsAvailableStatusRequest{
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
	resp, err := c.db.Get(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//check if present
	if !resp.IsAvailable {
		return ErrAlreadyUnavailable
	}

	err = c.db.UpdateIsAvailableStatus(ctx, &port.UpdateIsAvailableStatusRequest{
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
	// resp, err := c.db.Get(ctx, req.)
	// if err != nil {
	// 	switch err {
	// 	case port.ErrSysNoRows:
	// 		return ErrIdNotFound
	// 	default:
	// 		return ErrUnknown
	// 	}
	// }

	return nil
}

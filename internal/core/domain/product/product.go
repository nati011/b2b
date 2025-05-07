package product

import (
	"context"
	"errors"
	"math"

	category "b2b.nati011.github.com/internal/core/domain/category"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/product"
)

var (
	ErrEmptyGetContent                                        = errors.New("oopsy, no product found")
	ErrAlreadyActive                                          = errors.New("oopsy, product already active")
	ErrAlreadyInactive                                        = errors.New("oopsy, product already inactive")
	ErrIdNotFound                                             = errors.New("oopsy, id not found")
	ErrNameNotSupplied                                        = errors.New("oopsy, name is not supplied")
	ErrDescNotSupplied                                        = errors.New("oopsy, description is not supplied")
	ErrNameDuplicate                                          = errors.New("oopsy, name duplicate")
	ErrImagesMustBeAtleastTwo                                 = errors.New("oopsy, images must be atleast two")
	ErrPriceNotSupplied                                       = errors.New("oopsy, price is not supplied")
	ErrAttributeValuesCannotBeEmpty                           = errors.New("oopsy, attribute values cannot be empty")
	ErrUnknown                                                = errors.New("oopsy, unkown error")
	ErrCategoryNotFound                                       = errors.New("oopsy, category not found")
	ErrPriceCannotBeNegative                                  = errors.New("oopsy, price cannot be negative")
	ErrStockUnavailable                                       = errors.New("oopsy, requested quantity greater than stock")
	ErrStockReservationQtyMustBeLessThanOrEqualToAvailableQty = errors.New("oopsy, reserved quantity cannot be more than available quantity")
	ErrFreeReservationQtyMustBeLessThanOrEqualToReservedQty   = errors.New("oopsy, free reservation quantity cannot be more than reserved quantity")
)

type Image struct {
	ImageUrl string
	BlurHash string
}

type CreateRequest struct {
	Name          string
	Desc          string
	ExternalID    string
	Images        []string
	Price         float64
	Attributes    map[string]string
	DistributorId int
	CategoryId    []int
}

type GetResponse struct {
	Id             int
	Name           string
	Desc           string
	ExternalID     string
	Images         []Image
	Price          float64
	Attributes     map[string]string
	DistributorId  int
	CategoryId     []int
	Stock          int
	AvailableStock int
	ReservedStock  int
	IsActive       bool
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name          string
	ExternalID    string
	DistributorId int
	CategoryId    []int
	PriceMin      int
	PriceMax      int
}

type UpdateRequest struct {
	Id         int
	Name       string
	ExternalID string
	Price      float64
	Desc       string
	Images     []string
	CategoryId []int
}

type GoodsReceivingRequest struct {
	Id     int
	Amount int
}

type DispatchRequest struct {
	Id     int
	Amount int
}

type GetProductsWithCategoriesRequest struct {
	List []int
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (id int, err error)
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	Update(ctx context.Context, req *UpdateRequest) (int, error)
	ReceiveGoods(ctx context.Context, req *GoodsReceivingRequest) error
	Dispatch(ctx context.Context, req *DispatchRequest) error
	Activate(ctx context.Context, id int) error
	Deactivate(ctx context.Context, id int) error
	IsActive(ctx context.Context, id int) (bool, error)
	Reserve(ctx context.Context, id int, qty int) error
	FreeReservation(ctx context.Context, id int, qty int) error
}

type ProductService struct {
	DB              port.DB
	CategoryService category.Provider
}

func NewProduct(db port.DB, categoryService category.Provider) Provider {
	return &ProductService{
		DB:              db,
		CategoryService: categoryService,
	}
}

func (p *ProductService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	//validate
	err := p.validateName(ctx, req.Name)
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
	err = create_validatePrice(int(req.Price))
	if err != nil {
		return 0, err
	}
	err = validateAttributes(req.Attributes)
	if err != nil {
		return 0, err
	}

	//validate categories
	for _, i := range req.CategoryId {
		_, err := p.CategoryService.Get(ctx, i)
		if err != nil {
			switch err {
			case category.ErrIdNotFound:
				return 0, ErrCategoryNotFound
			default:
				return 0, ErrUnknown
			}
		}
	}

	// FIXME: optimize, perhaps relegate operation to frontend
	images, err := generateBlurHash(req.Images)
	if err != nil {
		return 0, err
	}

	//create
	id, err := p.DB.Create(ctx, &port.CreateRequest{
		Name:          req.Name,
		Desc:          req.Desc,
		ExternalID:    req.ExternalID,
		Images:        images,
		Price:         req.Price,
		Attributes:    req.Attributes,
		CategoryId:    req.CategoryId,
		DistributorId: req.DistributorId,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	return id, nil
}

func (p *ProductService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := p.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	if resp.Id != id {
		return GetResponse{}, ErrIdNotFound
	}
	var images []Image
	for _, value := range resp.Images {
		image := Image{
			ImageUrl: value.ImageUrl,
			BlurHash: value.BlurHash,
		}
		images = append(images, image)
	}
	response := GetResponse{
		Id:            resp.Id,
		Name:          resp.Name,
		Desc:          resp.Desc,
		ExternalID:    resp.ExternalID,
		Images:        images,
		Price:         resp.Price,
		Attributes:    resp.Attributes,
		DistributorId: resp.DistributorId,
		CategoryId:    resp.CategoryId,
		Stock:         resp.Stock,
		IsActive:      resp.IsActive,
	}
	return response, nil
}

func (p *ProductService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := GetAllResponse{}
	if req.Name != "" {
		resp_getByName, err := p.DB.GetByName(ctx, &port.GetByNameRequest{
			Name: req.Name,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp_getByName.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}

			resp.List = append(resp.List, response)
		}
	}

	if req.ExternalID != "" {
		resp_getByExtId, err := p.DB.GetByExternalId(ctx, &port.GetByExternalIdRequest{
			ExternalId: req.ExternalID,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_getByExtId.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}
			resp.List = append(resp.List, response)
		}
	}

	if req.CategoryId != nil {
		//validate category
		for _, i := range req.CategoryId {
			_, err := p.CategoryService.Get(ctx, i)
			if err != nil {
				switch err {
				case category.ErrIdNotFound:
					return GetAllResponse{}, ErrCategoryNotFound
				default:
					return GetAllResponse{}, ErrUnknown
				}
			}
		}

		resp_getByCategoryId, err := p.DB.GetByCategory(ctx, &port.GetByCategoryRequest{
			CategoryId: req.CategoryId,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_getByCategoryId.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}

			resp.List = append(resp.List, response)
		}
	}

	if req.DistributorId != 0 {
		resp_getByDistId, err := p.DB.GetByDistributorId(ctx, &port.GetByDistributorIdRequest{
			DistributorId: req.DistributorId,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_getByDistId.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}

			resp.List = append(resp.List, response)
		}
	}

	if req.PriceMax != 0 && req.PriceMin != 0 {
		resp_getByName, err := p.DB.GetByPriceRange(ctx, &port.GetByPriceRangeRequest{
			PriceMin: req.PriceMin,
			PriceMax: req.PriceMax,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_getByName.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}
			resp.List = append(resp.List, response)
		}
	} else if req.PriceMax != 0 && req.PriceMin == 0 {
		resp_getByName, err := p.DB.GetByPriceRange(ctx, &port.GetByPriceRangeRequest{
			PriceMin: 0,
			PriceMax: req.PriceMax,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_getByName.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}
			resp.List = append(resp.List, response)
		}
	} else if req.PriceMax == 0 && req.PriceMin != 0 {
		resp_getByName, err := p.DB.GetByPriceRange(ctx, &port.GetByPriceRangeRequest{
			PriceMin: req.PriceMin,
			PriceMax: math.MaxInt,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_getByName.List {
			var images []Image
			for _, value := range i.Images {
				image := Image{
					ImageUrl: value.ImageUrl,
					BlurHash: value.BlurHash,
				}
				images = append(images, image)
			}
			response := GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}
			resp.List = append(resp.List, response)
		}
	}
	if len(resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}

	return resp, nil
}

func (p *ProductService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := p.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}

	resp_val := GetAllResponse{}
	for _, i := range resp.List {
		var images []Image
		for _, value := range i.Images {
			image := Image{
				ImageUrl: value.ImageUrl,
				BlurHash: value.BlurHash,
			}
			images = append(images, image)
		}
		response := GetResponse{
			Id:            i.Id,
			Name:          i.Name,
			Desc:          i.Desc,
			ExternalID:    i.ExternalID,
			Images:        images,
			Price:         i.Price,
			Attributes:    i.Attributes,
			DistributorId: i.DistributorId,
			CategoryId:    i.CategoryId,
			Stock:         i.Stock,
			IsActive:      i.IsActive,
		}
		resp_val.List = append(resp_val.List, response)
	}
	return resp_val, nil
}

func (p *ProductService) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate Id
	_, err := p.Get(ctx, req.Id)
	if err != nil {
		return 0, ErrIdNotFound
	}

	//update
	if req.ExternalID != "" {
		err = p.DB.UpdateExternalID(ctx, &port.UpdateExternalIDRequest{
			Id:         req.Id,
			ExternalId: req.ExternalID,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, ErrUnknown
			}

		}
	}

	if req.Name != "" {
		err = p.validateName(ctx, req.Name)
		if err != nil {
			return 0, err
		}
		err = p.DB.UpdateName(ctx, &port.UpdateNameRequest{
			Id:   req.Id,
			Name: req.Name,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, ErrUnknown
			}
		}
	}

	if req.Price != 0 {
		err = update_validatePrice(int(req.Price))
		if err != nil {
			return 0, err
		}
		err = p.DB.UpdatePrice(ctx, &port.UpdatePriceRequest{
			Id:    req.Id,
			Price: req.Price,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, ErrUnknown
			}

		}
	}

	if req.Desc != "" {
		err = validateDesc(req.Desc)
		if err != nil {
			return 0, err
		}
		err = p.DB.UpdateDesc(ctx, &port.UpdateDescRequest{
			Id:   req.Id,
			Desc: req.Desc,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, ErrUnknown
			}
		}
	}

	if req.Images != nil {
		err = validateImages(req.Images)
		if err != nil {
			return 0, err
		}

		images, err := generateBlurHash(req.Images)
		if err != nil {
			return 0, ErrUnknown
		}

		err = p.DB.UpdateImages(ctx, &port.UpdateImagesRequest{
			Id:     req.Id,
			Images: images,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, ErrUnknown
			}
		}
	}

	if req.CategoryId != nil {
		//validate categories
		for _, i := range req.CategoryId {
			_, err := p.CategoryService.Get(ctx, i)
			if err != nil {
				switch err {
				case category.ErrIdNotFound:
					return 0, ErrCategoryNotFound
				default:
					return 0, ErrUnknown
				}
			}
		}
		err = p.DB.UpdateCategoryId(ctx, &port.UpdateCategoryIdRequest{
			Id:         req.Id,
			CategoryId: req.CategoryId,
		})
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, ErrUnknown
			}
		}
	}
	return req.Id, nil
}

func (p *ProductService) ReceiveGoods(ctx context.Context, req *GoodsReceivingRequest) error {
	//validate Id
	_, err := p.Get(ctx, req.Id)
	if err != nil {
		return ErrIdNotFound
	}

	err = p.DB.GoodsReceiving(ctx, &port.GoodsReceivingRequest{
		Id:     req.Id,
		Amount: req.Amount,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (p *ProductService) Dispatch(ctx context.Context, req *DispatchRequest) error {
	//validate Id
	prod, err := p.Get(ctx, req.Id)
	if err != nil {
		return ErrIdNotFound
	}

	//validate stock quantity
	if prod.Stock < req.Amount {
		return ErrStockUnavailable
	}
	err = p.DB.Dispatch(ctx, &port.DispatchRequest{
		Id:     req.Id,
		Amount: req.Amount,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (p *ProductService) Activate(ctx context.Context, id int) error {
	//validate Id
	resp, err := p.Get(ctx, id)
	if err != nil {
		return ErrIdNotFound
	}

	//check if alreay active
	if resp.IsActive {
		return ErrAlreadyActive
	}

	err = p.DB.UpdateActiveStatus(ctx, &port.UpdateActiveStatusRequest{
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

func (p *ProductService) Deactivate(ctx context.Context, id int) error {
	//validate Id
	resp, err := p.Get(ctx, id)
	if err != nil {
		return ErrIdNotFound
	}

	//check if already inactive
	if !resp.IsActive {
		return ErrAlreadyInactive
	}

	err = p.DB.UpdateActiveStatus(ctx, &port.UpdateActiveStatusRequest{
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

func (p *ProductService) IsActive(ctx context.Context, id int) (bool, error) {
	//validate Id
	resp, err := p.Get(ctx, id)
	if err != nil {
		return false, ErrIdNotFound
	}
	return resp.IsActive, nil
}

func (p *ProductService) Reserve(ctx context.Context, id int, qty int) error {
	//validate Id
	product, err := p.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}
	//validate if qty is less than or equal to available qty
	if product.AvailableStock <= qty {
		return ErrStockReservationQtyMustBeLessThanOrEqualToAvailableQty
	}

	err = p.DB.Reserve(ctx, &port.ReserveRequest{
		Id:     id,
		Amount: qty,
	})
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (p *ProductService) FreeReservation(ctx context.Context, id int, qty int) error {
	//validate Id
	product, err := p.Get(ctx, id)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	//validate if qty is less than or equal to reserved qty
	if qty > product.ReservedStock {
		return ErrFreeReservationQtyMustBeLessThanOrEqualToReservedQty
	}
	err = p.DB.FreeReservation(ctx, &port.FreeReservedRequest{
		Id:     id,
		Amount: qty,
	})
	if err != nil {
		return ErrIdNotFound
	}
	return nil
}

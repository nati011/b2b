package product

import (
	"context"
	"errors"
)

var (
	ErrEmptyGetContent              = errors.New("oopsy, no product found")
	ErrAlreadyActive                = errors.New("oopsy, product already active")
	ErrAlreadyInactive              = errors.New("oopsy, product already inactive")
	ErrIdNotFound                   = errors.New("oopsy, id not found")
	ErrNameNotSupplied              = errors.New("oopsy, name is not supplied")
	ErrDescNotSupplied              = errors.New("oopsy, description is not supplied")
	ErrNameDuplicate                = errors.New("oopsy, name duplicate")
	ErrImagesMustBeAtleastTwo       = errors.New("oopsy, images must be atleast two")
	ErrPriceNotSupplied             = errors.New("oopsy, price is not supplied")
	ErrAttributeValuesCannotBeEmpty = errors.New("oopsy, attribute values cannot be empty")
	ErrDuplicateNameNotAllowed      = errors.New("oopsy, duplicate name not allowed")
	ErrPriceCannotBeZero            = errors.New("oopsy, price cannot be zero")
)

type CreateProductRequest struct {
	Name          string
	Desc          string
	ExternalID    string
	Images        []string
	Price         float64
	Attributes    map[string]string
	DistributorId int
	CategoryId    []int
}

type GetProductResponse struct {
	Id            int
	Name          string
	Desc          string
	ExternalID    string
	Images        []string
	Price         float64
	Attributes    map[string]string
	DistributorId int
	CategoryId    []int
	Stock         int
	IsActive      bool
}

type GetAllResponse struct {
	List []GetProductResponse
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
	Name       string
	ExternalID string
	CategoryId int
	Price      int
}

type GoodsReceivingRequest struct {
	Id     int
	Amount int
}

type DispatchRequest struct {
	Id     int
	Amount int
}

type Provider interface {
	Create(ctx context.Context, req *CreateProductRequest) (id int, err error)
	Get(ctx context.Context, id int) (GetProductResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	Update(ctx context.Context, req *UpdateRequest) (int, error)
	ReceiveGoods(ctx context.Context, req *GoodsReceivingRequest) error
	Dispatch(ctx context.Context, req *DispatchRequest) error
	Activate(ctx context.Context, id int) error
	Deactivate(ctx context.Context, id int) error
}

type ProductService struct {
}

func NewProduct() Provider {
	return &ProductService{}
}

func (p *ProductService) Create(ctx context.Context, req *CreateProductRequest) (int, error) {
	return 0, nil
}

func (p *ProductService) Get(ctx context.Context, id int) (GetProductResponse, error) {
	return GetProductResponse{}, nil
}

func (p *ProductService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

func (p *ProductService) GetAll(ctx context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

func (p *ProductService) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	return 0, nil
}

func (p *ProductService) ReceiveGoods(ctx context.Context, req *GoodsReceivingRequest) error {
	return nil
}

func (p *ProductService) Dispatch(ctx context.Context, req *DispatchRequest) error {
	return nil
}

func (p *ProductService) Activate(ctx context.Context, id int) error {
	return nil
}

func (p *ProductService) Deactivate(ctx context.Context, id int) error {
	return nil
}

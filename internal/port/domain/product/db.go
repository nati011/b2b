package product

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type Image struct {
	ImageUrl string
	BlurHash string
}

type CreateRequest struct {
	Name          string
	Desc          string
	ExternalID    string
	Images        []Image
	Price         float64
	Attributes    map[string]string
	DistributorId int
	CategoryId    []int
}

type GetResponse struct {
	Id            int
	Name          string
	Desc          string
	ExternalID    string
	Images        []Image
	Price         float64
	Attributes    map[string]string
	DistributorId int
	CategoryId    []int
	Stock         int
	IsActive      bool
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByNameRequest struct {
	Name string
}

type GetByExternalIdRequest struct {
	ExternalId string
}

type GetByDistributorIdRequest struct {
	DistributorId int
}

type GetByCategoryRequest struct {
	CategoryId []int
}

type GetByPriceRangeRequest struct {
	PriceMin int
	PriceMax int
}

type UpdateNameRequest struct {
	Id   int
	Name string
}

type UpdateExternalIDRequest struct {
	Id         int
	ExternalId string
}

type UpdatePriceRequest struct {
	Id    int
	Price float64
}

type UpdateDescRequest struct {
	Id   int
	Desc string
}

type UpdateImagesRequest struct {
	Id     int
	Images []Image
}

type UpdateActiveStatusRequest struct {
	Id     int
	Status bool
}

type UpdateCategoryIdRequest struct {
	Id         int
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

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByName(ctx context.Context, req *GetByNameRequest) (GetAllResponse, error)
	GetByExternalId(ctx context.Context, req *GetByExternalIdRequest) (GetAllResponse, error)
	GetByDistributorId(ctx context.Context, req *GetByDistributorIdRequest) (GetAllResponse, error)
	GetByCategory(ctx context.Context, req *GetByCategoryRequest) (GetAllResponse, error)
	GetByPriceRange(ctx context.Context, req *GetByPriceRangeRequest) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	UpdateName(ctx context.Context, req *UpdateNameRequest) error
	UpdateExternalID(ctx context.Context, req *UpdateExternalIDRequest) error
	UpdatePrice(ctx context.Context, req *UpdatePriceRequest) error
	UpdateDesc(ctx context.Context, req *UpdateDescRequest) error
	UpdateImages(ctx context.Context, req *UpdateImagesRequest) error
	UpdateActiveStatus(ctx context.Context, req *UpdateActiveStatusRequest) error
	UpdateCategoryId(ctx context.Context, req *UpdateCategoryIdRequest) error
	GoodsReceiving(ctx context.Context, req *GoodsReceivingRequest) error
	Dispatch(ctx context.Context, req *DispatchRequest) error
}

type DB interface {
	Reader
	Writer
}

package configurableProduct

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type CreateRequest struct {
	Name              string
	Desc              string
	ExternalId        string
	AttributeKeys     []string
	Products          []int
	Images            []Image
	IsAvailableStatus bool
}

type PriceRangeResponse struct {
	Min int
	Max int
}

type Image struct {
	ImageUrl string
	BlurHash string
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
	Images        []Image
}

type GetAllResponse struct {
	List []GetResponse
}

type UpdateNameRequest struct {
	Id   int
	Name string
}

type UpdateDescRequest struct {
	Id   int
	Desc string
}

type UpdateExternalIdRequest struct {
	Id         int
	ExternalId string
}

type UpdateProductRequest struct {
	Id         int
	ProductIds []int
}

type UpdateIsAvailableStatusRequest struct {
	Id     int
	Status bool
}

type UpdateImagesRequest struct {
	Id     int
	Images []Image
}

type UpdateAttributes struct {
	Id            int
	AttributeKeys []string
}

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByName(ctx context.Context, name string) (GetAllResponse, error)
	GetByExternalId(ctx context.Context, extId string) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	UpdateName(ctx context.Context, req *UpdateNameRequest) error
	UpdateDesc(ctx context.Context, req *UpdateDescRequest) error
	UpdateExternalId(ctx context.Context, req *UpdateExternalIdRequest) error
	UpdateIsAvailableStatus(ctx context.Context, req *UpdateIsAvailableStatusRequest) error
	UpdateProducts(ctx context.Context, req *UpdateProductRequest) error
	UpdateImages(ctx context.Context, req *UpdateImagesRequest) error
	UpdateAttributes(ctx context.Context, req *UpdateAttributes) error
}

type DB interface {
	Reader
	Writer
}

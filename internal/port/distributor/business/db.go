package distributor

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type CreateBusinessInformation struct {
	Name          string `json:"name"`
	Tin           string `json:"tin"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	GeneralZone   string `json:"general_zone"`
	Region        string `json:"region"`
	Woreda        string `json:"woreda"`
	DistributorId int    `json:"distributorId"`
}

type CreateBusinessResponse struct {
	BusinessId int
}
type GetResponse struct {
	Id            int
	Name          string
	Tin           int
	DistributorId int
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetById(ctx context.Context, req int) (GetResponse, error)
	GetByDistributorId(ctx context.Context, distributorId int) (GetResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateBusinessInformation) (CreateBusinessResponse, error)
	Update(ctx context.Context, req *CreateBusinessInformation) (CreateBusinessResponse, error)
}

type DB interface {
	Reader
	Writer
}

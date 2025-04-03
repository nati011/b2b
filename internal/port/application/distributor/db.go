package distributor

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type UpdateRequest struct {
	Id         int
	FirstName  string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}
type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Woreda      string
	UserId      int
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByIdRequest struct {
	Id int
}

type CreateRequest struct {
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Woreda      string
	UserId      int
}

type BusinessLocation struct {
	GeneralZone string
	Region      string
	Woreda      string
}

type GetByParamRequest struct {
	Id   int
	Name string
	Tin  string
}

type CreateBusinessInformation struct {
	Name          string
	Tin           int
	Latitude      string
	Longitude     string
	GeneralZone   string
	Region        string
	Woreda        string
	DistributorId int
}

type CreateBusinessResponse struct {
	BusinessId int
}
type GetBusinessResponse struct {
	Id            int
	Name          string
	Tin           string
	DistributorId int
}

type GetAllBusinessResponse struct {
	List []GetResponse
}

type Reader interface {
	GetById(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetBusinessAll(ctx context.Context) (GetAllResponse, error)
	GetBusinessById(ctx context.Context, req int) (GetBusinessResponse, error)
	GetBusinessByDistributorId(ctx context.Context, distributorId int) (GetBusinessResponse, error)
	GetByParam(ctx context.Context, req GetByParamRequest) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Update(ctx context.Context, req *UpdateRequest) (int, error)
	CreateBusiness(ctx context.Context, req *CreateBusinessInformation) (CreateBusinessResponse, error)
}

type DB interface {
	Reader
	Writer
}

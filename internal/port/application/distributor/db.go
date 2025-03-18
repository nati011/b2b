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

type CreateRequest struct {
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
	Password   string
}
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
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByIdRequest struct {
	Id int
}

type RegisterDistributorRequest struct {
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	DOB             time.Time `json:"date_of_birth"`
	PhoneNumber     string    `json:"phone_number"`
	ConfirmPassword string    `json:"confirmed_password"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Username        string    `json:"username"`
	ExternalId      string    `json:"external_id"`
}

type BusinessLocation struct {
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
}

type UpdateBusinessRequest struct {
	Id            int              `json:"id"`
	DistributorId int              `json:"distributorId"`
	Name          string           `json:"name"`
	Tin           int              `json:"tin"`
	Location      BusinessLocation `json:"location"`
}

type RegisterDistributorResponse struct {
	Id      int    `json:"distributorId"`
	Message string `json:"message"`
}

type GetByParamRequest struct {
	Id    int
	Name  string
	Email string
}

type CreateBusinessInformation struct {
	Name          string `json:"name"`
	Tin           int    `json:"tin"`
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
type GetBusinessResponse struct {
	Id            int
	Name          string
	Tin           int
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
	GetByDistributorId(ctx context.Context, distributorId int) (GetBusinessResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	CreateBusiness(ctx context.Context, req *CreateBusinessInformation) (CreateBusinessResponse, error)
	UpdateBusiness(ctx context.Context, req *UpdateBusinessRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}

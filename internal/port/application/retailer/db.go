package retailer

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

type RegisterRetailerRequest struct {
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	DOB             time.Time `json:"date_of_birth"`
	PhoneNumber     string    `json:"phone_number"`
	ConfirmPassword string    `json:"confirmed_password"`
	FirstName       string    `json:"full_name"`
	Username        string    `json:"username"`
	ExternalId      string    `json:"external_id"`
}

type BusinessLocation struct {
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
}

type UpdateBusinessRequest struct {
	Id         int              `json:"id"`
	RetailerId int              `json:"retailerId"`
	Name       string           `json:"name"`
	Tin        int              `json:"tin"`
	Location   BusinessLocation `json:"location"`
}

type RegisterRetailerResponse struct {
	Id      int    `json:"retailerId"`
	Message string `json:"message"`
}

type GetByParamRequest struct {
	Id    int
	Name  string
	Email string
}

type CreateBusinessInformation struct {
	Name        string `json:"name"`
	Tin         int    `json:"tin"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
	RetailerId  int    `json:"retailerId"`
}

type CreateBusinessResponse struct {
	BusinessId int
}
type GetBusinessResponse struct {
	Id         int
	Name       string
	Tin        int
	RetailerId int
}

type GetAllBusinessResponse struct {
	List []GetResponse
}

type Reader interface {
	GetById(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetBusinessAll(ctx context.Context) (GetAllResponse, error)
	GetBusinessById(ctx context.Context, req int) (GetBusinessResponse, error)
	GetByRetailerId(ctx context.Context, retailerId int) (GetBusinessResponse, error)
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

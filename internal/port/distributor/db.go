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
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
	Password   string
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
	DistributorId int              `json:"distributorId"`
	Name          string           `json:"name"`
	Tin           int              `json:"tin"`
	Region        BusinessLocation `json:"region"`
}

type RegisterDistributorResponse struct {
	DistributorId int    `json:"distributorId"`
	Message       string `json:"message"`
}

type GetByParamRequest struct {
	Id    int
	Name  string
	Email string
}

type Reader interface {
	GetById(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}

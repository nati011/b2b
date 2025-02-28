package distributor

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type MutationRequest struct {
	Name          string
	Tin           int
	DistributorId int
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

type GetByIdRequest struct {
	Id int
}

type GetByDistributorId struct {
	Distributor int
}

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetById(ctx context.Context, req *GetByIdRequest) (GetResponse, error)
	GetByDistributorId(ctx context.Context, req *GetByIdRequest) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *MutationRequest) (int, error)
	Update(ctx context.Context, req *MutationRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}

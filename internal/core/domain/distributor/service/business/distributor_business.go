package distributor

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/distributor/business"
)

var (
	ErrEmptyGetContent = errors.New("oopsy, no distributor found")
	ErrUnknown         = errors.New("oopsy, unkown error")
	SUCCESS_MESSAGE    = "Ahoy!"
)

type Provider interface {
	Create(ctx context.Context, req *port.CreateBusinessInformation) (response port.CreateBusinessResponse, err error)
	GetAll(ctx context.Context) (port.GetAllResponse, error)
	GetByDistributor(ctx context.Context, distributorId int) (port.GetResponse, error)
	GetById(ctx context.Context, id int) (port.GetResponse, error)
}

type DistributorService struct {
	db port.DB
}

func (d *DistributorService) Create(ctx context.Context, req *port.CreateBusinessInformation) (resp port.CreateBusinessResponse, err error) {
	panic("unimplemented")
}

func (d *DistributorService) GetAll(ctx context.Context) (resp port.GetAllResponse, err error) {
	panic("unimplemented")
}

func (d *DistributorService) GetByDistributor(ctx context.Context, distributorId int) (resp port.GetResponse, err error) {
	panic("unimplemented")
}
func (d *DistributorService) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	panic("unimplemented")
}
func NewDistributorBusinessService(db port.DB) Provider {
	return &DistributorService{
		db: db,
	}
}

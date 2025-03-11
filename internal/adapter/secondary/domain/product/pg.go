package product

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/domain/product"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByName(ctx context.Context, req *port.GetByNameRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, req *port.GetByExternalIdRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByDistributorId(ctx context.Context, req *port.GetByDistributorIdRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByCategory(ctx context.Context, req *port.GetByCategoryRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByPriceRange(ctx context.Context, req *port.GetByPriceRangeRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	return nil
}

func (p *Postgres) UpdateExternalID(ctx context.Context, req *port.UpdateExternalIDRequest) error {
	return nil
}

func (p *Postgres) UpdatePrice(ctx context.Context, req *port.UpdatePriceRequest) error {
	return nil
}

func (p *Postgres) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) error {
	return nil
}

func (p *Postgres) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	return nil
}

func (p *Postgres) UpdateActiveStatus(ctx context.Context, req *port.UpdateActiveStatusRequest) error {
	return nil
}

func (p *Postgres) UpdateCategoryId(ctx context.Context, req *port.UpdateCategoryIdRequest) error {
	return nil
}

func (p *Postgres) GoodsReceiving(ctx context.Context, req *port.GoodsReceivingRequest) error {
	return nil
}

func (p *Postgres) Dispatch(ctx context.Context, req *port.DispatchRequest) error {
	return nil
}

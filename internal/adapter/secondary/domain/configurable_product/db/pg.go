package configurable_product

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/domain/configurable_product"
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

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	return nil
}

func (p *Postgres) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) error {
	return nil
}

func (p *Postgres) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	return nil
}

func (p *Postgres) UpdateIsAvailableStatus(ctx context.Context, req *port.UpdateIsAvailableStatusRequest) error {
	return nil
}

func (p *Postgres) UpdateProducts(ctx context.Context, req *port.UpdateProductRequest) error {
	return nil
}

func (p *Postgres) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	return nil
}

func (p *Postgres) UpdateAttributes(ctx context.Context, req *port.UpdateAttributes) error {
	return nil
}

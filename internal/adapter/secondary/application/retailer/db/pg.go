package db

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/application/retailer"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		db: db,
	}
}

func (m *Postgres) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	return 0, nil
}

func (m *Postgres) CreateUserAgent(ctx context.Context, id int) error {
	return nil
}

func (m *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	return nil
}

func (m *Postgres) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	return nil
}

func (m *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_retailer_by_id($1);"
	err := m.db.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Name, &response.Tin, &response.Latitude, &response.Longitude, &response.GeneralZone, &response.Region, &response.Woreda)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	return response, nil
}

func (m *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (m *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (m *Postgres) GetByTin(ctx context.Context, tin string) (port.GetResponse, error) {
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Postgres) GetAllUserAgents(ctx context.Context, id int) (port.GetAllUserResponse, error) {
	return port.GetAllUserResponse{}, nil
}

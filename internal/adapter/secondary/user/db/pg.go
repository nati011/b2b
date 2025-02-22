package db

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/user"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetByEmail(ctx context.Context, email string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByPhone(ctx context.Context, phone string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByUsername(ctx context.Context, username string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByActiveStatus(ctx context.Context, status bool) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByIsActiveStatus(ctx context.Context, status bool) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetAll(context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(context.Context, *port.CreateRequest) (int, error) {
	return 0, nil
}

func (m *Postgres) CreateAndActivate(context.Context, *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}

func (p *Postgres) UpdateFullName(context.Context, *port.UpdateFullNameRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateEmail(context.Context, *port.UpdateEmailRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateDOB(context.Context, *port.UpdateDOBRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateIsActiveStatus(context.Context, *port.UpdateIsActiveRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdatePhone(context.Context, *port.UpdatePhoneRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateUsername(context.Context, *port.UpdateUsernameRequest) (int, error) {
	return 0, nil
}

func (m *Postgres) AssignRole(ctx context.Context, id int, roleId int) error {
	return nil
}

func (m *Postgres) RemoveAssignedRole(ctx context.Context, id int, roleId int) error {
	return nil
}

func (m *Postgres) GetAllAssignedRole(ctx context.Context, id int) (port.GetAllAssignedRoleResponse, error) {
	return port.GetAllAssignedRoleResponse{}, nil
}

func (m *Postgres) HasAccessToResource(ctx context.Context, id int, resourceId int) (bool, error) {
	return false, nil
}

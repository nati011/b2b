package email

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/application/email-template"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{Pool: db}
}

func (p *PostgresReaderWriter) Create(ctx context.Context, req *CreateRequest) (CreateResponse, error) {
	query := "SELECT * FROM templates AS t WHERE t.name = $1;"
	var name string
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Name,
	)
	if err != nil {
		log.Fatalf("unable to execute search query: %q", err)
		return CreateResponse{}, err
	}

	rows.Row.Scan(&name)
	log.Println("name=", name)
	return CreateResponse{Name: name}, nil
}

func (p *Postgres) Update(ctx context.Context, req *port.UpdateRequest) error {
	return nil
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

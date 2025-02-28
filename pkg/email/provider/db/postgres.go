package email

import (
	"context"
	"database/sql"
	"log"
)

type PostgresReaderWriter struct {
	Pool *sql.DB
}

func NewPostgresReaderWriter(pool *sql.DB) ReaderWriter {
	return &PostgresReaderWriter{Pool: pool}
}

func (p *PostgresReaderWriter) Create(ctx context.Context, req *CreateRequest) (CreateResponse, error) {
	query := "SELECT * FROM templates AS t WHERE t.name = $1;"
	var name string
	err := p.Pool.QueryRowContext(ctx, query, req.Name).Scan(&name)
	if err != nil {
		log.Fatalf("unable to execute search query: %q", err)
		return CreateResponse{}, err
	}
	log.Println("name=", name)
	return CreateResponse{Name: name}, nil
}

func (p *PostgresReaderWriter) Update(ctx context.Context, req *UpdateRequest) (UpdateResponse, error) {
	return UpdateResponse{}, nil
}

func (p *PostgresReaderWriter) Get(ctx context.Context, r string) (GetResponse, error) {
	return GetResponse{}, nil
}

func (p *PostgresReaderWriter) GetAll(ctx context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

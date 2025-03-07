package db

import (
	"context"

<<<<<<< HEAD
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
=======
	"database/sql"
>>>>>>> dec51710 (+ resplve sql.db issue)

	port "b2b.nati011.github.com/internal/port/distributor"
)

type Postgres struct {
<<<<<<< HEAD
	db *pgxpool.Pool
=======
	db *sql.DB
>>>>>>> dec51710 (+ resplve sql.db issue)
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	print("Here_______________________________")
	var resourceId int
<<<<<<< HEAD
	query := "SELECT * FROM public.create_distributor_user($1, $2, $3, $4, $5, $6);"
	err := p.db.QueryRow(ctx, query,
=======
	query := "SELECT * FROM public.create_user($1, $2, $3, $4, $5, $6);"

	err := p.db.QueryRowContext(ctx, query,
>>>>>>> dec51710 (+ resplve sql.db issue)
		req.FullName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
		req.ExternalId,
	).Scan(&resourceId)
	print("HERE2________________")
	print(err, query)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_distributors();"

	err := p.db.QueryRowContext(ctx, query).Scan()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func (p *Postgres) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_distributor_by_id($1);"

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response.Id)
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

<<<<<<< HEAD
func NewPostgres(db *pgxpool.Pool) port.DB {
=======
func NewPostgres(db *sql.DB) port.DB {
>>>>>>> dec51710 (+ resplve sql.db issue)
	return &Postgres{
		db: db,
	}
}

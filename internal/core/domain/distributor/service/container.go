package distributor

import (
<<<<<<< HEAD
	db "b2b.nati011.github.com/internal/adapter/secondary/distributor/db"
	"github.com/jackc/pgx/v5/pgxpool"
=======
	"database/sql"

	adapter "b2b.nati011.github.com/internal/adapter/secondary/distributor/db"
>>>>>>> dec51710 (+ resplve sql.db issue)
)

type Container struct {
	//exportable
	DistributorProvider Provider
<<<<<<< HEAD
	dbPool              *pgxpool.Pool
}

func NewContainer(*pgxpool.Pool) *Container {
=======
}

func NewContainer(db *sql.DB) *Container {
>>>>>>> dec51710 (+ resplve sql.db issue)
	c := new(Container)
	c.initDistributorProvider(db)
	return c
}

func (c *Container) initDistributorProvider(db *sql.DB) {
	postgres := adapter.NewPostgres(
		db,
	)
	c.DistributorProvider = NewDistributorService(postgres)
}

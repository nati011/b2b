package distributor

import (
	db "b2b.nati011.github.com/internal/adapter/secondary/distributor/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	//exportable
	DistributorProvider Provider
	dbPool              *pgxpool.Pool
}

func NewContainer(*pgxpool.Pool) *Container {
	c := new(Container)
	c.initDistributorProvider()
	return c
}

func (c *Container) initDistributorProvider() {
	postgres := db.NewPostgres(
		c.dbPool,
	)
	c.DistributorProvider = NewDistributorService(postgres)
}

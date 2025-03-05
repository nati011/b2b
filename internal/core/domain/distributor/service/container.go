package distributor

import (
	db "b2b.nati011.github.com/internal/adapter/secondary/distributor/db"
	"github.com/jackc/pgx/v5"
)

type Container struct {
	//exportable
	DistributorProvider Provider
	dbPool              *pgx.Conn
}

func NewContainer(*pgx.Conn) *Container {
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

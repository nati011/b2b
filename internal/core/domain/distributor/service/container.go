package distributor

import (
	"database/sql"

	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
)

type Container struct {
	//exportable
	DistributorProvider Provider
}

func NewContainer(db *sql.DB) *Container {
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

package test_container

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/retailer/db"
	user "b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

type TestContainer struct {
	UserService     user.Provider
	RetailerService retailer.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.UserService = user.NewIntegrationTestContainer(db).UserService
	container.RetailerService = retailer.NewRetailerService(container.UserService, db_adapter.NewPostgres(db, &config.Pagination{
		Limit:  10,
		Offset: 0,
	}))

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.UserService = user.NewIntegrationTestContainer(db).UserService
	t.RetailerService = retailer.NewRetailerService(t.UserService, db_adapter.NewPostgres(db, &config.Pagination{
		Limit:  10,
		Offset: 0,
	}))
}

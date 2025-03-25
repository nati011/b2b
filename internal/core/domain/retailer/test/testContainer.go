package test_container

import (
	"database/sql"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/retailer/db"
	user "b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

type TestContainer struct {
	UserService     user.Provider
	RetailerService retailer.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.UserService = user.NewTestContainer().UserService
	container.RetailerService = retailer.NewRetailerService(container.UserService, db_adapter.NewPostgres(db))

	return container
}

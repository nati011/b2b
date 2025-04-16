package test

import (
	"database/sql"

	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	"b2b.nati011.github.com/internal/core/application/resource"
)

type TestContainer struct {
	ResourceService resource.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	c := TestContainer{}
	c.ResourceService = resource.NewResource(adapter.NewPostgres(db))
	return c
}

func (t *TestContainer) TeardownIntegrationTestContainer(db *sql.DB) {
}

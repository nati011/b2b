package core

import (
	"database/sql"

	resource_db_port "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	"b2b.nati011.github.com/internal/core/application/resource"
)

type Container struct {
	ResourceService resource.Provider
}

func NewContainer(db *sql.DB) *Container {
	container := Container{}
	container.InitResourceService(db)

	return &container
}

func (m *Container) InitResourceService(db *sql.DB) {
	m.ResourceService = resource.NewResource(resource_db_port.NewPostgres(db))
}

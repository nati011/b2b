package test

import (
	"database/sql"

	port "b2b.nati011.github.com/internal/adapter/secondary/domain/config"
	"b2b.nati011.github.com/internal/core/domain/config"
)

type IntegrationTestContainer struct {
	ConfigService config.Provider
}

func NewIntegrationTestContainer(DB *sql.DB) IntegrationTestContainer {
	return IntegrationTestContainer{
		ConfigService: config.NewConfig(port.NewPostgres(DB)),
	}
}

func (t *IntegrationTestContainer) Teardown() {
	t.ConfigService = config.NewConfig(port.NewMock())
}

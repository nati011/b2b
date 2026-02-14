package template

import (
	"database/sql"

	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"
	"b2b.nati011.github.com/internal/core/application/template"
)

type TestContainer struct {
	TemplateService template.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	c := TestContainer{}
	c.TemplateService = template.NewTemplateService(adapter.NewPostgres(db))
	return c
}

func (t *TestContainer) Teardown() {
}

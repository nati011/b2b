package template

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"
	"b2b.nati011.github.com/internal/core/application/template"
)

type TestContainer struct {
	TemplateService template.Provider
}

func NewIntegrationTestContainer() TestContainer {
	c := TestContainer{}
	c.TemplateService = template.NewTemplateService(adapter.NewMock())
	return c
}

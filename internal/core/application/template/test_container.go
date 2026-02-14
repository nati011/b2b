package template

import adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"

type TestContainer struct {
	TemplateService Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}
	c.TemplateService = NewTemplateService(adapter.NewMock())
	return c
}

func (t *TestContainer) Teardown() {
	t.TemplateService = NewTemplateService(adapter.NewMock())
}

package resource

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
)

type TestContainer struct {
	ResourceService Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}
	c.ResourceService = NewResource(adapter.NewMock())
	return c
}

func (t *TestContainer) Teardown() {
	t.ResourceService = NewResource(adapter.NewMock())
}

package config

import (
	port "b2b.nati011.github.com/internal/adapter/secondary/domain/config"
)

type TestContainer struct {
	ConfigService Provider
}

func NewTestContainer() TestContainer {
	return TestContainer{
		ConfigService: NewConfig(port.NewMock()),
	}
}

func (t *TestContainer) Teardown() {
	t.ConfigService = NewConfig(port.NewMock())
}

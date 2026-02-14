package payment

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
)

type TestContainer struct {
	Service Provider
}

func NewTestContainer() TestContainer {
	return TestContainer{
		Service: NewPaymentService(adapter.NewMock()),
	}
}

func (t *TestContainer) Teardown() {
	t.Service = NewPaymentService(adapter.NewMock())
}

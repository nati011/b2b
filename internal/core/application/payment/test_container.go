package payment

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
)

type TestContainer struct {
	service Provider
}

func NewTestContainer() TestContainer {
	return TestContainer{
		service: NewPaymentService(adapter.NewMock()),
	}
}

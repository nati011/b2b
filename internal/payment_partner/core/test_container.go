package payment_partner

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment_partner/db"
)

type TestContainer struct {
	PartnerService Provider
}

func NewIntegrationTestContainer() TestContainer {
	c := TestContainer{}
	c.PartnerService = NewPartner(adapter.NewMock())
	return c
}

func (t *TestContainer) Teardown() {
	t.PartnerService = NewPartner(adapter.NewMock())
}

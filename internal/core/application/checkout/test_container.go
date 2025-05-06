package checkout

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
)

type TestContainer struct {
	CheckoutService Provider
	PartnerService  partner.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.CheckoutService = NewCheckoutService(
		adapter.NewMock(),
		container.PartnerService,
		transaction.NewPackageIntegrationTestContainer().TransactionService,
		"https://example.com",
		"https://example.com",
	)

	return container
}

func (t *TestContainer) TearDown() {

	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.CheckoutService = NewCheckoutService(
		adapter.NewMock(),
		t.PartnerService,
		transaction.NewPackageIntegrationTestContainer().TransactionService,
		"https://example.com",
		"https://example.com",
	)
}

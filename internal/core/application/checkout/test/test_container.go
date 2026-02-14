package test

import (
	"database/sql"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
)

type TestContainer struct {
	CheckoutService    checkout.Provider
	TransactionService transaction.Provider
	PartnerService     partner.Provider
}

// (DB port.DB, partner partner.Provider, order order.Provider, transaction transaction.Provider, frontendUrl string, baseUrl string)

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.CheckoutService = checkout.NewCheckoutService(
		payment.NewTestContainer().Service,
		container.PartnerService,
		container.TransactionService,
		"https://example.com",
		"https://example.com",
	)

	return container
}

func (t *TestContainer) TearDown() {
	t.CheckoutService = checkout.NewCheckoutService(
		payment.NewTestContainer().Service,
		t.PartnerService,
		t.TransactionService,
		"https://example.com",
		"https://example.com",
	)
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.CheckoutService = checkout.NewCheckoutService(
		payment.NewTestContainer().Service,
		container.PartnerService,
		container.TransactionService,
		"https://example.com",
		"https://example.com",
	)
	return container
}

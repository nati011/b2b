package test

import (
	port "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	CheckoutService    checkout.Provider
	TransactionService transaction.Provider
	PartnerService     partner.Provider
	OrderService       order.Provider
	ProductService     product.Provider
}

// (DB port.DB, partner partner.Provider, order order.Provider, transaction transaction.Provider, frontendUrl string, baseUrl string)

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	orderContainer := order.NewPackageIntegrationTestContainer()
	container.OrderService = orderContainer.OrderService
	container.ProductService = orderContainer.ProductService

	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.OrderService = orderContainer.OrderService
	container.ProductService = orderContainer.ProductService
	container.CheckoutService = checkout.NewCheckoutService(
		port.NewMock(),
		container.PartnerService,
		container.TransactionService,
		"https://example.com",
		"https://example.com",
	)

	return container
}

func (t *TestContainer) TearDown() {
	orderContainer := order.NewPackageIntegrationTestContainer()
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	t.OrderService = orderContainer.OrderService
	t.ProductService = orderContainer.ProductService
	t.CheckoutService = checkout.NewCheckoutService(
		port.NewMock(),
		t.PartnerService,
		t.TransactionService,
		"https://example.com",
		"https://example.com",
	)
}

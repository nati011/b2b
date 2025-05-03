package checkout

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
	product "b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	CheckoutService Provider
	PartnerService  partner.Provider
	ProductService  product.Provider
	OrderService    order.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	orderContainer := order.NewPackageIntegrationTestContainer()
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.OrderService = orderContainer.OrderService
	container.ProductService = orderContainer.ProductService
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
	t.OrderService = order.NewPackageIntegrationTestContainer().OrderService
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	orderContainer := order.NewPackageIntegrationTestContainer()
	t.OrderService = orderContainer.OrderService
	t.ProductService = orderContainer.ProductService
	t.CheckoutService = NewCheckoutService(
		adapter.NewMock(),
		t.PartnerService,
		transaction.NewPackageIntegrationTestContainer().TransactionService,
		"https://example.com",
		"https://example.com",
	)
}

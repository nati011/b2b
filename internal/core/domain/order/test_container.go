package order

import (
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	OrderService   Provider
	InvoiceService invoice.Provider
	ProductService product.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.InvoiceService = invoice.NewInvoice(
		invoice_db.NewMock(),
	)
	container.ProductService = product.NewPackageIntegrationTestContainer().ProductService
	container.OrderService = NewOrderService(
		order_db.NewMock(),
		container.InvoiceService,
		container.ProductService,
	)

	return container
}

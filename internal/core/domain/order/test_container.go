package order

import (
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/invoice"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/order"
	"b2b.nati011.github.com/internal/core/domain/invoice"
)

type TestContainer struct {
	OrderService   Provider
	InvoiceService invoice.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.InvoiceService = invoice.NewInvoice(
		invoice_db.NewMock(),
	)
	container.OrderService = NewOrderService(
		order_db.NewMock(),
		container.InvoiceService,
	)

	return container
}

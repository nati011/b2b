package order

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	product_test "b2b.nati011.github.com/internal/core/domain/product/test"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	retailer_test "b2b.nati011.github.com/internal/core/domain/retailer/test"
)

type TestContainer struct {
	OrderService    order.Provider
	InvoiceService  invoice.Provider
	ProductService  product.Provider
	RetailerService retailer.Provider
	CheckoutService checkout.Provider
	PartnerService  partner.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.InvoiceService = invoice.NewInvoice(
		invoice_db.NewPostgres(db, config.NewPaginationBuilder().Build()),
	)

	checkout_container := checkout.NewPackageIntegrationTestContainer()

	container.PartnerService = checkout_container.PartnerService
	container.RetailerService = retailer_test.NewDBIntegrationTestContainer(db).RetailerService
	container.CheckoutService = checkout_container.CheckoutService
	container.ProductService = product_test.NewDBIntegrationTestContainer(db).ProductService
	container.OrderService = order.NewOrderService(
		order_db.NewPostgres(db, &config.Pagination{
			Limit:  10,
			Offset: 0,
		}),
		container.InvoiceService,
		container.ProductService,
		container.RetailerService,
		container.CheckoutService,
	)

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.InvoiceService = invoice.NewInvoice(
		invoice_db.NewPostgres(db, config.NewPaginationBuilder().Build()),
	)
	t.RetailerService = retailer.NewPackageIntegrationTestContainer().RetailerService
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.ProductService = product_test.NewDBIntegrationTestContainer(db).ProductService
	t.OrderService = order.NewOrderService(
		order_db.NewPostgres(db, &config.Pagination{
			Limit:  10,
			Offset: 0,
		}),
		t.InvoiceService,
		t.ProductService,
		t.RetailerService,
		t.CheckoutService,
	)
}

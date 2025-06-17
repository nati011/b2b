package order

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributor_test "b2b.nati011.github.com/internal/core/domain/distributor/test"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	product_test "b2b.nati011.github.com/internal/core/domain/product/test"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	retailer_test "b2b.nati011.github.com/internal/core/domain/retailer/test"
)

type TestContainer struct {
	OrderService       order.Provider
	RenderService      render.Provider
	InvoiceService     invoice.Provider
	ProductService     product.Provider
	RetailerService    retailer.Provider
	CheckoutService    checkout.Provider
	PartnerService     partner.Provider
	DistributorService distributor.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.RenderService = render.NewMock()
	container.InvoiceService = invoice.NewInvoice(
		invoice_db.NewPostgres(db, config.NewPaginationBuilder().Build()),
		container.RenderService,
	)

	checkout_container := checkout.NewPackageIntegrationTestContainer()
	container.DistributorService = distributor_test.NewDBIntegrationTestContainer(db).DistributorService
	container.PartnerService = checkout_container.PartnerService
	container.RetailerService = retailer_test.NewDBIntegrationTestContainer(db).RetailerService
	container.CheckoutService = checkout_container.CheckoutService
	container.ProductService = product.NewProduct(
		product_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		category.NewCategory(category_db.NewMock()),
		container.DistributorService,
	)
	container.OrderService = order.NewOrderService(
		order_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		container.InvoiceService,
		container.ProductService,
		container.RetailerService,
		container.CheckoutService,
	)

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.RenderService = render.NewMock()
	t.InvoiceService = invoice.NewInvoice(
		invoice_db.NewPostgres(db, config.NewPaginationBuilder().Build()),
		t.RenderService,
	)
	t.RetailerService = retailer.NewPackageIntegrationTestContainer().RetailerService
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.ProductService = product_test.NewDBIntegrationTestContainer(db).ProductService
	t.OrderService = order.NewOrderService(
		order_db.NewPostgres(db, config.NewPaginationBuilder().Build()),
		t.InvoiceService,
		t.ProductService,
		t.RetailerService,
		t.CheckoutService,
	)
}

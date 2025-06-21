package order_expiry

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	payment_db "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	config_db "b2b.nati011.github.com/internal/adapter/secondary/domain/config"
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	transaction "b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/category"
	config_module "b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributor_test "b2b.nati011.github.com/internal/core/domain/distributor/test"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
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
	ConfigService      config_module.Provider
	PaymentService     payment.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
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
	container.PaymentService = payment.NewPaymentService(payment_db.NewPostgres(db))
	container.CheckoutService = checkout.NewCheckoutService(
		container.PaymentService,
		container.PartnerService,
		transaction.NewPackageIntegrationTestContainer().TransactionService,
		"https://example.com",
		"https://example.com",
	)
	container.ProductService = product.NewProduct(
		product_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		category.NewCategory(category_db.NewMock()),
		container.DistributorService,
	)
	container.ConfigService = config_module.NewConfig(config_db.NewPostgres(db))
	container.OrderService = order.NewOrderService(
		order_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		container.InvoiceService,
		container.ProductService,
		container.RetailerService,
		container.CheckoutService,
		container.ConfigService,
	)

	return container
}

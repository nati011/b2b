package order

import (
	"database/sql"

	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	config_db "b2b.nati011.github.com/internal/adapter/secondary/domain/config"
	distributor_db "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/category"
	config_module "b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
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
	CategoryService    category.Provider
	UserService        user.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.RenderService = render.NewMock()
	container.InvoiceService = invoice.NewInvoice(
		invoice_db.NewMock(),
		container.RenderService,
	)
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.UserService = user.NewTestContainer().UserService
	container.DistributorService = distributor.NewDistributorService(
		container.UserService,
		distributor_db.NewMock(),
		distributorApproval.NewTestContainer().DistributorApprovalService)
	container.ProductService = product.NewProduct(
		product_db.NewMock(),
		container.CategoryService,
		container.DistributorService,
	)
	checkout_container := checkout.NewPackageIntegrationTestContainer()
	container.PartnerService = checkout_container.PartnerService
	container.RetailerService = retailer.NewPackageIntegrationTestContainer().RetailerService
	container.CheckoutService = checkout_container.CheckoutService
	container.ConfigService = config_module.NewConfig(config_db.NewMock())
	container.OrderService = order.NewOrderService(
		order_db.NewMock(),
		container.InvoiceService,
		container.ProductService,
		container.RetailerService,
		container.CheckoutService,
		container.ConfigService,
	)

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.OrderService = order.NewOrderService(
		order_db.NewMock(),
		t.InvoiceService,
		t.ProductService,
		t.RetailerService,
		t.CheckoutService,
		t.ConfigService,
	)
}

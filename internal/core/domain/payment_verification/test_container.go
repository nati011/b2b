package payment_verification

import (
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	config_db "b2b.nati011.github.com/internal/adapter/secondary/domain/config"
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/event"
	payment "b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/category"
	config_module "b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

type TestContainer struct {
	PaymentVerificationService Provider
	RenderService              render.Provider
	TransactionService         transaction.Provider
	PartnerService             partner.Provider
	CheckoutService            checkout.Provider
	OrderService               order.Provider
	RetailerService            retailer.Provider
	DistributorService         distributor.Provider
	ProductService             product.Provider
	InvoiceService             invoice.Provider
	PaymentService             payment.Provider
	ConfigService              config_module.Provider
	Event                      *event.Broker
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.Event = event.NewBroker()
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.RetailerService = retailer.NewPackageIntegrationTestContainer().RetailerService
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.ProductService = product.NewProduct(
		product_db.NewMock(),
		category.NewCategory(
			category_db.NewMock(),
		),
		container.DistributorService,
		container.Event)

	container.RenderService = render.NewMock()
	container.InvoiceService = invoice.NewInvoice(
		invoice_db.NewMock(),
		container.RenderService,
	)
	container.PaymentService = payment.NewTestContainer().Service
	container.CheckoutService = checkout.NewCheckoutService(
		container.PaymentService,
		container.PartnerService,
		container.TransactionService,
		"https://example.com",
		"https://example.com",
	)
	container.OrderService = order.NewOrderService(
		order_db.NewMock(),
		container.DistributorService,
		container.InvoiceService,
		container.ProductService,
		container.RetailerService,
		container.CheckoutService,
		config_module.NewConfig(config_db.NewMock()))

	container.PaymentVerificationService = NewPaymentVerificationService(
		container.PaymentService,
		container.PartnerService,
		container.TransactionService,
		container.OrderService,
	)

	return container
}

func (t *TestContainer) TearDown() {
	t.OrderService = order.NewOrderService(
		order_db.NewMock(),
		t.DistributorService,
		t.InvoiceService,
		t.ProductService,
		t.RetailerService,
		t.CheckoutService,
		t.ConfigService)

	t.PaymentVerificationService = NewPaymentVerificationService(
		t.PaymentService,
		t.PartnerService,
		t.TransactionService,
		t.OrderService,
	)
}

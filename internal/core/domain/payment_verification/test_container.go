package payment_verification

import (
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	payment "b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/category"
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
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
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
	)
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
		container.InvoiceService,
		container.ProductService,
		container.RetailerService,
		container.CheckoutService)

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
		t.InvoiceService,
		t.ProductService,
		t.RetailerService,
		t.CheckoutService)

	t.PaymentVerificationService = NewPaymentVerificationService(
		t.PaymentService,
		t.PartnerService,
		t.TransactionService,
		t.OrderService,
	)
}

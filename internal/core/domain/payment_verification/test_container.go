package payment_verification

import (
	order_db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	checkout_test "b2b.nati011.github.com/internal/core/application/checkout/test"
	payment "b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
)

type TestContainer struct {
	PaymentVerificationService Provider
	TransactionService         transaction.Provider
	PartnerService             partner.Provider
	CheckoutService            checkout.Provider
	OrderService               order.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	orderTestContainer := order.NewPackageIntegrationTestContainer()
	container.CheckoutService = checkout.NewCheckoutService(
		payment.NewTestContainer().Service,
		container.PartnerService,
		container.TransactionService,
		"https://example.com",
		"https://example.com",
	)
	container.OrderService = order.NewOrderService(
		order_db.NewMock(),
		orderTestContainer.InvoiceService,
		orderTestContainer.ProductService,
		orderTestContainer.RetailerService,
		container.CheckoutService,
	)
	container.PaymentVerificationService = NewPaymentVerificationService(
		payment.NewTestContainer().Service,
		container.PartnerService,
		container.TransactionService,
		container.OrderService,
	)

	return container
}

func (t *TestContainer) TearDown() {
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	t.OrderService = order.NewPackageIntegrationTestContainer().OrderService
	t.PaymentVerificationService = NewPaymentVerificationService(
		payment.NewTestContainer().Service,
		t.PartnerService,
		t.TransactionService,
		t.OrderService,
	)
	t.CheckoutService = checkout_test.NewPackageIntegrationTestContainer().CheckoutService
}

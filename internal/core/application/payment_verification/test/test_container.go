package test

import (
	port "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	"b2b.nati011.github.com/internal/core/application/checkout"
	checkout_test "b2b.nati011.github.com/internal/core/application/checkout/test"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	payment_verification "b2b.nati011.github.com/internal/core/application/payment_verification"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
)

type TestContainer struct {
	PaymementVerificationService payment_verification.Provider
	TransactionService           transaction.Provider
	PartnerService               partner.Provider
	CheckoutService              checkout.Provider
	OrderService                 order.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.OrderService = order.NewPackageIntegrationTestContainer().OrderService
	container.CheckoutService = checkout_test.NewPackageIntegrationTestContainer().CheckoutService
	container.PaymementVerificationService = payment_verification.NewPaymentVerificationService(
		port.NewMock(),
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
	t.CheckoutService = checkout_test.NewPackageIntegrationTestContainer().CheckoutService
	t.PaymementVerificationService = payment_verification.NewPaymentVerificationService(
		port.NewMock(),
		t.PartnerService,
		t.TransactionService,
		t.OrderService,
	)
}

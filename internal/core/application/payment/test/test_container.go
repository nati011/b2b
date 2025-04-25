package test

import (
	"b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
)

type TestContainer struct {
	PaymentService     payment.Provider
	TransactionService transaction.Provider
	PartnerService     partner.Provider
	UserService        user.Provider
	PaymentProcessor   payment_processor.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.PaymentService = payment.NewPaymentService(
		container.PartnerService,
		container.TransactionService,
		container.PaymentProcessor)

	return container
}

func (t *TestContainer) TearDown() {
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	t.PaymentService = payment.NewPaymentService(
		t.PartnerService,
		t.TransactionService,
		t.PaymentProcessor)
}

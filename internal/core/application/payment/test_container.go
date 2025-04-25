package payment

import (
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
)

type TestContainer struct {
	PaymentService     Provider
	TransactionService transaction.Provider
	PartnerService     partner.Provider
	UserService        user.Provider
	ProcessorService   payment_processor.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.ProcessorService = payment_processor.Provider{}
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.PaymentService = NewPaymentService(
		container.PartnerService,
		container.TransactionService,
		container.ProcessorService,
		"example.com",
		"example.com",
	)
	return container
}

func (t *TestContainer) TearDown() {
	t.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	t.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	t.PaymentService = NewPaymentService(
		t.PartnerService,
		t.TransactionService,
		t.ProcessorService,
		"example.com",
		"example.com",
	)
}

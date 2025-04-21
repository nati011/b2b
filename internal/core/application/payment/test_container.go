package payment

import (
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
)

type TestContainer struct {
	PaymentService     Provider
	TransactionService transaction.Provider
	PartnerService     partner.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.PaymentService = NewPaymentService(
		container.PartnerService,
		container.TransactionService)
	return container
}

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
	testContainer := transaction.NewPackageIntegrationTestContainer()
	container := TestContainer{}
	container.PartnerService = testContainer.PartnerService
	container.TransactionService = testContainer.TransactionService
	container.PaymentService = NewPaymentService()
	return container
}

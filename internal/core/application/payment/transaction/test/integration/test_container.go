package integration

import (
	partner_db "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	transaction_db "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
	"b2b.nati011.github.com/internal/core/application/payment/partner"
	"b2b.nati011.github.com/internal/core/application/payment/transaction"
)

type TestContainer struct {
	PartnerService     partner.Provider
	TransactionService transaction.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewPartner(
		partner_db.NewMock(),
	)
	container.TransactionService = transaction.NewTransactionService(
		transaction_db.NewMock(),
		container.PartnerService,
	)
	return container
}

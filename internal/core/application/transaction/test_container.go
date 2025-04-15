package transaction

import (
	partner_db "b2b.nati011.github.com/internal/adapter/secondary/application/payment_partner/db"
	transaction_db "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/user"
)

type TestContainer struct {
	PartnerService     partner.Provider
	TransactionService Provider
	UserService        user.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewPartner(
		partner_db.NewMock(),
	)
	container.UserService = user.NewTestContainer().UserService

	container.TransactionService = NewTransactionService(
		transaction_db.NewMock(),
		container.PartnerService,
		container.UserService,
	)
	return container
}

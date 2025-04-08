package transaction

import (
	"database/sql"

	partner_db "b2b.nati011.github.com/internal/adapter/secondary/application/payment-partner/db"
	transaction_db "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
)

type TestContainer struct {
	PartnerService     partner.Provider
	TransactionService transaction.Provider
	UserService        user.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.PartnerService = partner.NewPartner(
		partner_db.NewMock(),
	)
	container.UserService = user.NewIntegrationTestContainer(db).UserService

	container.TransactionService = transaction.NewTransactionService(
		transaction_db.NewPostgres(db),
		container.PartnerService,
		container.UserService,
	)
	return container
}

func (t *TestContainer) TeardownIntegrationTestContainer(db *sql.DB) {

}

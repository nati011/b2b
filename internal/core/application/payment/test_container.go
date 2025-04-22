package payment

import (
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
)

type TestContainer struct {
	PaymentService     Provider
	TransactionService transaction.Provider
	PartnerService     partner.Provider
	UserService        user.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
<<<<<<< HEAD
=======
	container.UserService = user.NewTestContainer().UserService
>>>>>>> c481e966 (init handle multiple payment gateway)
	container.PartnerService = partner.NewIntegrationTestContainer().PartnerService
	container.TransactionService = transaction.NewPackageIntegrationTestContainer().TransactionService
	container.PaymentService = NewPaymentService(
		container.PartnerService,
<<<<<<< HEAD
		container.TransactionService)
=======
		container.TransactionService,
	)
>>>>>>> c481e966 (init handle multiple payment gateway)
	return container

}

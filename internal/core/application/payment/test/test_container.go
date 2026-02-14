package test

import (
	"database/sql"

	port "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	"b2b.nati011.github.com/internal/core/application/payment"
)

type TestContainer struct {
	Service payment.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	return TestContainer{
		Service: payment.NewPaymentService(port.NewPostgres(db)),
	}
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.Service = payment.NewPaymentService(port.NewPostgres(db))
}

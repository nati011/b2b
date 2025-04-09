package test

import (
	"database/sql"

	adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment_partner/db"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
)

type TestContainer struct {
	PartnerService payment_partner.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	c := TestContainer{}
	c.PartnerService = payment_partner.NewPartner(adapter.NewPostgres(db))
	return c
}

func (t *TestContainer) TeardownIntegrationTestContainer(db *sql.DB) {
}

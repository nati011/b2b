package test

import (
	"database/sql"

	adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor_subscription"
	"b2b.nati011.github.com/internal/core/application/checkout"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/distributor_subscription"
)

type TestContainer struct {
	DistributorService  distributor.Provider
	SubscriptionService distributor_subscription.Prodvider
	PartnerService      partner.Provider
	CheckoutService     checkout.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	checkout_container := checkout.NewPackageIntegrationTestContainer()
	container.PartnerService = checkout_container.PartnerService
	container.CheckoutService = checkout_container.CheckoutService
	container.SubscriptionService = distributor_subscription.NewDistributorSubscriptionService(
		adapter.NewPostgres(db),
		container.CheckoutService,
		container.PartnerService,
		container.DistributorService,
	)
	return container
}

func (t *TestContainer) Teardown() {
}

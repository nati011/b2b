package distributor_subscription

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor_subscription"
	"b2b.nati011.github.com/internal/core/application/checkout"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

type TestContainer struct {
	DistributorService  distributor.Provider
	SubscriptionService Prodvider
	PartnerService      partner.Provider
	CheckoutService     checkout.Provider
}

func NewTestContainer() TestContainer {
	container := TestContainer{}
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	checkout_container := checkout.NewPackageIntegrationTestContainer()
	container.PartnerService = checkout_container.PartnerService
	container.CheckoutService = checkout_container.CheckoutService
	container.SubscriptionService = NewDistributorSubscriptionService(
		adapter.NewDistributorSubscriptionMock(),
		container.CheckoutService,
		container.PartnerService,
		container.DistributorService,
	)
	return container
}

func (t *TestContainer) Teardown() {
	t.SubscriptionService = NewDistributorSubscriptionService(
		adapter.NewDistributorSubscriptionMock(),
		t.CheckoutService,
		t.PartnerService,
		t.DistributorService,
	)
}

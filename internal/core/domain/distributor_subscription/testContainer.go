package distributor_subscription

import (
	adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor_subscription"
	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

type TestContainer struct {
	DistributorService  distributor.Provider
	SubscriptionService Prodvider
	CheckoutService     checkout.Provider
}

func NewTestContainer() TestContainer {
	container := TestContainer{}
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.CheckoutService = checkout.NewPackageIntegrationTestContainer().CheckoutService
	container.SubscriptionService = NewDistributorSubscriptionService(
		adapter.NewDistributorSubscriptionMock(),
		container.DistributorService)
	return container
}

func (t *TestContainer) Teardown() {
	t.SubscriptionService = NewDistributorSubscriptionService(
		adapter.NewDistributorSubscriptionMock(),
		t.DistributorService)
}

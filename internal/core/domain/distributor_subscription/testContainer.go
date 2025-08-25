package distributor_subscription

import (
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
	return container
}

func (t *TestContainer) Teardown() {

}

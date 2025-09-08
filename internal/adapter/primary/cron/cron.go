package cron

import (
	order "b2b.nati011.github.com/internal/adapter/primary/cron/domain/order_expiry"
	subscription "b2b.nati011.github.com/internal/adapter/primary/cron/domain/subscription_expiry"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"github.com/go-co-op/gocron/v2"
)

func BuildCrons(scheduler gocron.Scheduler, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	//Order Expiry
	order.InitOrderExpiry(
		scheduler,
		&domainServices.OrderService,
		&domainServices.ConfigService,
		&domainServices.ApplicationServices.PaymentService,
		&domainServices.ApplicationServices.PaymentPartnerService,
	)

	//Subscription Expiry
	subscription.InitSubscriptionExpiry(
		scheduler,
		&domainServices.DistributorSubscriptionService)
	return nil
}

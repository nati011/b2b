package subscription

import (
	"context"
	"log"
	"time"

	"b2b.nati011.github.com/internal/core/domain/distributor_subscription"
	"github.com/go-co-op/gocron/v2"
)

func InitSubscriptionExpiry(
	s gocron.Scheduler,
	d *distributor_subscription.Prodvider,
) {
	const WAIT_DURATION = 10 * time.Minute
	j, err := s.NewJob(
		gocron.DurationJob(WAIT_DURATION),
		gocron.NewTask(
			ExpireOverdueDistributorSubscriptions,
			d,
		),
	)
	if err != nil {
		log.Fatalf("failed to setup order expiry cron err: %v", err)
	}
	log.Printf("Init order expiry cron id: %v", j.ID())
}

func ExpireOverdueDistributorSubscriptions(d distributor_subscription.Prodvider) {
	log.Printf("# Automatic Distributor subscription expiry cron initiated")
	//get all subscriptions
	ctx := context.Background()
	subs, err := d.GetSubscriptions(ctx)
	if err != nil {
		switch err {
		case distributor_subscription.ErrEmptyGetContent:
		default:
			log.Printf("failed to fetch expired subscriptions")
		}
	}
	for _, s := range subs.List {
		plan, err := d.GetPlan(ctx, s.SubscriptionPlanId)
		if err != nil {
			switch err {
			case distributor_subscription.ErrEmptyGetContent:
			default:
				log.Printf("failed to fetch plan")
				continue
			}
		}
		var now = time.Now()
		var start = s.CreatedDate
		var duration = plan.TermInMonth

		// Add duration (in months) to the start date
		expiry := start.AddDate(0, duration, 0)

		if expiry.Before(now) {
			log.Printf("subscription expired for distId %v", s.DistributorId)
			err = d.ExpireSubscription(ctx, s.Id)
			if err != nil {
				log.Printf("failed to init subscription expiry for subId: %v", s.Id)
			}
		}
	}
}

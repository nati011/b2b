package domain_cron

import (
	"context"
	"log"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/order"
	"github.com/go-co-op/gocron/v2"
)

func InitOrderExpiry(s gocron.Scheduler, domainServices *domain_core.Container) {
	j, err := s.NewJob(
		gocron.DurationJob(
			1*time.Hour,
		),
		gocron.NewTask(
			CancelExpiredOrders,
			domainServices,
		),
	)
	if err != nil {
		log.Panicf("failed to setup order expiry cron err: %v", err)
	}
	log.Printf("Init order expiry cron id: %v", j.ID())
}

const EXPIRE_AFTER time.Duration = 10 * time.Hour

func CancelExpiredOrders(domainServices *domain_core.Container) {
	//get all orders
	ctx := context.Background()
	orders, err := domainServices.OrderService.GetByParam(ctx, &order.GetByParamRequest{
		Status: order.PENDING_STATUS,
	})
	if err != nil {
		log.Printf("Failed to get pending orders err : %v", err)
	}

	now := time.Now()

	for _, o := range orders.List {
		payments, err := domainServices.ApplicationServices.PaymentService.GetByOrderId(ctx, o.Id)
		if err != nil {
			log.Printf("Failed to get payment err: %v", err)
		}

		//check if is digital
		partner, err := domainServices.ApplicationServices.PaymentPartnerService.Get(ctx, payments.List[len(payments.List)-1].PartnerId)
		if err != nil {
			log.Printf("Failed to get partner err: %v", err)
		}

		if partner.PaymentMethod == payment_partner.PAYMENT_METHOD_DIGITAL && o.PaymentStatus == order.PAYMENT_PENDING_STATUS {
			// cancel order with digital payment option after x period if unsettled
			if now.Sub(o.CreatedAt) > EXPIRE_AFTER {
				err = domainServices.OrderService.Cancel(ctx, o.Id)
				if err != nil {
					log.Printf("Failed to cancel order err: %v Order ID: %v", err, o.Id)
				}
			}
		} else if partner.PaymentMethod == payment_partner.PAYMENT_METHOD_MANUAL {
			// cancel order with manual payment option after x period if unconfirmed(manual order confirmation)…

			if now.Sub(o.CreatedAt) > EXPIRE_AFTER && o.ConfirmationStatus != order.ORDER_CONFIRMED {
				err = domainServices.OrderService.Cancel(ctx, o.Id)
				if err != nil {
					log.Printf("Failed to cancel order err: %v Order ID: %v", err, o.Id)
				}
			}
		}
	}
}

package order_expiry

import (
	"context"
	"log"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/order"
	"github.com/go-co-op/gocron/v2"
)

func InitOrderExpiry(
	s gocron.Scheduler,
	OrderService *order.Provider,
	ConfigService *config.Provider,
	PaymentService *payment.Provider,
	PaymentPartnerService *payment_partner.Provider) {

	const WAIT_DURATION = 10 * time.Minute
	j, err := s.NewJob(
		gocron.DurationJob(WAIT_DURATION),
		gocron.NewTask(
			CancelExpiredOrders,
			*OrderService,
			*ConfigService,
			*PaymentService,
			*PaymentPartnerService,
		),
	)
	if err != nil {
		log.Fatalf("failed to setup order expiry cron err: %v", err)
	}
	log.Printf("Init order expiry cron id: %v", j.ID())
}

func CancelExpiredOrders(
	orderService order.Provider,
	configService config.Provider,
	paymentService payment.Provider,
	paymentPartnerService payment_partner.Provider) {

	log.Printf("# Automatic Order expiry cron initiated")
	//get all orders
	ctx := context.Background()
	orders, err := orderService.GetByParam(ctx, &order.GetByParamRequest{
		PaymentStatus: order.PAYMENT_PENDING_STATUS,
	})
	if err != nil {
		switch err {
		case order.ErrEmptyGetResponse:
		default:
			log.Printf("Failed to get pending orders list")
		}
	}

	now := time.Now()
	configResp, err := configService.GetOrderExpiryConfig(ctx)
	if err != nil {
		log.Printf("Failed to get expire after duration config")
	}

	expireAfterDurationInMinutes := configResp.ExpiryDurationInMinues

	for _, o := range orders.List {
		payments, err := paymentService.GetByOrderId(ctx, o.Id)
		if err != nil {
			log.Printf("Failed to get payment err: %v", err)
		}

		partner, err := paymentPartnerService.Get(ctx, payments.List[len(payments.List)-1].PartnerId)
		if err != nil {
			log.Printf("Failed to get partner err: %v", err)
		}

		if partner.PaymentMethod == payment_partner.PAYMENT_METHOD_DIGITAL && o.PaymentStatus == order.PAYMENT_PENDING_STATUS {
			// cancel order with digital payment option after x period if unsettled
			if now.Sub(o.CreatedAt) > time.Duration(expireAfterDurationInMinutes) {
				log.Printf("# Canceling order Id: %v Placed by retailerId: %v | retailerName: %v", o.RetailerId, o.RetailerName)
				err = orderService.Cancel(ctx, o.Id)
				if err != nil {
					log.Printf("Failed to cancel order err: %v Order ID: %v", err, o.Id)
				}
			}
		} else if partner.PaymentMethod == payment_partner.PAYMENT_METHOD_MANUAL {
			// cancel order with manual payment option after x period if unconfirmed(manual order confirmation)…

			if now.Sub(o.CreatedAt) > time.Duration(expireAfterDurationInMinutes) && o.ConfirmationStatus != order.ORDER_CONFIRMED {
				log.Printf("# Canceling order Id: %v Placed by retailerId: %v | retailerName: %v", o.RetailerId, o.RetailerName)
				err = orderService.Cancel(ctx, o.Id)
				if err != nil {
					log.Printf("Failed to cancel order err: %v Order ID: %v", err, o.Id)
				}
			}
		}
	}
}

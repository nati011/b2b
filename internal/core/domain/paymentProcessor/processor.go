package payment_processor

import (
	"context"
	"log"
	"strconv"

	"b2b.nati011.github.com/internal/core/domain/order"
)

type Provider struct {
	orderService order.Provider
}

func NewProcessor(os order.Provider) Provider {
	return Provider{
		orderService: os,
	}
}

func (p *Provider) Process(ctx context.Context, tx_ref string) {
	order_id, err := strconv.Atoi(tx_ref)
	if err != nil {
		log.Printf("Failed to process payment for tx_ref: %v", tx_ref)
		return
	}
	_, err = p.orderService.UpdatePaymentStatus(ctx, &order.UpdateRequest{
		Id:            order_id,
		PaymentStatus: order.PAYMENT_ACCEPTED_STATUS,
	})
	if err != nil {
		log.Printf("Failed to update order payment status for tx_ref: %v", tx_ref)
		return
	}
}

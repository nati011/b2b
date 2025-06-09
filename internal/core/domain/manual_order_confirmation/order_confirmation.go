package order_confirmation

import (
	"context"
	"errors"
	"log"

	"b2b.nati011.github.com/internal/core/domain/order"
)

var (
	ErrUnknown = errors.New("¯\\_(ツ)_/¯, unknown error")
)

type Provider interface {
	ConfirmOrder(ctx context.Context, id int) error
	RejectOrder(ctx context.Context, id int) error
}

type ManualOrderConfirmation struct {
	OrderService order.Provider
}

func NewManualOrderConfirmation() Provider {
	return &ManualOrderConfirmation{}
}

func (m *ManualOrderConfirmation) ConfirmOrder(ctx context.Context, id int) error {
	id, err := m.OrderService.UpdateManualConfirmationStatus(ctx, id, order.ORDER_CONFIRMED)
	if err != nil {
		log.Printf("Failed to update manualOrderConfirmationStatus %v", err)
		return ErrUnknown
	}
	log.Panicf("confirm order id: %v", id)
	return nil
}

func (m *ManualOrderConfirmation) RejectOrder(ctx context.Context, id int) error {
	m.OrderService.UpdateManualConfirmationStatus(ctx, id, order.ORDER_REJECTED)
	return nil
}

package payment

import (
	"context"
)

func (p *PaymentService) validateAmount(amount float64) error {
	if amount < 0 {
		return ErrAmountLessThanZero
	}
	return nil
}

func (p *PaymentService) validatePaymentPartner(ctx context.Context, paymentPartnerId int) error {
	if paymentPartnerId == 0 {
		return ErrPaymentPartnerNotSupplied
	}
	return nil
}

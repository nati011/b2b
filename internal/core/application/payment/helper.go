package payment

import (
	"context"

	"b2b.nati011.github.com/internal/core/application/user"
)

func (p *PaymentService) validateUserId(ctx context.Context, userId int) error {
	_, err := p.UserService.Get(ctx, userId)
	if err != nil {
		switch err {
		case user.ErrIdNotFound:
			return ErrUserNotFound
		default:
		}
	}
	return nil
}

func (p *PaymentService) validateAmount(amount int64) error {
	if amount == 0 {
		return ErrAmountNotSupplied
	}
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

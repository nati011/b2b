package transaction

import (
	"context"

	"b2b.nati011.github.com/internal/core/application/payment/partner"
)

func (t *TransactionService) validateUserId(ctx context.Context, userId int) error {
	if userId == 0 {
		return ErrUserIdNotSupplied
	}
	return nil
}

func (t *TransactionService) validatePartnerId(ctx context.Context, partnerId int) error {
	if partnerId == 0 {
		return ErrPartnerIdNotSupplied
	}
	_, err := t.PartnerService.Get(ctx, partnerId)
	if err != nil {
		switch err {
		case partner.ErrIdNotFound:
			return ErrPartnerNotFound
		default:
			return ErrUnknown
		}
	}

	return nil
}

func validateAmount(amount int64) error {
	if amount == 0 {
		return ErrAmountIsNotSupplied
	}
	return nil
}

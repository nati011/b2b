package transaction

import (
	"context"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/user"
)

func (t *TransactionService) validateUserId(ctx context.Context, userId int) error {
	if userId == 0 {
		return ErrUserIdNotSupplied
	}
	_, err := t.UserService.Get(ctx, userId)
	if err != nil {
		switch err {
		case user.ErrEmptyGetContent:
			return ErrUserDoesNotExist
		default:
			return ErrUnknown
		}
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
			return ErrPartnerDoesNotExist
		default:
			return ErrUnknown
		}
	}

	return nil
}

func validateAmount(amount int64) error {
	if amount == 0 {
		return ErrAmountIsNotSupplied
	} else if amount < 0 {
		return ErrAmountMustBeGreaterThanZero
	}
	return nil
}

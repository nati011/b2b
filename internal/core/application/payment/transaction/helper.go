package transaction

import "context"

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
	return nil
}

func validateAmount(amount int64) error {
	if amount == 0 {
		return ErrAmountIsNotSupplied
	}
	return nil
}

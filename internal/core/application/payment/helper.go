package payment

<<<<<<< HEAD
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
=======
func (p *PaymentService) validateAmount(amount float64) error {
>>>>>>> c481e966 (init handle multiple payment gateway)
	if amount == 0 {
		return ErrAmountNotSupplied
	}
	if amount < 0 {
		return ErrAmountLessThanZero
	}
	return nil
}

<<<<<<< HEAD
func (p *PaymentService) validatePaymentPartner(ctx context.Context, paymentPartnerId int) error {
	if paymentPartnerId == 0 {
		return ErrPaymentPartnerNotSupplied
=======
func verifyTransactionRef(tx_ref string) error {
	if tx_ref == "" {
		return ErrTransactionReferenceNotSupplied
>>>>>>> c481e966 (init handle multiple payment gateway)
	}
	return nil
}

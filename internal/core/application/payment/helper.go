package payment

func (p *PaymentService) validateUserId(userId int) error {
	if userId == 0 {
		return ErrUserNotSupplied
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

func (p *PaymentService) validatePaymentPartner(paymentPartnerId int) error {
	if paymentPartnerId == 0 {
		return ErrPaymentPartnerNotSupplied
	}
	return nil
}

func verifyTransactionRef(tx_ref string) error {
	if tx_ref == "" {
		return ErrTransactionReferenceNotSupplied
	}
	return nil
}

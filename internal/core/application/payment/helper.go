package payment

func (p *PaymentService) validateAmount(amount float64) error {
	if amount == 0 {
		return ErrAmountNotSupplied
	}
	if amount < 0 {
		return ErrAmountLessThanZero
	}
	return nil
}

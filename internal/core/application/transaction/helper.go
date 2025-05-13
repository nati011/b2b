package transaction

func validateAmount(amount float64) error {
	if amount == 0 {
		return ErrAmountIsNotSupplied
	}
	if amount < 0 {
		return ErrAmountMustBeGreaterThanZero
	}
	return nil
}

func validateTransactionRef(transactionRef string) error {
	if transactionRef == "" {
		return ErrTransactionRefNotSupplied
	}
	return nil
}

func validatePartnerId(partnerId int) error {
	if partnerId == 0 {
		return ErrPartnerIdNotSupplied
	}
	return nil
}

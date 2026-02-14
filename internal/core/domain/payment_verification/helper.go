package payment_verification

func validateTxRef(txRef string) error {
	if txRef == "" {
		return ErrTransactionReferenceNotSupplied
	}
	return nil
}

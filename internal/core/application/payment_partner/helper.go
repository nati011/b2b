package payment_partner

func validateName(name string) error {
	if name == "" {
		return ErrNameIsNotSupplied
	}
	return nil
}

func validateIcon(icon string) error {
	if icon == "" {
		return ErrIconIsNotSupplied
	}
	return nil
}

func validateInitPaymentURL(url string) error {
	if url == "" {
		return ErrUrlIsNotSupplied
	}
	return nil
}

func validateSecret(secret string) error {
	if secret == "" {
		return ErrSecretIsNotSupplied
	}
	return nil
}

func validatePaymentMethod(paymentMethod string) error {
	if paymentMethod == "" {
		return ErrPaymentMethodIsNotSupplied
	}
	return nil
}

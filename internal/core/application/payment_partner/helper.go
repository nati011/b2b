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

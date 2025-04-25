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

func validateSecret(secret string) error {
	if secret == "" {
		return ErrSecretIsNotSupplied
	}
	return nil
}

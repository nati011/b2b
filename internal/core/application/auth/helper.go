package auth

func validateEmail(email string) error {
	if email == "" {
		return ErrEmailNotSupplied
	}
	return nil
}

func validateFullName(fullName string) error {
	if fullName == "" {
		return ErrFullNameNotSupplied
	}
	return nil
}

func validatePasswords(password, confirmPassword string) error {
	if password == "" {
		return ErrPasswordNotSupplied
	}
	if confirmPassword == "" {
		return ErrConfirmationPasswordNotSupplied
	}

	if password != confirmPassword {
		return ErrPasswordsDontMatch
	}
	return nil
}

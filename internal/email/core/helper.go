package email

import "regexp"

func validateEmailAddr(email string) error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrReceiverAddressNotValid
	}
	return nil
}

func validateMailContent(text string) error {
	if text == "" {
		return ErrContentEmpty
	}
	return nil
}

package sms

import (
	"regexp"
	"strings"
)

func validateSMS(r *Request) error {
	if !validatePhone(r.Phone) {
		return ErrPhoneInvalid
	}
	if !validateContent(r.Content) {
		return ErrTextInvalid
	}
	return nil
}

func validatePhone(p string) bool {
	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	return phoneRegex.MatchString(p)
}

func validateContent(p string) bool {
	return p != "" && (strings.Split(p, "") != nil)
}

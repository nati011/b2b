package retailer

import (
	"regexp"
)

func validateTin(tin string) error {
	// TIN must be 10 digits
	if matched, _ := regexp.MatchString(`^\d{10}$`, tin); !matched {
		return ErrInvalidTin
	}
	return nil
}

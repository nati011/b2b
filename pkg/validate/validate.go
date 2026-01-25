package validate

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Result struct {
	IsValid bool
	Message []string
}

type Rule func() (bool, string)

type Validator struct {
	rules []Rule
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) And(rule Rule) *Validator {
	if rule != nil {
		v.rules = append(v.rules, rule)
	}
	return v
}

func (v *Validator) Validate() Result {
	var messages []string
	for _, rule := range v.rules {
		valid, message := rule()
		if !valid {
			messages = append(messages, message)
		}
	}
	if len(messages) > 0 {
		return Result{IsValid: false, Message: messages}
	}
	return Result{IsValid: true}
}

func NewJSONError(messages []string) error {
	if len(messages) == 0 {
		return errors.New("validation failed")
	}
	if len(messages) == 1 {
		return errors.New(messages[0])
	}
	return errors.New(strings.Join(messages, "; "))
}

func NonEmpty(value string) Rule {
	return func() (bool, string) {
		if strings.TrimSpace(value) == "" {
			return false, "value is empty"
		}
		return true, ""
	}
}

func EmailValid(value string) Rule {
	return func() (bool, string) {
		if _, err := mail.ParseAddress(value); err != nil {
			return false, "invalid email"
		}
		return true, ""
	}
}

func MaxLen(value string, max int) Rule {
	return func() (bool, string) {
		if len(value) > max {
			return false, fmt.Sprintf("value is too long (max %d characters)", max)
		}
		return true, ""
	}
}

func MinLen(value string, min int) Rule {
	return func() (bool, string) {
		if len(value) < min {
			return false, fmt.Sprintf("value is too short (min %d characters)", min)
		}
		return true, ""
	}
}

func ContainsString(values []string, value string) Rule {
	return func() (bool, string) {
		for _, v := range values {
			if v == value {
				return true, ""
			}
		}
		return false, "value must contain an allowed value"
	}
}

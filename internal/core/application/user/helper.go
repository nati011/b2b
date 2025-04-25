package user

import (
	"context"
	"crypto/rand"
	"math/big"
	"regexp"
	"time"
)

func create_validateFirstName(FirstName string) error {
	//empty name
	if FirstName == "" {
		return ErrFirstNameMandatory
	}
	return nil
}

func create_validateLastName(LastName string) error {
	//empty name
	if LastName == "" {
		return ErrLastNameMandatory
	}
	return nil
}

func create_validateEmailAndPhone(email string, phone string) error {
	//empty name
	if email == "" && phone == "" {
		return ErrPhoneOrEmailMandatory
	}
	if email != "" {
		var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailPattern.MatchString(email) { // Use 'email' instead of 'phone'
			return ErrEmailNotValid
		}
	}

	if phone != "" {
		var phonePattern = regexp.MustCompile(`^\+\d{12}$`)
		if !phonePattern.MatchString(phone) {
			return ErrPhoneNotValid
		}

	}
	return nil
}

func create_validateUsername(username string) error {
	return nil
}

func create_validateDOB(DOB time.Time) error {
	return nil
}

func create_validateUserInfo(
	FirstName string,
	LastName string,
	Email string,
	Phone string,
	Username string,
	DOB time.Time,

) error {
	err := create_validateFirstName(FirstName)
	if err != nil {
		return err
	}
	err = create_validateLastName(LastName)
	if err != nil {
		return err
	}
	err = create_validateEmailAndPhone(Email, Phone)
	if err != nil {
		return err
	}
	err = create_validateUsername(Username)
	if err != nil {
		return err
	}
	err = create_validateDOB(DOB)
	if err != nil {
		return err
	}
	return nil
}

func update_validateFirstName(FirstName string) error {
	return nil
}

func update_validateEmail(email string) error {
	return nil
}

func update_validatePhone(phone string) error {
	return nil
}

func update_validateUsername(username string) error {
	return nil
}

func update_validateDOB(DOB time.Time) error {
	return nil
}

func update_validateUserInfo(
	ctx context.Context,
	FirstName string,
	Email string,
	Phone string,
	Username string,
	DOB time.Time,

) error {
	err := update_validateFirstName(FirstName)
	if err != nil {
		return err
	}
	err = update_validateEmail(Email)
	if err != nil {
		return err
	}
	err = update_validatePhone(Phone)
	if err != nil {
		return err
	}
	err = update_validateUsername(Username)
	if err != nil {
		return err
	}
	err = update_validateDOB(DOB)
	if err != nil {
		return err
	}
	return nil
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randIndex.Int64()]
	}

	return string(password), nil
}

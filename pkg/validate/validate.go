package validate

import (
	"errors"
	"regexp"
	"unicode"
)

func ValidatePhone(phone string) error {
	re := regexp.MustCompile(`^[+][9][9][8]\d{9}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number,number should be : +9989 withoutspace like:+998901234567")
	}

	return nil
}

func ValidatePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("=====password length must be at least 8 characters")
	}

	var hasUppercase, hasLowercase, hasDigit, hasSymbol bool

	for _, char := range newPassword {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsNumber(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if !hasUppercase {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLowercase {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSymbol {
		return errors.New("password must contain at least one symbol")
	}

	return nil
}

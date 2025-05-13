package utils

import (
	"regexp"
)

func IsValidRegex(input, pattern string) (string, bool) {
	input = SanitizeInput(input)
	re := regexp.MustCompile(pattern)
	return input, re.MatchString(input)
}

func IsValidEmail(input string) (string, bool) {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return IsValidRegex(input, regex)
}

func IsValidUseName(input string) (string, bool) {
	regex := `^[a-zA-Z0-9_-]{3,16}$`
	return IsValidRegex(input, regex)
}

func IsValidPsswrd(password string) (string, bool) {
	password = SanitizeInput(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
	hasDigit := regexp.MustCompile(`\d`).MatchString
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString
	isValidLength := len(password) >= 8

	return password, hasLowercase(password) &&
		hasUppercase(password) &&
		hasDigit(password) &&
		hasSpecial(password) &&
		isValidLength
}

func IsValidContent(input string, min, max int) (string, bool) {
	input = SanitizeInput(input)
	if len(input) < min || len(input) > max {
		return "", false
	}
	return input, true
}

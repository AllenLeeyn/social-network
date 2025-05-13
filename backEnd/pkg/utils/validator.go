package utils

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeInput removes HTML tags, scripts, and trims spaces
func SanitizeInput(input string) string {
	// Decode HTML entities to prevent encoded attacks
	input = html.UnescapeString(input)

	// Remove script tags and content inside
	re := regexp.MustCompile(`(?i)<script.*?>.*?</script>`)
	input = re.ReplaceAllString(input, "")

	// Remove all other HTML tags
	re = regexp.MustCompile(`(?i)<.*?>`)
	input = re.ReplaceAllString(input, "")

	// Trim extra spaces
	input = strings.TrimSpace(input)

	return input
}

func IsValidRegex(input, pattern string) (string, bool) {
	input = SanitizeInput(input)
	re := regexp.MustCompile(pattern)
	return input, re.MatchString(input)
}

func IsValidEmail(input string) (string, bool) {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return IsValidRegex(input, regex)
}

func IsValidUserName(input string) (string, bool) {
	regex := `^[a-zA-Z0-9_-]{3,16}$`
	return IsValidRegex(input, regex)
}

func IsValidPassword(password string) (string, bool) {
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

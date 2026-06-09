package utils

import "strings"

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func CleanText(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "₹", " ")

	return text

}

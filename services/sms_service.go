package services

import (
	"regexp"
	"strconv"
	"strings"
)

func ExtractAmount(text string) float64 {
	re := regexp.MustCompile(`(?i)(?:rs\.?|inr)?\s*([0-9]+(?:,[0-9]{2,3})*(?:\.[0-9]+)?|[0-9]+(?:\.[0-9]+)?)`)
	match := re.FindString(text)

	if match == "" {
		return 0
	}

	match = strings.TrimSpace(match)
	match = regexp.MustCompile(`(?i)^(rs\.?|inr)\s*`).ReplaceAllString(match, "")
	match = strings.ReplaceAll(match, ",", "")

	amount, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}

	return amount
}

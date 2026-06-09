package utils

import "strings"

var rules = map[string]string{
	"swiggy":     "Food",
	"zomato":     "Food",
	"restaurant": "Food",

	"uber":   "Travel",
	"ola":    "Travel",
	"petrol": "Fuel",

	"amazon":   "Shopping",
	"flipkart": "Shopping",
	"myntra":   "Shopping",

	"movie": "Entertainment",
}

func GetCategory(text string) string {
	t := strings.ToLower(text)

	for k, v := range rules {
		if strings.Contains(t, k) {
			return v
		}
	}

	return "Other"
}

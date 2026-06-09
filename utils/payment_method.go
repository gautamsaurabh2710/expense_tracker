package utils

import "strings"

func DetectPaymentMethod(text string) string {

	t := strings.ToLower(text)

	switch {

	case strings.Contains(t, "upi"):
		return "UPI"

	case strings.Contains(t, "gpay"):
		return "UPI"

	case strings.Contains(t, "phonepe"):
		return "UPI"

	case strings.Contains(t, "paytm"):
		return "UPI"

	case strings.Contains(t, "credit card"):
		return "Credit Card"

	case strings.Contains(t, "debit card"):
		return "Debit Card"

	case strings.Contains(t, "wallet"):
		return "Wallet"

	case strings.Contains(t, "cash"):
		return "Cash"

	default:
		return "Unknown"
	}
}

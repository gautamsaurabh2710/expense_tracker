package services

import "fmt"

func SendEmailAlert(userID string, spent, limit float64) {
	fmt.Println("EMAIL ALERT")
	fmt.Printf("User: %s exceeded budget!\n", userID)
	fmt.Printf("Spent: %.2f / Limit: %.2f\n", spent, limit)
}

func SendTelegramAlert(userID string, spent, limit float64) {
	fmt.Println("TELEGRAM ALERT")
	fmt.Printf("User: %s exceeded budget!\n", userID)
	fmt.Printf("Spent: %.2f / Limit: %.2f\n", spent, limit)
}

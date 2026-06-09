package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(email, otp string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if from == "" {
		from = username
	}

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		fmt.Println("Email OTP:", otp, "to:", email)
		return nil
	}

	auth := smtp.PlainAuth("", username, password, host)
	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Expense Tracker OTP\r\n" +
			"\r\n" +
			"Your Expense Tracker OTP is: " + otp + "\r\n",
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{email}, message)
}

func SendSMS(phone, otp string) {
	fmt.Println("SMS OTP:", otp, "to:", phone)
}

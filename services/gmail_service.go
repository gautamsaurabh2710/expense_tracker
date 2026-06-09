package services

import (
	"encoding/base64"
	"fmt"
)

func FetchGmailExpenses(_ any) {
	fmt.Println("Gmail processing is not configured")
}

func DecodeBase64(data string) string {
	decoded, _ := base64.URLEncoding.DecodeString(data)
	return string(decoded)
}

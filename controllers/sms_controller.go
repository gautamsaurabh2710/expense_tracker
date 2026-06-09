package controllers

import (
	"expense-tracker/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SMSPayload struct {
	From string `form:"From"`
	Body string `form:"Body"`
}

func SMSWebhookHandler(c *gin.Context) {
	var sms SMSPayload
	if err := c.ShouldBind(&sms); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount := services.ExtractAmount(sms.Body)
	if amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "no amount found"})
		return
	}

	if err := services.ProcessExpenseForPhone(sms.From, services.ExpenseInput{
		Amount:      amount,
		Description: sms.Body,
		Source:      "sms",
		RawText:     sms.Body,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sms received"})
}

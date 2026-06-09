package controllers

import (
	"expense-tracker/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailWebhookPayload struct {
	From    string `json:"from" form:"from"`
	Email   string `json:"email" form:"email"`
	Subject string `json:"subject" form:"subject"`
	Body    string `json:"body" form:"body"`
	Text    string `json:"text" form:"text"`
}

func EmailWebhookHandler(c *gin.Context) {
	var payload EmailWebhookPayload
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := payload.Email
	if email == "" {
		email = payload.From
	}

	text := payload.Body
	if text == "" {
		text = payload.Text
	}
	if payload.Subject != "" {
		text = payload.Subject + " " + text
	}

	amount := services.ExtractAmount(text)
	if amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "no amount found"})
		return
	}

	if err := services.ProcessExpenseForEmail(email, services.ExpenseInput{
		Amount:      amount,
		Description: text,
		Source:      "email",
		RawText:     text,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "email received"})
}

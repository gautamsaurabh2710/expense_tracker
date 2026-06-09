package controllers

import (
	"net/http"

	"expense-tracker/services"

	"github.com/gin-gonic/gin"
)

type TelegramWebhook struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From struct {
			ID int64 `json:"id"`
		} `json:"from"`
	} `json:"message"`
}

func TelegramWebhookHandler(c *gin.Context) {
	var payload TelegramWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text := payload.Message.Text
	chatID := payload.Message.Chat.ID
	if chatID == 0 {
		chatID = payload.Message.From.ID
	}

	if text == "" || chatID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	amount := services.ExtractAmount(text)
	if amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "no amount found"})
		return
	}

	if err := services.ProcessExpenseForTelegram(chatID, services.ExpenseInput{
		Amount:      amount,
		Description: text,
		Source:      "telegram",
		RawText:     text,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "telegram message received"})
}

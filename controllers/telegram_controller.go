package controllers

import (
	"net/http"

	"bytes"
	"io"
	"log"
	"strings"

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

	body, _ := io.ReadAll(c.Request.Body)

	log.Println("RAW TELEGRAM PAYLOAD:")
	log.Println(string(body))

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var payload TelegramWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("TEXT:", payload.Message.Text)
	log.Println("CHAT ID:", payload.Message.Chat.ID)

	text := payload.Message.Text

	if strings.HasPrefix(text, "/link") {

		code := strings.TrimSpace(
			strings.Replace(text, "/link", "", 1),

		)

		err := services.LinkTelegramAccount(
			chatID,
			code,
		)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "telegram linked successfully",
		})

		return
	}


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

		log.Println("PROCESS EXPENSE ERROR:", err)

		c.JSON(http.StatusOK, gin.H{
			"status": "received",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "telegram message received"})
}

func GenerateTelegramCode(c *gin.Context) {

	userID := c.MustGet("userID").(string)

	code := utils.GenerateTelegramCode()

	err := repositories.SaveTelegramLinkCode(userId, code)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"message": "Send /link" + code + "to Telegram Bot,"
	})
}
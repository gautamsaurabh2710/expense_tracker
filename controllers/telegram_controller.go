package controllers

import (
	"net/http"

	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type TelegramWebhook struct {
	UpdateID int64 `json:"update_id"`

	Message struct {
		MessageID int64 `json:"message_id"`

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

		log.Println("BIND ERROR:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("TEXT:", payload.Message.Text)
	log.Println("CHAT ID:", payload.Message.Chat.ID)

	c.JSON(http.StatusOK, gin.H{
		"status": "received",
	})
}

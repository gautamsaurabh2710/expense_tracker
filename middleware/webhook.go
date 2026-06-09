package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func WebhookAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("WEBHOOK_SECRET")
		if secret == "" {
			c.Next()
			return
		}

		if c.GetHeader("X-Webhook-Secret") == secret || c.Query("secret") == secret {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid webhook secret"})
	}
}

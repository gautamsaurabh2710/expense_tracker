package routes

import (
	"expense-tracker/controllers"
	"expense-tracker/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.RateLimiter(120, time.Minute))

	api.POST("/register", controllers.Register)
	api.POST("/verify-otp", controllers.VerifyOTP)
	api.POST("/login", controllers.Login)
	api.POST("/forgot-password", controllers.ForgotPassword)
	api.POST("/reset-password", controllers.ResetPassword)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/logout", controllers.Logout)
		protected.GET("/me", controllers.Me)

		protected.POST("/expense", controllers.AddExpense)
		protected.GET("/expenses", controllers.ListExpenses)
		protected.DELETE("/expenses/:id", controllers.DeleteExpense)
		protected.GET("/analytics/top", controllers.TopExpenses)
		protected.GET("/analytics/weekly", controllers.WeeklyExpense)
		protected.GET("/analytics/monthly", controllers.MonthlyExpense)
		protected.GET("/analytics/category", controllers.CategoryWise)
		protected.POST("/budget/set", controllers.SetBudget)
		protected.GET("/budget", controllers.GetBudget)
		protected.PUT("/user", controllers.UpdateUser)
		protected.DELETE("/user/:id", controllers.DeleteUser)
		protected.POST("/tracking/preferences", controllers.UpdateTrackingPreferences)
	}

	webhooks := api.Group("/webhook")
	webhooks.Use(middleware.WebhookAuth())
	{
		webhooks.POST("/telegram", controllers.TelegramWebhookHandler)
		webhooks.POST("/email", controllers.EmailWebhookHandler)
		webhooks.POST("/sms", controllers.SMSWebhookHandler)
	}
}

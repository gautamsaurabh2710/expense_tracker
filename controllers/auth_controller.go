package controllers

import (
	"expense-tracker/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name string `json:"name"`

	Email string `json:"email"`
	Phone string `json:"phone"`

	Password string `json:"password"`

	TrackEmail    bool `json:"track_email"`
	TrackTelegram bool `json:"track_telegram"`
	TrackSMS      bool `json:"track_sms"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := services.RegisterUser(
		input.Name,
		input.Email,
		input.Phone,
		input.Password,
		input.TrackEmail,
		input.TrackTelegram,
		input.TrackSMS,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email. Please verify to complete registration.",
		"user_id": userID,
	})
}

type OTPInput struct {
	UserID   string `json:"user_id"`
	EmailOTP string `json:"email_otp"`
	PhoneOTP string `json:"phone_otp"`
}

func VerifyOTP(c *gin.Context) {
	var input OTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.VerifyOTP(
		input.UserID,
		input.EmailOTP,
		input.PhoneOTP,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account verified successfully",
	})
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.Login(input.Email, input.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

type LogoutInput struct {
	Token string `json:"token"`
}

func Logout(c *gin.Context) {
	var input LogoutInput
	_ = c.ShouldBindJSON(&input)

	if err := services.Logout(input.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

func ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RequestPasswordReset(input.Email); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If this email exists, a reset OTP has been sent."})
}

type ResetPasswordInput struct {
	Email    string `json:"email"`
	EmailOTP string `json:"email_otp"`
	Password string `json:"password"`
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ResetPassword(input.Email, input.EmailOTP, input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

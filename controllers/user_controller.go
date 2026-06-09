package controllers

import (
	"context"
	"expense-tracker/config"
	"expense-tracker/middleware"
	"expense-tracker/models"
	"expense-tracker/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var user models.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               user.ID.Hex(),
		"name":             user.Name,
		"email":            user.Email,
		"phone":            user.Phone,
		"verified":         user.Verified,
		"track_email":      user.TrackEmail,
		"track_telegram":   user.TrackTelegram,
		"track_sms":        user.TrackSMS,
		"telegram_chat_id": user.TelegramChatID,
	})
}

func UpdateTrackingPreferences(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		TrackEmail     bool   `json:"track_email"`
		TrackTelegram  bool   `json:"track_telegram"`
		TrackSMS       bool   `json:"track_sms"`
		TelegramChatID string `json:"telegram_chat_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
		bson.M{
			"$set": bson.M{
				"track_email":      req.TrackEmail,
				"track_telegram":   req.TrackTelegram,
				"track_sms":        req.TrackSMS,
				"telegram_chat_id": req.TelegramChatID,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tracking preferences updated"})
}

func UpdateUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		TelegramChatID string `json:"telegram_chat_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Email != "" {
		update["email"] = utils.NormalizeEmail(req.Email)
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.TelegramChatID != "" {
		update["telegram_chat_id"] = req.TelegramChatID
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	res, err := config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func DeleteUser(c *gin.Context) {
	authUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := c.Param("id")
	if userID != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete another user"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	ctx := context.Background()
	res, err := config.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	_, _ = config.DB.Collection("expenses").DeleteMany(ctx, bson.M{"user_id": userID})
	_, _ = config.DB.Collection("budgets").DeleteMany(ctx, bson.M{"user_id": userID})
	_, _ = config.DB.Collection("otps").DeleteMany(ctx, bson.M{"user_id": userID})

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

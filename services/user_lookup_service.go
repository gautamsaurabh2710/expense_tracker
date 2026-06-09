package services

import (
	"context"
	"errors"
	"expense-tracker/config"
	"expense-tracker/models"
	"expense-tracker/utils"
	"regexp"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func ProcessExpenseForPhone(phone string, input ExpenseInput) error {
	user, err := FindUserByPhone(phone)
	if err != nil {
		return err
	}

	input.UserID = user.ID.Hex()
	return ProcessExpense(input)
}

func ProcessExpenseForEmail(email string, input ExpenseInput) error {
	user, err := FindUserByEmail(email)
	if err != nil {
		return err
	}

	input.UserID = user.ID.Hex()
	return ProcessExpense(input)
}

func ProcessExpenseForTelegram(chatID int64, input ExpenseInput) error {
	user, err := FindUserByTelegramChatID(chatID)
	if err != nil {
		return err
	}

	input.UserID = user.ID.Hex()
	return ProcessExpense(input)
}

func FindUserByEmail(email string) (models.User, error) {
	var user models.User
	normalized := utils.NormalizeEmail(email)
	if normalized == "" {
		return user, errors.New("email is required")
	}

	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{"email": normalized}).
		Decode(&user)
	if err != nil {
		return user, errors.New("no user found for this email")
	}

	return user, nil
}

func FindUserByPhone(phone string) (models.User, error) {
	var user models.User
	normalized := NormalizePhone(phone)
	if normalized == "" {
		return user, errors.New("phone number is required")
	}

	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{
			"phone": bson.M{"$regex": normalized + "$"},
		}).
		Decode(&user)
	if err != nil {
		return user, errors.New("no user found for this phone number")
	}

	return user, nil
}

func FindUserByTelegramChatID(chatID int64) (models.User, error) {
	var user models.User
	if chatID == 0 {
		return user, errors.New("telegram chat id is required")
	}

	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{"telegram_chat_id": strconv.FormatInt(chatID, 10)}).
		Decode(&user)
	if err != nil {
		return user, errors.New("no user found for this telegram chat id")
	}

	return user, nil
}

func NormalizePhone(phone string) string {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	if len(digits) > 10 {
		digits = digits[len(digits)-10:]
	}
	return digits
}

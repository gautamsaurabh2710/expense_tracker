package services

import (
	"context"
	"errors"
	"strconv"

	"expense-tracker/config"
	"expense-tracker/repositories"

	"go.mongodb.org/mongo-driver/bson"
)

func LinkTelegramAccount(chatID int64, code string) error {
	user, err := repositories.FindUserByTelegramCode(code)

	if err != nil {
		return errors.New("invalid code")
	}

	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{
			"_id": user.ID,
		},
		bson.M{
			"$set": bson.M{
				"telegram_chat_id": strconv.FormatInt(chatID, 10),
				"track_telegram":   true,
			},

			"$unset": bson.M{
				"telegram_link_code": "",
			},
		},
	)

	return err
}

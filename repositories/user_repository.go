package repositories

import (
	"context"
	"expense-tracker/config"
	"expense-tracker/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveTelegramLinkCode(userID string, code string) error {

	log.Println("USER ID:", userID)
	log.Println("CODE:", code)

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return err
	}

	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"telegram_link_code": code,
			},
		},
	)

	log.Println("MATCHED:", result.MatchedCount)
	log.Println("MODIFIED:", result.ModifiedCount)
	log.Println("UPDATE ERROR:", err)

	return err
}

func FindUserByTelegramCode(code string) (models.User, error) {

	var user models.User

	err := config.DB.Collection("users").
		FindOne(
			context.Background(),
			bson.M{
				"telegram_link_code": code,
			},
		).
		Decode(&user)

	return user, err
}

package services

import (
	"context"
	"expense-tracker/config"
	"expense-tracker/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetMonthlySpent(userID string) (float64, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	collection := config.DB.Collection("expenses")

	pipeline := mongoPipeline(
		bson.D{{"$match", bson.M{
			"user_id": userID,
			"created_at": bson.M{
				"$gte": start,
			},
		}}},
		bson.D{{"$group", bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	)

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}

	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return utils.BsonFloat(result[0], "total"), nil
}

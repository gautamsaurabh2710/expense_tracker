package services

import (
	"context"
	"expense-tracker/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TopExpenses(userID string) ([]bson.M, error) {
	collection := config.DB.Collection("expenses")

	pipeline := mongoPipeline(
		bson.D{{"$match", bson.M{"user_id": userID}}},
		bson.D{{"$sort", bson.M{"amount": -1}}},
		bson.D{{"$limit", 5}},
	)

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func WeeklyExpense(userID string) ([]bson.M, error) {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)

	collection := config.DB.Collection("expenses")

	pipeline := mongoPipeline(
		bson.D{{"$match", bson.M{
			"user_id": userID,
			"created_at": bson.M{
				"$gte": weekAgo,
				"$lte": now,
			},
		}}},
		bson.D{{"$group", bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	)

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func MonthlyExpense(userID string) ([]bson.M, error) {
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
		return nil, err
	}

	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func CategoryWiseExpense(userID string) ([]bson.M, error) {
	collection := config.DB.Collection("expenses")

	pipeline := mongoPipeline(
		bson.D{{"$match", bson.M{"user_id": userID}}},
		bson.D{{"$group", bson.M{
			"_id":   "$category",
			"total": bson.M{"$sum": "$amount"},
		}}},
		bson.D{{"$sort", bson.M{"total": -1}}},
	)

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func PaymentMethodWise(userID string) ([]bson.M, error) {
	collection := config.DB.Collection("expenses")

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{
				"user_id": userID,
			}},
		},
		{
			{"$group", bson.M{
				"_id": "$payment_method",
				"total": bson.M{
					"$sum": "$amount",
				},
			}},
		},
	}

	cursor, err := collection.Aggregate(
		context.Background(),
		pipeline,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var results []bson.M

	if err := cursor.All(
		context.Background(),
		&results,
	); err != nil {
		return nil, err
	}

	return results, nil
}

func mongoPipeline(stages ...bson.D) mongo.Pipeline {
	var pipeline mongo.Pipeline

	for _, s := range stages {
		pipeline = append(pipeline, s)
	}

	return pipeline
}

package services

import (
	"context"
	"expense-tracker/config"
	"expense-tracker/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetBudget(userID string, limit float64) error {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	filter := bson.M{
		"user_id": userID,
		"month":   month,
		"year":    year,
	}

	update := bson.M{
		"$set": bson.M{
			"limit":      limit,
			"updated_at": time.Now(),
		},
		"$setOnInsert": bson.M{
			"user_id":    userID,
			"month":      month,
			"year":       year,
			"created_at": time.Now(),
		},
	}

	_, err := config.DB.Collection("budgets").UpdateOne(
		context.Background(),
		filter,
		update,
		options.Update().SetUpsert(true),
	)

	return err
}

func GetBudget(userID string) (models.Budget, error) {
	var budget models.Budget

	now := time.Now()

	err := config.DB.Collection("budgets").
		FindOne(context.Background(), bson.M{
			"user_id": userID,
			"month":   int(now.Month()),
			"year":    now.Year(),
		}).Decode(&budget)

	return budget, err
}

type BudgetSummary struct {
	Limit     float64 `json:"limit"`
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
	Month     int     `json:"month"`
	Year      int     `json:"year"`
}

func GetBudgetSummary(userID string) (BudgetSummary, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	summary := BudgetSummary{Month: int(now.Month()), Year: now.Year()}

	budget, err := GetBudget(userID)
	if err != nil && err != mongo.ErrNoDocuments {
		return summary, err
	}
	if err == nil {
		summary.Limit = budget.Limit
	}

	cursor, err := config.DB.Collection("expenses").Aggregate(context.Background(), mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id":    userID,
			"created_at": bson.M{"$gte": start},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return summary, err
	}
	defer cursor.Close(context.Background())

	var rows []bson.M
	if err := cursor.All(context.Background(), &rows); err != nil {
		return summary, err
	}
	if len(rows) > 0 {
		if total, ok := rows[0]["total"].(float64); ok {
			summary.Spent = total
		}
	}

	summary.Remaining = summary.Limit - summary.Spent
	return summary, nil
}

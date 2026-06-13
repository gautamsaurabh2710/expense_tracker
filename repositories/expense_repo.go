package repositories

import (
	"context"
	"expense-tracker/config"
	"expense-tracker/models"
	"regexp"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateExpense(exp models.Expense) error {
	result, err := config.DB.Collection("expenses").
		InsertOne(context.Background(), exp)

    log.Println("INSERT RESULT:", result.InsertedID)
	log.Println("INSERT ERROR:", err)

	return err
}

func GetExpenses(userID string) ([]models.Expense, error) {
	var expenses []models.Expense

	cursor, err := config.DB.Collection("expenses").
		Find(context.Background(), bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))

	if err != nil {
		return expenses, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var e models.Expense
		cursor.Decode(&e)
		expenses = append(expenses, e)
	}

	return expenses, nil

}

type ExpenseFilters struct {
	Category      string
	PaymentMethod string
	Source        string
	Search        string
	From          time.Time
	To            time.Time
	MinAmount     float64
	MaxAmount     float64
	Limit         int64
}

func ListExpenses(userID string, filters ExpenseFilters) ([]models.Expense, error) {
	query := bson.M{"user_id": userID}

	if filters.Category != "" {
		query["category"] = bson.M{"$regex": "^" + regexp.QuoteMeta(filters.Category) + "$", "$options": "i"}
	}
	if filters.PaymentMethod != "" {
		query["payment_method"] = bson.M{"$regex": "^" + regexp.QuoteMeta(filters.PaymentMethod) + "$", "$options": "i"}
	}
	if filters.Source != "" {
		query["source"] = bson.M{"$regex": "^" + regexp.QuoteMeta(filters.Source) + "$", "$options": "i"}
	}
	if filters.Search != "" {
		query["description"] = bson.M{"$regex": filters.Search, "$options": "i"}
	}
	if !filters.From.IsZero() || !filters.To.IsZero() {
		dateQuery := bson.M{}
		if !filters.From.IsZero() {
			dateQuery["$gte"] = filters.From
		}
		if !filters.To.IsZero() {
			dateQuery["$lte"] = filters.To
		}
		query["created_at"] = dateQuery
	}
	if filters.MinAmount > 0 || filters.MaxAmount > 0 {
		amountQuery := bson.M{}
		if filters.MinAmount > 0 {
			amountQuery["$gte"] = filters.MinAmount
		}
		if filters.MaxAmount > 0 {
			amountQuery["$lte"] = filters.MaxAmount
		}
		query["amount"] = amountQuery
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if filters.Limit > 0 {
		opts.SetLimit(filters.Limit)
	}

	cursor, err := config.DB.Collection("expenses").Find(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	expenses := []models.Expense{}
	if err := cursor.All(context.Background(), &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}

func DeleteExpense(userID, expenseID string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return false, err
	}

	res, err := config.DB.Collection("expenses").DeleteOne(
		context.Background(),
		bson.M{"_id": objectID, "user_id": userID},
	)
	if err != nil {
		return false, err
	}

	return res.DeletedCount > 0, nil
}

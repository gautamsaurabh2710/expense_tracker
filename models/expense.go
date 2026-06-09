package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID string `bson:"user_id" json:"-"`

	Amount float64 `bson:"amount" json:"amount"`

	Category string `bson:"category" json:"category"`

	PaymentMethod string `bson:"payment_method" json:"payment_method"`

	Description string `bson:"description" json:"description"`

	Source string `bson:"source" json:"source"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

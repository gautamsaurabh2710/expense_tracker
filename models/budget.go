package models

import "time"

type Budget struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Month     int       `bson:"month" json:"month"`
	Year      int       `bson:"year" json:"year"`
	Limit     float64   `bson:"limit" json:"limit"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

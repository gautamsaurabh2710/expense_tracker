package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Phone     string             `bson:"phone,omitempty" json:"phone"`
	Password  string             `bson:"password,omitempty"`
	Verified  bool               `bson:"verified" json:"verified"`
	LoginType string             `bson:"login_type" json:"login_type"`

	TrackEmail       bool   `bson:"track_email" json:"track_email"`
	TrackTelegram    bool   `bson:"track_telegram" json:"track_telegram"`
	TrackSMS         bool   `bson:"track_sms" json:"track_sms"`
	TelegramChatID   string `bson:"telegram_chat_id,omitempty" json:"telegram_chat_id"`
	TelegramLinkCode string `bson:"telegram_link_code,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
}

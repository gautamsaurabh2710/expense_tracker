package models

import "time"

type OTP struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	EmailOTP  string    `bson:"email_otp" json:"email_otp"`
	PhoneOTP  string    `bson:"phone_otp" json:"phone_otp"`
	EmailOK   bool      `bson:"email_ok" json:"email_ok"`
	PhoneOK   bool      `bson:"phone_ok" json:"phone_ok"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
}

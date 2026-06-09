package services

import (
	"context"
	"expense-tracker/config"
	"expense-tracker/models"
	"expense-tracker/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//REGISTER USER + SEND OTP

func RegisterUser(name, email, phone, password string, trackEmail, trackTelegram, trackSMS bool) (string, error) {
	normalizedEmail := utils.NormalizeEmail(email)

	var existing models.User
	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{"email": normalizedEmail}).
		Decode(&existing)
	if err == nil {
		return "", fmt.Errorf("email already registered")
	}
	if err != mongo.ErrNoDocuments {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	user := models.User{
		Name:          name,
		Email:         normalizedEmail,
		Phone:         phone,
		Password:      string(hashedPassword),
		Verified:      false,
		LoginType:     "email",
		TrackEmail:    trackEmail,
		TrackTelegram: trackTelegram,
		TrackSMS:      trackSMS,
		CreatedAt:     time.Now(),
	}

	res, err := config.DB.Collection("users").InsertOne(context.Background(), user)

	if err != nil {
		return "", err
	}

	userID := res.InsertedID.(primitive.ObjectID).Hex()

	//generate OTPs

	emailOTP := utils.GenerateOTP()

	otp := models.OTP{
		UserID: userID,

		EmailOTP:  emailOTP,
		PhoneOTP:  "",
		EmailOK:   false,
		PhoneOK:   true,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if _, err := config.DB.Collection("otps").InsertOne(context.Background(), otp); err != nil {
		return "", err
	}

	//SEND OTPs

	if err := utils.SendEmail(email, emailOTP); err != nil {
		return "", err
	}

	return userID, nil

}

//VERIFY OTP

func VerifyOTP(userID, emailOTP, phoneOTP string) error {

	var otp models.OTP

	err := config.DB.Collection("otps").
		FindOne(context.Background(), bson.M{"user_id": userID}).
		Decode(&otp)

	if err != nil {
		return err
	}

	if otp.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("OTP expires")
	}

	update := bson.M{}

	if otp.EmailOTP == emailOTP {
		update["email_ok"] = true
	}

	if len(update) == 0 {
		return fmt.Errorf("invalid OTP")
	}

	_, err = config.DB.Collection("otps").UpdateOne(
		context.Background(),
		bson.M{"user_id": userID},
		bson.M{"$set": update},
	)

	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(otp.UserID)
	if err != nil {
		return err
	}

	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"verified": true}},
	)
	if err != nil {
		return err
	}

	return nil

}

//LOGIN + JWT

func Login(email, password string) (string, error) {
	normalizedEmail := utils.NormalizeEmail(email)

	var user models.User
	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{"email": normalizedEmail}).
		Decode(&user)

	if err != nil {
		return "", fmt.Errorf("user not found")

	}

	if !user.Verified {
		return "", fmt.Errorf("user not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {

		return "", fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func Logout(_ string) error {
	return nil
}

func RequestPasswordReset(email string) error {
	normalizedEmail := utils.NormalizeEmail(email)

	var user models.User
	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{"email": normalizedEmail}).
		Decode(&user)
	if err != nil {
		return fmt.Errorf("if this email exists, a reset OTP has been sent")
	}

	otpCode := utils.GenerateOTP()
	otp := models.OTP{
		UserID:    user.ID.Hex(),
		EmailOTP:  otpCode,
		PhoneOTP:  "",
		EmailOK:   false,
		PhoneOK:   true,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	_, _ = config.DB.Collection("otps").DeleteMany(context.Background(), bson.M{"user_id": user.ID.Hex()})
	if _, err := config.DB.Collection("otps").InsertOne(context.Background(), otp); err != nil {
		return err
	}

	if err := utils.SendEmail(normalizedEmail, otpCode); err != nil {
		return err
	}

	return nil
}

func ResetPassword(email, emailOTP, password string) error {
	normalizedEmail := utils.NormalizeEmail(email)
	if password == "" {
		return fmt.Errorf("password is required")
	}

	var user models.User
	err := config.DB.Collection("users").
		FindOne(context.Background(), bson.M{"email": normalizedEmail}).
		Decode(&user)
	if err != nil {
		return fmt.Errorf("invalid reset request")
	}

	var otp models.OTP
	err = config.DB.Collection("otps").
		FindOne(context.Background(), bson.M{"user_id": user.ID.Hex()}).
		Decode(&otp)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}
	if otp.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("OTP expired")
	}
	if otp.EmailOTP != emailOTP {
		return fmt.Errorf("invalid OTP")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"password": string(hashedPassword), "verified": true}},
	)
	if err != nil {
		return err
	}

	_, _ = config.DB.Collection("otps").DeleteMany(context.Background(), bson.M{"user_id": user.ID.Hex()})
	return nil
}

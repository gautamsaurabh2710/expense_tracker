package repositories



func SaveTelegramLinkCode(userID string, code string) error {
	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return err
	}

	_, err = config.DBCollection("users").UpdateOne(
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
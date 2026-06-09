package services

func CheckBudget(userID string) error {

	// 1. Get budget
	budget, err := GetBudget(userID)
	if err != nil {
		return nil // no budget set
	}

	// 2. Get monthly spending
	spent, err := GetMonthlySpent(userID)
	if err != nil {
		return err
	}

	// 3. Compare
	if spent >= budget.Limit {

		// 4. Trigger alerts
		SendEmailAlert(userID, spent, budget.Limit)
		SendTelegramAlert(userID, spent, budget.Limit)
	}

	return nil
}

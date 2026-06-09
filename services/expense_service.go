package services

import (
	"expense-tracker/models"
	"expense-tracker/repositories"
	"expense-tracker/utils"
	"time"
)

func AddExpense(userID string, amount float64, description, source string) error {

	category := utils.GetCategory(description)

	expense := models.Expense{
		UserID:        userID,
		Amount:        amount,
		Category:      category,
		PaymentMethod: utils.DetectPaymentMethod(description),
		Description:   description,
		Source:        source,
		CreatedAt:     time.Now(),
	}

	return repositories.CreateExpense(expense)
}

type ExpenseFilters = repositories.ExpenseFilters

func ListExpenses(userID string, filters ExpenseFilters) ([]models.Expense, error) {
	return repositories.ListExpenses(userID, filters)
}

func DeleteExpense(userID, expenseID string) (bool, error) {
	return repositories.DeleteExpense(userID, expenseID)
}

package services

import (
	"expense-tracker/models"
	"expense-tracker/repositories"
	"expense-tracker/utils"
	"fmt"
	"log"
	"time"
)

type ExpenseInput struct {
	UserID        string
	Amount        float64
	Description   string
	Source        string
	RawText       string
	PaymentMethod string
}

func ProcessExpense(input ExpenseInput) error {

	log.Println("PROCESSING EXPENSE")
	log.Println("USER ID:", input.UserID)
	log.Println("AMOUNT:", input.Amount)

	if input.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	//CLEAN DATA
	cleanText := utils.CleanText(input.Description + " " + input.RawText)

	//CATEGORY DETECTION

	category := utils.GetCategory(cleanText)

	//BUILD EXPENSE OBJECT

	paymentMethod := utils.DetectPaymentMethod(cleanText)

	expense := models.Expense{
		UserID: input.UserID,

		Amount: input.Amount,

		Category: category,

		PaymentMethod: paymentMethod,

		Description: input.Description,

		Source: input.Source,

		CreatedAt: time.Now(),
	}

	log.Println("EXPENSE OBJECT:", expense)

	//SAVE TO DB
	if err := repositories.CreateExpense(expense); err != nil {
		return err
	}

	log.Println("EXPENSE SAVED SUCCESSFULLY")

	go CheckBudget(input.UserID)
	return nil
}

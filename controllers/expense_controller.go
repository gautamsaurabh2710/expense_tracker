package controllers

import (
	"expense-tracker/middleware"
	"expense-tracker/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ExpenseRequest struct {
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func AddExpense(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.ProcessExpense(services.ExpenseInput{
		UserID:      userID,
		Amount:      req.Amount,
		Description: req.Description,
		Source:      "web",
		RawText:     req.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense added"})
}

func ListExpenses(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filters := services.ExpenseFilters{
		Category:      c.Query("category"),
		PaymentMethod: c.Query("payment_method"),
		Source:        c.Query("source"),
		Search:        c.Query("search"),
		Limit:         100,
	}

	if value := c.Query("from"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
			return
		}
		filters.From = parsed
	}
	if value := c.Query("to"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
			return
		}
		filters.To = parsed.Add(24*time.Hour - time.Nanosecond)
	}
	if value := c.Query("min_amount"); value != "" {
		amount, err := strconv.ParseFloat(value, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid minimum amount"})
			return
		}
		filters.MinAmount = amount
	}
	if value := c.Query("max_amount"); value != "" {
		amount, err := strconv.ParseFloat(value, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maximum amount"})
			return
		}
		filters.MaxAmount = amount
	}
	if value := c.Query("limit"); value != "" {
		limit, err := strconv.ParseInt(value, 10, 64)
		if err != nil || limit < 1 || limit > 500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 500"})
			return
		}
		filters.Limit = limit
	}

	expenses, err := services.ListExpenses(userID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func DeleteExpense(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deleted, err := services.DeleteExpense(userID, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense deleted"})
}

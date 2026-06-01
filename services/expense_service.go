package services

import (
	"errors"
	"fmt"

	"github.com/SadikMR/go-expense-tracker-api/models"
	"github.com/SadikMR/go-expense-tracker-api/utils"
)

var (
	ErrExpenseNotFound = errors.New("expense not found")
	ErrForbidden       = errors.New("access denied")
)

// CreateExpenseInput holds the validated fields for creating an expense.
type CreateExpenseInput struct {
	Title       string
	Amount      float64
	Category    string
	Note        string
	ExpenseDate string
}

// UpdateExpenseInput holds the validated fields for updating an expense.
type UpdateExpenseInput struct {
	Title       string
	Amount      float64
	Category    string
	Note        string
	ExpenseDate string
}

// CreateExpense validates input and creates a new expense for the given user.
func CreateExpense(userID int, input CreateExpenseInput) (*models.Expense, error) {
	if !utils.ValidateRequired(input.Title, input.Category, input.ExpenseDate) {
		return nil, fmt.Errorf("CreateExpense: title, category, and expense_date are required")
	}
	if !utils.ValidatePositiveAmount(input.Amount) {
		return nil, fmt.Errorf("CreateExpense: amount must be greater than zero")
	}
	if !utils.ValidateCategory(input.Category) {
		return nil, fmt.Errorf("CreateExpense: invalid category %q", input.Category)
	}
	if !utils.ValidateDate(input.ExpenseDate) {
		return nil, fmt.Errorf("CreateExpense: expense_date must be in YYYY-MM-DD format")
	}

	id, err := models.NextExpenseID()
	if err != nil {
		return nil, fmt.Errorf("CreateExpense: %w", err)
	}

	expense := &models.Expense{
		ID:          id,
		UserID:      userID,
		Title:       input.Title,
		Amount:      input.Amount,
		Category:    input.Category,
		Note:        input.Note,
		ExpenseDate: input.ExpenseDate,
	}

	if err := models.CreateExpense(expense); err != nil {
		return nil, fmt.Errorf("CreateExpense: %w", err)
	}
	return expense, nil
}

// ListExpenses returns expenses for a user, capped by limit.
func ListExpenses(userID, limit int) ([]models.Expense, error) {
	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ListExpenses: %w", err)
	}

	if limit > 0 && limit < len(expenses) {
		return expenses[:limit], nil
	}
	return expenses, nil
}

// GetExpense returns a single expense owned by the user.
// Returns ErrExpenseNotFound if it does not exist.
func GetExpense(id, userID int) (*models.Expense, error) {
	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("GetExpense: %w", err)
	}
	if expense == nil {
		return nil, ErrExpenseNotFound
	}
	return expense, nil
}

// UpdateExpense validates input and updates an existing expense owned by the user.
func UpdateExpense(id, userID int, input UpdateExpenseInput) (*models.Expense, error) {
	expense, err := GetExpense(id, userID)
	if err != nil {
		return nil, err
	}

	if !utils.ValidateRequired(input.Title, input.Category, input.ExpenseDate) {
		return nil, fmt.Errorf("UpdateExpense: title, category, and expense_date are required")
	}
	if !utils.ValidatePositiveAmount(input.Amount) {
		return nil, fmt.Errorf("UpdateExpense: amount must be greater than zero")
	}
	if !utils.ValidateCategory(input.Category) {
		return nil, fmt.Errorf("UpdateExpense: invalid category %q", input.Category)
	}
	if !utils.ValidateDate(input.ExpenseDate) {
		return nil, fmt.Errorf("UpdateExpense: expense_date must be in YYYY-MM-DD format")
	}

	expense.Title = input.Title
	expense.Amount = input.Amount
	expense.Category = input.Category
	expense.Note = input.Note
	expense.ExpenseDate = input.ExpenseDate

	if err := models.UpdateExpense(expense); err != nil {
		return nil, fmt.Errorf("UpdateExpense: %w", err)
	}
	return expense, nil
}

// DeleteExpense removes an expense owned by the user.
func DeleteExpense(id, userID int) error {
	_, err := GetExpense(id, userID)
	if err != nil {
		return err
	}

	if err := models.DeleteExpense(id, userID); err != nil {
		return fmt.Errorf("DeleteExpense: %w", err)
	}
	return nil
}

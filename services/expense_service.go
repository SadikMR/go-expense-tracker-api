package services

import (
	"errors"
	"fmt"
	"sort"

	"github.com/SadikMR/go-expense-tracker-api/models"
	"github.com/SadikMR/go-expense-tracker-api/utils"
)

var (
	ErrExpenseNotFound = errors.New("expense not found")
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

// ListExpensesInput holds the optional filter and sort parameters for listing expenses.
type ListExpensesInput struct {
	DateFrom  string
	DateTo    string
	SortBy    string
	SortOrder string
	Limit     int
}

// filterByDate filters expenses by optional date range using direct YYYY-MM-DD string comparison.
// Reused by both ListExpenses and GetSummary — single source of truth for date filtering.
func filterByDate(expenses []models.Expense, dateFrom, dateTo string) []models.Expense {
	filtered := make([]models.Expense, 0, len(expenses))
	for _, e := range expenses {
		if dateFrom != "" && e.ExpenseDate < dateFrom {
			continue
		}
		if dateTo != "" && e.ExpenseDate > dateTo {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
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

// ListExpenses returns filtered and sorted expenses for a user.
func ListExpenses(userID int, input ListExpensesInput) ([]models.Expense, error) {
	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ListExpenses: %w", err)
	}

	filtered := filterByDate(expenses, input.DateFrom, input.DateTo)

	sortBy := input.SortBy
	sortOrder := input.SortOrder
	if sortBy == "" {
		sortBy = "expense_date"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sort.Slice(filtered, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "amount":
			less = filtered[i].Amount < filtered[j].Amount
		default:
			less = filtered[i].ExpenseDate < filtered[j].ExpenseDate
		}
		if sortOrder == "asc" {
			return less
		}
		return !less
	})

	if input.Limit > 0 && input.Limit < len(filtered) {
		return filtered[:input.Limit], nil
	}
	return filtered, nil
}

// GetExpense returns a single expense owned by the user.
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

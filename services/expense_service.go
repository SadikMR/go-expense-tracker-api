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

// CreateExpenseInput holds validated fields for creating an expense.
type CreateExpenseInput struct {
	Title       string
	Amount      float64
	Category    string
	Note        string
	ExpenseDate string
}

// UpdateExpenseInput holds validated fields for updating an expense.
// Fields are pointers so omitted values are not treated as zero values.
type UpdateExpenseInput struct {
	Title       *string
	Amount      *float64
	Category    *string
	Note        *string
	ExpenseDate *string
}

// ListExpensesInput holds optional filter, category, and sort parameters.
type ListExpensesInput struct {
	DateFrom  string
	DateTo    string
	Category  string
	SortBy    string
	SortOrder string
	Limit     int
}

// filterByDate filters expenses by optional date range.
// Reused by both ListExpenses and GetSummary — single source of truth.
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

// filterByCategory filters expenses by an optional category.
func filterByCategory(expenses []models.Expense, category string) []models.Expense {
	if category == "" {
		return expenses
	}

	filtered := make([]models.Expense, 0, len(expenses))
	for _, e := range expenses {
		if e.Category == category {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// validateExpenseInput validates expense fields and returns a clean user-facing error.
// Validation errors use errors.New (no wrap) so controllers can distinguish them
// from infrastructure errors (which use fmt.Errorf with %w).
func validateExpenseInput(title string, amount float64, category, expenseDate string) error {
	if !utils.ValidateRequired(title) {
		return errors.New("Title is required")
	}
	if !utils.ValidatePositiveAmount(amount) {
		return errors.New("Amount must be greater than zero")
	}
	if !utils.ValidateRequired(category) {
		return errors.New("Category is required")
	}
	if !utils.ValidateCategory(category) {
		return errors.New("Invalid category")
	}
	if !utils.ValidateRequired(expenseDate) {
		return errors.New("Expense date is required")
	}
	if !utils.ValidateDate(expenseDate) {
		return errors.New("Expense date must be in YYYY-MM-DD format")
	}
	return nil
}

func validateExpenseUpdateInput(input UpdateExpenseInput) error {
	if input.Title == nil && input.Amount == nil && input.Category == nil && input.Note == nil && input.ExpenseDate == nil {
		return errors.New("At least one field must be provided")
	}
	if input.Title != nil {
		if !utils.ValidateRequired(*input.Title) {
			return errors.New("Title is required")
		}
	}
	if input.Amount != nil {
		if !utils.ValidatePositiveAmount(*input.Amount) {
			return errors.New("Amount must be greater than zero")
		}
	}
	if input.Category != nil {
		if !utils.ValidateRequired(*input.Category) {
			return errors.New("Category is required")
		}
		if !utils.ValidateCategory(*input.Category) {
			return errors.New("Invalid category")
		}
	}
	if input.ExpenseDate != nil {
		if !utils.ValidateRequired(*input.ExpenseDate) {
			return errors.New("Expense date is required")
		}
		if !utils.ValidateDate(*input.ExpenseDate) {
			return errors.New("Expense date must be in YYYY-MM-DD format")
		}
	}
	return nil
}

// CreateExpense validates input and creates a new expense for the given user.
func CreateExpense(userID int, input CreateExpenseInput) (*models.Expense, error) {
	if err := validateExpenseInput(input.Title, input.Amount, input.Category, input.ExpenseDate); err != nil {
		return nil, err
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
	filtered = filterByCategory(filtered, input.Category)

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

	if err := validateExpenseUpdateInput(input); err != nil {
		return nil, err
	}

	if input.Title != nil {
		expense.Title = *input.Title
	}
	if input.Amount != nil {
		expense.Amount = *input.Amount
	}
	if input.Category != nil {
		expense.Category = *input.Category
	}
	if input.Note != nil {
		expense.Note = *input.Note
	}
	if input.ExpenseDate != nil {
		expense.ExpenseDate = *input.ExpenseDate
	}

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

package models

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SadikMR/go-expense-tracker-api/utils"
)

var expenseCSVHeader = []string{
	"id", "user_id", "title", "amount",
	"category", "note", "expense_date", "created_at",
}

func expensesCSVPath() string {
	if utils.AppConfig.ExpensesCSVPath != "" {
		return utils.AppConfig.ExpensesCSVPath
	}
	return filepath.Join(dataDir(), "expenses.csv")
}

// Expense represents a single expense record belonging to a user.
type Expense struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
	CreatedAt   string  `json:"created_at"`
}

// EnsureExpensesCSV creates the expenses CSV file with its header if it does not exist.
func EnsureExpensesCSV() error {
	return utils.EnsureCSVExists(expensesCSVPath(), expenseCSVHeader)
}

// GetExpensesByUserID returns all expenses belonging to a specific user.
func GetExpensesByUserID(userID int) ([]Expense, error) {
	rows, err := utils.ReadCSV(expensesCSVPath())
	if err != nil {
		return nil, fmt.Errorf("GetExpensesByUserID: %w", err)
	}

	expenses := make([]Expense, 0, len(rows))
	for _, row := range rows {
		e, err := rowToExpense(row)
		if err != nil {
			continue
		}
		if e.UserID == userID {
			expenses = append(expenses, e)
		}
	}
	return expenses, nil
}

// GetExpenseByID returns a single expense by ID scoped to the user.
// Returns nil, nil when not found.
func GetExpenseByID(id, userID int) (*Expense, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("GetExpenseByID: %w", err)
	}
	for _, e := range expenses {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, nil
}

// CreateExpense appends a new expense record to the CSV.
func CreateExpense(e *Expense) error {
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := utils.AppendCSV(expensesCSVPath(), expenseToRow(e)); err != nil {
		return fmt.Errorf("CreateExpense: %w", err)
	}
	return nil
}

// UpdateExpense rewrites the CSV replacing the matching expense row.
func UpdateExpense(updated *Expense) error {
	rows, err := utils.ReadCSV(expensesCSVPath())
	if err != nil {
		return fmt.Errorf("UpdateExpense: %w", err)
	}

	replaced := false
	for i, row := range rows {
		e, err := rowToExpense(row)
		if err != nil {
			continue
		}
		if e.ID == updated.ID && e.UserID == updated.UserID {
			rows[i] = expenseToRow(updated)
			replaced = true
			break
		}
	}

	if !replaced {
		return fmt.Errorf("UpdateExpense: expense id=%d not found", updated.ID)
	}

	return utils.RewriteCSV(expensesCSVPath(), expenseCSVHeader, rows)
}

// DeleteExpense rewrites the CSV excluding the matching expense row.
func DeleteExpense(id, userID int) error {
	rows, err := utils.ReadCSV(expensesCSVPath())
	if err != nil {
		return fmt.Errorf("DeleteExpense: %w", err)
	}

	filtered := make([][]string, 0, len(rows))
	deleted := false
	for _, row := range rows {
		e, err := rowToExpense(row)
		if err != nil {
			continue
		}
		if e.ID == id && e.UserID == userID {
			deleted = true
			continue // skip this row — effectively deletes it
		}
		filtered = append(filtered, row)
	}

	if !deleted {
		return fmt.Errorf("DeleteExpense: expense id=%d not found", id)
	}

	return utils.RewriteCSV(expensesCSVPath(), expenseCSVHeader, filtered)
}

// NextExpenseID returns the next available expense ID.
func NextExpenseID() (int, error) {
	rows, err := utils.ReadCSV(expensesCSVPath())
	if err != nil {
		return 0, fmt.Errorf("NextExpenseID: %w", err)
	}
	max := 0
	for _, row := range rows {
		e, err := rowToExpense(row)
		if err != nil {
			continue
		}
		if e.ID > max {
			max = e.ID
		}
	}
	return max + 1, nil
}

func rowToExpense(row []string) (Expense, error) {
	if len(row) < 8 {
		return Expense{}, fmt.Errorf("rowToExpense: expected 8 fields, got %d", len(row))
	}
	id, err := strconv.Atoi(row[0])
	if err != nil {
		return Expense{}, fmt.Errorf("rowToExpense: invalid id %q: %w", row[0], err)
	}
	userID, err := strconv.Atoi(row[1])
	if err != nil {
		return Expense{}, fmt.Errorf("rowToExpense: invalid user_id %q: %w", row[1], err)
	}
	amount, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return Expense{}, fmt.Errorf("rowToExpense: invalid amount %q: %w", row[3], err)
	}
	return Expense{
		ID:          id,
		UserID:      userID,
		Title:       row[2],
		Amount:      amount,
		Category:    row[4],
		Note:        row[5],
		ExpenseDate: row[6],
		CreatedAt:   row[7],
	}, nil
}

func expenseToRow(e *Expense) []string {
	return []string{
		strconv.Itoa(e.ID),
		strconv.Itoa(e.UserID),
		e.Title,
		strconv.FormatFloat(e.Amount, 'f', 2, 64),
		e.Category,
		e.Note,
		e.ExpenseDate,
		e.CreatedAt,
	}
}

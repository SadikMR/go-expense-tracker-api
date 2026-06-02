package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/utils"
)

func stringPtr(value string) *string {
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "data")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("failed to locate repo root from %q", dir)
	return ""
}

func expenseCSVPath(t *testing.T) string {
	return filepath.Join(repoRoot(t), "data", "expenses.csv")
}

func backupExpensesCSV(t *testing.T) ([]byte, bool) {
	t.Helper()
	data, err := os.ReadFile(expenseCSVPath(t))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}
		t.Fatalf("failed to back up expenses.csv: %v", err)
	}
	return data, true
}

func restoreExpensesCSV(t *testing.T, backup []byte, existed bool) {
	t.Helper()
	path := expenseCSVPath(t)
	if existed {
		if err := os.WriteFile(path, backup, 0644); err != nil {
			t.Fatalf("failed to restore expenses.csv: %v", err)
		}
		return
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove expenses.csv: %v", err)
	}
}

func writeExpensesCSV(t *testing.T, rows [][]string) {
	t.Helper()
	head := []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
	if err := utils.RewriteCSV(expenseCSVPath(t), head, rows); err != nil {
		t.Fatalf("failed to write expenses.csv: %v", err)
	}
}

func TestValidateExpenseUpdateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   UpdateExpenseInput
		wantErr string
	}{
		{
			name:    "no fields provided",
			input:   UpdateExpenseInput{},
			wantErr: "At least one field must be provided",
		},
		{
			name:    "title empty",
			input:   UpdateExpenseInput{Title: stringPtr("")},
			wantErr: "Title is required",
		},
		{
			name:    "invalid amount",
			input:   UpdateExpenseInput{Amount: floatPtr(0)},
			wantErr: "Amount must be greater than zero",
		},
		{
			name:    "invalid category",
			input:   UpdateExpenseInput{Category: stringPtr("Bad")},
			wantErr: "Invalid category",
		},
		{
			name:    "invalid expense date",
			input:   UpdateExpenseInput{ExpenseDate: stringPtr("2025/06/01")},
			wantErr: "Expense date must be in YYYY-MM-DD format",
		},
		{
			name:  "valid partial title",
			input: UpdateExpenseInput{Title: stringPtr("Updated")},
		},
		{
			name:  "valid amount and date",
			input: UpdateExpenseInput{Amount: floatPtr(250), ExpenseDate: stringPtr("2025-06-15")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExpenseUpdateInput(tt.input)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestUpdateExpensePartialFields(t *testing.T) {
	backup, existed := backupExpensesCSV(t)
	t.Cleanup(func() {
		restoreExpensesCSV(t, backup, existed)
	})

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("failed to change working directory to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	writeExpensesCSV(t, [][]string{
		{"1", "100", "Old Title", "100.00", "Food", "Old note", "2025-06-01", "2025-01-01T00:00:00Z"},
	})

	expense, err := UpdateExpense(1, 100, UpdateExpenseInput{Title: stringPtr("New Title")})
	if err != nil {
		t.Fatalf("unexpected error updating expense title: %v", err)
	}
	if expense.Title != "New Title" {
		t.Fatalf("expected title to be updated, got %q", expense.Title)
	}
	if expense.Amount != 100 {
		t.Fatalf("expected amount to remain unchanged, got %f", expense.Amount)
	}
	if expense.Category != "Food" {
		t.Fatalf("expected category to remain unchanged, got %q", expense.Category)
	}

	expense, err = UpdateExpense(1, 100, UpdateExpenseInput{Amount: floatPtr(250)})
	if err != nil {
		t.Fatalf("unexpected error updating expense amount: %v", err)
	}
	if expense.Amount != 250 {
		t.Fatalf("expected amount to be updated, got %f", expense.Amount)
	}
	if expense.Title != "New Title" {
		t.Fatalf("expected title to remain updated, got %q", expense.Title)
	}

	err = func() error {
		_, err := UpdateExpense(1, 100, UpdateExpenseInput{Category: stringPtr("Bad")})
		return err
	}()
	if err == nil || err.Error() != "Invalid category" {
		t.Fatalf("expected invalid category error, got %v", err)
	}
}

func TestCreateExpenseSuccess(t *testing.T) {
	backup, existed := backupExpensesCSV(t)
	t.Cleanup(func() {
		restoreExpensesCSV(t, backup, existed)
	})

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("failed to change working directory to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	writeExpensesCSV(t, [][]string{})

	input := CreateExpenseInput{
		Title:       "Lunch",
		Amount:      12.50,
		Category:    "Food",
		Note:        "Office lunch",
		ExpenseDate: "2025-06-01",
	}

	expense, err := CreateExpense(100, input)
	if err != nil {
		t.Fatalf("CreateExpense failed: %v", err)
	}

	if expense.Title != "Lunch" {
		t.Fatalf("expected title Lunch, got %q", expense.Title)
	}
	if expense.Amount != 12.50 {
		t.Fatalf("expected amount 12.50, got %v", expense.Amount)
	}
	if expense.UserID != 100 {
		t.Fatalf("expected userID 100, got %d", expense.UserID)
	}
	if expense.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestCreateExpenseValidationErrors(t *testing.T) {
	backup, existed := backupExpensesCSV(t)
	t.Cleanup(func() {
		restoreExpensesCSV(t, backup, existed)
	})

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("failed to change working directory to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	writeExpensesCSV(t, [][]string{})

	tests := []struct {
		name    string
		input   CreateExpenseInput
		wantErr string
	}{
		{
			name: "empty title",
			input: CreateExpenseInput{
				Title:       "",
				Amount:      10.0,
				Category:    "Food",
				ExpenseDate: "2025-06-01",
			},
			wantErr: "Title is required",
		},
		{
			name: "zero amount",
			input: CreateExpenseInput{
				Title:       "Test",
				Amount:      0,
				Category:    "Food",
				ExpenseDate: "2025-06-01",
			},
			wantErr: "Amount must be greater than zero",
		},
		{
			name: "negative amount",
			input: CreateExpenseInput{
				Title:       "Test",
				Amount:      -5.0,
				Category:    "Food",
				ExpenseDate: "2025-06-01",
			},
			wantErr: "Amount must be greater than zero",
		},
		{
			name: "empty category",
			input: CreateExpenseInput{
				Title:       "Test",
				Amount:      10.0,
				Category:    "",
				ExpenseDate: "2025-06-01",
			},
			wantErr: "Category is required",
		},
		{
			name: "invalid category",
			input: CreateExpenseInput{
				Title:       "Test",
				Amount:      10.0,
				Category:    "InvalidCat",
				ExpenseDate: "2025-06-01",
			},
			wantErr: "Invalid category",
		},
		{
			name: "empty date",
			input: CreateExpenseInput{
				Title:       "Test",
				Amount:      10.0,
				Category:    "Food",
				ExpenseDate: "",
			},
			wantErr: "Expense date is required",
		},
		{
			name: "invalid date format",
			input: CreateExpenseInput{
				Title:       "Test",
				Amount:      10.0,
				Category:    "Food",
				ExpenseDate: "06/01/2025",
			},
			wantErr: "Expense date must be in YYYY-MM-DD format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateExpense(100, tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestListExpensesService(t *testing.T) {
	backup, existed := backupExpensesCSV(t)
	t.Cleanup(func() {
		restoreExpensesCSV(t, backup, existed)
	})

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("failed to change working directory to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	writeExpensesCSV(t, [][]string{
		{"1", "100", "Lunch", "10.00", "Food", "", "2025-06-01", "2025-01-01T00:00:00Z"},
		{"2", "100", "Uber", "20.00", "Transport", "", "2025-06-02", "2025-01-02T00:00:00Z"},
		{"3", "200", "Movie", "15.00", "Entertainment", "", "2025-06-01", "2025-01-03T00:00:00Z"},
	})

	tests := []struct {
		name      string
		userID    int
		input     ListExpensesInput
		wantCount int
		validate  func(*testing.T, []interface{})
	}{
		{
			name:      "list all user expenses",
			userID:    100,
			input:     ListExpensesInput{},
			wantCount: 2,
		},
		{
			name:   "filter by date from",
			userID: 100,
			input: ListExpensesInput{
				DateFrom: "2025-06-02",
			},
			wantCount: 1,
		},
		{
			name:   "filter by date range",
			userID: 100,
			input: ListExpensesInput{
				DateFrom: "2025-06-01",
				DateTo:   "2025-06-01",
			},
			wantCount: 1,
		},
		{
			name:   "sort by amount desc",
			userID: 100,
			input: ListExpensesInput{
				SortBy:    "amount",
				SortOrder: "desc",
			},
			wantCount: 2,
		},
		{
			name:   "limit results",
			userID: 100,
			input: ListExpensesInput{
				Limit: 1,
			},
			wantCount: 1,
		},
		{
			name:      "other user expenses not included",
			userID:    200,
			input:     ListExpensesInput{},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenses, err := ListExpenses(tt.userID, tt.input)
			if err != nil {
				t.Fatalf("ListExpenses failed: %v", err)
			}
			if len(expenses) != tt.wantCount {
				t.Fatalf("expected %d expenses, got %d", tt.wantCount, len(expenses))
			}
		})
	}
}

func TestDeleteExpenseService(t *testing.T) {
	backup, existed := backupExpensesCSV(t)
	t.Cleanup(func() {
		restoreExpensesCSV(t, backup, existed)
	})

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repoRoot(t)); err != nil {
		t.Fatalf("failed to change working directory to repo root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	writeExpensesCSV(t, [][]string{
		{"1", "100", "Lunch", "10.00", "Food", "", "2025-06-01", "2025-01-01T00:00:00Z"},
		{"2", "100", "Uber", "20.00", "Transport", "", "2025-06-02", "2025-01-02T00:00:00Z"},
	})

	tests := []struct {
		name      string
		id        int
		userID    int
		wantError bool
		validate  func(*testing.T)
	}{
		{
			name:      "delete existing expense",
			id:        1,
			userID:    100,
			wantError: false,
			validate: func(t *testing.T) {
				expenses, err := ListExpenses(100, ListExpensesInput{})
				if err != nil {
					t.Fatalf("verify list failed: %v", err)
				}
				if len(expenses) != 1 {
					t.Fatalf("expected 1 expense after delete, got %d", len(expenses))
				}
			},
		},
		{
			name:      "delete non-existent expense",
			id:        999,
			userID:    100,
			wantError: true,
		},
		{
			name:      "delete other user expense",
			id:        1,
			userID:    200,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteExpense(tt.id, tt.userID)
			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantError && tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

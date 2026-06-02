package models

import (
	"os"
	"path/filepath"
	"testing"
)

func expenseFilePath(t *testing.T) string {
	return filepath.Join("data", "expenses.csv")
}

func writeExpensesCSV(t *testing.T, path string, rows [][]string) {
	t.Helper()
	head := []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
	testWriteCSV(t, path, head, rows)
}

func TestNextExpenseID(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	path := expenseFilePath(t)
	backup, existed := testBackupFile(t, path)
	t.Cleanup(func() {
		testRestoreFile(t, path, backup, existed)
	})

	writeExpensesCSV(t, path, [][]string{
		{"10", "100", "Lunch", "12.50", "Food", "note", "2025-06-01", "2025-01-01T00:00:00Z"},
	})

	id, err := NextExpenseID()
	if err != nil {
		t.Fatalf("NextExpenseID failed: %v", err)
	}
	if id != 11 {
		t.Fatalf("expected next expense ID 11, got %d", id)
	}
}

func TestGetExpenseByID(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	path := expenseFilePath(t)
	backup, existed := testBackupFile(t, path)
	t.Cleanup(func() {
		testRestoreFile(t, path, backup, existed)
	})

	writeExpensesCSV(t, path, [][]string{
		{"1", "100", "Dinner", "20.00", "Food", "evening meal", "2025-06-02", "2025-01-02T00:00:00Z"},
	})

	expense, err := GetExpenseByID(1, 100)
	if err != nil {
		t.Fatalf("GetExpenseByID failed: %v", err)
	}
	if expense == nil {
		t.Fatal("expected expense, got nil")
	}
	if expense.Title != "Dinner" {
		t.Fatalf("expected title Dinner, got %q", expense.Title)
	}

	nullExpense, err := GetExpenseByID(2, 100)
	if err != nil {
		t.Fatalf("GetExpenseByID failed: %v", err)
	}
	if nullExpense != nil {
		t.Fatal("expected nil expense for missing ID")
	}
}

func TestCreateExpenseAppends(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	path := expenseFilePath(t)
	backup, existed := testBackupFile(t, path)
	t.Cleanup(func() {
		testRestoreFile(t, path, backup, existed)
	})

	writeExpensesCSV(t, path, [][]string{
		{"1", "100", "Dinner", "20.00", "Food", "evening meal", "2025-06-02", "2025-01-02T00:00:00Z"},
	})

	exp := &Expense{ID: 2, UserID: 100, Title: "Breakfast", Amount: 8.75, Category: "Food", Note: "morning", ExpenseDate: "2025-06-03"}
	if err := CreateExpense(exp); err != nil {
		t.Fatalf("CreateExpense failed: %v", err)
	}

	created, err := GetExpenseByID(2, 100)
	if err != nil {
		t.Fatalf("GetExpenseByID failed: %v", err)
	}
	if created == nil || created.Title != "Breakfast" {
		t.Fatalf("expected created expense Breakfast, got %+v", created)
	}
}

func TestUpdateExpense(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

	tests := []struct {
		name      string
		initial   *Expense
		update    *Expense
		wantError bool
		validate  func(*testing.T, *Expense)
	}{
		{
			name: "update existing expense",
			initial: &Expense{
				ID:          1,
				UserID:      100,
				Title:       "Old",
				Amount:      10.0,
				Category:    "Food",
				ExpenseDate: "2025-06-01",
			},
			update: &Expense{
				ID:          1,
				UserID:      100,
				Title:       "New",
				Amount:      20.0,
				Category:    "Transport",
				ExpenseDate: "2025-06-02",
			},
			validate: func(t *testing.T, e *Expense) {
				if e.Title != "New" {
					t.Fatalf("expected title New, got %q", e.Title)
				}
				if e.Amount != 20.0 {
					t.Fatalf("expected amount 20.0, got %v", e.Amount)
				}
			},
		},
		{
			name: "update non-existent expense",
			initial: &Expense{
				ID:          1,
				UserID:      100,
				Title:       "Test",
				Amount:      10.0,
				Category:    "Food",
				ExpenseDate: "2025-06-01",
			},
			update: &Expense{
				ID:          999,
				UserID:      100,
				Title:       "New",
				Amount:      20.0,
				Category:    "Transport",
				ExpenseDate: "2025-06-02",
			},
			wantError: true,
		},
		{
			name: "update different user expense",
			initial: &Expense{
				ID:          1,
				UserID:      100,
				Title:       "Old",
				Amount:      10.0,
				Category:    "Food",
				ExpenseDate: "2025-06-01",
			},
			update: &Expense{
				ID:          1,
				UserID:      200,
				Title:       "New",
				Amount:      20.0,
				Category:    "Transport",
				ExpenseDate: "2025-06-02",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "data/expenses.csv"
			testBackupFile(t, path)
			head := []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
			testWriteCSV(t, path, head, [][]string{
				{"1", "100", tt.initial.Title, "10.00", tt.initial.Category, "", "2025-06-01", "2025-01-01T00:00:00Z"},
			})

			err := UpdateExpense(tt.update)
			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.wantError && tt.validate != nil {
				exp, err := GetExpenseByID(tt.update.ID, tt.update.UserID)
				if err != nil {
					t.Fatalf("verify read failed: %v", err)
				}
				if exp == nil {
					t.Fatal("expense not found after update")
				}
				tt.validate(t, exp)
			}
		})
	}
}

func TestDeleteExpense(t *testing.T) {
	cleanup := testEnterTempDir(t)
	t.Cleanup(cleanup)

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
				exp, err := GetExpenseByID(1, 100)
				if err != nil {
					t.Fatalf("verify read failed: %v", err)
				}
				if exp != nil {
					t.Fatal("expected nil after delete")
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
			name:      "delete different user expense",
			id:        1,
			userID:    200,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "data/expenses.csv"
			testBackupFile(t, path)
			head := []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
			testWriteCSV(t, path, head, [][]string{
				{"1", "100", "Test", "10.00", "Food", "", "2025-06-01", "2025-01-01T00:00:00Z"},
			})

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

func TestEnsureExpensesCSV(t *testing.T) {
	tests := []struct {
		name       string
		priorExist bool
		wantError  bool
	}{
		{
			name:       "create when not exists",
			priorExist: false,
			wantError:  false,
		},
		{
			name:       "skip when exists",
			priorExist: true,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := testEnterTempDir(t)
			t.Cleanup(cleanup)

			path := "data/expenses.csv"
			if tt.priorExist {
				if err := os.Mkdir("data", 0755); err != nil && !os.IsExist(err) {
					t.Fatalf("mkdir failed: %v", err)
				}
				if err := RewriteCSV(path, []string{"id"}, [][]string{}); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := EnsureExpensesCSV()
			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if _, err := os.Stat(path); err != nil {
				t.Fatalf("csv file not created: %v", err)
			}
		})
	}
}

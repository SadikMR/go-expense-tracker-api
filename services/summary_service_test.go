package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SadikMR/go-expense-tracker-api/utils"
)

func enterSummaryTempDir(t *testing.T) func() {
	t.Helper()
	tmp, err := os.MkdirTemp("", "summary-service-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tmp, "data"), 0755); err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("failed to create temp data dir: %v", err)
	}
	orig, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("failed to change working directory to temp dir: %v", err)
	}
	return func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
		if err := os.RemoveAll(tmp); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}
}

func writeSummaryExpensesCSV(t *testing.T, rows [][]string) {
	t.Helper()
	head := []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
	if err := utils.RewriteCSV("data/expenses.csv", head, rows); err != nil {
		t.Fatalf("failed to write summary expenses csv: %v", err)
	}
}

func TestGetSummaryTotalsAndCategories(t *testing.T) {
	cleanup := enterSummaryTempDir(t)
	t.Cleanup(cleanup)

	writeSummaryExpensesCSV(t, [][]string{
		{"1", "100", "Lunch", "10.00", "Food", "", "2025-06-01", "2025-01-01T00:00:00Z"},
		{"2", "100", "Uber", "20.00", "Transport", "", "2025-06-01", "2025-01-01T00:00:00Z"},
		{"3", "100", "Coffee", "5.00", "Food", "", "2025-06-02", "2025-01-02T00:00:00Z"},
		{"4", "200", "Movie", "15.00", "Entertainment", "", "2025-06-01", "2025-01-01T00:00:00Z"},
	})

	result, err := GetSummary(100, "", "")
	if err != nil {
		t.Fatalf("GetSummary failed: %v", err)
	}
	if result.TotalAmount != 35 {
		t.Fatalf("expected total amount 35, got %v", result.TotalAmount)
	}
	if result.TotalCount != 3 {
		t.Fatalf("expected total count 3, got %d", result.TotalCount)
	}

	categoryTotals := map[string]CategorySummary{}
	for _, summary := range result.ByCategory {
		categoryTotals[summary.Category] = summary
	}

	food, ok := categoryTotals["Food"]
	if !ok || food.Total != 15 || food.Count != 2 {
		t.Fatalf("expected Food summary 15/2, got %+v", food)
	}
	transport, ok := categoryTotals["Transport"]
	if !ok || transport.Total != 20 || transport.Count != 1 {
		t.Fatalf("expected Transport summary 20/1, got %+v", transport)
	}
}

func TestGetSummaryDateFiltering(t *testing.T) {
	cleanup := enterSummaryTempDir(t)
	t.Cleanup(cleanup)

	writeSummaryExpensesCSV(t, [][]string{
		{"1", "100", "Lunch", "10.00", "Food", "", "2025-06-01", "2025-01-01T00:00:00Z"},
		{"2", "100", "Coffee", "5.00", "Food", "", "2025-06-02", "2025-01-02T00:00:00Z"},
		{"3", "100", "Burger", "12.00", "Food", "", "2025-06-03", "2025-01-03T00:00:00Z"},
	})

	result, err := GetSummary(100, "2025-06-02", "2025-06-02")
	if err != nil {
		t.Fatalf("GetSummary failed: %v", err)
	}
	if result.TotalAmount != 5 {
		t.Fatalf("expected total amount 5, got %v", result.TotalAmount)
	}
	if result.TotalCount != 1 {
		t.Fatalf("expected total count 1, got %d", result.TotalCount)
	}
	if len(result.ByCategory) != 1 || result.ByCategory[0].Category != "Food" {
		t.Fatalf("expected single Food category result, got %+v", result.ByCategory)
	}
}

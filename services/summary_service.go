package services

import (
	"fmt"

	"github.com/SadikMR/go-expense-tracker-api/models"
)

// CategorySummary holds the aggregated total and count for a single category.
type CategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// SummaryResult holds the full expense summary for a user.
type SummaryResult struct {
	TotalAmount float64           `json:"total_amount"`
	TotalCount  int               `json:"total_count"`
	DateFrom    string            `json:"date_from,omitempty"`
	DateTo      string            `json:"date_to,omitempty"`
	ByCategory  []CategorySummary `json:"by_category"`
}

// GetSummary returns a summary of expenses for a user with optional date filtering.
// Reuses filterByDate from expense_service — single source of truth for date logic.
func GetSummary(userID int, dateFrom, dateTo string) (*SummaryResult, error) {
	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("GetSummary: %w", err)
	}

	filtered := filterByDate(expenses, dateFrom, dateTo)

	// Aggregate totals and counts per category using maps
	categoryTotals := make(map[string]float64)
	categoryCounts := make(map[string]int)

	result := &SummaryResult{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	for _, e := range filtered {
		result.TotalAmount += e.Amount
		result.TotalCount++
		categoryTotals[e.Category] += e.Amount
		categoryCounts[e.Category]++
	}

	// Build by_category slice from maps
	byCategory := make([]CategorySummary, 0, len(categoryTotals))
	for cat, total := range categoryTotals {
		byCategory = append(byCategory, CategorySummary{
			Category: cat,
			Total:    total,
			Count:    categoryCounts[cat],
		})
	}
	result.ByCategory = byCategory

	return result, nil
}

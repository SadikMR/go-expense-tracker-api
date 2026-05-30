package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCSV reads all data rows from a CSV file, skipping the header row.
// Returns an empty slice (not an error) when the file has no data rows.
func ReadCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadCSV: open %s: %w", path, err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ReadCSV: parse %s: %w", path, err)
	}
	if len(rows) <= 1 {
		return [][]string{}, nil
	}
	return rows[1:], nil
}

// AppendCSV appends a single row to an existing CSV file.
func AppendCSV(path string, row []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("AppendCSV: open %s: %w", path, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(row); err != nil {
		return fmt.Errorf("AppendCSV: write %s: %w", path, err)
	}
	return nil
}

// RewriteCSV overwrites an entire CSV file with the given header and rows.
// This is the correct pattern for update and delete operations on CSV files.
func RewriteCSV(path string, header []string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("RewriteCSV: create %s: %w", path, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(header); err != nil {
		return fmt.Errorf("RewriteCSV: write header %s: %w", path, err)
	}
	if err := w.WriteAll(rows); err != nil {
		return fmt.Errorf("RewriteCSV: write rows %s: %w", path, err)
	}
	return nil
}

// EnsureCSVExists creates the CSV file with its header row if it does not exist.
// Safe to call on every startup.
func EnsureCSVExists(path string, header []string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return RewriteCSV(path, header, [][]string{})
	}
	return nil
}
